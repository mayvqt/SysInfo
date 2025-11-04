//go:build darwin
// +build darwin

package collector

import "github.com/mayvqt/sysinfo/internal/types"

// collectMemoryModulesPlatform implements macOS-specific memory module collection
func collectMemoryModulesPlatform() []types.MemoryModule {
	// TODO: Implement using system_profiler SPMemoryDataType
	return make([]types.MemoryModule, 0)
}
