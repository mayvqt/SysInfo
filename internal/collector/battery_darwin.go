//go:build darwin

package collector

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
)

// CollectBattery collects battery information on macOS
func CollectBattery() (*types.BatteryData, error) {
	data := &types.BatteryData{
		Present:   false,
		Batteries: []types.BatteryInfo{},
		OnBattery: false,
	}

	// Use pmset to get battery information
	pmsetOutput, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return data, nil // No battery information available
	}

	pmsetStr := string(pmsetOutput)

	// Check if running on battery
	if strings.Contains(pmsetStr, "'Battery Power'") || strings.Contains(pmsetStr, "discharging") {
		data.OnBattery = true
	}

	// Parse pmset output for basic info
	battery := types.BatteryInfo{
		Name:          "InternalBattery",
		TimeToEmpty:   -1,
		TimeToFull:    -1,
		TimeRemaining: -1,
	}

	// Extract charge percentage from pmset
	// Format: "Now drawing from 'Battery Power'  -InternalBattery-0 (id=1234567)  95%; discharging; 4:30 remaining present: true"
	chargeRegex := regexp.MustCompile(`(\d+)%`)
	if match := chargeRegex.FindStringSubmatch(pmsetStr); len(match) > 1 {
		if charge, err := strconv.ParseFloat(match[1], 64); err == nil {
			battery.ChargeLevel = charge
		}
	}

	// Extract time remaining from pmset
	// Format can be "4:30 remaining" or "0:30 remaining" or "(no estimate)"
	timeRegex := regexp.MustCompile(`(\d+):(\d+) remaining`)
	if match := timeRegex.FindStringSubmatch(pmsetStr); len(match) > 2 {
		hours, _ := strconv.Atoi(match[1])
		minutes, _ := strconv.Atoi(match[2])
		totalMinutes := int64(hours*60 + minutes)
		battery.TimeRemaining = totalMinutes
	}

	// Determine state from pmset
	if strings.Contains(pmsetStr, "charging") {
		battery.State = "Charging"
		battery.IsCharging = true
		if battery.TimeRemaining > 0 {
			battery.TimeToFull = battery.TimeRemaining
			battery.TimeRemaining = battery.TimeToFull
		}
	} else if strings.Contains(pmsetStr, "discharging") {
		battery.State = "Discharging"
		battery.IsDischarging = true
		if battery.TimeRemaining > 0 {
			battery.TimeToEmpty = battery.TimeRemaining
		}
	} else if strings.Contains(pmsetStr, "charged") || strings.Contains(pmsetStr, "finishing charge") {
		battery.State = "Full"
	} else if strings.Contains(pmsetStr, "AC attached; not charging") {
		battery.State = "Not charging"
	} else {
		battery.State = "Idle"
	}

	// Use ioreg to get detailed battery information
	ioregOutput, err := exec.Command("ioreg", "-r", "-c", "AppleSmartBattery").Output()
	if err == nil {
		ioregStr := string(ioregOutput)

		// Parse ioreg output
		battery.CycleCount = parseIoregUint64(ioregStr, "CycleCount")
		battery.Capacity = parseIoregUint64(ioregStr, "DesignCapacity")
		battery.CapacityFull = parseIoregUint64(ioregStr, "MaxCapacity")
		battery.CapacityNow = parseIoregUint64(ioregStr, "CurrentCapacity")

		// Voltage information (in millivolts)
		if voltage := parseIoregUint64(ioregStr, "Voltage"); voltage > 0 {
			battery.Voltage = float64(voltage) / 1000.0 // Convert mV to V
		}

		// Current (in milliamps, negative when discharging)
		if current := parseIoregInt64(ioregStr, "InstantAmperage"); current != 0 {
			battery.Current = current
			// Calculate power from current and voltage
			if battery.Voltage > 0 {
				battery.PowerNow = uint64(absInt64(current)) * uint64(battery.Voltage*1000) / 1000 // mW
			}
		}

		// Temperature (in hundredths of a degree Kelvin)
		if temp := parseIoregUint64(ioregStr, "Temperature"); temp > 0 {
			// Convert from centikelvins to Celsius
			battery.Temperature = (float64(temp) / 100.0) - 273.15
		}

		// Manufacturer
		if manufacturer := parseIoregString(ioregStr, "Manufacturer"); manufacturer != "" {
			battery.Vendor = manufacturer
		}

		// Device name
		if deviceName := parseIoregString(ioregStr, "DeviceName"); deviceName != "" {
			battery.Model = deviceName
		}

		// Serial number
		if serial := parseIoregString(ioregStr, "Serial"); serial != "" {
			battery.SerialNumber = serial
		}

		// Battery technology
		battery.Technology = "Li-ion" // macOS batteries are typically Li-ion or Li-poly

		// Calculate health if we have both capacities
		if battery.Capacity > 0 && battery.CapacityFull > 0 {
			battery.Health = float64(battery.CapacityFull) / float64(battery.Capacity) * 100.0
			if battery.Health > 100.0 {
				battery.Health = 100.0
			}
		}

		// Recalculate charge level from ioreg data if more accurate
		if battery.CapacityFull > 0 && battery.CapacityNow > 0 {
			chargeFromCapacity := float64(battery.CapacityNow) / float64(battery.CapacityFull) * 100.0
			// Use the ioreg value if pmset didn't give us a percentage
			if battery.ChargeLevel == 0 || chargeFromCapacity > 0 {
				battery.ChargeLevel = chargeFromCapacity
				if battery.ChargeLevel > 100.0 {
					battery.ChargeLevel = 100.0
				}
			}
		}

		// Recalculate time estimates if we have power data
		if battery.PowerNow > 0 && battery.TimeRemaining < 0 {
			if battery.IsDischarging && battery.CapacityNow > 0 {
				// Time = capacity / power (both in mWh and mW)
				hours := float64(battery.CapacityNow) / float64(battery.PowerNow)
				battery.TimeToEmpty = int64(hours * 60)
				battery.TimeRemaining = battery.TimeToEmpty
			} else if battery.IsCharging && battery.CapacityFull > battery.CapacityNow {
				energyNeeded := battery.CapacityFull - battery.CapacityNow
				hours := float64(energyNeeded) / float64(battery.PowerNow)
				battery.TimeToFull = int64(hours * 60)
				battery.TimeRemaining = battery.TimeToFull
			}
		}
	}

	// Only add battery if we got some data
	if battery.ChargeLevel > 0 || battery.CapacityNow > 0 {
		data.Present = true
		data.Batteries = append(data.Batteries, battery)
		data.TotalCapacity = battery.Capacity
	}

	return data, nil
}

