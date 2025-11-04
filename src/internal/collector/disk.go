package collector

import (
	"fmt"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/mayvqt/sysinfo/internal/utils"
	"github.com/shirou/gopsutil/v3/disk"
)

// CollectDisk gathers disk and partition information
func CollectDisk(includeSMART bool) (*types.DiskData, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	data := &types.DiskData{
		Partitions: make([]types.PartitionInfo, 0),
		IOStats:    make([]types.DiskIOStat, 0),
	}

	// Collect partition information
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // Skip partitions we can't access
		}

		partInfo := types.PartitionInfo{
			Device:         partition.Device,
			MountPoint:     partition.Mountpoint,
			FSType:         partition.Fstype,
			Total:          usage.Total,
			Free:           usage.Free,
			Used:           usage.Used,
			UsedPercent:    usage.UsedPercent,
			TotalFormatted: utils.FormatBytes(usage.Total),
			UsedFormatted:  utils.FormatBytes(usage.Used),
			FreeFormatted:  utils.FormatBytes(usage.Free),
			InodesTotal:    usage.InodesTotal,
			InodesUsed:     usage.InodesUsed,
			InodesFree:     usage.InodesFree,
		}

		data.Partitions = append(data.Partitions, partInfo)
	}

	// Collect I/O statistics
	ioCounters, err := disk.IOCounters()
	if err == nil {
		for name, io := range ioCounters {
			ioStat := types.DiskIOStat{
				Name:       name,
				ReadCount:  io.ReadCount,
				WriteCount: io.WriteCount,
				ReadBytes:  io.ReadBytes,
				WriteBytes: io.WriteBytes,
				ReadTime:   io.ReadTime,
				WriteTime:  io.WriteTime,
				IoTime:     io.IoTime,
			}
			data.IOStats = append(data.IOStats, ioStat)
		}
	}

	// Collect SMART data if requested
	if includeSMART {
		data.SMARTData = CollectSMART()
	}

	return data, nil
}

// CollectSMART gathers SMART data from drives
// Note: This is a placeholder. Full SMART implementation requires platform-specific code
// or external libraries like smartmontools
func CollectSMART() []types.SMARTInfo {
	// This is a basic implementation
	// For production use, you would need to:
	// 1. Call smartctl or similar tools
	// 2. Parse the output
	// 3. Handle platform differences

	smartData := make([]types.SMARTInfo, 0)

	// Platform-specific SMART collection would go here
	// For now, return empty to avoid errors

	return smartData
}
