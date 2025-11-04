//go:build darwin
// +build darwin

package collector

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mayvqt/sysinfo/src/internal/types"
)

// SmartctlOutput represents the JSON output from smartctl (same as Linux)
type SmartctlOutputDarwin struct {
	Device struct {
		Name     string `json:"name"`
		InfoName string `json:"info_name"`
		Type     string `json:"type"`
		Protocol string `json:"protocol"`
	} `json:"device"`
	ModelFamily   string              `json:"model_family"`
	ModelName     string              `json:"model_name"`
	SerialNumber  string              `json:"serial_number"`
	UserCapacity  UserCapacityDarwin  `json:"user_capacity"`
	SmartStatus   SmartStatusDarwin   `json:"smart_status"`
	Temperature   TemperatureDarwin   `json:"temperature"`
	PowerOnTime   PowerOnTimeDarwin   `json:"power_on_time"`
	AtaSmartAttrs AtaSmartAttrsDarwin `json:"ata_smart_attributes"`
	NvmeSmartLog  NvmeSmartLogDarwin  `json:"nvme_smart_health_information_log"`
}

type UserCapacityDarwin struct {
	Blocks uint64 `json:"blocks"`
	Bytes  uint64 `json:"bytes"`
}

type SmartStatusDarwin struct {
	Passed bool `json:"passed"`
}

type TemperatureDarwin struct {
	Current int `json:"current"`
}

type PowerOnTimeDarwin struct {
	Hours uint64 `json:"hours"`
}

type AtaSmartAttrsDarwin struct {
	Table []SmartAttributeDarwin `json:"table"`
}

type SmartAttributeDarwin struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Worst      int    `json:"worst"`
	Threshold  int    `json:"thresh"`
	RawValue   int64  `json:"raw"`
	RawString  string `json:"raw"`
	WhenFailed string `json:"when_failed"`
}

type NvmeSmartLogDarwin struct {
	Temperature      int    `json:"temperature"`
	PowerOnHours     uint64 `json:"power_on_hours"`
	DataUnitsRead    uint64 `json:"data_units_read"`
	DataUnitsWritten uint64 `json:"data_units_written"`
}

// collectSMARTPlatform implements macOS-specific SMART data collection
func collectSMARTPlatform() []types.SMARTInfo {
	smartData := make([]types.SMARTInfo, 0)

	// Check if smartctl is available
	_, err := exec.LookPath("smartctl")
	if err != nil {
		// smartctl not available, return empty
		// User needs to install smartmontools: brew install smartmontools
		return smartData
	}

	// Get list of devices
	devices := getDarwinDiskDevices()

	for _, device := range devices {
		info := collectDeviceSMARTDarwin(device)
		if info != nil {
			smartData = append(smartData, *info)
		}
	}

	return smartData
}

// getDarwinDiskDevices returns a list of disk devices to check
func getDarwinDiskDevices() []string {
	devices := make([]string, 0)

	// Try smartctl --scan first
	cmd := exec.Command("smartctl", "--scan")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) > 0 {
				devices = append(devices, fields[0])
			}
		}
	}

	// If scan didn't work, try common device paths
	if len(devices) == 0 {
		commonDevices := []string{
			"/dev/disk0", "/dev/disk1", "/dev/disk2",
		}
		for _, dev := range commonDevices {
			cmd := exec.Command("smartctl", "-i", dev)
			if err := cmd.Run(); err == nil {
				devices = append(devices, dev)
			}
		}
	}

	return devices
}

// collectDeviceSMARTDarwin collects SMART data for a specific device on macOS
func collectDeviceSMARTDarwin(device string) *types.SMARTInfo {
	// Run smartctl with JSON output
	cmd := exec.Command("smartctl", "-a", "-j", device)
	output, err := cmd.Output()
	if err != nil {
		// Even if smartctl returns non-zero, it might still have data
		if len(output) == 0 {
			return nil
		}
	}

	var smartOutput SmartctlOutputDarwin
	if err := json.Unmarshal(output, &smartOutput); err != nil {
		return nil
	}

	info := &types.SMARTInfo{
		Device:      device,
		ModelFamily: smartOutput.ModelFamily,
		DeviceModel: smartOutput.ModelName,
		Serial:      smartOutput.SerialNumber,
		Capacity:    smartOutput.UserCapacity.Bytes,
		Healthy:     smartOutput.SmartStatus.Passed,
		Attributes:  make(map[string]string),
	}

	// Extract temperature
	if smartOutput.Temperature.Current > 0 {
		info.Temperature = smartOutput.Temperature.Current
	}

	// Extract power-on hours
	if smartOutput.PowerOnTime.Hours > 0 {
		info.PowerOnHours = smartOutput.PowerOnTime.Hours
	}

	// For NVMe devices (including Apple Silicon SSDs), use NVMe-specific data
	if smartOutput.NvmeSmartLog.Temperature > 0 {
		info.Temperature = smartOutput.NvmeSmartLog.Temperature
		info.PowerOnHours = smartOutput.NvmeSmartLog.PowerOnHours
		info.Attributes["Data_Units_Read"] = fmt.Sprintf("%d", smartOutput.NvmeSmartLog.DataUnitsRead)
		info.Attributes["Data_Units_Written"] = fmt.Sprintf("%d", smartOutput.NvmeSmartLog.DataUnitsWritten)
	}

	// Parse ATA SMART attributes
	for _, attr := range smartOutput.AtaSmartAttrs.Table {
		info.Attributes[attr.Name] = fmt.Sprintf("%d", attr.RawValue)
		info.Attributes[attr.Name+"_Current"] = fmt.Sprintf("%d", attr.Value)
		info.Attributes[attr.Name+"_Worst"] = fmt.Sprintf("%d", attr.Worst)
		info.Attributes[attr.Name+"_Threshold"] = fmt.Sprintf("%d", attr.Threshold)

		if attr.WhenFailed != "" && attr.WhenFailed != "-" {
			info.Healthy = false
		}

		// Extract common values
		switch attr.ID {
		case 9: // Power-on hours
			info.PowerOnHours = uint64(attr.RawValue)
		case 194: // Temperature
			info.Temperature = int(attr.RawValue)
		}
	}

	return info
}
