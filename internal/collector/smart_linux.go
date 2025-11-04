//go:build linux
// +build linux

package collector

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
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
	RawValue   int64  `json:"raw_value"`
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
		Device:          device,
		ModelFamily:     smartOutput.ModelFamily,
		DeviceModel:     smartOutput.ModelName,
		Serial:          smartOutput.SerialNumber,
		Capacity:        smartOutput.UserCapacity.Bytes,
		Healthy:         smartOutput.SmartStatus.Passed,
		Attributes:      make(map[string]string),
		DetailedAttribs: make([]types.SMARTAttribute, 0),
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

	// Parse ATA SMART attributes with detailed information
	failingAttrs := make([]string, 0)
	warningAttrs := make([]string, 0)

	for _, attr := range smartOutput.AtaSmartAttrs.Table {
		info.Attributes[attr.Name] = fmt.Sprintf("%d", attr.RawValue)
		info.Attributes[attr.Name+"_Current"] = fmt.Sprintf("%d", attr.Value)
		info.Attributes[attr.Name+"_Worst"] = fmt.Sprintf("%d", attr.Worst)
		info.Attributes[attr.Name+"_Threshold"] = fmt.Sprintf("%d", attr.Threshold)

		// Create detailed attribute
		detailedAttr := types.SMARTAttribute{
			ID:         uint8(attr.ID),
			Name:       attr.Name,
			Value:      uint8(attr.Value),
			Worst:      uint8(attr.Worst),
			Threshold:  uint8(attr.Threshold),
			RawValue:   uint64(attr.RawValue),
			RawString:  attr.RawString,
			WhenFailed: attr.WhenFailed,
			Type:       "Old_age", // smartctl doesn't always provide this
			Updated:    "Always",
		}
		info.DetailedAttribs = append(info.DetailedAttribs, detailedAttr)

		// Check for failures
		if attr.WhenFailed != "" && attr.WhenFailed != "-" {
			info.Healthy = false
			if attr.WhenFailed == "FAILING_NOW" || attr.WhenFailed == "now" {
				failingAttrs = append(failingAttrs, fmt.Sprintf("%s (Value: %d, Threshold: %d)",
					attr.Name, attr.Value, attr.Threshold))
			}
		}

		// Check for critical attributes with non-zero values
		criticalAttrs := map[string]bool{
			"Reallocated_Sector_Ct":  true,
			"Current_Pending_Sector": true,
			"Offline_Uncorrectable":  true,
			"Reported_Uncorrect":     true,
		}
		if criticalAttrs[attr.Name] && attr.RawValue > 0 {
			warningAttrs = append(warningAttrs, fmt.Sprintf("%s = %d", attr.Name, attr.RawValue))
		}

		// Extract common values
		switch attr.ID {
		case 9: // Power-on hours
			info.PowerOnHours = uint64(attr.RawValue)
		case 12: // Power cycle count
			info.PowerCycleCount = uint64(attr.RawValue)
		case 194: // Temperature
			info.Temperature = int(attr.RawValue)
		}
	}

	// Create health assessment
	if len(failingAttrs) > 0 || len(warningAttrs) > 0 || !smartOutput.SmartStatus.Passed {
		info.HealthAssessment = &types.SMARTHealthStatus{
			Passed:            smartOutput.SmartStatus.Passed,
			FailingAttributes: failingAttrs,
			WarningAttributes: warningAttrs,
		}

		if len(failingAttrs) > 0 {
			info.HealthAssessment.OverallAssessment = "FAIL"
		} else if len(warningAttrs) > 0 {
			info.HealthAssessment.OverallAssessment = "WARN"
		} else {
			info.HealthAssessment.OverallAssessment = "PASS"
		}

		// Temperature assessment
		if info.Temperature > 70 {
			info.HealthAssessment.TemperatureStatus = "CRITICAL"
		} else if info.Temperature > 60 {
			info.HealthAssessment.TemperatureStatus = "HIGH"
		} else if info.Temperature > 45 {
			info.HealthAssessment.TemperatureStatus = "WARM"
		} else if info.Temperature > 0 {
			info.HealthAssessment.TemperatureStatus = "NORMAL"
		}
	}

	return info
}
