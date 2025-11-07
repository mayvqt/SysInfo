package collector

import (
	"testing"
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
	}
}

// BenchmarkCollectGPU benchmarks GPU collection
func BenchmarkCollectGPU(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CollectGPU()
	}
}