// parseIoregUint64 extracts a uint64 value from ioreg output
func parseIoregUint64(output, key string) uint64 {
	pattern := fmt.Sprintf(`"%s" = (\d+)`, regexp.QuoteMeta(key))
	re := regexp.MustCompile(pattern)
	if match := re.FindStringSubmatch(output); len(match) > 1 {
		if val, err := strconv.ParseUint(match[1], 10, 64); err == nil {
			return val
		}
	}
	return 0
}

// parseIoregInt64 extracts an int64 value from ioreg output (can be negative)
func parseIoregInt64(output, key string) int64 {
	pattern := fmt.Sprintf(`"%s" = (-?\d+)`, regexp.QuoteMeta(key))
	re := regexp.MustCompile(pattern)
	if match := re.FindStringSubmatch(output); len(match) > 1 {
		if val, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			return val
		}
	}
	return 0
}

// parseIoregString extracts a string value from ioreg output
func parseIoregString(output, key string) string {
	pattern := fmt.Sprintf(`"%s" = "([^"]+)"`, regexp.QuoteMeta(key))
	re := regexp.MustCompile(pattern)
	if match := re.FindStringSubmatch(output); len(match) > 1 {
		return match[1]
	}
	return ""
}

// absInt64 returns the absolute value of an int64
func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
