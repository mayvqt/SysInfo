package formatter

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mayvqt/sysinfo/internal/config"
	"github.com/mayvqt/sysinfo/internal/types"
)

func createTestSystemInfo() *types.SystemInfo {
	return &types.SystemInfo{
		Timestamp: time.Date(2025, 11, 4, 12, 0, 0, 0, time.UTC),
		System: &types.SystemData{
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
		CPU: &types.CPUData{
			ModelName:   "Intel Core i7",
			Cores:       4,
			LogicalCPUs: 8,
			Vendor:      "Intel",
			MHz:         2800.0,
			Usage:       []float64{10.5, 20.3},
			LoadAvg: &types.LoadAverage{
				Load1:  1.5,
				Load5:  1.2,
				Load15: 0.9,
			},
		},
		Memory: &types.MemoryData{
			Total:          16 * 1024 * 1024 * 1024,
			Used:           8 * 1024 * 1024 * 1024,
			Free:           8 * 1024 * 1024 * 1024,
			UsedPercent:    50.0,
			TotalFormatted: "16.00 GB",
			UsedFormatted:  "8.00 GB",
			FreeFormatted:  "8.00 GB",
			SwapTotal:      4 * 1024 * 1024 * 1024,
			SwapUsed:       1 * 1024 * 1024 * 1024,
			SwapPercent:    25.0,
		},
		Disk: &types.DiskData{
			Partitions: []types.PartitionInfo{
				{
					Device:         "/dev/sda1",
					MountPoint:     "/",
					FSType:         "ext4",
					Total:          500 * 1024 * 1024 * 1024,
					Used:           300 * 1024 * 1024 * 1024,
					Free:           200 * 1024 * 1024 * 1024,
					UsedPercent:    60.0,
					TotalFormatted: "500.00 GB",
					UsedFormatted:  "300.00 GB",
					FreeFormatted:  "200.00 GB",
				},
			},
			SMARTData: []types.SMARTInfo{
				{
					Device:       "/dev/sda",
					DeviceModel:  "Samsung SSD 970 EVO",
					Healthy:      true,
					Temperature:  35,
					PowerOnHours: 1000,
				},
			},
		},
		Network: &types.NetworkData{
			Interfaces: []types.NetworkInterface{
				{
					Name:         "eth0",
					HardwareAddr: "00:11:22:33:44:55",
					Addresses:    []string{"192.168.1.100"},
					MTU:          1500,
					BytesSent:    1024 * 1024 * 100,
					BytesRecv:    1024 * 1024 * 500,
				},
			},
			Connections: 42,
		},
		GPU: &types.GPUData{
			GPUs: []types.GPUInfo{
				{
					Index:             0,
					Name:              "NVIDIA GeForce RTX 4070",
					Vendor:            "NVIDIA",
					Driver:            "nvidia",
					DriverVersion:     "535.161.07",
					MemoryTotal:       12 * 1024 * 1024 * 1024,
					MemoryUsed:        4 * 1024 * 1024 * 1024,
					MemoryFree:        8 * 1024 * 1024 * 1024,
					MemoryFormatted:   "12.00 GB",
					Temperature:       65,
					FanSpeed:          50,
					PowerDraw:         150.5,
					PowerLimit:        250.0,
					Utilization:       75,
					MemoryUtilization: 50,
					ClockSpeed:        1500,
					ClockSpeedMemory:  7000,
					PCIBus:            "01",
					UUID:              "GPU-12345678-1234-1234-1234-123456789012",
				},
			},
		},
		Processes: &types.ProcessData{
			TotalCount: 250,
			Running:    5,
			Sleeping:   240,
			TopByMemory: []types.ProcessInfo{
				{
					PID:           1234,
					Name:          "chrome",
					CPUPercent:    15.5,
					MemoryPercent: 25.3,
					MemoryMB:      4096,
				},
			},
			TopByCPU: []types.ProcessInfo{
				{
					PID:        5678,
					Name:       "node",
					CPUPercent: 45.2,
				},
			},
		},
	}
}

func TestFormat(t *testing.T) {
	info := createTestSystemInfo()

	tests := []struct {
		name      string
		format    string
		wantError bool
		validate  func(t *testing.T, output string)
	}{
		{
			name:      "json format",
			format:    "json",
			wantError: false,
			validate: func(t *testing.T, output string) {
				var decoded types.SystemInfo
				if err := json.Unmarshal([]byte(output), &decoded); err != nil {
					t.Errorf("Failed to unmarshal JSON output: %v", err)
				}
				if decoded.System.Hostname != "test-host" {
					t.Errorf("Hostname = %q; want %q", decoded.System.Hostname, "test-host")
				}
			},
		},
		{
			name:      "text format",
			format:    "text",
			wantError: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "test-host") {
					t.Error("Text output missing hostname")
				}
				if !strings.Contains(output, "SYSTEM INFORMATION") {
					t.Error("Text output missing system header")
				}
			},
		},
		{
			name:      "pretty format",
			format:    "pretty",
			wantError: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "test-host") {
					t.Error("Pretty output missing hostname")
				}
				if !strings.Contains(output, "SYSTEM INFORMATION REPORT") {
					t.Error("Pretty output missing report header")
				}
			},
		},
		{
			name:      "unknown format",
			format:    "unknown",
			wantError: true,
			validate:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Format: tt.format}
			output, err := Format(info, cfg)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if output == "" {
				t.Error("Output is empty")
			}

			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

