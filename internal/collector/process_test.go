package collector

import (
	"testing"
)

// TestCollectProcesses verifies basic process collection works
func TestCollectProcesses(t *testing.T) {
	data, err := CollectProcesses()
	if err != nil {
		t.Fatalf("CollectProcesses failed: %v", err)
	}

	if data == nil {
		t.Fatal("CollectProcesses returned nil data")
	}

	// There should always be processes running
	if data.TotalCount == 0 {
		t.Error("TotalCount is 0 (expected at least 1)")
	}

	// Running + Sleeping + Others should be <= TotalCount
	// (there may be other states like zombie, stopped, etc.)
	if data.Running+data.Sleeping > data.TotalCount {
		t.Errorf("Running (%d) + Sleeping (%d) > TotalCount (%d)",
			data.Running, data.Sleeping, data.TotalCount)
	}

	t.Logf("Processes: Total=%d, Running=%d, Sleeping=%d",
		data.TotalCount, data.Running, data.Sleeping)

	// Check top processes by memory
	if len(data.TopByMemory) > 0 {
		t.Logf("Top %d processes by memory:", len(data.TopByMemory))
		for i, proc := range data.TopByMemory {
			if proc.PID <= 0 {
				t.Errorf("TopByMemory[%d] has invalid PID: %d", i, proc.PID)
			}
			if proc.Name == "" {
				t.Errorf("TopByMemory[%d] has empty name", i)
			}
			if proc.MemoryPercent < 0 || proc.MemoryPercent > 100 {
				t.Errorf("TopByMemory[%d] MemoryPercent = %f is out of range [0, 100]",
					i, proc.MemoryPercent)
			}

			t.Logf("  [%d] %s (PID %d): %d MB (%.1f%%)",
				i+1, proc.Name, proc.PID, proc.MemoryMB, proc.MemoryPercent)

			// List should be sorted by memory (descending)
			if i > 0 && proc.MemoryMB > data.TopByMemory[i-1].MemoryMB {
				t.Errorf("TopByMemory not sorted: [%d] %d MB > [%d] %d MB",
					i, proc.MemoryMB, i-1, data.TopByMemory[i-1].MemoryMB)
			}
		}
	} else {
		t.Log("No top memory processes available")
	}

	// Check top processes by CPU
	if len(data.TopByCPU) > 0 {
		t.Logf("Top %d processes by CPU:", len(data.TopByCPU))
		for i, proc := range data.TopByCPU {
			if proc.PID <= 0 {
				t.Errorf("TopByCPU[%d] has invalid PID: %d", i, proc.PID)
			}
			if proc.Name == "" {
				t.Errorf("TopByCPU[%d] has empty name", i)
			}
			// CPU percent can be > 100 on multi-core systems
			if proc.CPUPercent < 0 {
				t.Errorf("TopByCPU[%d] CPUPercent = %f is negative", i, proc.CPUPercent)
			}

			t.Logf("  [%d] %s (PID %d): %.1f%%",
				i+1, proc.Name, proc.PID, proc.CPUPercent)

			// List should be sorted by CPU (descending)
			if i > 0 && proc.CPUPercent > data.TopByCPU[i-1].CPUPercent {
				t.Errorf("TopByCPU not sorted: [%d] %.1f%% > [%d] %.1f%%",
					i, proc.CPUPercent, i-1, data.TopByCPU[i-1].CPUPercent)
			}
		}
	} else {
		t.Log("No top CPU processes available")
	}
}

func TestCollectProcessesStatusCounts(t *testing.T) {
	data, err := CollectProcesses()
	if err != nil {
		t.Fatalf("CollectProcesses failed: %v", err)
	}

	// Counts should be non-negative
	if data.Running < 0 {
		t.Error("Running count is negative")
	}
	if data.Sleeping < 0 {
		t.Error("Sleeping count is negative")
	}

	// Note: On Windows, Running and Sleeping may both be 0 because
	// the process status field may use different states
	// As long as TotalCount > 0 (verified in TestCollectProcesses), this is acceptable
	t.Logf("Process status: Running=%d, Sleeping=%d, Total=%d",
		data.Running, data.Sleeping, data.TotalCount)
}

func TestCollectProcessesTopListSizes(t *testing.T) {
	data, err := CollectProcesses()
	if err != nil {
		t.Fatalf("CollectProcesses failed: %v", err)
	}

	// Top lists should not exceed the expected limit (typically 10)
	maxTop := 10
	if len(data.TopByMemory) > maxTop {
		t.Errorf("TopByMemory has %d entries (expected <= %d)", len(data.TopByMemory), maxTop)
	}
	if len(data.TopByCPU) > maxTop {
		t.Errorf("TopByCPU has %d entries (expected <= %d)", len(data.TopByCPU), maxTop)
	}
}

func BenchmarkCollectProcesses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CollectProcesses()
	}
}
