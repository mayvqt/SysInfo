//go:build linux
// +build linux

package collector

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mayvqt/sysinfo/src/internal/types"
)

// SmartctlOutput represents the JSON output from smartctl
type SmartctlOutput struct {
	Device struct {
		Name     string `json:"name"`
		InfoName string `json:"info_name"`
		Type     string `json:"type"`
		Protocol string `json:"protocol"`
	} `json:"device"`
	ModelFamily   string        `json:"model_family"`
	ModelName     string        `json:"model_name"`
	SerialNumber  string        `json:"serial_number"`
	UserCapacity  UserCapacity  `json:"user_capacity"`
	SmartStatus   SmartStatus   `json:"smart_status"`
	Temperature   Temperature   `json:"temperature"`
	PowerOnTime   PowerOnTime   `json:"power_on_time"`
	AtaSmartAttrs AtaSmartAttrs `json:"ata_smart_attributes"`
	NvmeSmartLog  NvmeSmartLog  `json:"nvme_smart_health_information_log"`
}

type UserCapacity struct {
	Blocks uint64 `json:"blocks"`
	Bytes  uint64 `json:"bytes"`
}

type SmartStatus struct {
	Passed bool `json:"passed"`
}

type Temperature struct {
	Current int `json:"current"`
}

type PowerOnTime struct {
	Hours uint64 `json:"hours"`
}

type AtaSmartAttrs struct {
	Table []SmartAttribute `json:"table"`
}

type SmartAttribute struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Worst      int    `json:"worst"`
	Threshold  int    `json:"thresh"`
	RawValue   int64  `json:"raw"`
	RawString  string `json:"raw"`
	WhenFailed string `json:"when_failed"`
}

type NvmeSmartLog struct {
	Temperature      int    `json:"temperature"`
	PowerOnHours     uint64 `json:"power_on_hours"`
	DataUnitsRead    uint64 `json:"data_units_read"`
	DataUnitsWritten uint64 `json:"data_units_written"`
}

// collectSMARTPlatform implements Linux-specific SMART data collection
func collectSMARTPlatform() []types.SMARTInfo {
	smartData := make([]types.SMARTInfo, 0)

	// Check if smartctl is available
	_, err := exec.LookPath("smartctl")
	if err != nil {
		// smartctl not available, return empty
		return smartData
	}

	// Get list of devices
	devices := getLinuxDiskDevices()

	for _, device := range devices {
		info := collectDeviceSMART(device)
		if info != nil {
			smartData = append(smartData, *info)
		}
	}

	return smartData
}

// getLinuxDiskDevices returns a list of disk devices to check
func getLinuxDiskDevices() []string {
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
			"/dev/sda", "/dev/sdb", "/dev/sdc", "/dev/sdd",
			"/dev/nvme0", "/dev/nvme1",
		}
		for _, dev := range commonDevices {
			// Check if device exists by trying a quick smartctl command
			cmd := exec.Command("smartctl", "-i", dev)
			if err := cmd.Run(); err == nil {
				devices = append(devices, dev)
			}
		}
	}

	return devices
}

// collectDeviceSMART collects SMART data for a specific device
func collectDeviceSMART(device string) *types.SMARTInfo {
	// Run smartctl with JSON output
	cmd := exec.Command("smartctl", "-a", "-j", device)
	output, err := cmd.Output()
	if err != nil {
		// Even if smartctl returns non-zero, it might still have data
		// smartctl returns non-zero for disks with warnings
		if len(output) == 0 {
			return nil
		}
	}

	var smartOutput SmartctlOutput
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

	// For NVMe devices, use NVMe-specific data
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
