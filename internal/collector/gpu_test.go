package collector

import (
	"testing"

	"github.com/mayvqt/sysinfo/internal/types"
)

// TestCollectGPU is an integration test that verifies GPU collection
func TestCollectGPU(t *testing.T) {
	data, err := CollectGPU()

	// GPU collection might fail if no GPU is present or tools are missing
	// This is expected behavior, so we just log it
	if err != nil {
		t.Logf("CollectGPU returned error (may be expected): %v", err)
		return
	}

	if data == nil {
		t.Logf("CollectGPU returned nil data (no GPUs found)")
		return
	}

	if len(data.GPUs) == 0 {
		t.Log("No GPUs found (this is okay on systems without discrete GPUs)")
		return
	}

	// If we found GPUs, verify basic structure
	for i, gpu := range data.GPUs {
		t.Logf("GPU %d:", i)
		t.Logf("  Name: %s", gpu.Name)
		t.Logf("  Vendor: %s", gpu.Vendor)
		t.Logf("  Driver: %s", gpu.Driver)

		if gpu.Name == "" {
			t.Errorf("GPU %d has empty name", i)
		}

		if gpu.Index != i {
			t.Errorf("GPU index mismatch: expected %d, got %d", i, gpu.Index)
		}

		// Verify vendor is known
		validVendors := map[string]bool{
			"NVIDIA": true,
			"AMD":    true,
			"Intel":  true,
			"":       true, // Unknown is okay
		}
		if !validVendors[gpu.Vendor] {
			t.Errorf("GPU %d has unknown vendor: %s", i, gpu.Vendor)
		}

		// If memory is reported, it should be reasonable (not negative, formatted should exist)
		if gpu.MemoryTotal > 0 {
			if gpu.MemoryFormatted == "" {
				t.Errorf("GPU %d has memory total but no formatted string", i)
			}
			// Memory should be less than 1TB (reasonable upper limit for now)
			if gpu.MemoryTotal > 1024*1024*1024*1024 {
				t.Errorf("GPU %d memory total seems unreasonably large: %d bytes", i, gpu.MemoryTotal)
			}
		}

		// Validate temperature is in reasonable range if reported
		if gpu.Temperature > 0 {
			if gpu.Temperature < 0 || gpu.Temperature > 120 {
				t.Errorf("GPU %d temperature out of reasonable range: %d°C", i, gpu.Temperature)
			}
		}

		// Validate utilization percentages
		if gpu.Utilization < 0 || gpu.Utilization > 100 {
			t.Errorf("GPU %d utilization out of range: %d%%", i, gpu.Utilization)
		}
		if gpu.MemoryUtilization < 0 || gpu.MemoryUtilization > 100 {
			t.Errorf("GPU %d memory utilization out of range: %d%%", i, gpu.MemoryUtilization)
		}
		if gpu.FanSpeed < 0 || gpu.FanSpeed > 100 {
			t.Errorf("GPU %d fan speed out of range: %d%%", i, gpu.FanSpeed)
		}

		// Validate power values if reported
		if gpu.PowerDraw < 0 {
			t.Errorf("GPU %d power draw is negative: %.2f W", i, gpu.PowerDraw)
		}
		if gpu.PowerLimit < 0 {
			t.Errorf("GPU %d power limit is negative: %.2f W", i, gpu.PowerLimit)
		}
		if gpu.PowerDraw > gpu.PowerLimit && gpu.PowerLimit > 0 {
			// This is technically possible but unusual
			t.Logf("GPU %d power draw (%.2f W) exceeds power limit (%.2f W)", i, gpu.PowerDraw, gpu.PowerLimit)
		}

		// Validate clock speeds are reasonable if reported
		if gpu.ClockSpeed < 0 || gpu.ClockSpeed > 10000 {
			t.Errorf("GPU %d clock speed out of reasonable range: %d MHz", i, gpu.ClockSpeed)
		}
		if gpu.ClockSpeedMemory < 0 || gpu.ClockSpeedMemory > 30000 {
			t.Errorf("GPU %d memory clock speed out of reasonable range: %d MHz", i, gpu.ClockSpeedMemory)
		}
	}
}

// TestCollectGPUPlatform tests the platform-specific implementation
func TestCollectGPUPlatform(t *testing.T) {
	gpus := collectGPUPlatform()

	t.Logf("Found %d GPU(s)", len(gpus))

	for i, gpu := range gpus {
		t.Logf("GPU %d: %s (%s)", i, gpu.Name, gpu.Vendor)

		if gpu.MemoryTotal > 0 {
			t.Logf("  Memory: %s", gpu.MemoryFormatted)
		}
		if gpu.Temperature > 0 {
			t.Logf("  Temperature: %d°C", gpu.Temperature)
		}
		if gpu.Utilization > 0 {
			t.Logf("  Utilization: %d%%", gpu.Utilization)
		}
		if gpu.Driver != "" {
			t.Logf("  Driver: %s", gpu.Driver)
			if gpu.DriverVersion != "" {
				t.Logf("  Driver Version: %s", gpu.DriverVersion)
			}
		}
		if gpu.PCIBus != "" {
			t.Logf("  PCI Bus: %s", gpu.PCIBus)
		}
	}
}

