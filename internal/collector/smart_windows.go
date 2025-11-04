//go:build windows
// +build windows

package collector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/yusufpapurcu/wmi"
)

// Win32_DiskDrive represents WMI disk drive data
type Win32_DiskDrive struct {
	DeviceID          string
	Model             string
	SerialNumber      string
	Size              uint64
	Status            string
	InterfaceType     string
	MediaType         string
	Partitions        uint32
	FirmwareRevision  string
	BytesPerSector    uint32
	SectorsPerTrack   uint32
	TracksPerCylinder uint32
	TotalCylinders    uint64
	TotalHeads        uint32
	TotalSectors      uint64
	TotalTracks       uint64
}

// MSStorageDriver_ATAPISmartData represents SMART data from WMI
type MSStorageDriver_ATAPISmartData struct {
	VendorSpecific []uint8
}

// MSStorageDriver_FailurePredictStatus represents disk health prediction
type MSStorageDriver_FailurePredictStatus struct {
	PredictFailure bool
	Reason         uint32
}

// MSStorageDriver_FailurePredictData represents SMART attributes
type MSStorageDriver_FailurePredictData struct {
	VendorSpecific []uint8
}

// MSStorageDriver_FailurePredictThresholds represents SMART thresholds
type MSStorageDriver_FailurePredictThresholds struct {
	VendorSpecific []uint8
}

// collectSMARTPlatform implements Windows-specific SMART data collection
func collectSMARTPlatform() []types.SMARTInfo {
	smartData := make([]types.SMARTInfo, 0)

	var drives []Win32_DiskDrive
	query := "SELECT * FROM Win32_DiskDrive"
	err := wmi.Query(query, &drives)
	if err != nil {
		return smartData
	}

	// Process each drive found
	for _, drive := range drives {
		info := types.SMARTInfo{
			Device:          drive.DeviceID,
			DeviceModel:     strings.TrimSpace(drive.Model),
			Serial:          strings.TrimSpace(drive.SerialNumber),
			FirmwareVersion: strings.TrimSpace(drive.FirmwareRevision),
			Capacity:        drive.Size,
			Healthy:         true, // Default to healthy
			Attributes:      make(map[string]string),
			DetailedAttribs: make([]types.SMARTAttribute, 0),
		}

		// Determine drive type and rotation rate from MediaType
		if drive.MediaType != "" {
			info.Attributes["MediaType"] = drive.MediaType
			// Try to determine if SSD
			mediaLower := strings.ToLower(drive.MediaType)
			if strings.Contains(mediaLower, "ssd") || strings.Contains(mediaLower, "solid") {
				info.RotationRate = 0
				info.Attributes["DriveType"] = "SSD"
			} else if strings.Contains(mediaLower, "fixed") {
				info.Attributes["DriveType"] = "HDD"
			}
		}

		// Add disk geometry information
		if drive.BytesPerSector > 0 {
			info.Attributes["BytesPerSector"] = fmt.Sprintf("%d", drive.BytesPerSector)
		}
		if drive.TotalSectors > 0 {
			info.Attributes["TotalSectors"] = fmt.Sprintf("%d", drive.TotalSectors)
		}
		if drive.TotalCylinders > 0 {
			info.Attributes["TotalCylinders"] = fmt.Sprintf("%d", drive.TotalCylinders)
		}
		if drive.TotalHeads > 0 {
			info.Attributes["TotalHeads"] = fmt.Sprintf("%d", drive.TotalHeads)
		}

		// Check disk health status
		if drive.Status != "" {
			info.Attributes["Status"] = drive.Status
			if drive.Status != "OK" {
				info.Healthy = false
			}
		}

		// Try to get SMART health prediction (may not be available on VMs or some drives)
		healthy, smartAvailable := checkDiskHealth(drive.DeviceID)
		if smartAvailable {
			if !healthy {
				info.Healthy = false
			}

			// Get SMART attributes if available
			attributes, detailedAttribs := getSMARTAttributes(drive.DeviceID)
			if len(attributes) > 0 {
				for k, v := range attributes {
					info.Attributes[k] = v
				}

				// Extract common SMART values
				if temp, ok := attributes["Temperature"]; ok {
					if tempInt, err := strconv.Atoi(temp); err == nil {
						info.Temperature = tempInt
					}
				}
				if poh, ok := attributes["Power_On_Hours"]; ok {
					if pohInt, err := strconv.ParseUint(poh, 10, 64); err == nil {
						info.PowerOnHours = pohInt
					}
				}
				if pcc, ok := attributes["Power_Cycle_Count"]; ok {
					if pccInt, err := strconv.ParseUint(pcc, 10, 64); err == nil {
						info.PowerCycleCount = pccInt
					}
				}
			}

			// Add detailed attributes
			if len(detailedAttribs) > 0 {
				info.DetailedAttribs = detailedAttribs

				// Perform health assessment based on detailed attributes
				info.HealthAssessment = assessDriveHealth(detailedAttribs, info.Temperature)
			}
		} else {
			info.Attributes["SMART"] = "Not Available"
		}

		// Add interface and media type information
		if drive.InterfaceType != "" {
			info.Attributes["Interface"] = drive.InterfaceType
		}
		if drive.MediaType != "" {
			info.Attributes["MediaType"] = drive.MediaType
		}

		smartData = append(smartData, info)
	}

	return smartData
}

