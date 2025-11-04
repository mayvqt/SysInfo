package collector

import (
	"fmt"
	"sort"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/shirou/gopsutil/v3/process"
)

// CollectProcesses gathers process information
func CollectProcesses() (*types.ProcessData, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	data := &types.ProcessData{
		TotalCount:  len(processes),
		TopByMemory: make([]types.ProcessInfo, 0),
		TopByCPU:    make([]types.ProcessInfo, 0),
	}

	processInfos := make([]types.ProcessInfo, 0)
	running := 0
	sleeping := 0

	for _, proc := range processes {
		name, _ := proc.Name()
		username, _ := proc.Username()
		cpuPercent, _ := proc.CPUPercent()
		memPercent, _ := proc.MemoryPercent()
		memInfo, _ := proc.MemoryInfo()
		status, _ := proc.Status()
		createTime, _ := proc.CreateTime()

		// Count status
		if len(status) > 0 {
			switch status[0] {
			case "R":
				running++
			case "S":
				sleeping++
			}
		}

		memMB := uint64(0)
		if memInfo != nil {
			memMB = memInfo.RSS / (1024 * 1024)
		}

		pInfo := types.ProcessInfo{
			PID:           proc.Pid,
			Name:          name,
			Username:      username,
			CPUPercent:    cpuPercent,
			MemoryPercent: memPercent,
			MemoryMB:      memMB,
			Status:        status[0],
			CreateTime:    createTime,
		}

		processInfos = append(processInfos, pInfo)
	}

	data.Running = running
	data.Sleeping = sleeping

	// Get top 10 by memory
	sortedByMem := make([]types.ProcessInfo, len(processInfos))
	copy(sortedByMem, processInfos)
	sort.Slice(sortedByMem, func(i, j int) bool {
		return sortedByMem[i].MemoryMB > sortedByMem[j].MemoryMB
	})
	if len(sortedByMem) > 10 {
		data.TopByMemory = sortedByMem[:10]
	} else {
		data.TopByMemory = sortedByMem
	}

	// Get top 10 by CPU
	sortedByCPU := make([]types.ProcessInfo, len(processInfos))
	copy(sortedByCPU, processInfos)
	sort.Slice(sortedByCPU, func(i, j int) bool {
		return sortedByCPU[i].CPUPercent > sortedByCPU[j].CPUPercent
	})
	if len(sortedByCPU) > 10 {
		data.TopByCPU = sortedByCPU[:10]
	} else {
		data.TopByCPU = sortedByCPU
	}

	return data, nil
}
