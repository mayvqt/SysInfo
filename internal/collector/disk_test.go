package collector

import (
	"testing"
)

// TestCollectDisk verifies basic disk collection works
func TestCollectDisk(t *testing.T) {
	data, err := CollectDisk(false) // Without SMART
	if err != nil {
		t.Fatalf("CollectDisk failed: %v", err)
	}

	if data == nil {
		t.Fatal("CollectDisk returned nil data")
	}

	// There should be at least one partition on any system
	if len(data.Partitions) == 0 {
		t.Error("No partitions found (expected at least one)")
	}

	// Verify partition data sanity
	for i, part := range data.Partitions {
		if part.Device == "" {
			t.Errorf("Partition[%d] has empty device", i)
		}

		if part.Total == 0 {
			t.Logf("Warning: Partition[%d] (%s) has zero total (may be virtual)", i, part.Device)
		}

		if part.UsedPercent < 0 || part.UsedPercent > 100 {
			t.Errorf("Partition[%d] UsedPercent = %f is out of range [0, 100]", i, part.UsedPercent)
		}

		// Formatted strings should not be empty if Total > 0
		if part.Total > 0 {
			if part.TotalFormatted == "" {
				t.Errorf("Partition[%d] TotalFormatted is empty", i)
			}
			if part.UsedFormatted == "" {
				t.Errorf("Partition[%d] UsedFormatted is empty", i)
			}
			if part.FreeFormatted == "" {
				t.Errorf("Partition[%d] FreeFormatted is empty", i)
			}
		}

		t.Logf("Partition: %s mounted at %s, Type: %s, Total: %s, Used: %.1f%%",
			part.Device, part.MountPoint, part.FSType, part.TotalFormatted, part.UsedPercent)
	}
}

func TestCollectDiskWithSMART(t *testing.T) {
	// SMART collection may require elevated privileges
	data, err := CollectDisk(true)
	if err != nil {
		t.Logf("CollectDisk with SMART failed (may be expected without privileges): %v", err)
		return
	}

	if data == nil {
		t.Fatal("CollectDisk returned nil data")
	}

	// SMART data may be empty if not available or no privileges
	if len(data.SMARTData) > 0 {
		t.Logf("Found %d SMART devices", len(data.SMARTData))
		for i, smart := range data.SMARTData {
			if smart.Device == "" {
				t.Errorf("SMARTData[%d] has empty device", i)
			}
			t.Logf("SMART Device: %s, Healthy: %v, Temp: %d°C",
				smart.Device, smart.Healthy, smart.Temperature)
		}
	} else {
		t.Log("No SMART data available (may require elevated privileges or unsupported on this system)")
	}
}

func TestCollectDiskPhysicalDisks(t *testing.T) {
	data, err := CollectDisk(false)
	if err != nil {
		t.Fatalf("CollectDisk failed: %v", err)
	}

	// Physical disks may or may not be available depending on platform/privileges
	if len(data.PhysicalDisks) > 0 {
		t.Logf("Found %d physical disks", len(data.PhysicalDisks))
		for i, disk := range data.PhysicalDisks {
			if disk.Name == "" {
				t.Errorf("PhysicalDisk[%d] has empty name", i)
			}
			// Size may be 0 on some platforms or if not available
			t.Logf("Physical Disk: %s, Model: %s, Size: %s, Type: %s",
				disk.Name, disk.Model, disk.SizeFormatted, disk.Type)
		}
	} else {
		t.Log("No physical disk information available (platform-dependent)")
	}
}

func TestCollectDiskIOStats(t *testing.T) {
	data, err := CollectDisk(false)
	if err != nil {
		t.Fatalf("CollectDisk failed: %v", err)
	}

	// IO stats may or may not be available
	if len(data.IOStats) > 0 {
		t.Logf("Found %d disk I/O stat entries", len(data.IOStats))
		for i, io := range data.IOStats {
			if io.Name == "" {
				t.Errorf("IOStat[%d] has empty name", i)
			}
			// Counters are uint64 and can't be negative, just verify they exist
			t.Logf("IO Stats: %s, Reads: %d, Writes: %d",
				io.Name, io.ReadCount, io.WriteCount)
		}
	} else {
		t.Log("No I/O statistics available (platform-dependent)")
	}
}

func BenchmarkCollectDisk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CollectDisk(false)
	}
}

func BenchmarkCollectDiskWithSMART(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CollectDisk(true)
	}
}
