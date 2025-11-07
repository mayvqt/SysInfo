//go:build darwin
// +build darwin

package collector

import (
	"encoding/xml"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
)

// SystemProfilerDisplays represents the XML output from system_profiler SPDisplaysDataType
type SystemProfilerDisplays struct {
	Displays []Display `xml:"array>dict"`
}

type Display struct {
	Items []DisplayItem `xml:"array>dict"`
}

type DisplayItem struct {
	Name          string `xml:"key>string"`
	ChipsetModel  string `xml:"dict>key[.='sppci_model']>string"`
	VRAMTotal     string `xml:"dict>key[.='sppci_vram']>string"`
	Vendor        string `xml:"dict>key[.='sppci_vendor']>string"`
	DeviceID      string `xml:"dict>key[.='sppci_device_id']>string"`
	BusType       string `xml:"dict>key[.='sppci_bus']>string"`
}

// collectGPUPlatform implements macOS-specific GPU data collection
func collectGPUPlatform() []types.GPUInfo {
	gpus := make([]types.GPUInfo, 0)

	// Use system_profiler to get display/GPU information
	cmd := exec.Command("system_profiler", "SPDisplaysDataType", "-xml")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to text parsing
		return collectGPUsFromSystemProfilerText()
	}

	// Parse XML output
	var displays SystemProfilerDisplays
	if err := xml.Unmarshal(output, &displays); err != nil {
		return collectGPUsFromSystemProfilerText()
	}

	index := 0
	for _, display := range displays.Displays {
		for _, item := range display.Items {
			gpuInfo := types.GPUInfo{
				Index: index,
			}

			// Parse chipset model
			if item.ChipsetModel != "" {
				gpuInfo.Name = item.ChipsetModel
			}

			// Determine vendor
			vendor := strings.ToLower(item.Vendor)
			if strings.Contains(vendor, "nvidia") {
				gpuInfo.Vendor = "NVIDIA"
			} else if strings.Contains(vendor, "amd") || strings.Contains(vendor, "ati") {
				gpuInfo.Vendor = "AMD"
			} else if strings.Contains(vendor, "intel") {
				gpuInfo.Vendor = "Intel"
			} else if strings.Contains(vendor, "apple") {
				gpuInfo.Vendor = "Apple"
			} else {
				gpuInfo.Vendor = item.Vendor
			}

			// Parse VRAM (format: "8 GB" or "8192 MB")
			if item.VRAMTotal != "" {
				gpuInfo.MemoryTotal = parseVRAM(item.VRAMTotal)
				if gpuInfo.MemoryTotal > 0 {
					gpuInfo.MemoryFormatted = item.VRAMTotal
				}
			}

			// Device ID as UUID substitute
			if item.DeviceID != "" {
				gpuInfo.UUID = item.DeviceID
			}

			// Try to get Metal information for Apple Silicon
			if gpuInfo.Vendor == "Apple" {
				enrichAppleSiliconGPU(&gpuInfo)
			}

			gpus = append(gpus, gpuInfo)
			index++
		}
	}

	return gpus
}

// collectGPUsFromSystemProfilerText parses text output as fallback
func collectGPUsFromSystemProfilerText() []types.GPUInfo {
	gpus := make([]types.GPUInfo, 0)

	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	output, err := cmd.Output()
	if err != nil {
		return gpus
	}

	lines := strings.Split(string(output), "\n")
	var currentGPU *types.GPUInfo
	index := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "Chipset Model:") {
			// New GPU found
			if currentGPU != nil {
				gpus = append(gpus, *currentGPU)
			}
			currentGPU = &types.GPUInfo{
				Index: index,
			}
			index++

			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				currentGPU.Name = strings.TrimSpace(parts[1])
				
				// Determine vendor from name
				nameLower := strings.ToLower(currentGPU.Name)
				if strings.Contains(nameLower, "nvidia") {
					currentGPU.Vendor = "NVIDIA"
				} else if strings.Contains(nameLower, "amd") || strings.Contains(nameLower, "radeon") {
					currentGPU.Vendor = "AMD"
				} else if strings.Contains(nameLower, "intel") {
					currentGPU.Vendor = "Intel"
				} else if strings.Contains(nameLower, "apple") {
					currentGPU.Vendor = "Apple"
				}
			}
		} else if currentGPU != nil {
			if strings.Contains(line, "VRAM") || strings.Contains(line, "Memory") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					vramStr := strings.TrimSpace(parts[1])
					currentGPU.MemoryTotal = parseVRAM(vramStr)
					if currentGPU.MemoryTotal > 0 {
						currentGPU.MemoryFormatted = vramStr
					}
				}
			} else if strings.Contains(line, "Vendor:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 && currentGPU.Vendor == "" {
					currentGPU.Vendor = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// Add last GPU
	if currentGPU != nil {
		gpus = append(gpus, *currentGPU)
	}

	return gpus
}

// enrichAppleSiliconGPU adds Apple Silicon specific information
func enrichAppleSiliconGPU(gpu *types.GPUInfo) {
	// For Apple Silicon, we can try to get more info from Metal
	// This would require additional commands or APIs
	// For now, just set driver to Metal
	gpu.Driver = "Metal"
}

// parseVRAM converts VRAM string to bytes
func parseVRAM(vramStr string) uint64 {
	vramStr = strings.TrimSpace(vramStr)
	vramStr = strings.ToUpper(vramStr)

	// Remove common suffixes to get the number
	var multiplier uint64 = 1
	if strings.Contains(vramStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		vramStr = strings.Replace(vramStr, "GB", "", -1)
	} else if strings.Contains(vramStr, "MB") {
		multiplier = 1024 * 1024
		vramStr = strings.Replace(vramStr, "MB", "", -1)
	}

	vramStr = strings.TrimSpace(vramStr)
	if val, err := strconv.ParseUint(vramStr, 10, 64); err == nil {
		return val * multiplier
	}

	return 0
}
