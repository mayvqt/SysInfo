package collector

import (
	"fmt"
	"os"
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
			fmt.Fprintf(os.Stderr, "Error collecting system info: %v\n", err)
		}
	}

	// Collect CPU information
	if cfg.ShouldCollect("cpu") {
		info.CPU, err = CollectCPU()
		if err != nil && cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Error collecting CPU info: %v\n", err)
		}
	}

	// Collect memory information
	if cfg.ShouldCollect("memory") {
		info.Memory, err = CollectMemory()
		if err != nil && cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Error collecting memory info: %v\n", err)
		}
	}

	// Collect disk information
	if cfg.ShouldCollect("disk") {
		info.Disk, err = CollectDisk(cfg.ShouldCollect("smart"))
		if err != nil && cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Error collecting disk info: %v\n", err)
		}
	}

	// Collect network information
	if cfg.ShouldCollect("network") {
		info.Network, err = CollectNetwork()
		if err != nil && cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Error collecting network info: %v\n", err)
		}
	}

	// Collect process information
	if cfg.ShouldCollect("process") {
		info.Processes, err = CollectProcesses()
		if err != nil && cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Error collecting process info: %v\n", err)
		}
	}

	// Collect GPU information
	if cfg.ShouldCollect("gpu") {
		info.GPU, err = CollectGPU()
		if err != nil && cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Error collecting GPU info: %v\n", err)
		}
	}

	return info, nil
}