// checkDiskHealth queries WMI for disk failure prediction
// Returns (healthy, smartAvailable)
func checkDiskHealth(deviceID string) (bool, bool) {
	// Convert deviceID format for WMI namespace
	// Example: \\.\PHYSICALDRIVE0 -> PhysicalDrive0
	devicePath := strings.ReplaceAll(deviceID, `\\.\`, "")
	devicePath = strings.ReplaceAll(devicePath, `\`, "")

	var status []MSStorageDriver_FailurePredictStatus
	query := fmt.Sprintf(`SELECT PredictFailure FROM MSStorageDriver_FailurePredictStatus WHERE InstanceName LIKE '%%%s%%'`, devicePath)

	err := wmi.QueryNamespace(query, &status, `root\wmi`)
	if err != nil || len(status) == 0 {
		return true, false // SMART not available, assume healthy
	}

	return !status[0].PredictFailure, true
} // getSMARTAttributes retrieves SMART attribute values
func getSMARTAttributes(deviceID string) (map[string]string, []types.SMARTAttribute) {
	attributes := make(map[string]string)
	detailedAttribs := make([]types.SMARTAttribute, 0)

	// Convert deviceID format
	devicePath := strings.ReplaceAll(deviceID, `\\.\`, "")
	devicePath = strings.ReplaceAll(devicePath, `\`, "")

	var data []MSStorageDriver_FailurePredictData
	query := fmt.Sprintf(`SELECT VendorSpecific FROM MSStorageDriver_FailurePredictData WHERE InstanceName LIKE '%%%s%%'`, devicePath)

	err := wmi.QueryNamespace(query, &data, `root\wmi`)
	if err != nil || len(data) == 0 {
		return attributes, detailedAttribs
	}

	// Get thresholds for comparison
	thresholds := getSMARTThresholds(deviceID)

	// Parse vendor-specific SMART data
	// The VendorSpecific field contains raw SMART data in a specific format
	// Each SMART attribute is 12 bytes: ID (1), Flags (2), Current (1), Worst (1), Raw (6), Reserved (1)
	vendorData := data[0].VendorSpecific
	if len(vendorData) < 362 { // SMART data should be at least 362 bytes
		return attributes, detailedAttribs
	}

	// SMART attributes start at offset 2
	for i := 2; i < len(vendorData)-12; i += 12 {
		if i+12 > len(vendorData) {
			break
		}

		id := vendorData[i]
		if id == 0 {
			continue // Skip empty entries
		}

		flags := uint16(vendorData[i+1]) | (uint16(vendorData[i+2]) << 8)
		current := vendorData[i+3]
		worst := vendorData[i+4]

		// Raw value is 6 bytes (little-endian)
		rawValue := uint64(vendorData[i+5]) |
			uint64(vendorData[i+6])<<8 |
			uint64(vendorData[i+7])<<16 |
			uint64(vendorData[i+8])<<24 |
			uint64(vendorData[i+9])<<32 |
			uint64(vendorData[i+10])<<40

		// Map common SMART attribute IDs to names
		attrName := getSMARTAttributeName(id)
		if attrName != "" {
			// Add to simple attributes map
			attributes[attrName] = fmt.Sprintf("%d", rawValue)
			attributes[attrName+"_Current"] = fmt.Sprintf("%d", current)
			attributes[attrName+"_Worst"] = fmt.Sprintf("%d", worst)

			// Get threshold for this attribute
			threshold := uint8(0)
			if thresh, ok := thresholds[id]; ok {
				threshold = thresh
			}

			// Determine attribute type based on flags
			attrType := "Old_age"
			if flags&0x01 != 0 {
				attrType = "Pre-fail"
			}

			updated := "Always"
			if flags&0x02 != 0 {
				updated = "Offline"
			}

			// Determine if failing
			whenFailed := "Never"
			if threshold > 0 && current <= threshold {
				whenFailed = "FAILING_NOW"
			} else if threshold > 0 && worst <= threshold {
				whenFailed = "In_the_past"
			}

			// Create human-readable raw string based on attribute
			rawString := formatRawValue(id, rawValue)

			// Add to detailed attributes
			detailedAttr := types.SMARTAttribute{
				ID:         id,
				Name:       attrName,
				Flag:       flags,
				Value:      current,
				Worst:      worst,
				Threshold:  threshold,
				RawValue:   rawValue,
				Type:       attrType,
				Updated:    updated,
				WhenFailed: whenFailed,
				RawString:  rawString,
			}
			detailedAttribs = append(detailedAttribs, detailedAttr)
		}
	}

	return attributes, detailedAttribs
}

// getSMARTAttributeName maps SMART attribute IDs to human-readable names
func getSMARTAttributeName(id uint8) string {
	smartAttrs := map[uint8]string{
		// Common attributes across all drives
		1:  "Read_Error_Rate",
		2:  "Throughput_Performance",
		3:  "Spin_Up_Time",
		4:  "Start_Stop_Count",
		5:  "Reallocated_Sector_Count",
		6:  "Read_Channel_Margin",
		7:  "Seek_Error_Rate",
		8:  "Seek_Time_Performance",
		9:  "Power_On_Hours",
		10: "Spin_Retry_Count",
		11: "Calibration_Retry_Count",
		12: "Power_Cycle_Count",
		13: "Soft_Read_Error_Rate",

		// SSD-specific attributes
		170: "Available_Reserved_Space",
		171: "SSD_Program_Fail_Count",
		172: "SSD_Erase_Fail_Count",
		173: "SSD_Wear_Leveling_Count",
		174: "Unexpected_Power_Loss_Count",
		175: "Power_Loss_Protection_Failure",
		176: "Erase_Fail_Count",
		177: "Wear_Range_Delta",
		178: "Used_Reserved_Block_Count_Chip",
		179: "Used_Reserved_Block_Count_Total",
		180: "Unused_Reserved_Block_Count_Total",
		181: "Program_Fail_Count_Total",
		182: "Erase_Fail_Count",
		183: "Runtime_Bad_Block",
		184: "End-to-End_Error",
		185: "Head_Stability",
		186: "Induced_Op-Vibration_Detection",
		187: "Reported_Uncorrectable_Errors",
		188: "Command_Timeout",
		189: "High_Fly_Writes",
		190: "Temperature_Difference",
		191: "G-Sense_Error_Rate",
		192: "Power-Off_Retract_Count",
		193: "Load_Cycle_Count",
		194: "Temperature",
		195: "Hardware_ECC_Recovered",
		196: "Reallocation_Event_Count",
		197: "Current_Pending_Sector",
		198: "Offline_Uncorrectable",
		199: "UDMA_CRC_Error_Count",
		200: "Multi_Zone_Error_Rate",
		201: "Soft_Read_Error_Rate",
		202: "Data_Address_Mark_Errors",
		203: "Run_Out_Cancel",
		204: "Soft_ECC_Correction",
		205: "Thermal_Asperity_Rate",
		206: "Flying_Height",
		207: "Spin_High_Current",
		208: "Spin_Buzz",
		209: "Offline_Seek_Performance",
		210: "Vibration_During_Write",
		211: "Vibration_During_Write_Chip",
		212: "Shock_During_Write",

		// Western Digital attributes
		220: "Disk_Shift",
		221: "G-Sense_Error_Rate",
		222: "Loaded_Hours",
		223: "Load_Retry_Count",
		224: "Load_Friction",
		225: "Load_Cycle_Count",
		226: "Load-in_Time",
		227: "Torque_Amplification_Count",
		228: "Power-Off_Retract_Count",

		// Seagate/Samsung/Crucial/Intel/Micron/SandForce SSD attributes
		230: "Drive_Life_Protection_Status",
		231: "SSD_Life_Left",
		232: "Available_Reserved_Space",
		233: "Media_Wearout_Indicator",
		234: "Average_Erase_Count",
		235: "Good_Block_Count",
		240: "Head_Flying_Hours",
		241: "Total_LBAs_Written",
		242: "Total_LBAs_Read",
		243: "Total_LBAs_Written_Expanded",
		244: "Total_LBAs_Read_Expanded",
		245: "Reserved_245",
		246: "Total_LBAs_Written_Expanded",
		247: "Reserved_247",
		248: "Reserved_248",
		249: "NAND_Writes_1GiB",
		250: "Read_Error_Retry_Rate",
		251: "Minimum_Spares_Remaining",
		252: "Newly_Added_Bad_Flash_Block",
		253: "Reserved_253",
		254: "Free_Fall_Protection",
	}

	if name, ok := smartAttrs[id]; ok {
		return name
	}
	// Return generic name for unknown attributes
	return fmt.Sprintf("Attribute_%d", id)
}

// getSMARTThresholds retrieves SMART threshold values
func getSMARTThresholds(deviceID string) map[uint8]uint8 {
	thresholds := make(map[uint8]uint8)

	// Convert deviceID format
	devicePath := strings.ReplaceAll(deviceID, `\\.\`, "")
	devicePath = strings.ReplaceAll(devicePath, `\`, "")

	var data []MSStorageDriver_FailurePredictThresholds
	query := fmt.Sprintf(`SELECT VendorSpecific FROM MSStorageDriver_FailurePredictThresholds WHERE InstanceName LIKE '%%%s%%'`, devicePath)

	err := wmi.QueryNamespace(query, &data, `root\wmi`)
	if err != nil || len(data) == 0 {
		return thresholds
	}

	// Parse threshold data (similar structure to attributes)
	vendorData := data[0].VendorSpecific
	if len(vendorData) < 362 {
		return thresholds
	}

	// Thresholds start at offset 2, same format as attributes
	for i := 2; i < len(vendorData)-12; i += 12 {
		if i+12 > len(vendorData) {
			break
		}

		id := vendorData[i]
		if id == 0 {
			continue
		}

		threshold := vendorData[i+1] // Threshold value
		thresholds[id] = threshold
	}

	return thresholds
}

// formatRawValue formats raw SMART values in human-readable format
func formatRawValue(attributeID uint8, rawValue uint64) string {
	switch attributeID {
	case 9: // Power On Hours
		days := rawValue / 24
		hours := rawValue % 24
		return fmt.Sprintf("%d hours (%d days, %d hours)", rawValue, days, hours)
	case 12: // Power Cycle Count
		return fmt.Sprintf("%d", rawValue)
	case 194: // Temperature
		// Temperature is usually in the lower byte
		temp := rawValue & 0xFF
		if temp > 0 && temp < 200 {
			return fmt.Sprintf("%d Celsius", temp)
		}
		return fmt.Sprintf("%d", rawValue)
	case 4, 192, 193: // Start/Stop, Retract, Load counts
		return fmt.Sprintf("%d", rawValue)
	case 5, 196, 197, 198: // Sector counts
		if rawValue > 0 {
			return fmt.Sprintf("%d sectors", rawValue)
		}
		return "0"
	case 241, 242: // Total LBAs Written/Read
		// Convert to approximate GB (512 byte sectors)
		gb := float64(rawValue) * 512.0 / 1000000000.0
		if gb > 1000 {
			return fmt.Sprintf("%.2f TB", gb/1000.0)
		}
		return fmt.Sprintf("%.2f GB", gb)
	case 231: // SSD Life Left
		return fmt.Sprintf("%d%%", rawValue)
	case 233: // Media Wearout Indicator
		return fmt.Sprintf("%d%%", rawValue)
	default:
		return fmt.Sprintf("%d", rawValue)
	}
}

// assessDriveHealth performs health assessment based on SMART attributes
func assessDriveHealth(attributes []types.SMARTAttribute, temperature int) *types.SMARTHealthStatus {
	health := &types.SMARTHealthStatus{
		Passed:            true,
		OverallAssessment: "PASS",
		FailingAttributes: make([]string, 0),
		WarningAttributes: make([]string, 0),
	}

	criticalAttributes := map[string]bool{
		"Reallocated_Sector_Count":      true,
		"Current_Pending_Sector":        true,
		"Offline_Uncorrectable":         true,
		"Reported_Uncorrectable_Errors": true,
		"SSD_Program_Fail_Count":        true,
		"SSD_Erase_Fail_Count":          true,
		"Runtime_Bad_Block":             true,
	}

	warningAttributes := map[string]bool{
		"Reallocation_Event_Count": true,
		"Spin_Retry_Count":         true,
		"UDMA_CRC_Error_Count":     true,
		"Command_Timeout":          true,
	}

	for _, attr := range attributes {
		// Check for failing attributes
		if attr.WhenFailed == "FAILING_NOW" {
			health.Passed = false
			health.OverallAssessment = "FAIL"
			health.FailingAttributes = append(health.FailingAttributes,
				fmt.Sprintf("%s (Value: %d, Threshold: %d)", attr.Name, attr.Value, attr.Threshold))
		}

		// Check critical attributes with non-zero raw values
		if criticalAttributes[attr.Name] && attr.RawValue > 0 {
			health.WarningAttributes = append(health.WarningAttributes,
				fmt.Sprintf("%s = %d", attr.Name, attr.RawValue))
			if health.OverallAssessment == "PASS" {
				health.OverallAssessment = "WARN"
			}
		}

		// Check warning attributes
		if warningAttributes[attr.Name] && attr.RawValue > 0 {
			health.WarningAttributes = append(health.WarningAttributes,
				fmt.Sprintf("%s = %d", attr.Name, attr.RawValue))
		}

		// Check SSD wear level
		if attr.Name == "SSD_Life_Left" || attr.Name == "Media_Wearout_Indicator" {
			health.PercentUsed = 100.0 - float64(attr.RawValue)
			if attr.RawValue < 10 {
				health.CriticalWarning = "SSD wear level critical (less than 10% life remaining)"
				if health.OverallAssessment != "FAIL" {
					health.OverallAssessment = "WARN"
				}
			} else if attr.RawValue < 20 {
				health.WarningAttributes = append(health.WarningAttributes,
					fmt.Sprintf("SSD life remaining: %d%%", attr.RawValue))
			}
		}
	}

	// Check temperature
	if temperature > 0 {
		if temperature >= 70 {
			health.TemperatureStatus = "CRITICAL"
			health.WarningAttributes = append(health.WarningAttributes,
				fmt.Sprintf("Temperature critical: %d°C", temperature))
			if health.OverallAssessment != "FAIL" {
				health.OverallAssessment = "WARN"
			}
		} else if temperature >= 60 {
			health.TemperatureStatus = "HIGH"
			health.WarningAttributes = append(health.WarningAttributes,
				fmt.Sprintf("Temperature high: %d°C", temperature))
		} else if temperature >= 45 {
			health.TemperatureStatus = "WARM"
		} else {
			health.TemperatureStatus = "NORMAL"
		}
	}

	return health
}
