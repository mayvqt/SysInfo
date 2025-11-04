package collector

import (
	"fmt"
	"time"

	"github.com/mayvqt/sysinfo/src/internal/types"
	"github.com/shirou/gopsutil/v3/host"
)

// CollectSystem gathers general system information
func CollectSystem() (*types.SystemData, error) {
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %w", err)
	}

	uptime := formatUptime(info.Uptime)

	return &types.SystemData{
		Hostname:        info.Hostname,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformFamily:  info.PlatformFamily,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
		Uptime:          info.Uptime,
		UptimeFormatted: uptime,
		BootTime:        info.BootTime,
		Procs:           info.Procs,
	}, nil
}

// formatUptime converts seconds to a human-readable format
func formatUptime(seconds uint64) string {
	duration := time.Duration(seconds) * time.Second
	days := duration / (24 * time.Hour)
	duration -= days * 24 * time.Hour
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
