//go:build linux

package collector

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/mayvqt/sysinfo/internal/utils"
)

// lsblkDevice represents the structure of lsblk JSON output
type lsblkDevice struct {
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Size      string        `json:"size"`
	Model     string        `json:"model"`
	Serial    string        `json:"serial"`
	Rota      string        `json:"rota"` // 1 for HDD, 0 for SSD
	Removable string        `json:"rm"`   // 1 for removable
	Tran      string        `json:"tran"` // transport (sata, nvme, usb, etc.)
	Children  []lsblkDevice `json:"children,omitempty"`
}

type lsblkOutput struct {
	BlockDevices []lsblkDevice `json:"blockdevices"`
}

func collectPhysicalDisksPlatform() []types.PhysicalDisk {
	// Try lsblk first (most reliable)
	lsblkDisks := collectDisksLsblk()
	if len(lsblkDisks) > 0 {
		return lsblkDisks
	}

	// Fallback to /sys/block parsing
	return collectDisksSysBlock()
}

// collectDisksLsblk uses lsblk to get physical disk information
func collectDisksLsblk() []types.PhysicalDisk {
	// Check if lsblk is available
	if _, err := exec.LookPath("lsblk"); err != nil {
		return nil
	}

	// Run lsblk with JSON output
	cmd := exec.Command("lsblk", "-J", "-b", "-d", "-o", "NAME,TYPE,SIZE,MODEL,SERIAL,ROTA,RM,TRAN")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var lsblkOut lsblkOutput
	if err := json.Unmarshal(output, &lsblkOut); err != nil {
		return nil
	}

	disks := make([]types.PhysicalDisk, 0)
	for _, device := range lsblkOut.BlockDevices {
		// Only process disk devices (not partitions, loops, etc.)
		if device.Type != "disk" {
			continue
		}

		disk := types.PhysicalDisk{
			Name:  "/dev/" + device.Name,
			Model: strings.TrimSpace(device.Model),
		}

		// Parse size
		if size, err := strconv.ParseUint(device.Size, 10, 64); err == nil {
			disk.Size = size
			disk.SizeFormatted = utils.FormatBytes(size)
		}

		// Set serial number
		disk.SerialNumber = strings.TrimSpace(device.Serial)

		// Determine disk type (HDD vs SSD)
		switch device.Rota {
		case "1":
			disk.Type = "HDD"
			// Try to get RPM from sysfs
			if rpm := getDiskRPM(device.Name); rpm > 0 {
				disk.RPM = rpm
			}
		case "0":
			// Check if it's NVMe
			if strings.HasPrefix(device.Name, "nvme") {
				disk.Type = "NVMe"
			} else {
				disk.Type = "SSD"
			}
		}

		// Set interface type
		disk.Interface = strings.ToUpper(device.Tran)
		if disk.Interface == "" {
			// Try to detect from device name
			if strings.HasPrefix(device.Name, "nvme") {
				disk.Interface = "NVMe"
			} else if strings.HasPrefix(device.Name, "sd") {
				disk.Interface = "SATA"
			} else if strings.HasPrefix(device.Name, "hd") {
				disk.Interface = "IDE"
			}
		}

		// Set removable flag
		disk.Removable = device.Removable == "1"

		disks = append(disks, disk)
	}

	return disks
}

// collectDisksSysBlock parses /sys/block for disk information (fallback)
func collectDisksSysBlock() []types.PhysicalDisk {
	disks := make([]types.PhysicalDisk, 0)

	sysBlockPath := "/sys/block"
	entries, err := os.ReadDir(sysBlockPath)
	if err != nil {
		return disks
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip loop devices, ram disks, etc.
		if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") || strings.HasPrefix(name, "dm-") {
			continue
		}

		devicePath := filepath.Join(sysBlockPath, name)

		// Check if it's a physical device (has a device directory)
		if _, err := os.Stat(filepath.Join(devicePath, "device")); os.IsNotExist(err) {
			continue
		}

		disk := types.PhysicalDisk{
			Name: "/dev/" + name,
		}

		// Read size
		if sizeStr, err := os.ReadFile(filepath.Join(devicePath, "size")); err == nil {
			// Size is in 512-byte sectors
			if sectors, err := strconv.ParseUint(strings.TrimSpace(string(sizeStr)), 10, 64); err == nil {
				disk.Size = sectors * 512
				disk.SizeFormatted = utils.FormatBytes(disk.Size)
			}
		}

		// Read model
		if modelBytes, err := os.ReadFile(filepath.Join(devicePath, "device", "model")); err == nil {
			disk.Model = strings.TrimSpace(string(modelBytes))
		}

		// Read rotational flag
		if rotaBytes, err := os.ReadFile(filepath.Join(devicePath, "queue", "rotational")); err == nil {
			if strings.TrimSpace(string(rotaBytes)) == "1" {
				disk.Type = "HDD"
				if rpm := getDiskRPM(name); rpm > 0 {
					disk.RPM = rpm
				}
			} else {
				if strings.HasPrefix(name, "nvme") {
					disk.Type = "NVMe"
				} else {
					disk.Type = "SSD"
				}
			}
		}

		// Read removable flag
		if rmBytes, err := os.ReadFile(filepath.Join(devicePath, "removable")); err == nil {
			disk.Removable = strings.TrimSpace(string(rmBytes)) == "1"
		}

		// Detect interface
		if strings.HasPrefix(name, "nvme") {
			disk.Interface = "NVMe"
		} else if strings.HasPrefix(name, "sd") {
			disk.Interface = "SATA"
		} else if strings.HasPrefix(name, "hd") {
			disk.Interface = "IDE"
		} else if strings.HasPrefix(name, "mmc") {
			disk.Interface = "MMC"
		}

		disks = append(disks, disk)
	}

	return disks
}

// getDiskRPM attempts to get disk RPM (for HDDs)
func getDiskRPM(deviceName string) uint32 {
	// Try smartctl if available
	if _, err := exec.LookPath("smartctl"); err == nil {
		cmd := exec.Command("smartctl", "-i", "/dev/"+deviceName)
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Rotation Rate:") {
					parts := strings.Fields(line)
					for i, part := range parts {
						if part == "rpm" && i > 0 {
							if rpm, err := strconv.ParseUint(parts[i-1], 10, 32); err == nil {
								return uint32(rpm)
							}
						}
					}
				}
			}
		}
	}

	// Common HDD speeds as fallback
	// If we can't determine, return 0 (unknown)
	return 0
}
