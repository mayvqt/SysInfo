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
		Partitions:    make([]types.PartitionInfo, 0),
		PhysicalDisks: make([]types.PhysicalDisk, 0),
		IOStats:       make([]types.DiskIOStat, 0),
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

			// Create physical disk info from IO counter data
			physicalDisk := types.PhysicalDisk{
				Name:          name,
				SizeFormatted: "N/A",
			}
			data.PhysicalDisks = append(data.PhysicalDisks, physicalDisk)
		}
	}

	// Get detailed physical disk information from platform-specific implementation
	physicalDisks := collectPhysicalDisksPlatform()
	if len(physicalDisks) > 0 {
		data.PhysicalDisks = physicalDisks
	}

	// Collect SMART data if requested
	if includeSMART {
		data.SMARTData = CollectSMART()
	}

	return data, nil
}

// CollectSMART gathers SMART data from drives
func CollectSMART() []types.SMARTInfo {
	// Call platform-specific implementation
	return collectSMARTPlatform()
}
