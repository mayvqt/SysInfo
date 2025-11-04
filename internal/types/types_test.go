package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSystemInfoMarshaling(t *testing.T) {
	now := time.Now()
	info := &SystemInfo{
		Timestamp: now,
		System: &SystemData{
			Hostname:        "test-host",
			OS:              "linux",
			Platform:        "ubuntu",
			PlatformFamily:  "debian",
			PlatformVersion: "22.04",
			KernelVersion:   "5.15.0",
			KernelArch:      "x86_64",
			Uptime:          3600,
			UptimeFormatted: "1h 0m 0s",
			BootTime:        1234567890,
			Procs:           150,
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal SystemInfo: %v", err)
	}

	// Test JSON unmarshaling
	var decoded SystemInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal SystemInfo: %v", err)
	}

	// Verify key fields
	if decoded.System.Hostname != "test-host" {
		t.Errorf("Hostname = %q; want %q", decoded.System.Hostname, "test-host")
	}
	if decoded.System.OS != "linux" {
		t.Errorf("OS = %q; want %q", decoded.System.OS, "linux")
	}
	if decoded.System.Uptime != 3600 {
		t.Errorf("Uptime = %d; want %d", decoded.System.Uptime, 3600)
	}
}

func TestCPUDataMarshaling(t *testing.T) {
	cpu := &CPUData{
		ModelName:   "Intel Core i7",
		Cores:       4,
		LogicalCPUs: 8,
		Vendor:      "Intel",
		Family:      "6",
		Model:       "142",
		Stepping:    10,
		MHz:         2800.0,
		MinMHz:      800.0,
		MaxMHz:      4200.0,
		CacheSize:   8192,
		Usage:       []float64{10.5, 20.3, 15.7, 8.2},
		LoadAvg: &LoadAverage{
			Load1:  1.5,
			Load5:  1.2,
			Load15: 0.9,
		},
		Flags:     []string{"fpu", "vme", "de"},
		Microcode: "0xb4",
	}

	data, err := json.Marshal(cpu)
	if err != nil {
		t.Fatalf("Failed to marshal CPUData: %v", err)
	}

	var decoded CPUData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal CPUData: %v", err)
	}

	if decoded.ModelName != "Intel Core i7" {
		t.Errorf("ModelName = %q; want %q", decoded.ModelName, "Intel Core i7")
	}
	if decoded.Cores != 4 {
		t.Errorf("Cores = %d; want %d", decoded.Cores, 4)
	}
	if len(decoded.Usage) != 4 {
		t.Errorf("len(Usage) = %d; want %d", len(decoded.Usage), 4)
	}
	if decoded.LoadAvg == nil {
		t.Error("LoadAvg is nil; want non-nil")
	} else {
		if decoded.LoadAvg.Load1 != 1.5 {
			t.Errorf("LoadAvg.Load1 = %f; want %f", decoded.LoadAvg.Load1, 1.5)
		}
	}
}

func TestMemoryDataMarshaling(t *testing.T) {
	mem := &MemoryData{
		Total:          16 * 1024 * 1024 * 1024,
		Available:      8 * 1024 * 1024 * 1024,
		Used:           8 * 1024 * 1024 * 1024,
		UsedPercent:    50.0,
		Free:           8 * 1024 * 1024 * 1024,
		TotalFormatted: "16.00 GB",
		UsedFormatted:  "8.00 GB",
		FreeFormatted:  "8.00 GB",
		SwapTotal:      4 * 1024 * 1024 * 1024,
		SwapUsed:       1 * 1024 * 1024 * 1024,
		SwapFree:       3 * 1024 * 1024 * 1024,
		SwapPercent:    25.0,
		Modules: []MemoryModule{
			{
				Locator:      "DIMM 0",
				Capacity:     8 * 1024 * 1024 * 1024,
				Speed:        3200,
				Type:         "DDR4",
				Manufacturer: "Samsung",
				PartNumber:   "M471A1K43CB1-CTD",
				SerialNumber: "12345678",
				FormFactor:   "SODIMM",
			},
		},
	}

	data, err := json.Marshal(mem)
	if err != nil {
		t.Fatalf("Failed to marshal MemoryData: %v", err)
	}

	var decoded MemoryData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal MemoryData: %v", err)
	}

	if decoded.Total != 16*1024*1024*1024 {
		t.Errorf("Total = %d; want %d", decoded.Total, 16*1024*1024*1024)
	}
	if decoded.UsedPercent != 50.0 {
		t.Errorf("UsedPercent = %f; want %f", decoded.UsedPercent, 50.0)
	}
	if len(decoded.Modules) != 1 {
		t.Errorf("len(Modules) = %d; want %d", len(decoded.Modules), 1)
	}
	if decoded.Modules[0].Locator != "DIMM 0" {
		t.Errorf("Module[0].Locator = %q; want %q", decoded.Modules[0].Locator, "DIMM 0")
	}
}

