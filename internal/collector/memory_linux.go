//go:build linux
// +build linux

package collector

import "github.com/mayvqt/sysinfo/internal/types"

// collectMemoryModulesPlatform implements Linux-specific memory module collection
func collectMemoryModulesPlatform() []types.MemoryModule {
	// TODO: Implement using dmidecode or /sys/devices/system/memory
	return make([]types.MemoryModule, 0)
}
