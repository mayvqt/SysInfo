//go:build linux
// +build linux

package collector

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/mayvqt/sysinfo/internal/utils"
)

// NvidiaSMILog represents the XML output from nvidia-smi
type NvidiaSMILog struct {
	GPUs []NvidiaGPU `xml:"gpu"`
}

type NvidiaGPU struct {
	ProductName string `xml:"product_name"`
	UUID        string `xml:"uuid"`
	PCIBus      string `xml:"pci>pci_bus"`
	Temperature struct {
		Current string `xml:"gpu_temp"`
	} `xml:"temperature"`
	Utilization struct {
		GPU    string `xml:"gpu_util"`
		Memory string `xml:"memory_util"`
	} `xml:"utilization"`
	FBMemory struct {
		Total string `xml:"total"`
		Used  string `xml:"used"`
		Free  string `xml:"free"`
	} `xml:"fb_memory_usage"`
	Power struct {
		Draw  string `xml:"power_draw"`
		Limit string `xml:"power_limit"`
	} `xml:"power_readings"`
	Clocks struct {
		Graphics string `xml:"graphics_clock"`
		Memory   string `xml:"mem_clock"`
	} `xml:"clocks"`
	FanSpeed      string `xml:"fan_speed"`
	DriverVersion string `xml:"driver_version"`
}

// collectGPUPlatform implements Linux-specific GPU data collection
func collectGPUPlatform() []types.GPUInfo {
	gpus := make([]types.GPUInfo, 0)

	// Try NVIDIA GPUs first (nvidia-smi)
	nvidiaGPUs := collectNvidiaGPUs()
	gpus = append(gpus, nvidiaGPUs...)

	// Try AMD GPUs (rocm-smi or lspci)
	amdGPUs := collectAMDGPUs()
	gpus = append(gpus, amdGPUs...)

	// Fallback to lspci for basic info if nothing else worked
	if len(gpus) == 0 {
		gpus = collectGPUsFromLspci()
	}

	return gpus
}

// collectNvidiaGPUs collects NVIDIA GPU information using nvidia-smi
func collectNvidiaGPUs() []types.GPUInfo {
	gpus := make([]types.GPUInfo, 0)

	// Check if nvidia-smi is available
	_, err := exec.LookPath("nvidia-smi")
	if err != nil {
		return gpus
	}

	// Try XML format first (more detailed)
	cmd := exec.Command("nvidia-smi", "-q", "-x")
	output, err := cmd.Output()
	if err == nil {
		var smiLog NvidiaSMILog
		if err := xml.Unmarshal(output, &smiLog); err == nil {
			for i, gpu := range smiLog.GPUs {
				gpuInfo := types.GPUInfo{
					Index:         i,
					Name:          gpu.ProductName,
					Vendor:        "NVIDIA",
					Driver:        "nvidia",
					DriverVersion: gpu.DriverVersion,
					UUID:          gpu.UUID,
					PCIBus:        gpu.PCIBus,
				}

				// Parse temperature
				if temp, err := strconv.Atoi(strings.TrimSpace(strings.Replace(gpu.Temperature.Current, "C", "", -1))); err == nil {
					gpuInfo.Temperature = temp
				}

				// Parse utilization
				if util, err := strconv.Atoi(strings.TrimSpace(strings.Replace(gpu.Utilization.GPU, "%", "", -1))); err == nil {
					gpuInfo.Utilization = util
				}
				if memUtil, err := strconv.Atoi(strings.TrimSpace(strings.Replace(gpu.Utilization.Memory, "%", "", -1))); err == nil {
					gpuInfo.MemoryUtilization = memUtil
				}

				// Parse memory (format: "12345 MiB")
				if total := parseMemoryMiB(gpu.FBMemory.Total); total > 0 {
					gpuInfo.MemoryTotal = total
					gpuInfo.MemoryFormatted = utils.FormatBytes(total)
				}
				if used := parseMemoryMiB(gpu.FBMemory.Used); used > 0 {
					gpuInfo.MemoryUsed = used
				}
				if free := parseMemoryMiB(gpu.FBMemory.Free); free > 0 {
					gpuInfo.MemoryFree = free
				}

				// Parse power (format: "123.45 W")
				if power := parsePowerWatts(gpu.Power.Draw); power > 0 {
					gpuInfo.PowerDraw = power
				}
				if limit := parsePowerWatts(gpu.Power.Limit); limit > 0 {
					gpuInfo.PowerLimit = limit
				}

				// Parse clocks (format: "1234 MHz")
				if clock := parseClockMHz(gpu.Clocks.Graphics); clock > 0 {
					gpuInfo.ClockSpeed = clock
				}
				if memClock := parseClockMHz(gpu.Clocks.Memory); memClock > 0 {
					gpuInfo.ClockSpeedMemory = memClock
				}

				// Parse fan speed
				if fan, err := strconv.Atoi(strings.TrimSpace(strings.Replace(gpu.FanSpeed, "%", "", -1))); err == nil {
					gpuInfo.FanSpeed = fan
				}

				gpus = append(gpus, gpuInfo)
			}
			return gpus
		}
	}

	// Fallback to CSV format for basic info
	cmd = exec.Command("nvidia-smi", "--query-gpu=index,name,temperature.gpu,utilization.gpu,memory.total,memory.used",
		"--format=csv,noheader,nounits")
	output, err = cmd.Output()
	if err != nil {
		return gpus
	}

	reader := csv.NewReader(strings.NewReader(string(output)))
	records, err := reader.ReadAll()
	if err != nil {
		return gpus
	}

	for _, record := range records {
		if len(record) < 6 {
			continue
		}

		gpuInfo := types.GPUInfo{
			Vendor: "NVIDIA",
			Driver: "nvidia",
		}

		if idx, err := strconv.Atoi(strings.TrimSpace(record[0])); err == nil {
			gpuInfo.Index = idx
		}
		gpuInfo.Name = strings.TrimSpace(record[1])
		if temp, err := strconv.Atoi(strings.TrimSpace(record[2])); err == nil {
			gpuInfo.Temperature = temp
		}
		if util, err := strconv.Atoi(strings.TrimSpace(record[3])); err == nil {
			gpuInfo.Utilization = util
		}
		if total, err := strconv.ParseUint(strings.TrimSpace(record[4]), 10, 64); err == nil {
			gpuInfo.MemoryTotal = total * 1024 * 1024 // Convert MiB to bytes
			gpuInfo.MemoryFormatted = utils.FormatBytes(gpuInfo.MemoryTotal)
		}
		if used, err := strconv.ParseUint(strings.TrimSpace(record[5]), 10, 64); err == nil {
			gpuInfo.MemoryUsed = used * 1024 * 1024 // Convert MiB to bytes
		}

		gpus = append(gpus, gpuInfo)
	}

	return gpus
}

