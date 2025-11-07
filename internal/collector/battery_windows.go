//go:build windows

package collector

import (
	"fmt"
	"math"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/yusufpapurcu/wmi"
)

// Win32Battery represents the Win32_Battery WMI class
type Win32Battery struct {
	Name                     string
	DeviceID                 string
	Description              string
	Chemistry                uint16
	DesignCapacity           uint32 // in mWh
	FullChargeCapacity       uint32 // in mWh
	DesignVoltage            uint64 // in mV
	BatteryStatus            uint16
	EstimatedChargeRemaining uint16 // Percentage
	EstimatedRunTime         uint32 // in minutes
	TimeToFullCharge         uint32 // in minutes
	PowerManagementSupported bool
}

// Win32PortableBattery represents the Win32_PortableBattery WMI class (more detailed)
type Win32PortableBattery struct {
	DeviceID            string
	Name                string
	Manufacturer        string
	ManufactureDate     string
	SerialNumber        string
	Chemistry           uint16
	DesignCapacity      uint32
	FullChargeCapacity  uint32
	DesignVoltage       uint64
	MaxRechargeTime     uint32
	TimeOnBattery       uint32
	TimeToFullCharge    uint32
	SmartBatteryVersion string
	Location            string
	CapacityMultiplier  uint16
}

// BatteryStaticData represents the BatteryStaticData WMI class
type BatteryStaticData struct {
	Tag              uint32
	ManufactureName  string
	DeviceName       string
	SerialNumber     string
	Chemistry        uint16
	DesignedCapacity uint32
	DefaultAlert1    uint32
	DefaultAlert2    uint32
	CriticalBias     uint32
	CycleCount       uint32
}

// BatteryFullChargedCapacity represents the BatteryFullChargedCapacity WMI class
type BatteryFullChargedCapacity struct {
	Tag                 uint32
	FullChargedCapacity uint32
}

// BatteryStatus represents the BatteryStatus WMI class
type BatteryStatus struct {
	Tag               uint32
	Voltage           uint32 // in mV
	ChargeRate        int32  // in mW (negative = discharging)
	DischargeRate     uint32 // in mW
	RemainingCapacity uint32 // in mWh
	PowerOnline       bool
	Charging          bool
	Discharging       bool
	Critical          bool
	Temperature       uint32 // in tenths of degrees Kelvin
}

