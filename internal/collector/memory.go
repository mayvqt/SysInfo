package collector

import (
	"fmt"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/mayvqt/sysinfo/internal/utils"
	"github.com/shirou/gopsutil/v3/mem"
)

// CollectMemory gathers memory information
func CollectMemory() (*types.MemoryData, error) {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		// Swap might not be available on all systems
		swap = &mem.SwapMemoryStat{}
	}

	data := &types.MemoryData{
		Total:          vmem.Total,
		Available:      vmem.Available,
		Used:           vmem.Used,
		UsedPercent:    vmem.UsedPercent,
		Free:           vmem.Free,
		TotalFormatted: utils.FormatBytes(vmem.Total),
		UsedFormatted:  utils.FormatBytes(vmem.Used),
		FreeFormatted:  utils.FormatBytes(vmem.Free),
		SwapTotal:      swap.Total,
		SwapUsed:       swap.Used,
		SwapFree:       swap.Free,
		SwapPercent:    swap.UsedPercent,
		Cached:         vmem.Cached,
		Buffers:        vmem.Buffers,
		Shared:         vmem.Shared,
	}

	// Try to collect physical memory module information
	modules := collectMemoryModules()
	if len(modules) > 0 {
		data.Modules = modules
	}

	return data, nil
}

// collectMemoryModules attempts to collect physical RAM module information
// This requires platform-specific implementation or external tools
func collectMemoryModules() []types.MemoryModule {
	// This is a placeholder - full implementation would require:
	// - Windows: WMI queries to Win32_PhysicalMemory
	// - Linux: dmidecode or /sys/devices/system/memory
	// - macOS: system_profiler SPMemoryDataType
	modules := make([]types.MemoryModule, 0)

	// Platform-specific collection would go here
	// For now, return empty to avoid errors

	return modules
}