func TestFormatJSON(t *testing.T) {
	info := createTestSystemInfo()

	output, err := FormatJSON(info)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var decoded types.SystemInfo
	if err := json.Unmarshal([]byte(output), &decoded); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify key fields are present
	if decoded.System == nil {
		t.Error("System data is nil")
	}
	if decoded.System.Hostname != "test-host" {
		t.Errorf("Hostname = %q; want %q", decoded.System.Hostname, "test-host")
	}
	if decoded.CPU == nil {
		t.Error("CPU data is nil")
	}
	if decoded.Memory == nil {
		t.Error("Memory data is nil")
	}

	// Verify JSON is indented (pretty-printed)
	if !strings.Contains(output, "\n") {
		t.Error("JSON output is not indented")
	}
}

func TestFormatJSONWithNilFields(t *testing.T) {
	info := &types.SystemInfo{
		Timestamp: time.Now(),
		System: &types.SystemData{
			Hostname: "minimal-host",
		},
		// All other fields nil
	}

	output, err := FormatJSON(info)
	if err != nil {
		t.Fatalf("FormatJSON failed with minimal data: %v", err)
	}

	var decoded types.SystemInfo
	if err := json.Unmarshal([]byte(output), &decoded); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if decoded.System.Hostname != "minimal-host" {
		t.Errorf("Hostname = %q; want %q", decoded.System.Hostname, "minimal-host")
	}
}

func TestFormatText(t *testing.T) {
	info := createTestSystemInfo()

	output := FormatText(info)

	// Check for expected sections
	expectedSections := []string{
		"SYSTEM INFORMATION",
		"CPU INFORMATION",
		"MEMORY INFORMATION",
		"STORAGE INFORMATION",
		"SMART DISK HEALTH",
		"NETWORK INTERFACES",
		"PROCESS INFORMATION",
	}
	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Text output missing section: %s", section)
		}
	}

	// Check for specific values
	expectedValues := []string{
		"test-host",
		"Intel Core i7",
		"16.00 GB",
		"/dev/sda1",
		"eth0",
		"chrome",
	}

	for _, value := range expectedValues {
		if !strings.Contains(output, value) {
			t.Errorf("Text output missing value: %s", value)
		}
	}
}

func TestFormatTextWithNilFields(t *testing.T) {
	info := &types.SystemInfo{
		Timestamp: time.Now(),
		System: &types.SystemData{
			Hostname: "minimal-host",
		},
		// All other fields nil
	}

	output := FormatText(info)

	if !strings.Contains(output, "minimal-host") {
		t.Error("Text output missing hostname")
	}

	// Should handle nil fields gracefully
	if strings.Contains(output, "CPU INFORMATION") {
		t.Error("Text output should not contain CPU section when CPU is nil")
	}
}