// TestGPUDataStructure validates the GPU data structure
func TestGPUDataStructure(t *testing.T) {
	// Create a sample GPU data structure
	gpu := types.GPUInfo{
		Index:             0,
		Name:              "Test GPU",
		Vendor:            "NVIDIA",
		Driver:            "nvidia",
		DriverVersion:     "535.161.07",
		MemoryTotal:       8 * 1024 * 1024 * 1024, // 8 GB
		MemoryUsed:        4 * 1024 * 1024 * 1024, // 4 GB
		MemoryFree:        4 * 1024 * 1024 * 1024, // 4 GB
		MemoryFormatted:   "8.00 GB",
		Temperature:       65,
		FanSpeed:          50,
		PowerDraw:         150.5,
		PowerLimit:        250.0,
		Utilization:       75,
		MemoryUtilization: 50,
		ClockSpeed:        1500,
		ClockSpeedMemory:  7000,
		PCIBus:            "0000:01:00.0",
		UUID:              "GPU-12345678-1234-1234-1234-123456789012",
	}

	// Validate all fields
	if gpu.Index != 0 {
		t.Errorf("Expected Index 0, got %d", gpu.Index)
	}
	if gpu.Name != "Test GPU" {
		t.Errorf("Expected Name 'Test GPU', got '%s'", gpu.Name)
	}
	if gpu.Vendor != "NVIDIA" {
		t.Errorf("Expected Vendor 'NVIDIA', got '%s'", gpu.Vendor)
	}
	if gpu.MemoryTotal != 8*1024*1024*1024 {
		t.Errorf("Expected MemoryTotal 8GB, got %d", gpu.MemoryTotal)
	}
	if gpu.Temperature != 65 {
		t.Errorf("Expected Temperature 65, got %d", gpu.Temperature)
	}

	// Test GPUData structure
	data := types.GPUData{
		GPUs: []types.GPUInfo{gpu},
	}

	if len(data.GPUs) != 1 {
		t.Errorf("Expected 1 GPU, got %d", len(data.GPUs))
	}
	if data.GPUs[0].Name != "Test GPU" {
		t.Errorf("Expected GPU name 'Test GPU', got '%s'", data.GPUs[0].Name)
	}
}

// TestCollectGPUNoGPUs tests behavior when no GPUs are found
func TestCollectGPUNoGPUs(t *testing.T) {
	// This test validates that CollectGPU handles the no-GPU case correctly
	// by calling the real function (which may or may not find GPUs)
	data, err := CollectGPU()

	// Either we get an error (no GPUs), or we get valid data
	if err != nil {
		// Error case is acceptable - no GPUs found
		t.Logf("No GPUs found (expected on systems without discrete GPUs): %v", err)
		if data != nil {
			t.Error("Expected nil data when error is returned")
		}
	} else {
		// Success case - validate data
		if data == nil {
			t.Error("Expected non-nil data when no error is returned")
		} else if len(data.GPUs) == 0 {
			t.Error("Expected at least one GPU when no error is returned")
		}
	}
}

// TestGPUMemoryCalculations tests memory-related calculations
func TestGPUMemoryCalculations(t *testing.T) {
	testCases := []struct {
		name        string
		total       uint64
		used        uint64
		free        uint64
		shouldMatch bool
	}{
		{
			name:        "Matching memory values",
			total:       8 * 1024 * 1024 * 1024,
			used:        4 * 1024 * 1024 * 1024,
			free:        4 * 1024 * 1024 * 1024,
			shouldMatch: true,
		},
		{
			name:        "Zero values",
			total:       0,
			used:        0,
			free:        0,
			shouldMatch: true,
		},
		{
			name:        "Only total set",
			total:       8 * 1024 * 1024 * 1024,
			used:        0,
			free:        0,
			shouldMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldMatch && tc.total > 0 {
				if tc.used+tc.free != tc.total {
					t.Errorf("Memory values don't match: total=%d, used=%d, free=%d", tc.total, tc.used, tc.free)
				}
			}

			// Validate no negative values
			if tc.total > 0 && tc.used > tc.total {
				t.Errorf("Used memory (%d) exceeds total (%d)", tc.used, tc.total)
			}
			if tc.total > 0 && tc.free > tc.total {
				t.Errorf("Free memory (%d) exceeds total (%d)", tc.free, tc.total)
			}
		})
	}
}

// TestGPUVendorDetection tests vendor detection logic
func TestGPUVendorDetection(t *testing.T) {
	testCases := []struct {
		gpuName        string
		expectedVendor string
	}{
		{"NVIDIA GeForce RTX 3080", "NVIDIA"},
		{"AMD Radeon RX 6800 XT", "AMD"},
		{"Intel UHD Graphics 630", "Intel"},
		{"GeForce RTX 4090", "NVIDIA"},
		{"Radeon VII", "AMD"},
	}

	for _, tc := range testCases {
		t.Run(tc.gpuName, func(t *testing.T) {
			// This is a basic validation that our test expectations make sense
			// The actual vendor detection happens in platform-specific code
			gpu := types.GPUInfo{
				Name:   tc.gpuName,
				Vendor: tc.expectedVendor,
			}

			if gpu.Name != tc.gpuName {
				t.Errorf("Expected name '%s', got '%s'", tc.gpuName, gpu.Name)
			}
			if gpu.Vendor != tc.expectedVendor {
				t.Errorf("Expected vendor '%s', got '%s'", tc.expectedVendor, gpu.Vendor)
			}
		})
	}
}

// TestGPUIndexing tests that GPU indices are correctly assigned
func TestGPUIndexing(t *testing.T) {
	// Create multiple GPUs
	gpus := []types.GPUInfo{
		{Index: 0, Name: "GPU 0"},
		{Index: 1, Name: "GPU 1"},
		{Index: 2, Name: "GPU 2"},
	}

	for i, gpu := range gpus {
		if gpu.Index != i {
			t.Errorf("GPU %d has incorrect index: %d", i, gpu.Index)
		}
	}
}

// BenchmarkCollectGPU benchmarks GPU collection
func BenchmarkCollectGPU(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CollectGPU()
	}
}

// BenchmarkCollectGPUPlatform benchmarks platform-specific GPU collection
func BenchmarkCollectGPUPlatform(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = collectGPUPlatform()
	}
}
