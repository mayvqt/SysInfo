//go:build darwin

package collector

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/mayvqt/sysinfo/internal/utils"
)

// diskutilListOutput represents diskutil list -plist output structure
type diskutilListOutput struct {
	AllDisks              []string `plist:"AllDisks"`
	AllDisksAndPartitions []struct {
		Content    string `plist:"Content"`
		DeviceID   string `plist:"DeviceIdentifier"`
		Size       uint64 `plist:"Size"`
		Partitions []struct {
			Content  string `plist:"Content"`
			DeviceID string `plist:"DeviceIdentifier"`
			Size     uint64 `plist:"Size"`
		} `plist:"Partitions,omitempty"`
	} `plist:"AllDisksAndPartitions"`
	VolumesFromDisks []string `plist:"VolumesFromDisks"`
	WholeDisks       []string `plist:"WholeDisks"`
}

// diskutilInfoOutput represents diskutil info -plist output structure
type diskutilInfoOutput struct {
	DeviceIdentifier   string `json:"DeviceIdentifier"`
	MediaName          string `json:"MediaName"`
	MediaType          string `json:"MediaType"`
	SolidState         bool   `json:"SolidState"`
	VirtualOrPhysical  string `json:"VirtualOrPhysical"`
	TotalSize          uint64 `json:"TotalSize"`
	DeviceBlockSize    uint64 `json:"DeviceBlockSize"`
	Protocol           string `json:"Protocol"` // SATA, NVMe, USB, etc.
	Removable          bool   `json:"Removable"`
	Internal           bool   `json:"Internal"`
	SMARTStatus        string `json:"SMARTStatus"`
}

func collectPhysicalDisksPlatform() []types.PhysicalDisk {
	disks := make([]types.PhysicalDisk, 0)

	// Check if diskutil is available
	if _, err := exec.LookPath("diskutil"); err != nil {
		return disks
	}

	// Get list of whole disks
	cmd := exec.Command("diskutil", "list", "-plist")
	output, err := cmd.Output()
	if err != nil {
		return disks
	}

	// Parse the plist output
	// Note: macOS uses plist format, but we'll use diskutil info with JSON for each disk
	// First, get the list of physical disks
	wholeDisks := getWholeDisks()

	for _, diskID := range wholeDisks {
		disk := getDiskInfo(diskID)
		if disk.Name != "" {
			disks = append(disks, disk)
		}
	}

	return disks
}

// getWholeDisks returns a list of whole disk identifiers
func getWholeDisks() []string {
	cmd := exec.Command("diskutil", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	disks := make([]string, 0)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		// Look for lines like "/dev/disk0 (internal, physical):"
		if strings.Contains(line, "/dev/disk") && (strings.Contains(line, "physical") || strings.Contains(line, "external")) {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				diskPath := strings.TrimSpace(parts[0])
				// Extract just the disk identifier (e.g., "disk0" from "/dev/disk0")
				diskID := strings.TrimPrefix(diskPath, "/dev/")
				disks = append(disks, diskID)
			}
		}
	}

	return disks
}

// getDiskInfo retrieves detailed information about a specific disk
func getDiskInfo(diskID string) types.PhysicalDisk {
	disk := types.PhysicalDisk{
		Name: "/dev/" + diskID,
	}

	// Use diskutil info with JSON output (available on recent macOS)
	cmd := exec.Command("diskutil", "info", "-plist", diskID)
	output, err := cmd.Output()
	if err != nil {
		return disk
	}

	// Parse plist output - we'll use a simpler text parsing approach
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Parse key-value pairs from plist
		if strings.Contains(line, "<key>MediaName</key>") && i+1 < len(lines) {
			modelLine := strings.TrimSpace(lines[i+1])
			disk.Model = extractPlistString(modelLine)
		} else if strings.Contains(line, "<key>TotalSize</key>") && i+1 < len(lines) {
			sizeLine := strings.TrimSpace(lines[i+1])
			if size := extractPlistInteger(sizeLine); size > 0 {
				disk.Size = size
				disk.SizeFormatted = utils.FormatBytes(size)
			}
		} else if strings.Contains(line, "<key>SolidState</key>") && i+1 < len(lines) {
			if strings.Contains(lines[i+1], "<true/>") {
				disk.Type = "SSD"
			} else {
				disk.Type = "HDD"
			}
		} else if strings.Contains(line, "<key>Protocol</key>") && i+1 < len(lines) {
			protocolLine := strings.TrimSpace(lines[i+1])
			protocol := extractPlistString(protocolLine)
			disk.Interface = strings.ToUpper(protocol)

			// Detect NVMe
			if strings.Contains(strings.ToUpper(protocol), "NVME") || strings.Contains(strings.ToUpper(protocol), "PCI") {
				disk.Type = "NVMe"
				disk.Interface = "NVMe"
			}
		} else if strings.Contains(line, "<key>Removable</key>") && i+1 < len(lines) {
			disk.Removable = strings.Contains(lines[i+1], "<true/>")
		} else if strings.Contains(line, "<key>Device / Media Name</key>") && i+1 < len(lines) {
			modelLine := strings.TrimSpace(lines[i+1])
			if disk.Model == "" {
				disk.Model = extractPlistString(modelLine)
			}
		}
	}

	// If type is still unknown, try to detect from disk name
	if disk.Type == "" {
		if strings.HasPrefix(diskID, "disk") {
			// Check if it's an SSD via system_profiler (slower but more accurate)
			if isSSD := checkSSDSystemProfiler(diskID); isSSD {
				disk.Type = "SSD"
			} else {
				disk.Type = "HDD"
			}
		}
	}

	return disk
}

// extractPlistString extracts a string value from a plist XML line
func extractPlistString(line string) string {
	// Line looks like: <string>APPLE SSD AP0512R</string>
	line = strings.TrimPrefix(line, "<string>")
	line = strings.TrimSuffix(line, "</string>")
	return strings.TrimSpace(line)
}

// extractPlistInteger extracts an integer value from a plist XML line
func extractPlistInteger(line string) uint64 {
	// Line looks like: <integer>512110190592</integer>
	line = strings.TrimPrefix(line, "<integer>")
	line = strings.TrimSuffix(line, "</integer>")
	line = strings.TrimSpace(line)
	if val, err := strconv.ParseUint(line, 10, 64); err == nil {
		return val
	}
	return 0
}

// checkSSDSystemProfiler uses system_profiler to check if a disk is SSD (fallback)
func checkSSDSystemProfiler(diskID string) bool {
	cmd := exec.Command("system_profiler", "SPSerialATADataType", "-json")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Simple check - SSDs usually have "Solid State" in the output
	return strings.Contains(string(output), "Solid State")
}