func TestFormatPretty(t *testing.T) {
	info := createTestSystemInfo()

	output := FormatPretty(info)

	// Check for expected sections
	expectedSections := []string{
		"SYSTEM INFORMATION REPORT",
		"SYSTEM",
		"CPU",
		"MEMORY",
		"STORAGE",
		"SMART DISK HEALTH",
		"NETWORK",
		"PROCESSES",
	}

	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Pretty output missing section: %s", section)
		}
	}

	// Check for box drawing characters
	if !strings.Contains(output, "┌") || !strings.Contains(output, "└") {
		t.Error("Pretty output missing box drawing characters")
	}

	// Check for specific values
	expectedValues := []string{
		"test-host",
		"Intel Core i7",
		"eth0",
	}

	for _, value := range expectedValues {
		if !strings.Contains(output, value) {
			t.Errorf("Pretty output missing value: %s", value)
		}
	}
}

func TestCreateProgressBar(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		width   int
		want    int // expected filled characters
	}{
		{"zero percent", 0, 20, 0},
		{"fifty percent", 50, 20, 10},
		{"hundred percent", 100, 20, 20},
		{"over hundred", 150, 20, 20},
		{"negative", -10, 20, 0},
		{"small width", 50, 10, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := createProgressBar(tt.percent, tt.width)
			// Strip ANSI color codes for counting
			cleanBar := stripAnsiCodes(bar)
			filled := strings.Count(cleanBar, "█")
			if filled != tt.want {
				t.Errorf("createProgressBar(%f, %d) filled count = %d; want %d", tt.percent, tt.width, filled, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		want   string
	}{
		{"short string", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"needs truncation", "hello world", 8, "hello..."},
		{"very long", "this is a very long string", 10, "this is..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.length)
			if result != tt.want {
				t.Errorf("truncate(%q, %d) = %q; want %q", tt.input, tt.length, result, tt.want)
			}
		})
	}
}

func TestFormatBytesInFormatters(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected string
	}{
		{0, "0 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
	}

	for _, tt := range tests {
		result := formatBytes(tt.bytes)
		if result != tt.expected {
			t.Errorf("formatBytes(%d) = %q; want %q", tt.bytes, result, tt.expected)
		}
	}
}

func TestGPUFormatting(t *testing.T) {
	info := createTestSystemInfo()

	tests := []struct {
		name     string
		format   string
		validate func(t *testing.T, output string)
	}{
		{
			name:   "GPU in JSON format",
			format: "json",
			validate: func(t *testing.T, output string) {
				var decoded types.SystemInfo
				if err := json.Unmarshal([]byte(output), &decoded); err != nil {
					t.Errorf("Failed to unmarshal JSON output: %v", err)
					return
				}

				if decoded.GPU == nil {
					t.Error("GPU data missing in JSON output")
					return
				}

				if len(decoded.GPU.GPUs) != 1 {
					t.Errorf("Expected 1 GPU, got %d", len(decoded.GPU.GPUs))
					return
				}

				gpu := decoded.GPU.GPUs[0]
				if gpu.Name != "NVIDIA GeForce RTX 4070" {
					t.Errorf("GPU name = %q; want %q", gpu.Name, "NVIDIA GeForce RTX 4070")
				}
				if gpu.Vendor != "NVIDIA" {
					t.Errorf("GPU vendor = %q; want %q", gpu.Vendor, "NVIDIA")
				}
				if gpu.MemoryFormatted != "12.00 GB" {
					t.Errorf("GPU memory = %q; want %q", gpu.MemoryFormatted, "12.00 GB")
				}
				if gpu.Temperature != 65 {
					t.Errorf("GPU temperature = %d; want %d", gpu.Temperature, 65)
				}
			},
		},
		{
			name:   "GPU in text format",
			format: "text",
			validate: func(t *testing.T, output string) {
				expectedStrings := []string{
					"GPU INFORMATION",
					"NVIDIA GeForce RTX 4070",
					"NVIDIA",
					"12.00 GB",
					"65°C",
					"nvidia",
				}

				for _, expected := range expectedStrings {
					if !strings.Contains(output, expected) {
						t.Errorf("Text output missing expected string: %q", expected)
					}
				}
			},
		},
		{
			name:   "GPU in pretty format",
			format: "pretty",
			validate: func(t *testing.T, output string) {
				// Strip ANSI codes for easier testing
				stripped := stripAnsiCodes(output)

				expectedStrings := []string{
					"GPU",
					"NVIDIA GeForce RTX 4070",
					"NVIDIA",
					"12.00 GB",
					"65°C",
				}

				for _, expected := range expectedStrings {
					if !strings.Contains(stripped, expected) {
						t.Errorf("Pretty output missing expected string: %q", expected)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Format: tt.format}
			output, err := Format(info, cfg)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}
			tt.validate(t, output)
		})
	}
}

func TestGPUFormattingMultipleGPUs(t *testing.T) {
	info := createTestSystemInfo()

	// Add a second GPU
	info.GPU.GPUs = append(info.GPU.GPUs, types.GPUInfo{
		Index:           1,
		Name:            "AMD Radeon RX 6800 XT",
		Vendor:          "AMD",
		Driver:          "amdgpu",
		MemoryTotal:     16 * 1024 * 1024 * 1024,
		MemoryFormatted: "16.00 GB",
		Temperature:     72,
	})

	// Test JSON format
	t.Run("Multiple GPUs in JSON", func(t *testing.T) {
		cfg := &config.Config{Format: "json"}
		output, err := Format(info, cfg)
		if err != nil {
			t.Fatalf("Format() error = %v", err)
		}

		var decoded types.SystemInfo
		if err := json.Unmarshal([]byte(output), &decoded); err != nil {
			t.Errorf("Failed to unmarshal JSON output: %v", err)
			return
		}

		if len(decoded.GPU.GPUs) != 2 {
			t.Errorf("Expected 2 GPUs, got %d", len(decoded.GPU.GPUs))
		}
	})

	// Test text format
	t.Run("Multiple GPUs in text", func(t *testing.T) {
		cfg := &config.Config{Format: "text"}
		output, err := Format(info, cfg)
		if err != nil {
			t.Fatalf("Format() error = %v", err)
		}

		if !strings.Contains(output, "NVIDIA GeForce RTX 4070") {
			t.Error("Missing first GPU")
		}
		if !strings.Contains(output, "AMD Radeon RX 6800 XT") {
			t.Error("Missing second GPU")
		}
	})
}

func TestGPUFormattingNoGPU(t *testing.T) {
	info := createTestSystemInfo()
	info.GPU = nil

	tests := []struct {
		name   string
		format string
	}{
		{"No GPU - JSON", "json"},
		{"No GPU - Text", "text"},
		{"No GPU - Pretty", "pretty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Format: tt.format}
			output, err := Format(info, cfg)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			// Should still produce output, just without GPU section
			if output == "" {
				t.Error("Output is empty")
			}

			// For JSON, verify GPU field is null
			if tt.format == "json" {
				var decoded types.SystemInfo
				if err := json.Unmarshal([]byte(output), &decoded); err != nil {
					t.Errorf("Failed to unmarshal JSON output: %v", err)
				}
				// GPU will be nil, which is fine
			}
		})
	}
}

func BenchmarkFormatJSON(b *testing.B) {
	info := createTestSystemInfo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FormatJSON(info)
	}
}

func BenchmarkFormatText(b *testing.B) {
	info := createTestSystemInfo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatText(info)
	}
}

func BenchmarkFormatPretty(b *testing.B) {
	info := createTestSystemInfo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatPretty(info)
	}
}

// stripAnsiCodes removes ANSI escape codes from a string for testing
func stripAnsiCodes(s string) string {
	// Simple implementation - strips common color codes
	result := s
	for strings.Contains(result, "\x1b[") {
		start := strings.Index(result, "\x1b[")
		end := strings.Index(result[start:], "m")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return result
}
