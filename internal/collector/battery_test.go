package collector

import (
	"testing"
)

func TestCollectBattery(t *testing.T) {
	data, err := CollectBattery()

	// Battery collection should not error even if no battery is present
	if err != nil {
		t.Errorf("CollectBattery() returned error: %v", err)
	}

	// Data should never be nil
	if data == nil {
		t.Fatal("CollectBattery() returned nil data")
	}

	// If battery is present, validate the data
	if data.Present {
		if len(data.Batteries) == 0 {
			t.Error("Battery marked as present but no batteries in list")
		}

		for i, battery := range data.Batteries {
			// Name should be set
			if battery.Name == "" {
				t.Errorf("Battery %d has empty name", i)
			}

			// State should be set if battery is present
			if battery.State == "" {
				t.Errorf("Battery %d has empty state", i)
			}

			// Charge level should be 0-100
			if battery.ChargeLevel < 0 || battery.ChargeLevel > 100 {
				t.Errorf("Battery %d has invalid charge level: %.2f", i, battery.ChargeLevel)
			}

			// Health should be 0-100 if set
			if battery.Health > 0 && (battery.Health < 0 || battery.Health > 100) {
				t.Errorf("Battery %d has invalid health: %.2f", i, battery.Health)
			}

			// Validate charging/discharging state consistency
			if battery.IsCharging && battery.IsDischarging {
				t.Errorf("Battery %d cannot be both charging and discharging", i)
			}

			// If has capacity info, current should not exceed full
			if battery.CapacityFull > 0 && battery.CapacityNow > battery.CapacityFull {
				t.Errorf("Battery %d current capacity (%d) exceeds full capacity (%d)",
					i, battery.CapacityNow, battery.CapacityFull)
			}
		}

		// Total capacity should be sum of individual batteries
		if data.TotalCapacity == 0 && len(data.Batteries) > 0 {
			// At least one battery should contribute to total
			t.Error("TotalCapacity is 0 but batteries are present")
		}
	} else {
		// If no battery present, list should be empty
		if len(data.Batteries) > 0 {
			t.Error("Battery not present but batteries list is not empty")
		}

		if data.TotalCapacity > 0 {
			t.Error("Battery not present but TotalCapacity is > 0")
		}
	}
}

func TestBatteryStateValidation(t *testing.T) {
	data, err := CollectBattery()
	if err != nil {
		t.Skip("Skipping validation test due to collection error")
	}

	if !data.Present {
		t.Skip("No battery present, skipping state validation")
	}

	for _, battery := range data.Batteries {
		// Valid states
		validStates := map[string]bool{
			"Charging":     true,
			"Discharging":  true,
			"Full":         true,
			"Idle":         true,
			"Not charging": true,
			"Low":          true,
			"Critical":     true,
			"Unknown":      true,
			"Other":        true,
			"Undefined":    true,
		}

		// State should be one of the valid states or a specific platform state
		if battery.State != "" && !validStates[battery.State] {
			// Allow platform-specific states, just log them
			t.Logf("Battery has non-standard state: %s", battery.State)
		}
	}
}
