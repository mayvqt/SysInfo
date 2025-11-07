//go:build windows
// +build windows

package collector

import (
	"encoding/csv"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/mayvqt/sysinfo/internal/utils"
	"github.com/yusufpapurcu/wmi"
)

// Win32_VideoController represents Windows WMI video controller
type Win32_VideoController struct {
	Name              string
	AdapterRAM        uint32
	DriverVersion     string
	VideoProcessor    string
	PNPDeviceID       string
	CurrentRefreshRate uint32
	VideoModeDescription string
	Status            string
}

// collectGPUPlatform implements Windows-specific GPU data collection
func collectGPUPlatform() []types.GPUInfo {
	gpus := make([]types.GPUInfo, 0)

	// Query WMI for video controller information
	var videoControllers []Win32_VideoController
	query := "SELECT Name, AdapterRAM, DriverVersion, VideoProcessor, PNPDeviceID, CurrentRefreshRate, VideoModeDescription, Status FROM Win32_VideoController"
	
	err := wmi.Query(query, &videoControllers)
	if err != nil {
		return gpus
	}

	for i, controller := range videoControllers {
		gpuInfo := types.GPUInfo{
			Index:         i,
			Name:          controller.Name,
			DriverVersion: controller.DriverVersion,
		}

		// Determine vendor from name or PNP device ID
		nameLower := strings.ToLower(controller.Name)
		pnpLower := strings.ToLower(controller.PNPDeviceID)
		
		if strings.Contains(nameLower, "nvidia") || strings.Contains(pnpLower, "nvidia") {
			gpuInfo.Vendor = "NVIDIA"
			gpuInfo.Driver = "nvidia"
		} else if strings.Contains(nameLower, "amd") || strings.Contains(nameLower, "radeon") || strings.Contains(pnpLower, "amd") {
			gpuInfo.Vendor = "AMD"
			gpuInfo.Driver = "amdgpu"
		} else if strings.Contains(nameLower, "intel") || strings.Contains(pnpLower, "intel") {
			gpuInfo.Vendor = "Intel"
			gpuInfo.Driver = "intel"
		} else if strings.Contains(nameLower, "microsoft") {
			gpuInfo.Vendor = "Microsoft"
			gpuInfo.Driver = "wddm"
		} else {
			gpuInfo.Vendor = "Unknown"
		}

		// Memory (AdapterRAM is in bytes)
		if controller.AdapterRAM > 0 {
			gpuInfo.MemoryTotal = uint64(controller.AdapterRAM)
			gpuInfo.MemoryFormatted = utils.FormatBytes(uint64(controller.AdapterRAM))
		}

		// Use PNP Device ID as UUID
		if controller.PNPDeviceID != "" {
			gpuInfo.UUID = controller.PNPDeviceID
		}

		// Extract PCI bus from PNP Device ID if available
		// Format: PCI\VEN_10DE&DEV_1234&SUBSYS_12345678&REV_A1\4&12345678&0&00E0
		if strings.HasPrefix(strings.ToUpper(controller.PNPDeviceID), "PCI\\") {
			parts := strings.Split(controller.PNPDeviceID, "\\")
			if len(parts) >= 2 {
				gpuInfo.PCIBus = parts[1]
			}
		}

		gpus = append(gpus, gpuInfo)
	}

	// Enrich NVIDIA GPUs with nvidia-smi if available (for advanced metrics)
	enrichNvidiaGPUsWindows(gpus)

	return gpus
}

