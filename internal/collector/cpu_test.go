package collector

import (
	"testing"
)

// TestCollectCPU verifies basic CPU collection works
func TestCollectCPU(t *testing.T) {
	data, err := CollectCPU()
	if err != nil {
		t.Fatalf("CollectCPU failed: %v", err)
	}

	if data == nil {
		t.Fatal("CollectCPU returned nil data")
	}

	// Basic sanity checks
	if data.ModelName == "" {
		t.Error("ModelName is empty")
	}

	if data.Cores <= 0 {
		t.Error("Cores should be > 0")
	}

	if data.LogicalCPUs <= 0 {
		t.Error("LogicalCPUs should be > 0")
	}

	if data.LogicalCPUs < data.Cores {
		t.Errorf("LogicalCPUs (%d) should be >= Cores (%d)", data.LogicalCPUs, data.Cores)
	}

	// Vendor should be present
	if data.Vendor == "" {
		t.Log("Warning: Vendor is empty (may be acceptable on some systems)")
	}

	// Usage may be empty if collection failed, but shouldn't panic
	t.Logf("CPU: %s, Cores: %d, Logical: %d", data.ModelName, data.Cores, data.LogicalCPUs)

	if len(data.Usage) > 0 {
		t.Logf("CPU Usage samples: %d", len(data.Usage))
		for i, usage := range data.Usage {
			if usage < 0 || usage > 100 {
				t.Errorf("CPU usage[%d] = %f is out of range [0, 100]", i, usage)
			}
		}
	}

	if data.LoadAvg != nil {
		t.Logf("Load Average: %.2f, %.2f, %.2f", data.LoadAvg.Load1, data.LoadAvg.Load5, data.LoadAvg.Load15)
		// Load averages should not be negative
		if data.LoadAvg.Load1 < 0 {
			t.Error("Load1 is negative")
		}
	}
}

func TestCollectCPUReturnsValidData(t *testing.T) {
	data, err := CollectCPU()
	if err != nil {
		t.Fatalf("CollectCPU failed: %v", err)
	}

	// Test that MHz is reasonable (should be > 0 on most systems)
	if data.MHz <= 0 {
		t.Log("Warning: MHz is 0 or negative (may be acceptable on some virtual systems)")
	}

	// Cache size should be non-negative
	if data.CacheSize < 0 {
		t.Error("CacheSize is negative")
	}

	// Flags might be empty on some systems but should not panic
	t.Logf("CPU flags count: %d", len(data.Flags))

	// Microcode may be empty
	if data.Microcode != "" {
		t.Logf("Microcode: %s", data.Microcode)
	}
}

func BenchmarkCollectCPU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CollectCPU()
	}
}
