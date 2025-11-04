//go:build darwin
// +build darwin

package collector

import (
	"testing"
)

func TestParseMemorySizeMac(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"8 GB", 8589934592},
		{"16 GB", 17179869184},
		{"4 GB", 4294967296},
		{"32 GB", 34359738368},
		{"1024 MB", 1073741824},
		{"512 KB", 524288},
		{"1 TB", 1099511627776},
		{"empty", 0},
		{"", 0},
		{"invalid", 0},
		{"1024", 0}, // Missing unit
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseMemorySizeMac(tt.input)
			if result != tt.expected {
				t.Errorf("parseMemorySizeMac(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseMemorySpeedMac(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"2400 MHz", 2400},
		{"2133 MHz", 2133},
		{"3200 MHz", 3200},
		{"1600 MHz", 1600},
		{"", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseMemorySpeedMac(tt.input)
			if result != tt.expected {
				t.Errorf("parseMemorySpeedMac(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseMemorySlot(t *testing.T) {
	tests := []struct {
		name     string
		input    DictData
		expected MemorySlot
	}{
		{
			name: "complete_slot",
			input: DictData{
				Keys: []string{
					"dimm_size",
					"dimm_type",
					"dimm_speed",
					"dimm_status",
					"dimm_manufacturer",
					"dimm_part_number",
					"dimm_serial_number",
				},
				Values: []string{
					"8 GB",
					"DDR4",
					"2400 MHz",
					"OK",
					"Crucial",
					"CT8G4DFS824A",
					"12345678",
				},
			},
			expected: MemorySlot{
				Size:         "8 GB",
				Type:         "DDR4",
				Speed:        "2400 MHz",
				Status:       "OK",
				Manufacturer: "Crucial",
				PartNumber:   "CT8G4DFS824A",
				SerialNumber: "12345678",
			},
		},
		{
			name: "empty_slot",
			input: DictData{
				Keys: []string{
					"dimm_size",
					"dimm_status",
				},
				Values: []string{
					"empty",
					"Empty",
				},
			},
			expected: MemorySlot{
				Size:   "empty",
				Status: "Empty",
			},
		},
		{
			name: "partial_data",
			input: DictData{
				Keys: []string{
					"dimm_size",
					"dimm_type",
				},
				Values: []string{
					"16 GB",
					"DDR3",
				},
			},
			expected: MemorySlot{
				Size: "16 GB",
				Type: "DDR3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMemorySlot(tt.input)
			if result.Size != tt.expected.Size {
				t.Errorf("Size = %q, expected %q", result.Size, tt.expected.Size)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %q, expected %q", result.Type, tt.expected.Type)
			}
			if result.Speed != tt.expected.Speed {
				t.Errorf("Speed = %q, expected %q", result.Speed, tt.expected.Speed)
			}
			if result.Status != tt.expected.Status {
				t.Errorf("Status = %q, expected %q", result.Status, tt.expected.Status)
			}
			if result.Manufacturer != tt.expected.Manufacturer {
				t.Errorf("Manufacturer = %q, expected %q", result.Manufacturer, tt.expected.Manufacturer)
			}
			if result.PartNumber != tt.expected.PartNumber {
				t.Errorf("PartNumber = %q, expected %q", result.PartNumber, tt.expected.PartNumber)
			}
			if result.SerialNumber != tt.expected.SerialNumber {
				t.Errorf("SerialNumber = %q, expected %q", result.SerialNumber, tt.expected.SerialNumber)
			}
		})
	}
}

func TestCollectMemoryModulesPlatformDarwin(t *testing.T) {
	// This test will only work if system_profiler is available
	modules := collectMemoryModulesPlatform()

	// Should not panic and should return a slice (may be empty)
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
		// Form factor should be either DIMM or SODIMM on macOS
		if mod.FormFactor != "DIMM" && mod.FormFactor != "SODIMM" {
			t.Logf("Module[%d] has unusual form factor: %s", i, mod.FormFactor)
		}
		t.Logf("Module[%d]: %s, %d bytes, %s, %d MHz, %s",
			i, mod.Locator, mod.Capacity, mod.Type, mod.Speed, mod.FormFactor)
	}

	t.Logf("Found %d memory modules on macOS", len(modules))
}

func TestParseSystemProfilerMemoryEmpty(t *testing.T) {
	// Test with empty XML
	slots := parseSystemProfilerMemory([]byte(""))
	if len(slots) != 0 {
		t.Errorf("Expected 0 slots from empty XML, got %d", len(slots))
	}

	// Test with minimal valid XML
	minimalXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<array></array>
</plist>`)

	slots = parseSystemProfilerMemory(minimalXML)
	if len(slots) != 0 {
		t.Errorf("Expected 0 slots from minimal XML, got %d", len(slots))
	}
}
