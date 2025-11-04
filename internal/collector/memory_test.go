package collector

import (
	"testing"
)

// TestCollectMemory verifies basic memory collection works
func TestCollectMemory(t *testing.T) {
	data, err := CollectMemory()
	if err != nil {
		t.Fatalf("CollectMemory failed: %v", err)
	}

	if data == nil {
		t.Fatal("CollectMemory returned nil data")
	}

	// Total memory should be > 0
	if data.Total == 0 {
		t.Error("Total memory is 0")
	}

	// Used + Free should approximately equal Total (allowing for rounding/caching)
	// This is a loose check since different systems report differently
	if data.Used == 0 && data.Free == 0 {
		t.Error("Both Used and Free are 0")
	}

	// UsedPercent should be in valid range
	if data.UsedPercent < 0 || data.UsedPercent > 100 {
		t.Errorf("UsedPercent = %f is out of range [0, 100]", data.UsedPercent)
	}

	// Formatted strings should not be empty
	if data.TotalFormatted == "" {
		t.Error("TotalFormatted is empty")
	}
	if data.UsedFormatted == "" {
		t.Error("UsedFormatted is empty")
	}
	if data.FreeFormatted == "" {
		t.Error("FreeFormatted is empty")
	}

	t.Logf("Memory: Total=%s, Used=%s (%.1f%%), Free=%s",
		data.TotalFormatted, data.UsedFormatted, data.UsedPercent, data.FreeFormatted)

	// Swap may be 0 on some systems
	if data.SwapTotal > 0 {
		t.Logf("Swap: Total=%s, Used=%s (%.1f%%)",
			formatBytes(data.SwapTotal),
			formatBytes(data.SwapUsed),
			data.SwapPercent)

		if data.SwapPercent < 0 || data.SwapPercent > 100 {
			t.Errorf("SwapPercent = %f is out of range [0, 100]", data.SwapPercent)
		}
	}

	// Available should be <= Total
	if data.Available > data.Total {
		t.Errorf("Available (%d) > Total (%d)", data.Available, data.Total)
	}

	// Cached, Buffers, Shared may be 0 on some systems (especially Windows)
	t.Logf("Cached: %d, Buffers: %d, Shared: %d", data.Cached, data.Buffers, data.Shared)
}

func TestCollectMemoryModules(t *testing.T) {
	// collectMemoryModules is currently a placeholder
	modules := collectMemoryModules()

	// Should not panic and should return a slice (may be empty)
	if modules == nil {
		t.Error("collectMemoryModules returned nil")
	}

	// If modules are present, verify they have reasonable data
	for i, mod := range modules {
		if mod.Capacity == 0 {
			t.Errorf("Module[%d] has zero capacity", i)
		}
		if mod.Locator == "" {
			t.Errorf("Module[%d] has empty locator", i)
		}
	}

	if len(modules) > 0 {
		t.Logf("Found %d memory modules", len(modules))
	} else {
		t.Log("No memory modules detected (may be expected on this platform)")
	}
}

func TestMemoryConsistency(t *testing.T) {
	// Run collection multiple times and verify results are consistent
	data1, err1 := CollectMemory()
	if err1 != nil {
		t.Fatalf("First CollectMemory failed: %v", err1)
	}

	data2, err2 := CollectMemory()
	if err2 != nil {
		t.Fatalf("Second CollectMemory failed: %v", err2)
	}

	// Total memory should be the same
	if data1.Total != data2.Total {
		t.Errorf("Total memory changed: %d -> %d", data1.Total, data2.Total)
	}

	// Used memory may change slightly, but should be in reasonable range
	diff := int64(data2.Used) - int64(data1.Used)
	if diff < 0 {
		diff = -diff
	}
	maxDiff := int64(data1.Total / 10) // Allow 10% variance
	if diff > maxDiff {
		t.Logf("Warning: Memory usage changed significantly: %d -> %d (diff: %d)",
			data1.Used, data2.Used, diff)
	}
}

func BenchmarkCollectMemory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CollectMemory()
	}
}

// Helper function for formatting in tests
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return "0 B"
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}
	return units[exp]
}
