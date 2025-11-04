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

	return &types.MemoryData{
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
	}, nil
}
