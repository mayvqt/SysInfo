//go:build linux

package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
)

const (
	powerSupplyPath = "/sys/class/power_supply"
)

// CollectBattery collects battery information on Linux
func CollectBattery() (*types.BatteryData, error) {
	data := &types.BatteryData{
		Present:   false,
		Batteries: []types.BatteryInfo{},
		OnBattery: false,
	}

	// Check if power_supply directory exists
	if _, err := os.Stat(powerSupplyPath); os.IsNotExist(err) {
		return data, nil // No battery information available
	}

	// Read all power supply devices
	entries, err := os.ReadDir(powerSupplyPath)
	if err != nil {
		return data, fmt.Errorf("failed to read power supply directory: %w", err)
	}

	var totalCapacity uint64
	acOnline := false

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		devicePath := filepath.Join(powerSupplyPath, entry.Name())
		deviceType, err := readSysFile(filepath.Join(devicePath, "type"))
		if err != nil {
			continue
		}

		deviceType = strings.TrimSpace(deviceType)

		// Check if AC adapter
		if deviceType == "Mains" {
			online, err := readSysFile(filepath.Join(devicePath, "online"))
			if err == nil && strings.TrimSpace(online) == "1" {
				acOnline = true
			}
			continue
		}

		// Only process battery devices
		if deviceType != "Battery" {
			continue
		}

		battery, err := readBatteryInfo(devicePath, entry.Name())
		if err != nil {
			continue
		}

		data.Batteries = append(data.Batteries, battery)
		data.Present = true
		totalCapacity += battery.Capacity
	}

	if len(data.Batteries) > 0 {
		data.TotalCapacity = totalCapacity
		data.OnBattery = !acOnline
	}

	return data, nil
}