// collectAMDGPUs collects AMD GPU information using rocm-smi
func collectAMDGPUs() []types.GPUInfo {
	gpus := make([]types.GPUInfo, 0)

	// Check if rocm-smi is available
	_, err := exec.LookPath("rocm-smi")
	if err != nil {
		return gpus
	}

	// rocm-smi --showproductname --showtemp --showuse --showmeminfo vram
	cmd := exec.Command("rocm-smi", "--showproductname", "--showtemp", "--showuse", "--showmeminfo", "vram", "--csv")
	output, err := cmd.Output()
	if err != nil {
		return gpus
	}

	// Parse CSV output
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) < 4 {
			continue
		}

		gpuInfo := types.GPUInfo{
			Index:  i - 1,
			Vendor: "AMD",
			Driver: "amdgpu",
		}

		// Parse fields (format varies, this is approximate)
		for _, field := range fields {
			field = strings.TrimSpace(field)
			if strings.Contains(strings.ToLower(field), "radeon") || strings.Contains(strings.ToLower(field), "rx") {
				gpuInfo.Name = field
			}
		}

		gpus = append(gpus, gpuInfo)
	}

	return gpus
}

// collectGPUsFromLspci uses lspci as a fallback for basic GPU info
func collectGPUsFromLspci() []types.GPUInfo {
	gpus := make([]types.GPUInfo, 0)

	cmd := exec.Command("lspci")
	output, err := cmd.Output()
	if err != nil {
		return gpus
	}

	lines := strings.Split(string(output), "\n")
	index := 0
	for _, line := range lines {
		lineLower := strings.ToLower(line)
		if strings.Contains(lineLower, "vga") || strings.Contains(lineLower, "3d controller") {
			gpuInfo := types.GPUInfo{
				Index: index,
			}

			// Extract PCI bus
			if parts := strings.Split(line, " "); len(parts) > 0 {
				gpuInfo.PCIBus = parts[0]
			}

			// Determine vendor and name
			if strings.Contains(lineLower, "nvidia") {
				gpuInfo.Vendor = "NVIDIA"
				gpuInfo.Name = extractGPUName(line, "NVIDIA")
			} else if strings.Contains(lineLower, "amd") || strings.Contains(lineLower, "ati") {
				gpuInfo.Vendor = "AMD"
				gpuInfo.Name = extractGPUName(line, "AMD")
			} else if strings.Contains(lineLower, "intel") {
				gpuInfo.Vendor = "Intel"
				gpuInfo.Name = extractGPUName(line, "Intel")
			} else {
				gpuInfo.Name = line
			}

			gpus = append(gpus, gpuInfo)
			index++
		}
	}

	return gpus
}

// Helper functions

func parseMemoryMiB(memStr string) uint64 {
	memStr = strings.TrimSpace(memStr)
	memStr = strings.Replace(memStr, "MiB", "", -1)
	memStr = strings.TrimSpace(memStr)
	if val, err := strconv.ParseUint(memStr, 10, 64); err == nil {
		return val * 1024 * 1024 // Convert MiB to bytes
	}
	return 0
}

func parsePowerWatts(powerStr string) float64 {
	powerStr = strings.TrimSpace(powerStr)
	powerStr = strings.Replace(powerStr, "W", "", -1)
	powerStr = strings.TrimSpace(powerStr)
	if val, err := strconv.ParseFloat(powerStr, 64); err == nil {
		return val
	}
	return 0
}

func parseClockMHz(clockStr string) int {
	clockStr = strings.TrimSpace(clockStr)
	clockStr = strings.Replace(clockStr, "MHz", "", -1)
	clockStr = strings.TrimSpace(clockStr)
	if val, err := strconv.Atoi(clockStr); err == nil {
		return val
	}
	return 0
}

func extractGPUName(line, vendor string) string {
	// Try to extract GPU name from lspci output
	if idx := strings.Index(line, ":"); idx > 0 {
		name := line[idx+1:]
		name = strings.TrimSpace(name)
		// Remove common prefixes
		name = strings.Replace(name, "VGA compatible controller:", "", 1)
		name = strings.Replace(name, "3D controller:", "", 1)
		name = strings.TrimSpace(name)
		return name
	}
	return fmt.Sprintf("%s GPU", vendor)
}
