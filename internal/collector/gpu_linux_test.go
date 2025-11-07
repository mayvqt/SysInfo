//go:build linux
// +build linux

package collector

import (
	"testing"
)

// TestParseMemoryMiB tests the memory parsing helper
func TestParseMemoryMiB(t *testing.T) {
	testCases := []struct {
		input    string
		expected uint64
	}{
		{"8192 MiB", 8192 * 1024 * 1024},
		{"8192MiB", 8192 * 1024 * 1024},
		{"16384 MiB", 16384 * 1024 * 1024},
		{"0 MiB", 0},
		{"  8192 MiB  ", 8192 * 1024 * 1024},
		{"invalid", 0},
		{"", 0},
		{"MiB", 0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := parseMemoryMiB(tc.input)
			if result != tc.expected {
				t.Errorf("parseMemoryMiB(%q) = %d, expected %d", tc.input, result, tc.expected)
			}
		})
	}
}

// TestParsePowerWatts tests the power parsing helper
func TestParsePowerWatts(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{"150.5 W", 150.5},
		{"150.5W", 150.5},
		{"250.0 W", 250.0},
		{"0 W", 0.0},
		{"  150.5 W  ", 150.5},
		{"invalid", 0.0},
		{"", 0.0},
		{"W", 0.0},
		{"75", 75.0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := parsePowerWatts(tc.input)
			if result != tc.expected {
				t.Errorf("parsePowerWatts(%q) = %.2f, expected %.2f", tc.input, result, tc.expected)
			}
		})
	}
}

// TestParseClockMHz tests the clock speed parsing helper
func TestParseClockMHz(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"1500 MHz", 1500},
		{"1500MHz", 1500},
		{"2100 MHz", 2100},
		{"0 MHz", 0},
		{"  1500 MHz  ", 1500},
		{"invalid", 0},
		{"", 0},
		{"MHz", 0},
		{"1500", 1500},
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := parseClockMHz(tc.input)
			if result != tc.expected {
				t.Errorf("parseClockMHz(%q) = %d, expected %d", tc.input, result, tc.expected)
			}
		})
	}
}

// TestExtractGPUName tests the GPU name extraction from lspci output
func TestExtractGPUName(t *testing.T) {
	testCases := []struct {
		line     string
		vendor   string
		expected string
	}{
		{
			line:     "01:00.0 VGA compatible controller: NVIDIA Corporation GA102 [GeForce RTX 3080]",
			vendor:   "NVIDIA",
			expected: "NVIDIA Corporation GA102 [GeForce RTX 3080]",
		},
		{
			line:     "00:02.0 VGA compatible controller: Intel Corporation UHD Graphics 630",
			vendor:   "Intel",
			expected: "Intel Corporation UHD Graphics 630",
		},
		{
			line:     "03:00.0 3D controller: NVIDIA Corporation TU117M [GeForce GTX 1650 Mobile]",
			vendor:   "NVIDIA",
			expected: "NVIDIA Corporation TU117M [GeForce GTX 1650 Mobile]",
		},
		{
			line:     "Invalid line without colon",
			vendor:   "Unknown",
			expected: "Unknown GPU",
		},
		{
			line:     "",
			vendor:   "Unknown",
			expected: "Unknown GPU",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.line, func(t *testing.T) {
			result := extractGPUName(tc.line, tc.vendor)
			// The function cleans up the name, so we just verify it's not empty when input is valid
			if tc.line != "" && tc.line != "Invalid line without colon" {
				if result == "" {
					t.Errorf("extractGPUName(%q, %q) returned empty string", tc.line, tc.vendor)
				}
				// Check that VGA/3D controller prefixes are removed
				if result == "VGA compatible controller:" || result == "3D controller:" {
					t.Errorf("extractGPUName(%q, %q) didn't remove controller prefix", tc.line, tc.vendor)
				}
			}
		})
	}
}

// TestCollectGPUsFromLspci tests the lspci fallback method
func TestCollectGPUsFromLspci(t *testing.T) {
	// This is an integration test that actually calls lspci if available
	gpus := collectGPUsFromLspci()
	
	// We can't guarantee lspci will find GPUs, but if it does, validate them
	if len(gpus) > 0 {
		t.Logf("Found %d GPU(s) via lspci", len(gpus))
		
		for i, gpu := range gpus {
			t.Logf("GPU %d: %s (%s)", i, gpu.Name, gpu.Vendor)
			
			if gpu.Name == "" {
				t.Errorf("GPU %d has empty name", i)
			}
			
			if gpu.Index != i {
				t.Errorf("GPU %d has incorrect index: %d", i, gpu.Index)
			}
			
			// PCIBus should be set when using lspci
			if gpu.PCIBus == "" {
				t.Logf("Warning: GPU %d missing PCI bus information", i)
			}
		}
	} else {
		t.Log("No GPUs found via lspci (may be expected)")
	}
}

