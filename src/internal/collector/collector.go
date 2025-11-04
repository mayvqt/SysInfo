package collector

import (
	"time"

	"github.com/mayvqt/sysinfo/internal/config"
	"github.com/mayvqt/sysinfo/internal/types"
)

// Collect gathers all requested system information
func Collect(cfg *config.Config) (*types.SystemInfo, error) {
	info := &types.SystemInfo{
		Timestamp: time.Now(),
	}

	var err error

	// Collect system information
	if cfg.ShouldCollect("system") {
		info.System, err = CollectSystem()
		if err != nil && cfg.Verbose {
			// Log error but continue
		}
	}

	// Collect CPU information
	if cfg.ShouldCollect("cpu") {
		info.CPU, err = CollectCPU()
		if err != nil && cfg.Verbose {
			// Log error but continue
		}
	}

	// Collect memory information
	if cfg.ShouldCollect("memory") {
		info.Memory, err = CollectMemory()
		if err != nil && cfg.Verbose {
			// Log error but continue
		}
	}

	// Collect disk information
	if cfg.ShouldCollect("disk") {
		info.Disk, err = CollectDisk(cfg.ShouldCollect("smart"))
		if err != nil && cfg.Verbose {
			// Log error but continue
		}
	}

	// Collect network information
	if cfg.ShouldCollect("network") {
		info.Network, err = CollectNetwork()
		if err != nil && cfg.Verbose {
			// Log error but continue
		}
	}

	// Collect process information
	if cfg.ShouldCollect("process") {
		info.Processes, err = CollectProcesses()
		if err != nil && cfg.Verbose {
			// Log error but continue
		}
	}

	return info, nil
}