// enrichNvidiaGPUsWindows uses nvidia-smi to get additional information for NVIDIA GPUs
func enrichNvidiaGPUsWindows(gpus []types.GPUInfo) {
	// Check if nvidia-smi is available (usually in C:\Program Files\NVIDIA Corporation\NVSMI\)
	_, err := exec.LookPath("nvidia-smi")
	if err != nil {
		// Try common installation path
		cmd := exec.Command("C:\\Program Files\\NVIDIA Corporation\\NVSMI\\nvidia-smi.exe", "--help")
		if err := cmd.Run(); err != nil {
			// nvidia-smi not available
			return
		}
	}

	// Get detailed NVIDIA GPU information using CSV format
	cmd := exec.Command("nvidia-smi", 
		"--query-gpu=index,name,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.used,memory.free,power.draw,power.limit,clocks.gr,clocks.mem,fan.speed,uuid",
		"--format=csv,noheader,nounits")
	
	output, err := cmd.Output()
	if err != nil {
		return
	}

	reader := csv.NewReader(strings.NewReader(string(output)))
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	// Create a map of NVIDIA GPUs by name for matching
	nvidiaGPUMap := make(map[string]*types.GPUInfo)
	for i := range gpus {
		if gpus[i].Vendor == "NVIDIA" {
			nvidiaGPUMap[gpus[i].Name] = &gpus[i]
		}
	}

	// Enrich with nvidia-smi data
	for _, record := range records {
		if len(record) < 14 {
			continue
		}

		name := strings.TrimSpace(record[1])
		
		// Find matching GPU in our list
		gpu, exists := nvidiaGPUMap[name]
		if !exists {
			// Try to match by index if name doesn't match exactly
			if idx, err := strconv.Atoi(strings.TrimSpace(record[0])); err == nil && idx < len(gpus) {
				if gpus[idx].Vendor == "NVIDIA" {
					gpu = &gpus[idx]
				}
			}
		}

		if gpu == nil {
			continue
		}

		// Parse temperature
		if temp, err := strconv.Atoi(strings.TrimSpace(record[2])); err == nil {
			gpu.Temperature = temp
		}

		// Parse GPU utilization
		if util, err := strconv.Atoi(strings.TrimSpace(record[3])); err == nil {
			gpu.Utilization = util
		}

		// Parse memory utilization
		if memUtil, err := strconv.Atoi(strings.TrimSpace(record[4])); err == nil {
			gpu.MemoryUtilization = memUtil
		}

		// Parse memory (convert MiB to bytes)
		if total, err := strconv.ParseUint(strings.TrimSpace(record[5]), 10, 64); err == nil {
			memBytes := total * 1024 * 1024
			gpu.MemoryTotal = memBytes
			gpu.MemoryFormatted = utils.FormatBytes(memBytes)
		}
		if used, err := strconv.ParseUint(strings.TrimSpace(record[6]), 10, 64); err == nil {
			gpu.MemoryUsed = used * 1024 * 1024
		}
		if free, err := strconv.ParseUint(strings.TrimSpace(record[7]), 10, 64); err == nil {
			gpu.MemoryFree = free * 1024 * 1024
		}

		// Parse power draw
		powerStr := strings.TrimSpace(record[8])
		if powerStr != "[N/A]" && powerStr != "" {
			if power, err := strconv.ParseFloat(powerStr, 64); err == nil {
				gpu.PowerDraw = power
			}
		}

		// Parse power limit
		limitStr := strings.TrimSpace(record[9])
		if limitStr != "[N/A]" && limitStr != "" {
			if limit, err := strconv.ParseFloat(limitStr, 64); err == nil {
				gpu.PowerLimit = limit
			}
		}

		// Parse clock speeds
		if clock, err := strconv.Atoi(strings.TrimSpace(record[10])); err == nil {
			gpu.ClockSpeed = clock
		}
		if memClock, err := strconv.Atoi(strings.TrimSpace(record[11])); err == nil {
			gpu.ClockSpeedMemory = memClock
		}

		// Parse fan speed
		fanStr := strings.TrimSpace(record[12])
		if fanStr != "[N/A]" && fanStr != "" {
			if fan, err := strconv.Atoi(fanStr); err == nil {
				gpu.FanSpeed = fan
			}
		}

		// Update UUID if available
		uuid := strings.TrimSpace(record[13])
		if uuid != "" && uuid != "[N/A]" {
			gpu.UUID = uuid
		}
	}
}

// Additional WMI queries for more detailed GPU info could include:
// - Win32_TemperatureProbe for temperature
// - Win32_PerfFormattedData_GPUPerformanceCounters for utilization
// - MSAcpi_ThermalZoneTemperature for thermal info