func TestDiskDataMarshaling(t *testing.T) {
	disk := &DiskData{
		Partitions: []PartitionInfo{
			{
				Device:         "/dev/sda1",
				MountPoint:     "/",
				FSType:         "ext4",
				Total:          500 * 1024 * 1024 * 1024,
				Free:           200 * 1024 * 1024 * 1024,
				Used:           300 * 1024 * 1024 * 1024,
				UsedPercent:    60.0,
				TotalFormatted: "500.00 GB",
				UsedFormatted:  "300.00 GB",
				FreeFormatted:  "200.00 GB",
				InodesTotal:    1000000,
				InodesUsed:     500000,
				InodesFree:     500000,
			},
		},
		PhysicalDisks: []PhysicalDisk{
			{
				Name:          "sda",
				Model:         "Samsung SSD 970 EVO",
				SerialNumber:  "S5H9NS0N123456",
				Size:          500 * 1024 * 1024 * 1024,
				SizeFormatted: "500.00 GB",
				Type:          "SSD",
				Interface:     "NVMe",
				Removable:     false,
			},
		},
		IOStats: []DiskIOStat{
			{
				Name:       "sda",
				ReadCount:  10000,
				WriteCount: 5000,
				ReadBytes:  1024 * 1024 * 100,
				WriteBytes: 1024 * 1024 * 50,
				ReadTime:   1000,
				WriteTime:  500,
				IoTime:     1500,
			},
		},
		SMARTData: []SMARTInfo{
			{
				Device:       "/dev/sda",
				Serial:       "S5H9NS0N123456",
				ModelFamily:  "Samsung 970 EVO",
				DeviceModel:  "Samsung SSD 970 EVO 500GB",
				Capacity:     500 * 1024 * 1024 * 1024,
				Healthy:      true,
				Temperature:  35,
				PowerOnHours: 1000,
				Attributes: map[string]string{
					"Reallocated_Sector_Count": "0",
					"Power_On_Hours":           "1000",
				},
			},
		},
	}

	data, err := json.Marshal(disk)
	if err != nil {
		t.Fatalf("Failed to marshal DiskData: %v", err)
	}

	var decoded DiskData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal DiskData: %v", err)
	}

	if len(decoded.Partitions) != 1 {
		t.Errorf("len(Partitions) = %d; want %d", len(decoded.Partitions), 1)
	}
	if decoded.Partitions[0].Device != "/dev/sda1" {
		t.Errorf("Partition[0].Device = %q; want %q", decoded.Partitions[0].Device, "/dev/sda1")
	}
	if len(decoded.SMARTData) != 1 {
		t.Errorf("len(SMARTData) = %d; want %d", len(decoded.SMARTData), 1)
	}
	if !decoded.SMARTData[0].Healthy {
		t.Error("SMARTData[0].Healthy = false; want true")
	}
}

func TestNetworkDataMarshaling(t *testing.T) {
	net := &NetworkData{
		Interfaces: []NetworkInterface{
			{
				Name:         "eth0",
				HardwareAddr: "00:11:22:33:44:55",
				Addresses:    []string{"192.168.1.100", "fe80::1"},
				Flags:        []string{"up", "broadcast", "running"},
				MTU:          1500,
				BytesSent:    1024 * 1024 * 100,
				BytesRecv:    1024 * 1024 * 500,
				PacketsSent:  10000,
				PacketsRecv:  50000,
				ErrorsIn:     0,
				ErrorsOut:    0,
				DropsIn:      5,
				DropsOut:     2,
			},
		},
		Connections: 42,
	}

	data, err := json.Marshal(net)
	if err != nil {
		t.Fatalf("Failed to marshal NetworkData: %v", err)
	}

	var decoded NetworkData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal NetworkData: %v", err)
	}

	if len(decoded.Interfaces) != 1 {
		t.Errorf("len(Interfaces) = %d; want %d", len(decoded.Interfaces), 1)
	}
	if decoded.Interfaces[0].Name != "eth0" {
		t.Errorf("Interface[0].Name = %q; want %q", decoded.Interfaces[0].Name, "eth0")
	}
	if decoded.Connections != 42 {
		t.Errorf("Connections = %d; want %d", decoded.Connections, 42)
	}
}

func TestProcessDataMarshaling(t *testing.T) {
	proc := &ProcessData{
		TotalCount: 250,
		Running:    5,
		Sleeping:   240,
		TopByMemory: []ProcessInfo{
			{
				PID:           1234,
				Name:          "chrome",
				Username:      "user",
				CPUPercent:    15.5,
				MemoryPercent: 25.3,
				MemoryMB:      4096,
				Status:        "running",
				CreateTime:    1234567890,
			},
		},
		TopByCPU: []ProcessInfo{
			{
				PID:           5678,
				Name:          "node",
				Username:      "user",
				CPUPercent:    45.2,
				MemoryPercent: 10.1,
				MemoryMB:      1024,
				Status:        "running",
				CreateTime:    1234567891,
			},
		},
	}

	data, err := json.Marshal(proc)
	if err != nil {
		t.Fatalf("Failed to marshal ProcessData: %v", err)
	}

	var decoded ProcessData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ProcessData: %v", err)
	}

	if decoded.TotalCount != 250 {
		t.Errorf("TotalCount = %d; want %d", decoded.TotalCount, 250)
	}
	if len(decoded.TopByMemory) != 1 {
		t.Errorf("len(TopByMemory) = %d; want %d", len(decoded.TopByMemory), 1)
	}
	if decoded.TopByMemory[0].Name != "chrome" {
		t.Errorf("TopByMemory[0].Name = %q; want %q", decoded.TopByMemory[0].Name, "chrome")
	}
}

func TestOmitemptyFields(t *testing.T) {
	// Test that omitempty works correctly for nil pointers
	info := &SystemInfo{
		Timestamp: time.Now(),
		// All other fields are nil
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal SystemInfo: %v", err)
	}

	// Should only contain timestamp field
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Only timestamp should be present
	if len(decoded) != 1 {
		t.Errorf("Expected 1 field (timestamp), got %d fields: %v", len(decoded), decoded)
	}

	if _, ok := decoded["timestamp"]; !ok {
		t.Error("timestamp field is missing")
	}
}