// CollectBattery collects battery information on Windows
func CollectBattery() (*types.BatteryData, error) {
	data := &types.BatteryData{
		Present:   false,
		Batteries: []types.BatteryInfo{},
		OnBattery: false,
	}

	// Query Win32_Battery for basic information
	var batteries []Win32Battery
	query := "SELECT * FROM Win32_Battery"
	err := wmi.Query(query, &batteries)
	if err != nil || len(batteries) == 0 {
		// No batteries found or WMI error
		return data, nil
	}

	data.Present = true

	// Query more detailed battery information from WMI namespace root\wmi
	var staticDataList []BatteryStaticData
	var statusList []BatteryStatus
	var fullCapacityList []BatteryFullChargedCapacity

	// These queries might fail on some systems, so we ignore errors
	_ = wmi.QueryNamespace("SELECT * FROM BatteryStaticData", &staticDataList, "root\\wmi")
	_ = wmi.QueryNamespace("SELECT * FROM BatteryStatus", &statusList, "root\\wmi")
	_ = wmi.QueryNamespace("SELECT * FROM BatteryFullChargedCapacity", &fullCapacityList, "root\\wmi")

	var totalCapacity uint64

	for i, winBatt := range batteries {
		battery := types.BatteryInfo{
			Name:          winBatt.Name,
			TimeToEmpty:   -1,
			TimeToFull:    -1,
			TimeRemaining: -1,
		}

		// Set battery technology/chemistry
		battery.Technology = getBatteryChemistry(winBatt.Chemistry)

		// Charge level from Win32_Battery
		battery.ChargeLevel = float64(winBatt.EstimatedChargeRemaining)

		// Design capacity (convert from mWh to mWh - already in correct unit)
		if winBatt.DesignCapacity > 0 {
			battery.Capacity = uint64(winBatt.DesignCapacity)
		}

		// Full charge capacity
		if winBatt.FullChargeCapacity > 0 {
			battery.CapacityFull = uint64(winBatt.FullChargeCapacity)
		}

		// Voltage
		if winBatt.DesignVoltage > 0 {
			battery.VoltageMin = float64(winBatt.DesignVoltage) / 1000.0 // mV to V
		}

		// Battery status
		battery.State, battery.IsCharging, battery.IsDischarging = getBatteryStatus(winBatt.BatteryStatus)

		// Time estimates from Win32_Battery
		if winBatt.EstimatedRunTime != math.MaxUint32 && winBatt.EstimatedRunTime > 0 {
			battery.TimeToEmpty = int64(winBatt.EstimatedRunTime)
			if battery.IsDischarging {
				battery.TimeRemaining = battery.TimeToEmpty
			}
		}
		if winBatt.TimeToFullCharge != math.MaxUint32 && winBatt.TimeToFullCharge > 0 {
			battery.TimeToFull = int64(winBatt.TimeToFullCharge)
			if battery.IsCharging {
				battery.TimeRemaining = battery.TimeToFull
			}
		}

		// Try to get more detailed information from root\wmi namespace
		if i < len(staticDataList) {
			staticData := staticDataList[i]
			if staticData.ManufactureName != "" {
				battery.Vendor = staticData.ManufactureName
			}
			if staticData.DeviceName != "" {
				battery.Model = staticData.DeviceName
			}
			if staticData.SerialNumber != "" {
				battery.SerialNumber = staticData.SerialNumber
			}
			if staticData.CycleCount > 0 {
				battery.CycleCount = uint64(staticData.CycleCount)
			}
			if staticData.DesignedCapacity > 0 {
				battery.Capacity = uint64(staticData.DesignedCapacity)
			}
		}

		// Get current status details
		if i < len(statusList) {
			status := statusList[i]

			// Current voltage
			if status.Voltage > 0 {
				battery.Voltage = float64(status.Voltage) / 1000.0 // mV to V
			}

			// Remaining capacity
			if status.RemainingCapacity > 0 {
				battery.CapacityNow = uint64(status.RemainingCapacity)
			}

			// Power consumption/charge rate
			if status.ChargeRate != 0 {
				if status.ChargeRate < 0 {
					// Discharging
					battery.PowerNow = uint64(-status.ChargeRate)
				} else {
					// Charging
					battery.PowerNow = uint64(status.ChargeRate)
				}
			} else if status.DischargeRate > 0 {
				battery.PowerNow = uint64(status.DischargeRate)
			}

			// Temperature (in tenths of Kelvin)
			if status.Temperature > 0 {
				battery.Temperature = (float64(status.Temperature) / 10.0) - 273.15
			}

			// Update charging/power status
			if status.PowerOnline {
				data.OnBattery = false
			} else {
				data.OnBattery = true
			}

			if status.Charging {
				battery.IsCharging = true
				battery.State = "Charging"
			} else if status.Discharging {
				battery.IsDischarging = true
				battery.State = "Discharging"
			}
		}

		// Get full charged capacity
		if i < len(fullCapacityList) {
			fullCap := fullCapacityList[i]
			if fullCap.FullChargedCapacity > 0 {
				battery.CapacityFull = uint64(fullCap.FullChargedCapacity)
			}
		}

		// Calculate charge level if we have capacity data
		if battery.CapacityFull > 0 && battery.CapacityNow > 0 {
			calculatedCharge := float64(battery.CapacityNow) / float64(battery.CapacityFull) * 100.0
			// Use calculated value if it seems more accurate
			if battery.ChargeLevel == 0 || (calculatedCharge > 0 && calculatedCharge <= 100) {
				battery.ChargeLevel = calculatedCharge
			}
		}

		// Calculate health
		if battery.Capacity > 0 && battery.CapacityFull > 0 {
			battery.Health = float64(battery.CapacityFull) / float64(battery.Capacity) * 100.0
			if battery.Health > 100.0 {
				battery.Health = 100.0
			}
		}

		// Recalculate time estimates if we have power data
		if battery.PowerNow > 0 {
			if battery.IsDischarging && battery.CapacityNow > 0 && battery.TimeToEmpty < 0 {
				hours := float64(battery.CapacityNow) / float64(battery.PowerNow)
				battery.TimeToEmpty = int64(hours * 60)
				battery.TimeRemaining = battery.TimeToEmpty
			} else if battery.IsCharging && battery.CapacityFull > battery.CapacityNow && battery.TimeToFull < 0 {
				energyNeeded := battery.CapacityFull - battery.CapacityNow
				hours := float64(energyNeeded) / float64(battery.PowerNow)
				battery.TimeToFull = int64(hours * 60)
				battery.TimeRemaining = battery.TimeToFull
			}
		}

		data.Batteries = append(data.Batteries, battery)
		totalCapacity += battery.Capacity
	}

	data.TotalCapacity = totalCapacity

	return data, nil
}

// getBatteryChemistry converts the WMI chemistry code to a string
func getBatteryChemistry(chemistry uint16) string {
	switch chemistry {
	case 1:
		return "Other"
	case 2:
		return "Unknown"
	case 3:
		return "Lead Acid"
	case 4:
		return "Nickel Cadmium"
	case 5:
		return "Nickel Metal Hydride"
	case 6:
		return "Lithium-ion"
	case 7:
		return "Zinc Air"
	case 8:
		return "Lithium Polymer"
	default:
		return fmt.Sprintf("Unknown (%d)", chemistry)
	}
}

// getBatteryStatus converts the WMI battery status code to state, charging, and discharging bools
func getBatteryStatus(status uint16) (string, bool, bool) {
	// BatteryStatus values:
	// 1 = Other, 2 = Unknown, 3 = Fully Charged, 4 = Low, 5 = Critical
	// 6 = Charging, 7 = Charging and High, 8 = Charging and Low, 9 = Charging and Critical
	// 10 = Undefined, 11 = Partially Charged

	switch status {
	case 1:
		return "Other", false, false
	case 2:
		return "Unknown", false, false
	case 3:
		return "Full", false, false
	case 4:
		return "Low", false, true
	case 5:
		return "Critical", false, true
	case 6, 7, 8, 9:
		return "Charging", true, false
	case 10:
		return "Undefined", false, false
	case 11:
		return "Idle", false, false
	default:
		return fmt.Sprintf("Unknown (%d)", status), false, false
	}
}
