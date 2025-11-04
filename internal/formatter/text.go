package formatter

import (
	"fmt"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
)

// FormatText formats the information as plain text
func FormatText(info *types.SystemInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Timestamp: %s\n\n", info.Timestamp.Format("2006-01-02 15:04:05")))

	// System information
	if info.System != nil {
		sb.WriteString("SYSTEM INFORMATION\n")
		sb.WriteString(fmt.Sprintf("Hostname: %s\n", info.System.Hostname))
		sb.WriteString(fmt.Sprintf("OS: %s\n", info.System.OS))
		sb.WriteString(fmt.Sprintf("Platform: %s %s\n", info.System.Platform, info.System.PlatformVersion))
		sb.WriteString(fmt.Sprintf("Platform Family: %s\n", info.System.PlatformFamily))
		sb.WriteString(fmt.Sprintf("Kernel: %s (%s)\n", info.System.KernelVersion, info.System.KernelArch))
		sb.WriteString(fmt.Sprintf("Uptime: %s\n", info.System.UptimeFormatted))
		sb.WriteString(fmt.Sprintf("Processes: %d\n\n", info.System.Procs))
	}

	// CPU information
	if info.CPU != nil {
		sb.WriteString("CPU INFORMATION\n")
		sb.WriteString(fmt.Sprintf("Model: %s\n", info.CPU.ModelName))
		sb.WriteString(fmt.Sprintf("Vendor: %s\n", info.CPU.Vendor))
		sb.WriteString(fmt.Sprintf("Physical Cores: %d\n", info.CPU.Cores))
		sb.WriteString(fmt.Sprintf("Logical CPUs: %d\n", info.CPU.LogicalCPUs))
		sb.WriteString(fmt.Sprintf("Frequency: %.2f MHz\n", info.CPU.MHz))
		if info.CPU.LoadAvg != nil {
			sb.WriteString(fmt.Sprintf("Load Average: %.2f, %.2f, %.2f\n",
				info.CPU.LoadAvg.Load1, info.CPU.LoadAvg.Load5, info.CPU.LoadAvg.Load15))
		}
		if len(info.CPU.Usage) > 0 {
			sb.WriteString("CPU Usage Per Core:\n")
			for i, usage := range info.CPU.Usage {
				sb.WriteString(fmt.Sprintf("  Core %d: %.2f%%\n", i, usage))
			}
		}
		sb.WriteString("\n")
	}

	// Memory information
	if info.Memory != nil {
		sb.WriteString("MEMORY INFORMATION\n")
		sb.WriteString(fmt.Sprintf("Total: %s\n", info.Memory.TotalFormatted))
		sb.WriteString(fmt.Sprintf("Used: %s (%.2f%%)\n", info.Memory.UsedFormatted, info.Memory.UsedPercent))
		sb.WriteString(fmt.Sprintf("Free: %s\n", info.Memory.FreeFormatted))
		if info.Memory.SwapTotal > 0 {
			sb.WriteString(fmt.Sprintf("Swap Total: %s\n", formatBytes(info.Memory.SwapTotal)))
			sb.WriteString(fmt.Sprintf("Swap Used: %s (%.2f%%)\n", formatBytes(info.Memory.SwapUsed), info.Memory.SwapPercent))
		}
		sb.WriteString("\n")
	}

	// Disk information
	if info.Disk != nil && len(info.Disk.Partitions) > 0 {
		sb.WriteString("DISK INFORMATION\n")
		for _, part := range info.Disk.Partitions {
			sb.WriteString(fmt.Sprintf("Device: %s\n", part.Device))
			sb.WriteString(fmt.Sprintf("  Mount Point: %s\n", part.MountPoint))
			sb.WriteString(fmt.Sprintf("  Type: %s\n", part.FSType))
			sb.WriteString(fmt.Sprintf("  Total: %s\n", part.TotalFormatted))
			sb.WriteString(fmt.Sprintf("  Used: %s (%.2f%%)\n", part.UsedFormatted, part.UsedPercent))
			sb.WriteString(fmt.Sprintf("  Free: %s\n", part.FreeFormatted))
		}
		sb.WriteString("\n")
	}

	// SMART disk health information
	if info.Disk != nil && len(info.Disk.SMARTData) > 0 {
		sb.WriteString("SMART DISK HEALTH\n")
		for _, smart := range info.Disk.SMARTData {
			deviceName := smart.Device
			if smart.DeviceModel != "" {
				deviceName = smart.DeviceModel
			}

			healthStatus := "HEALTHY"
			if !smart.Healthy {
				healthStatus = "WARNING"
			}

			sb.WriteString(fmt.Sprintf("Device: %s [%s]\n", deviceName, healthStatus))

			if smart.Serial != "" {
				sb.WriteString(fmt.Sprintf("  Serial: %s\n", smart.Serial))
			}
			if smart.ModelFamily != "" {
				sb.WriteString(fmt.Sprintf("  Model Family: %s\n", smart.ModelFamily))
			}
			if smart.Capacity > 0 {
				sb.WriteString(fmt.Sprintf("  Capacity: %s\n", formatBytes(smart.Capacity)))
			}
			if smart.Temperature > 0 {
				sb.WriteString(fmt.Sprintf("  Temperature: %d°C\n", smart.Temperature))
			}
			if smart.PowerOnHours > 0 {
				days := smart.PowerOnHours / 24
				sb.WriteString(fmt.Sprintf("  Power-On Hours: %d (%d days)\n", smart.PowerOnHours, days))
			}

			// Display key SMART attributes
			if len(smart.Attributes) > 0 {
				criticalAttrs := []string{
					"Reallocated_Sector_Count",
					"Current_Pending_Sector",
					"Offline_Uncorrectable",
					"UDMA_CRC_Error_Count",
					"SMART",
					"Status",
				}

				hasShownAttrs := false
				for _, attrName := range criticalAttrs {
					if val, ok := smart.Attributes[attrName]; ok {
						if !hasShownAttrs {
							sb.WriteString("  Attributes:\n")
							hasShownAttrs = true
						}
						sb.WriteString(fmt.Sprintf("    %s: %s\n", attrName, val))
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	// Network information
	if info.Network != nil && len(info.Network.Interfaces) > 0 {
		sb.WriteString("NETWORK INTERFACES\n")
		for _, iface := range info.Network.Interfaces {
			sb.WriteString(fmt.Sprintf("Interface: %s\n", iface.Name))
			if iface.HardwareAddr != "" {
				sb.WriteString(fmt.Sprintf("  MAC: %s\n", iface.HardwareAddr))
			}
			if len(iface.Addresses) > 0 {
				sb.WriteString(fmt.Sprintf("  Addresses: %s\n", strings.Join(iface.Addresses, ", ")))
			}
			if len(iface.Flags) > 0 {
				sb.WriteString(fmt.Sprintf("  Flags: %s\n", strings.Join(iface.Flags, ", ")))
			}
			sb.WriteString(fmt.Sprintf("  MTU: %d\n", iface.MTU))
			if iface.BytesSent > 0 || iface.BytesRecv > 0 {
				sb.WriteString(fmt.Sprintf("  Bytes Sent: %s\n", formatBytes(iface.BytesSent)))
				sb.WriteString(fmt.Sprintf("  Bytes Received: %s\n", formatBytes(iface.BytesRecv)))
			}
		}
		sb.WriteString("\n")
	}

	// Process information
	if info.Processes != nil {
		sb.WriteString("PROCESS INFORMATION\n")
		sb.WriteString(fmt.Sprintf("Total: %d (Running: %d, Sleeping: %d)\n",
			info.Processes.TotalCount, info.Processes.Running, info.Processes.Sleeping))

		if len(info.Processes.TopByMemory) > 0 {
			sb.WriteString("\nTop Processes by Memory:\n")
			for i, proc := range info.Processes.TopByMemory {
				if i >= 5 {
					break
				}
				sb.WriteString(fmt.Sprintf("  %s (PID %d): %d MB (%.2f%%)\n",
					proc.Name, proc.PID, proc.MemoryMB, proc.MemoryPercent))
			}
		}

		if len(info.Processes.TopByCPU) > 0 {
			sb.WriteString("\nTop Processes by CPU:\n")
			for i, proc := range info.Processes.TopByCPU {
				if i >= 5 {
					break
				}
				sb.WriteString(fmt.Sprintf("  %s (PID %d): %.2f%%\n",
					proc.Name, proc.PID, proc.CPUPercent))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(div), units[exp])
}