// readBatteryInfo reads battery information from a specific battery device
func readBatteryInfo(devicePath, name string) (types.BatteryInfo, error) {
	battery := types.BatteryInfo{
		Name:          name,
		TimeToEmpty:   -1,
		TimeToFull:    -1,
		TimeRemaining: -1,
	}

	// Read status
	if status, err := readSysFile(filepath.Join(devicePath, "status")); err == nil {
		battery.State = strings.TrimSpace(status)
		battery.IsCharging = battery.State == "Charging"
		battery.IsDischarging = battery.State == "Discharging"
	}

	// Read manufacturer/vendor
	if vendor, err := readSysFile(filepath.Join(devicePath, "manufacturer")); err == nil {
		battery.Vendor = strings.TrimSpace(vendor)
	}

	// Read model
	if model, err := readSysFile(filepath.Join(devicePath, "model_name")); err == nil {
		battery.Model = strings.TrimSpace(model)
	}

	// Read serial number
	if serial, err := readSysFile(filepath.Join(devicePath, "serial_number")); err == nil {
		battery.SerialNumber = strings.TrimSpace(serial)
	}

	// Read technology
	if tech, err := readSysFile(filepath.Join(devicePath, "technology")); err == nil {
		battery.Technology = strings.TrimSpace(tech)
	}

	// Read cycle count
	if cycleStr, err := readSysFile(filepath.Join(devicePath, "cycle_count")); err == nil {
		if cycle, err := strconv.ParseUint(strings.TrimSpace(cycleStr), 10, 64); err == nil {
			battery.CycleCount = cycle
		}
	}

	// Read capacity information (in µWh - microwatt-hours)
	// Convert to mWh (milliwatt-hours) by dividing by 1000
	if capacityStr, err := readSysFile(filepath.Join(devicePath, "energy_full_design")); err == nil {
		if capacity, err := strconv.ParseUint(strings.TrimSpace(capacityStr), 10, 64); err == nil {
			battery.Capacity = capacity / 1000 // Convert µWh to mWh
		}
	} else if capacityStr, err := readSysFile(filepath.Join(devicePath, "charge_full_design")); err == nil {
		// Some batteries report in µAh (microampere-hours)
		if capacity, err := strconv.ParseUint(strings.TrimSpace(capacityStr), 10, 64); err == nil {
			// Read voltage to convert to energy
			if voltageStr, err := readSysFile(filepath.Join(devicePath, "voltage_now")); err == nil {
				if voltage, err := strconv.ParseUint(strings.TrimSpace(voltageStr), 10, 64); err == nil {
					// energy (µWh) = charge (µAh) * voltage (µV) / 1000000
					// then convert to mWh
					battery.Capacity = (capacity * voltage) / 1000000 / 1000
				}
			}
		}
	}

	// Read full capacity
	if fullStr, err := readSysFile(filepath.Join(devicePath, "energy_full")); err == nil {
		if full, err := strconv.ParseUint(strings.TrimSpace(fullStr), 10, 64); err == nil {
			battery.CapacityFull = full / 1000
			battery.EnergyFull = full / 1000
		}
	} else if fullStr, err := readSysFile(filepath.Join(devicePath, "charge_full")); err == nil {
		if full, err := strconv.ParseUint(strings.TrimSpace(fullStr), 10, 64); err == nil {
			if voltageStr, err := readSysFile(filepath.Join(devicePath, "voltage_now")); err == nil {
				if voltage, err := strconv.ParseUint(strings.TrimSpace(voltageStr), 10, 64); err == nil {
					battery.CapacityFull = (full * voltage) / 1000000 / 1000
					battery.EnergyFull = battery.CapacityFull
				}
			}
		}
	}

	// Read current capacity
	if nowStr, err := readSysFile(filepath.Join(devicePath, "energy_now")); err == nil {
		if now, err := strconv.ParseUint(strings.TrimSpace(nowStr), 10, 64); err == nil {
			battery.CapacityNow = now / 1000
			battery.EnergyNow = now / 1000
		}
	} else if nowStr, err := readSysFile(filepath.Join(devicePath, "charge_now")); err == nil {
		if now, err := strconv.ParseUint(strings.TrimSpace(nowStr), 10, 64); err == nil {
			if voltageStr, err := readSysFile(filepath.Join(devicePath, "voltage_now")); err == nil {
				if voltage, err := strconv.ParseUint(strings.TrimSpace(voltageStr), 10, 64); err == nil {
					battery.CapacityNow = (now * voltage) / 1000000 / 1000
					battery.EnergyNow = battery.CapacityNow
				}
			}
		}
	}

	// Calculate charge level percentage
	if battery.CapacityFull > 0 {
		battery.ChargeLevel = float64(battery.CapacityNow) / float64(battery.CapacityFull) * 100.0
		if battery.ChargeLevel > 100.0 {
			battery.ChargeLevel = 100.0
		}
	}

	// Calculate health percentage
	if battery.Capacity > 0 && battery.CapacityFull > 0 {
		battery.Health = float64(battery.CapacityFull) / float64(battery.Capacity) * 100.0
		if battery.Health > 100.0 {
			battery.Health = 100.0
		}
	}

	// Read power consumption/charge rate
	if powerStr, err := readSysFile(filepath.Join(devicePath, "power_now")); err == nil {
		if power, err := strconv.ParseUint(strings.TrimSpace(powerStr), 10, 64); err == nil {
			battery.PowerNow = power / 1000 // Convert µW to mW
		}
	} else if currentStr, err := readSysFile(filepath.Join(devicePath, "current_now")); err == nil {
		// Calculate power from current and voltage
		if current, err := strconv.ParseInt(strings.TrimSpace(currentStr), 10, 64); err == nil {
			if voltageStr, err := readSysFile(filepath.Join(devicePath, "voltage_now")); err == nil {
				if voltage, err := strconv.ParseUint(strings.TrimSpace(voltageStr), 10, 64); err == nil {
					// power (µW) = abs(current (µA)) * voltage (µV) / 1000000
					absCurrent := current
					if absCurrent < 0 {
						absCurrent = -absCurrent
					}
					battery.PowerNow = uint64(absCurrent) * voltage / 1000000 / 1000 // Convert to mW
					battery.Current = current / 1000                                 // Convert µA to mA
				}
			}
		}
	}

	// Read voltage
	if voltageStr, err := readSysFile(filepath.Join(devicePath, "voltage_now")); err == nil {
		if voltage, err := strconv.ParseUint(strings.TrimSpace(voltageStr), 10, 64); err == nil {
			battery.Voltage = float64(voltage) / 1000000.0 // Convert µV to V
		}
	}

	// Read minimum voltage
	if minVoltageStr, err := readSysFile(filepath.Join(devicePath, "voltage_min_design")); err == nil {
		if minVoltage, err := strconv.ParseUint(strings.TrimSpace(minVoltageStr), 10, 64); err == nil {
			battery.VoltageMin = float64(minVoltage) / 1000000.0
		}
	}

	// Calculate time remaining
	if battery.PowerNow > 0 {
		if battery.IsDischarging && battery.CapacityNow > 0 {
			// Time to empty = energy remaining / power consumption
			// energy in mWh, power in mW, result in hours
			hours := float64(battery.CapacityNow) / float64(battery.PowerNow)
			battery.TimeToEmpty = int64(hours * 60) // Convert to minutes
			battery.TimeRemaining = battery.TimeToEmpty
		} else if battery.IsCharging && battery.CapacityFull > battery.CapacityNow {
			// Time to full = energy needed / charge rate
			energyNeeded := battery.CapacityFull - battery.CapacityNow
			hours := float64(energyNeeded) / float64(battery.PowerNow)
			battery.TimeToFull = int64(hours * 60)
			battery.TimeRemaining = battery.TimeToFull
		}
	}

	return battery, nil
}

// readSysFile reads a single value from a sysfs file
func readSysFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
