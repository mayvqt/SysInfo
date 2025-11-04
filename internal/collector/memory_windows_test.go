//go:build windows
// +build windows

package collector

import (
	"testing"
)

func TestGetMemoryType(t *testing.T) {
	tests := []struct {
		code     uint16
		expected string
	}{
		{0, "Unknown"},
		{20, "DDR"},
		{21, "DDR2"},
		{24, "DDR3"},
		{26, "DDR4"},
		{34, "DDR5"},
		{27, "LPDDR"},
		{28, "LPDDR2"},
		{29, "LPDDR3"},
		{30, "LPDDR4"},
		{35, "LPDDR5"},
		{32, "HBM (High Bandwidth Memory)"},
		{33, "HBM2 (High Bandwidth Memory Generation 2)"},
		{999, "Unknown"}, // Unknown code
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getMemoryType(tt.code)
			if result != tt.expected {
				t.Errorf("getMemoryType(%d) = %q, expected %q", tt.code, result, tt.expected)
			}
		})
	}
}

func TestGetFormFactor(t *testing.T) {
	tests := []struct {
		code     uint16
		expected string
	}{
		{0, "Unknown"},
		{8, "DIMM"},
		{12, "SODIMM"},
		{24, "FB-DIMM"},
		{11, "RIMM"},
		{7, "SIMM"},
		{13, "SRIMM"},
		{999, "Unknown"}, // Unknown code
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getFormFactor(tt.code)
			if result != tt.expected {
				t.Errorf("getFormFactor(%d) = %q, expected %q", tt.code, result, tt.expected)
			}
		})
	}
}

func TestGetLocator(t *testing.T) {
	tests := []struct {
		name     string
		mem      Win32_PhysicalMemory
		expected string
	}{
		{
			name: "device_locator",
			mem: Win32_PhysicalMemory{
				DeviceLocator: "DIMM0",
				BankLabel:     "BANK 0",
				Tag:           "Physical Memory 0",
			},
			expected: "DIMM0",
		},
		{
			name: "bank_label",
			mem: Win32_PhysicalMemory{
				DeviceLocator: "",
				BankLabel:     "BANK 1",
				Tag:           "Physical Memory 1",
			},
			expected: "BANK 1",
		},
		{
			name: "tag",
			mem: Win32_PhysicalMemory{
				DeviceLocator: "",
				BankLabel:     "",
				Tag:           "Physical Memory 2",
			},
			expected: "Physical Memory 2",
		},
		{
			name: "unknown",
			mem: Win32_PhysicalMemory{
				DeviceLocator: "",
				BankLabel:     "",
				Tag:           "",
			},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLocator(tt.mem)
			if result != tt.expected {
				t.Errorf("getLocator() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestCollectMemoryModulesPlatformWindows(t *testing.T) {
	// This test will call the actual WMI implementation
	modules := collectMemoryModulesPlatform()

	// Should not panic and should return a slice (may be empty on some systems)
	if modules == nil {
		t.Error("collectMemoryModulesPlatform returned nil")
	}

	// If modules are found, validate them
	for i, mod := range modules {
		if mod.Capacity == 0 {
			t.Errorf("Module[%d] has zero capacity", i)
		}
		if mod.Locator == "" {
			t.Errorf("Module[%d] has empty locator", i)
		}

		// Log the module details
		t.Logf("Module[%d]: %s", i, mod.Locator)
		t.Logf("  Capacity: %d bytes (%.2f GB)", mod.Capacity, float64(mod.Capacity)/(1024*1024*1024))
		t.Logf("  Type: %s", mod.Type)
		t.Logf("  Speed: %d MHz", mod.Speed)
		t.Logf("  Form Factor: %s", mod.FormFactor)

		if mod.Manufacturer != "" {
			t.Logf("  Manufacturer: %s", mod.Manufacturer)
		}
		if mod.PartNumber != "" {
			t.Logf("  Part Number: %s", mod.PartNumber)
		}
		if mod.SerialNumber != "" {
			t.Logf("  Serial Number: %s", mod.SerialNumber)
		}
	}

	if len(modules) > 0 {
		t.Logf("Successfully collected %d memory modules on Windows", len(modules))
	} else {
		t.Log("No memory modules detected (may be expected in virtualized environments)")
	}
}

func TestMemoryModuleFieldValidation(t *testing.T) {
	// Integration test to ensure all fields are properly populated
	modules := collectMemoryModulesPlatform()

	if len(modules) == 0 {
		t.Skip("No memory modules to validate")
	}

	for i, mod := range modules {
		t.Run("module_"+mod.Locator, func(t *testing.T) {
			// Capacity should always be > 0
			if mod.Capacity == 0 {
				t.Error("Capacity is 0")
			}

			// Locator should never be empty
			if mod.Locator == "" {
				t.Error("Locator is empty")
			}

			// Type might be "Unknown" but shouldn't be empty if populated
			if mod.Type != "" && mod.Type == "Unknown" {
				t.Logf("Module[%d] has Unknown type", i)
			}

			// Form Factor should be one of the known types if populated
			validFormFactors := map[string]bool{
				"DIMM": true, "SODIMM": true, "RIMM": true, "SIMM": true,
				"FB-DIMM": true, "Unknown": true, "Other": true, "SRIMM": true,
			}
			if mod.FormFactor != "" && !validFormFactors[mod.FormFactor] {
				t.Logf("Module[%d] has unusual form factor: %s", i, mod.FormFactor)
			}
		})
	}
}
