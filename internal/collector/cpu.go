package collector

import (
	"fmt"
	"time"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
)

// CollectCPU gathers CPU information
func CollectCPU() (*types.CPUData, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}

	if len(cpuInfo) == 0 {
		return nil, fmt.Errorf("no CPU information available")
	}

	cores, err := cpu.Counts(false)
	if err != nil {
		cores = 0
	}

	logicalCPUs, err := cpu.Counts(true)
	if err != nil {
		logicalCPUs = 0
	}

	// Get CPU usage per core
	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		percentages = []float64{}
	}

	data := &types.CPUData{
		ModelName:   cpuInfo[0].ModelName,
		Cores:       int32(cores),
		LogicalCPUs: int32(logicalCPUs),
		Vendor:      cpuInfo[0].VendorID,
		Family:      cpuInfo[0].Family,
		Model:       cpuInfo[0].Model,
		Stepping:    cpuInfo[0].Stepping,
		MHz:         cpuInfo[0].Mhz,
		CacheSize:   cpuInfo[0].CacheSize,
		Usage:       percentages,
		Flags:       cpuInfo[0].Flags,
		Microcode:   cpuInfo[0].Microcode,
	}

	// Get load average (Unix-like systems)
	loadAvg, err := load.Avg()
	if err == nil {
		data.LoadAvg = &types.LoadAverage{
			Load1:  loadAvg.Load1,
			Load5:  loadAvg.Load5,
			Load15: loadAvg.Load15,
		}
	}

	return data, nil
}
