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
	DeviceID      string
	Model         string
	SerialNumber  string
	Size          uint64
	Status        string
	InterfaceType string
	MediaType     string
	Partitions    uint32
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
	query := "SELECT DeviceID, Model, SerialNumber, Size, Status, InterfaceType, MediaType, Partitions FROM Win32_DiskDrive"
	err := wmi.Query(query, &drives)
	if err != nil {
		return smartData
	}

	// Process each drive found
	for _, drive := range drives {
		info := types.SMARTInfo{
			Device:      drive.DeviceID,
			DeviceModel: strings.TrimSpace(drive.Model),
			Serial:      strings.TrimSpace(drive.SerialNumber),
			Capacity:    drive.Size,
			Healthy:     true, // Default to healthy
			Attributes:  make(map[string]string),
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
			attributes := getSMARTAttributes(drive.DeviceID)
			if attributes != nil && len(attributes) > 0 {
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
func getSMARTAttributes(deviceID string) map[string]string {
	attributes := make(map[string]string)

	// Convert deviceID format
	devicePath := strings.ReplaceAll(deviceID, `\\.\`, "")
	devicePath = strings.ReplaceAll(devicePath, `\`, "")

	var data []MSStorageDriver_FailurePredictData
	query := fmt.Sprintf(`SELECT VendorSpecific FROM MSStorageDriver_FailurePredictData WHERE InstanceName LIKE '%%%s%%'`, devicePath)

	err := wmi.QueryNamespace(query, &data, `root\wmi`)
	if err != nil || len(data) == 0 {
		return attributes
	}

	// Parse vendor-specific SMART data
	// The VendorSpecific field contains raw SMART data in a specific format
	// Each SMART attribute is 12 bytes: ID (1), Flags (2), Current (1), Worst (1), Raw (6), Reserved (1)
	vendorData := data[0].VendorSpecific
	if len(vendorData) < 362 { // SMART data should be at least 362 bytes
		return attributes
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
			attributes[attrName] = fmt.Sprintf("%d", rawValue)
			attributes[attrName+"_Current"] = fmt.Sprintf("%d", current)
			attributes[attrName+"_Worst"] = fmt.Sprintf("%d", worst)
		}
	}

	return attributes
}

// getSMARTAttributeName maps SMART attribute IDs to human-readable names
func getSMARTAttributeName(id uint8) string {
	smartAttrs := map[uint8]string{
		1:   "Read_Error_Rate",
		5:   "Reallocated_Sector_Count",
		9:   "Power_On_Hours",
		12:  "Power_Cycle_Count",
		192: "Power-Off_Retract_Count",
		193: "Load_Cycle_Count",
		194: "Temperature",
		196: "Reallocation_Event_Count",
		197: "Current_Pending_Sector",
		198: "Offline_Uncorrectable",
		199: "UDMA_CRC_Error_Count",
		200: "Multi_Zone_Error_Rate",
		241: "Total_LBAs_Written",
		242: "Total_LBAs_Read",
	}

	if name, ok := smartAttrs[id]; ok {
		return name
	}
	return ""
}