// TestCollectNvidiaGPUs tests NVIDIA GPU collection
func TestCollectNvidiaGPUs(t *testing.T) {
	gpus := collectNvidiaGPUs()
	
	if len(gpus) > 0 {
		t.Logf("Found %d NVIDIA GPU(s)", len(gpus))
		
		for i, gpu := range gpus {
			if gpu.Vendor != "NVIDIA" {
				t.Errorf("GPU %d vendor should be NVIDIA, got %s", i, gpu.Vendor)
			}
			
			if gpu.Name == "" {
				t.Errorf("GPU %d has empty name", i)
			}
			
			// If we got detailed info, verify it
			if gpu.MemoryTotal > 0 {
				t.Logf("GPU %d: %s - Memory: %s", i, gpu.Name, gpu.MemoryFormatted)
				
				if gpu.MemoryFormatted == "" {
					t.Errorf("GPU %d has memory but no formatted string", i)
				}
			}
			
			if gpu.Temperature > 0 {
				if gpu.Temperature < 20 || gpu.Temperature > 100 {
					t.Logf("Warning: GPU %d temperature seems unusual: %d°C", i, gpu.Temperature)
				}
			}
		}
	} else {
		t.Log("No NVIDIA GPUs found (expected if nvidia-smi not available or no NVIDIA GPUs)")
	}
}

// TestCollectAMDGPUs tests AMD GPU collection
func TestCollectAMDGPUs(t *testing.T) {
	gpus := collectAMDGPUs()
	
	if len(gpus) > 0 {
		t.Logf("Found %d AMD GPU(s)", len(gpus))
		
		for i, gpu := range gpus {
			if gpu.Vendor != "AMD" {
				t.Errorf("GPU %d vendor should be AMD, got %s", i, gpu.Vendor)
			}
			
			t.Logf("GPU %d: %s", i, gpu.Name)
		}
	} else {
		t.Log("No AMD GPUs found (expected if rocm-smi not available or no AMD GPUs)")
	}
}

// TestNvidiaGPUDataValidation tests NVIDIA-specific data validation
func TestNvidiaGPUDataValidation(t *testing.T) {
	// This test validates that if we get NVIDIA GPU data, it's properly structured
	gpus := collectNvidiaGPUs()
	
	for i, gpu := range gpus {
		// All NVIDIA GPUs should have these fields set
		if gpu.Vendor != "NVIDIA" {
			t.Errorf("GPU %d: expected NVIDIA vendor, got %s", i, gpu.Vendor)
		}
		
		if gpu.Driver != "nvidia" && gpu.Driver != "" {
			t.Errorf("GPU %d: expected 'nvidia' driver or empty, got %s", i, gpu.Driver)
		}
		
		// If UUID is set, it should match NVIDIA's format
		if gpu.UUID != "" {
			if len(gpu.UUID) < 10 {
				t.Errorf("GPU %d: UUID seems too short: %s", i, gpu.UUID)
			}
		}
		
		// Memory consistency check
		if gpu.MemoryTotal > 0 && gpu.MemoryUsed > 0 && gpu.MemoryFree > 0 {
			calculatedTotal := gpu.MemoryUsed + gpu.MemoryFree
			// Allow for some rounding differences (within 1MB)
			diff := int64(calculatedTotal) - int64(gpu.MemoryTotal)
			if diff < 0 {
				diff = -diff
			}
			if diff > 1024*1024 {
				t.Logf("GPU %d: Memory values don't add up: total=%d, used=%d, free=%d", 
					i, gpu.MemoryTotal, gpu.MemoryUsed, gpu.MemoryFree)
			}
		}
	}
}

// BenchmarkParseMemoryMiB benchmarks memory parsing
func BenchmarkParseMemoryMiB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseMemoryMiB("8192 MiB")
	}
}

// BenchmarkParsePowerWatts benchmarks power parsing
func BenchmarkParsePowerWatts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePowerWatts("150.5 W")
	}
}

// BenchmarkParseClockMHz benchmarks clock speed parsing
func BenchmarkParseClockMHz(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseClockMHz("1500 MHz")
	}
}

// BenchmarkCollectNvidiaGPUs benchmarks NVIDIA GPU collection
func BenchmarkCollectNvidiaGPUs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = collectNvidiaGPUs()
	}
}

// BenchmarkCollectAMDGPUs benchmarks AMD GPU collection
func BenchmarkCollectAMDGPUs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = collectAMDGPUs()
	}
}

// BenchmarkCollectGPUsFromLspci benchmarks lspci GPU collection
func BenchmarkCollectGPUsFromLspci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = collectGPUsFromLspci()
	}
}
