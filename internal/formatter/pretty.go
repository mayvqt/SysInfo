package formatter

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/mayvqt/sysinfo/internal/types"
)

// FormatPretty formats the information with colors and tables
func FormatPretty(info *types.SystemInfo) string {
	var sb strings.Builder

	// Color definitions
	headerColor := color.New(color.FgCyan, color.Bold)
	labelColor := color.New(color.FgGreen)
	valueColor := color.New(color.FgWhite)

	// Timestamp
	sb.WriteString(headerColor.Sprintf("═══════════════════════════════════════════════════════════════\n"))
	sb.WriteString(headerColor.Sprintf("  SYSTEM INFORMATION REPORT\n"))
	sb.WriteString(headerColor.Sprintf("═══════════════════════════════════════════════════════════════\n"))
	sb.WriteString(fmt.Sprintf("%s %s\n\n", labelColor.Sprint("Timestamp:"), valueColor.Sprint(info.Timestamp.Format("2006-01-02 15:04:05"))))

	// System information
	if info.System != nil {
		sb.WriteString(headerColor.Sprintf("┌─ SYSTEM ─────────────────────────────────────────────────────┐\n"))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Hostname:"), valueColor.Sprint(info.System.Hostname)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("OS:"), valueColor.Sprint(info.System.OS)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s %s\n", labelColor.Sprint("Platform:"), valueColor.Sprint(info.System.Platform), valueColor.Sprint(info.System.PlatformVersion)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Kernel:"), valueColor.Sprintf("%s (%s)", info.System.KernelVersion, info.System.KernelArch)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Uptime:"), valueColor.Sprint(info.System.UptimeFormatted)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Processes:"), valueColor.Sprintf("%d", info.System.Procs)))
		sb.WriteString(headerColor.Sprintf("└──────────────────────────────────────────────────────────────┘\n\n"))
	}

	// CPU information
	if info.CPU != nil {
		sb.WriteString(headerColor.Sprintf("┌─ CPU ────────────────────────────────────────────────────────┐\n"))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Model:"), valueColor.Sprint(info.CPU.ModelName)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Vendor:"), valueColor.Sprint(info.CPU.Vendor)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Physical Cores:"), valueColor.Sprintf("%d", info.CPU.Cores)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Logical CPUs:"), valueColor.Sprintf("%d", info.CPU.LogicalCPUs)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Frequency:"), valueColor.Sprintf("%.2f MHz", info.CPU.MHz)))

		if info.CPU.CacheSize > 0 {
			sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Cache Size:"), valueColor.Sprintf("%d KB", info.CPU.CacheSize)))
		}

		if info.CPU.Microcode != "" {
			sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Microcode:"), valueColor.Sprint(info.CPU.Microcode)))
		}

		if info.CPU.LoadAvg != nil {
			sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Load Average:"),
				valueColor.Sprintf("%.2f, %.2f, %.2f", info.CPU.LoadAvg.Load1, info.CPU.LoadAvg.Load5, info.CPU.LoadAvg.Load15)))
		}

		if len(info.CPU.Usage) > 0 {
			sb.WriteString(fmt.Sprintf("│ %-20s\n", labelColor.Sprint("Core Usage:")))
			for i, usage := range info.CPU.Usage {
				bar := createProgressBar(usage, 20)
				sb.WriteString(fmt.Sprintf("│   Core %-2d: %s %s\n", i, bar, valueColor.Sprintf("%.1f%%", usage)))
			}
		}
		sb.WriteString(headerColor.Sprintf("└──────────────────────────────────────────────────────────────┘\n\n"))
	}

	// Memory information
	if info.Memory != nil {
		sb.WriteString(headerColor.Sprintf("┌─ MEMORY ─────────────────────────────────────────────────────┐\n"))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Total:"), valueColor.Sprint(info.Memory.TotalFormatted)))

		memBar := createProgressBar(info.Memory.UsedPercent, 30)
		sb.WriteString(fmt.Sprintf("│ %-20s %s %s\n", labelColor.Sprint("Used:"),
			memBar, valueColor.Sprintf("%s (%.1f%%)", info.Memory.UsedFormatted, info.Memory.UsedPercent)))
		sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Free:"), valueColor.Sprint(info.Memory.FreeFormatted)))

		if info.Memory.Cached > 0 {
			sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Cached:"), valueColor.Sprint(formatBytes(info.Memory.Cached))))
		}
		if info.Memory.Buffers > 0 {
			sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Buffers:"), valueColor.Sprint(formatBytes(info.Memory.Buffers))))
		}

		if info.Memory.SwapTotal > 0 {
			sb.WriteString(fmt.Sprintf("│ %-20s %s\n", labelColor.Sprint("Swap Total:"), valueColor.Sprint(formatBytes(info.Memory.SwapTotal))))
			swapBar := createProgressBar(info.Memory.SwapPercent, 30)
			sb.WriteString(fmt.Sprintf("│ %-20s %s %s\n", labelColor.Sprint("Swap Used:"),
				swapBar, valueColor.Sprintf("%s (%.1f%%)", formatBytes(info.Memory.SwapUsed), info.Memory.SwapPercent)))
		}

		if len(info.Memory.Modules) > 0 {
			sb.WriteString(fmt.Sprintf("│\n│ %s\n", labelColor.Sprint("Physical Modules:")))
			for _, module := range info.Memory.Modules {
				sb.WriteString(fmt.Sprintf("│   %s\n", valueColor.Sprintf("%s: %s", module.Locator, formatBytes(module.Capacity))))
				if module.Speed > 0 {
					sb.WriteString(fmt.Sprintf("│     Speed: %s, Type: %s\n", valueColor.Sprintf("%d MHz", module.Speed), valueColor.Sprint(module.Type)))
				}
			}
		}

		sb.WriteString(headerColor.Sprintf("└──────────────────────────────────────────────────────────────┘\n\n"))
	}

	// Disk information
	if info.Disk != nil {
		sb.WriteString(headerColor.Sprintf("┌─ STORAGE ────────────────────────────────────────────────────┐\n"))

		// Physical disks information first (the actual hardware)
		if len(info.Disk.PhysicalDisks) > 0 {
			sb.WriteString(fmt.Sprintf("│ %s\n", labelColor.Sprint("Physical Disks:")))
			sb.WriteString("│\n")
			for _, disk := range info.Disk.PhysicalDisks {
				diskType := disk.Type
				if disk.Type == "" {
					diskType = "Unknown"
				}

				// Show disk name and type
				sb.WriteString(fmt.Sprintf("│ %s [%s]", valueColor.Sprint(disk.Name), valueColor.Sprint(diskType)))
				if disk.Interface != "" {
					sb.WriteString(fmt.Sprintf(" %s", color.New(color.FgCyan).Sprint(disk.Interface)))
				}
				sb.WriteString("\n")

				// Show model
				if disk.Model != "" {
					sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Model:"), valueColor.Sprint(disk.Model)))
				}

				// Show size
				if disk.SizeFormatted != "" && disk.SizeFormatted != "N/A" {
					sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Size:"), valueColor.Sprint(disk.SizeFormatted)))
				}

				// Show serial number
				if disk.SerialNumber != "" {
					sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Serial:"), valueColor.Sprint(disk.SerialNumber)))
				}

				// Show RPM for HDDs
				if disk.RPM > 0 {
					sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("RPM:"), valueColor.Sprintf("%d", disk.RPM)))
				}

				// Show removable status
				if disk.Removable {
					sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Removable:"), color.New(color.FgYellow).Sprint("Yes")))
				}

				sb.WriteString("│\n")
			}
		}

		// Mounted partitions (filter out loop devices and snaps for cleaner output)
		if len(info.Disk.Partitions) > 0 {
			// Filter significant partitions
			var significantPartitions []types.PartitionInfo
			for _, part := range info.Disk.Partitions {
				// Skip loop devices (snap mounts) and very small partitions
				if strings.HasPrefix(part.Device, "/dev/loop") {
					continue
				}
				// Skip if squashfs (usually snaps)
				if part.FSType == "squashfs" {
					continue
				}
				significantPartitions = append(significantPartitions, part)
			}

			if len(significantPartitions) > 0 {
				sb.WriteString(fmt.Sprintf("│ %s\n", labelColor.Sprint("Mounted Partitions:")))
				sb.WriteString("│\n")
				for _, part := range significantPartitions {
					sb.WriteString(fmt.Sprintf("│ %s", valueColor.Sprintf("%s", part.Device)))
					if part.MountPoint != "" {
						sb.WriteString(fmt.Sprintf(" → %s", valueColor.Sprintf("%s", part.MountPoint)))
					}
					sb.WriteString("\n")

					sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Type:"), valueColor.Sprint(part.FSType)))
					sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Total:"), valueColor.Sprint(part.TotalFormatted)))

					diskBar := createProgressBar(part.UsedPercent, 28)
					sb.WriteString(fmt.Sprintf("│   %-18s %s %s\n", labelColor.Sprint("Used:"),
						diskBar, valueColor.Sprintf("%s (%.1f%%)", part.UsedFormatted, part.UsedPercent)))
					sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Free:"), valueColor.Sprint(part.FreeFormatted)))
					sb.WriteString("│\n")
				}
			}
		}

		sb.WriteString(headerColor.Sprintf("└──────────────────────────────────────────────────────────────┘\n\n"))
	}

	// SMART disk health information
	if info.Disk != nil && len(info.Disk.SMARTData) > 0 {
		sb.WriteString(headerColor.Sprintf("┌─ SMART DISK HEALTH ──────────────────────────────────────────┐\n"))
		for _, smart := range info.Disk.SMARTData {
			deviceName := smart.Device
			if smart.DeviceModel != "" {
				deviceName = smart.DeviceModel
			}

			healthStatus := "HEALTHY"
			healthColor := color.New(color.FgGreen, color.Bold)
			if !smart.Healthy {
				healthStatus = "WARNING"
				healthColor = color.New(color.FgRed, color.Bold)
			}

			sb.WriteString(fmt.Sprintf("│ %s [%s]\n", valueColor.Sprint(deviceName), healthColor.Sprint(healthStatus)))

			if smart.Serial != "" {
				sb.WriteString(fmt.Sprintf("│   %-20s %s\n", labelColor.Sprint("Serial:"), valueColor.Sprint(smart.Serial)))
			}
			if smart.Capacity > 0 {
				sb.WriteString(fmt.Sprintf("│   %-20s %s\n", labelColor.Sprint("Capacity:"), valueColor.Sprint(formatBytes(smart.Capacity))))
			}
			if smart.Temperature > 0 {
				tempColor := valueColor
				if smart.Temperature > 50 {
					tempColor = color.New(color.FgYellow)
				}
				if smart.Temperature > 60 {
					tempColor = color.New(color.FgRed)
				}
				sb.WriteString(fmt.Sprintf("│   %-20s %s\n", labelColor.Sprint("Temperature:"), tempColor.Sprintf("%d°C", smart.Temperature)))
			}
			if smart.PowerOnHours > 0 {
				days := smart.PowerOnHours / 24
				sb.WriteString(fmt.Sprintf("│   %-20s %s (%s days)\n",
					labelColor.Sprint("Power-On Hours:"),
					valueColor.Sprintf("%d", smart.PowerOnHours),
					valueColor.Sprintf("%d", days)))
			}

			// Display key SMART attributes
			if len(smart.Attributes) > 0 {
				// Show critical attributes
				criticalAttrs := []string{
					"Reallocated_Sector_Count",
					"Current_Pending_Sector",
					"Offline_Uncorrectable",
					"UDMA_CRC_Error_Count",
				}

				hasShownAttrs := false
				for _, attrName := range criticalAttrs {
					if val, ok := smart.Attributes[attrName]; ok && val != "0" {
						if !hasShownAttrs {
							sb.WriteString(fmt.Sprintf("│   %s\n", labelColor.Sprint("Critical Attributes:")))
							hasShownAttrs = true
						}
						warnColor := color.New(color.FgYellow)
						sb.WriteString(fmt.Sprintf("│     %s: %s\n", attrName, warnColor.Sprint(val)))
					}
				}
			}

			sb.WriteString("│\n")
		}
		sb.WriteString(headerColor.Sprintf("└──────────────────────────────────────────────────────────────┘\n\n"))
	}

	// Network information
	if info.Network != nil && len(info.Network.Interfaces) > 0 {
		sb.WriteString(headerColor.Sprintf("┌─ NETWORK ────────────────────────────────────────────────────┐\n"))
		for _, iface := range info.Network.Interfaces {
			sb.WriteString(fmt.Sprintf("│ %s\n", valueColor.Sprint(iface.Name)))
			if iface.HardwareAddr != "" {
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("MAC:"), valueColor.Sprint(iface.HardwareAddr)))
			}
			if len(iface.Addresses) > 0 {
				for i, addr := range iface.Addresses {
					if i == 0 {
						sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("IP:"), valueColor.Sprint(addr)))
					} else {
						sb.WriteString(fmt.Sprintf("│   %-18s %s\n", "", valueColor.Sprint(addr)))
					}
				}
			}
			if iface.BytesSent > 0 || iface.BytesRecv > 0 {
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Sent:"), valueColor.Sprint(formatBytes(iface.BytesSent))))
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Received:"), valueColor.Sprint(formatBytes(iface.BytesRecv))))
			}
			sb.WriteString("│\n")
		}
		sb.WriteString(headerColor.Sprintf("└──────────────────────────────────────────────────────────────┘\n\n"))
	}

	// Process information
	if info.Processes != nil {
		sb.WriteString(headerColor.Sprintf("┌─ PROCESSES ──────────────────────────────────────────────────┐\n"))
		sb.WriteString(fmt.Sprintf("│ %-20s %s (Running: %s, Sleeping: %s)\n",
			labelColor.Sprint("Total:"),
			valueColor.Sprintf("%d", info.Processes.TotalCount),
			valueColor.Sprintf("%d", info.Processes.Running),
			valueColor.Sprintf("%d", info.Processes.Sleeping)))

		if len(info.Processes.TopByMemory) > 0 {
			sb.WriteString(fmt.Sprintf("│\n│ %s\n", labelColor.Sprint("Top by Memory:")))
			for i, proc := range info.Processes.TopByMemory {
				if i >= 5 {
					break
				}
				sb.WriteString(fmt.Sprintf("│   %s\n", valueColor.Sprintf("%-30s %6d MB  %.1f%%",
					truncate(proc.Name, 30), proc.MemoryMB, proc.MemoryPercent)))
			}
		}

		if len(info.Processes.TopByCPU) > 0 {
			sb.WriteString(fmt.Sprintf("│\n│ %s\n", labelColor.Sprint("Top by CPU:")))
			for i, proc := range info.Processes.TopByCPU {
				if i >= 5 {
					break
				}
				sb.WriteString(fmt.Sprintf("│   %s\n", valueColor.Sprintf("%-30s %6.1f%%",
					truncate(proc.Name, 30), proc.CPUPercent)))
			}
		}

		sb.WriteString(headerColor.Sprintf("└──────────────────────────────────────────────────────────────┘\n"))
	}

	// GPU information
	if info.GPU != nil && len(info.GPU.GPUs) > 0 {
		sb.WriteString("\n")
		sb.WriteString(headerColor.Sprintf("┌─ GPU ────────────────────────────────────────────────────────┐\n"))
		for _, gpu := range info.GPU.GPUs {
			sb.WriteString(fmt.Sprintf("│ %s\n", valueColor.Sprintf("GPU %d: %s", gpu.Index, gpu.Name)))

			if gpu.Vendor != "" {
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Vendor:"), valueColor.Sprint(gpu.Vendor)))
			}

			if gpu.Driver != "" {
				driverStr := gpu.Driver
				if gpu.DriverVersion != "" {
					driverStr = fmt.Sprintf("%s (v%s)", gpu.Driver, gpu.DriverVersion)
				}
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Driver:"), valueColor.Sprint(driverStr)))
			}

			if gpu.MemoryTotal > 0 {
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Memory:"), valueColor.Sprint(gpu.MemoryFormatted)))
				if gpu.MemoryUsed > 0 {
					usedPercent := float64(gpu.MemoryUsed) / float64(gpu.MemoryTotal) * 100
					memBar := createProgressBar(usedPercent, 28)
					sb.WriteString(fmt.Sprintf("│   %-18s %s %s\n", labelColor.Sprint("Memory Used:"),
						memBar, valueColor.Sprintf("%s (%.1f%%)", formatBytes(gpu.MemoryUsed), usedPercent)))
				}
			}

			if gpu.Temperature > 0 {
				tempColor := valueColor
				if gpu.Temperature > 70 {
					tempColor = color.New(color.FgRed)
				} else if gpu.Temperature > 60 {
					tempColor = color.New(color.FgYellow)
				}
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Temperature:"), tempColor.Sprintf("%d°C", gpu.Temperature)))
			}

			if gpu.Utilization > 0 {
				utilBar := createProgressBar(float64(gpu.Utilization), 28)
				sb.WriteString(fmt.Sprintf("│   %-18s %s %s\n", labelColor.Sprint("GPU Utilization:"),
					utilBar, valueColor.Sprintf("%d%%", gpu.Utilization)))
			}

			if gpu.MemoryUtilization > 0 {
				memUtilBar := createProgressBar(float64(gpu.MemoryUtilization), 28)
				sb.WriteString(fmt.Sprintf("│   %-18s %s %s\n", labelColor.Sprint("Mem Utilization:"),
					memUtilBar, valueColor.Sprintf("%d%%", gpu.MemoryUtilization)))
			}

			if gpu.PowerDraw > 0 {
				powerStr := fmt.Sprintf("%.1f W", gpu.PowerDraw)
				if gpu.PowerLimit > 0 {
					powerStr = fmt.Sprintf("%.1f / %.1f W", gpu.PowerDraw, gpu.PowerLimit)
				}
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Power Draw:"), valueColor.Sprint(powerStr)))
			}

			if gpu.ClockSpeed > 0 {
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Clock Speed:"), valueColor.Sprintf("%d MHz", gpu.ClockSpeed)))
			}

			if gpu.ClockSpeedMemory > 0 {
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Memory Clock:"), valueColor.Sprintf("%d MHz", gpu.ClockSpeedMemory)))
			}

			if gpu.FanSpeed > 0 {
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("Fan Speed:"), valueColor.Sprintf("%d%%", gpu.FanSpeed)))
			}

			if gpu.PCIBus != "" {
				sb.WriteString(fmt.Sprintf("│   %-18s %s\n", labelColor.Sprint("PCI Bus:"), valueColor.Sprint(gpu.PCIBus)))
			}

			sb.WriteString("│\n")
		}
		sb.WriteString(headerColor.Sprintf("└──────────────────────────────────────────────────────────────┘\n"))
	}

	return sb.String()
}

// createProgressBar creates a text progress bar
func createProgressBar(percent float64, width int) string {
	filled := int(percent / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	// Color the bar based on usage
	if percent > 90 {
		return color.New(color.FgRed).Sprint(bar)
	} else if percent > 70 {
		return color.New(color.FgYellow).Sprint(bar)
	}
	return color.New(color.FgGreen).Sprint(bar)
}

// truncate truncates a string to the specified length
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// (removed) createTable was unused — kept tablewriter usage available if needed in future
