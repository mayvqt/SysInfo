//go:build linux
// +build linux

package collector

import (
	"testing"

	"github.com/mayvqt/sysinfo/internal/types"
)

func TestParseDmidecodeOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		checks   []func(t *testing.T, modules []types.MemoryModule)
	}{
		{
			name: "single_module",
			input: `Memory Device
	Size: 8192 MB
	Locator: DIMM0
	Bank Locator: BANK 0
	Type: DDR4
	Speed: 2400 MT/s
	Manufacturer: Crucial
	Serial Number: 12345678
	Part Number: CT8G4DFS824A
	Form Factor: DIMM
	Configured Memory Speed: 2133 MT/s`,
			expected: 1,
			checks: []func(t *testing.T, modules []types.MemoryModule){
				func(t *testing.T, modules []types.MemoryModule) {
					if modules[0].Capacity != 8589934592 { // 8GB in bytes
						t.Errorf("Expected capacity 8589934592, got %d", modules[0].Capacity)
					}
					if modules[0].Locator != "DIMM0" {
						t.Errorf("Expected locator DIMM0, got %s", modules[0].Locator)
					}
					if modules[0].Type != "DDR4" {
						t.Errorf("Expected type DDR4, got %s", modules[0].Type)
					}
					if modules[0].Speed != 2133 {
						t.Errorf("Expected speed 2133, got %d", modules[0].Speed)
					}
					if modules[0].Manufacturer != "Crucial" {
						t.Errorf("Expected manufacturer Crucial, got %s", modules[0].Manufacturer)
					}
					if modules[0].FormFactor != "DIMM" {
						t.Errorf("Expected form factor DIMM, got %s", modules[0].FormFactor)
					}
				},
			},
		},
		{
			name: "multiple_modules",
			input: `Memory Device
	Size: 16384 MB
	Locator: DIMM0
	Type: DDR4
	Speed: 3200 MT/s
	Manufacturer: Corsair
	Form Factor: DIMM

Memory Device
	Size: 16384 MB
	Locator: DIMM1
	Type: DDR4
	Speed: 3200 MT/s
	Manufacturer: Corsair
	Form Factor: DIMM`,
			expected: 2,
			checks: []func(t *testing.T, modules []types.MemoryModule){
				func(t *testing.T, modules []types.MemoryModule) {
					if len(modules) != 2 {
						t.Fatalf("Expected 2 modules, got %d", len(modules))
					}
					if modules[0].Locator != "DIMM0" {
						t.Errorf("Expected first module locator DIMM0, got %s", modules[0].Locator)
					}
					if modules[1].Locator != "DIMM1" {
						t.Errorf("Expected second module locator DIMM1, got %s", modules[1].Locator)
					}
				},
			},
		},
		{
			name: "empty_slots",
			input: `Memory Device
	Size: No Module Installed
	Locator: DIMM0

Memory Device
	Size: 8192 MB
	Locator: DIMM1
	Type: DDR3
	Speed: 1600 MT/s`,
			expected: 1,
			checks: []func(t *testing.T, modules []types.MemoryModule){
				func(t *testing.T, modules []types.MemoryModule) {
					if modules[0].Locator != "DIMM1" {
						t.Errorf("Expected locator DIMM1, got %s", modules[0].Locator)
					}
				},
			},
		},
		{
			name:     "empty_output",
			input:    "",
			expected: 0,
		},
		{
			name: "unknown_values",
			input: `Memory Device
	Size: 4096 MB
	Locator: DIMM0
	Type: Unknown
	Speed: Unknown
	Manufacturer: Unknown
	Serial Number: Unknown
	Part Number: NO DIMM
	Form Factor: Unknown`,
			expected: 1,
			checks: []func(t *testing.T, modules []types.MemoryModule){
				func(t *testing.T, modules []types.MemoryModule) {
					if modules[0].Type != "" {
						t.Errorf("Expected empty type for Unknown, got %s", modules[0].Type)
					}
					if modules[0].Manufacturer != "" {
						t.Errorf("Expected empty manufacturer for Unknown, got %s", modules[0].Manufacturer)
					}
					if modules[0].FormFactor != "" {
						t.Errorf("Expected empty form factor for Unknown, got %s", modules[0].FormFactor)
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modules := parseDmidecodeOutput(tt.input)
			if len(modules) != tt.expected {
				t.Errorf("Expected %d modules, got %d", tt.expected, len(modules))
			}
			for _, check := range tt.checks {
				check(t, modules)
			}
		})
	}
}

func TestParseMemorySize(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"8192 MB", 8589934592},
		{"8 GB", 8589934592},
		{"16384 MB", 17179869184},
		{"16 GB", 17179869184},
		{"1024 KB", 1048576},
		{"1 TB", 1099511627776},
		{"No Module Installed", 0},
		{"Unknown", 0},
		{"", 0},
		{"invalid", 0},
		{"1024", 0}, // Missing unit
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseMemorySize(tt.input)
			if result != tt.expected {
				t.Errorf("parseMemorySize(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseMemorySpeed(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"2400 MT/s", 2400},
		{"2400 MHz", 2400},
		{"3200 MT/s", 3200},
		{"1600 MHz", 1600},
		{"Unknown", 0},
		{"", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseMemorySpeed(tt.input)
			if result != tt.expected {
				t.Errorf("parseMemorySpeed(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCollectMemoryModulesPlatformLinux(t *testing.T) {
	// This test will only work if dmidecode is available and has permissions
	modules := collectMemoryModulesPlatform()

	// Should not panic and should return a slice (may be empty if no permissions)
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
		t.Logf("Module[%d]: %s, %d bytes, %s, %d MHz",
			i, mod.Locator, mod.Capacity, mod.Type, mod.Speed)
	}

	t.Logf("Found %d memory modules on Linux", len(modules))
}
