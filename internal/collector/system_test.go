package collector

import (
	"testing"
)

// TestCollectSystem verifies basic system collection works
func TestCollectSystem(t *testing.T) {
	data, err := CollectSystem()
	if err != nil {
		t.Fatalf("CollectSystem failed: %v", err)
	}

	if data == nil {
		t.Fatal("CollectSystem returned nil data")
	}

	// Hostname should always be available
	if data.Hostname == "" {
		t.Error("Hostname is empty")
	}

	// OS should always be available
	if data.OS == "" {
		t.Error("OS is empty")
	}

	// Platform should be available
	if data.Platform == "" {
		t.Error("Platform is empty")
	}

	// Uptime formatted string should not be empty
	if data.UptimeFormatted == "" {
		t.Error("UptimeFormatted is empty")
	}

	t.Logf("System Info: %s on %s (%s)", data.Hostname, data.OS, data.Platform)
}

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		name    string
		seconds uint64
		want    string
	}{
		{"zero", 0, "0m"},
		{"30 seconds", 30, "0m"},
		{"1 minute", 60, "1m"},
		{"90 seconds", 90, "1m"},
		{"1 hour", 3600, "1h 0m"},
		{"1 hour 30 min", 5400, "1h 30m"},
		{"1 day", 86400, "1d 0h 0m"},
		{"1 day 2 hours 30 min", 95400, "1d 2h 30m"},
		{"7 days", 604800, "7d 0h 0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatUptime(tt.seconds)
			if got != tt.want {
				t.Errorf("formatUptime(%d) = %q; want %q", tt.seconds, got, tt.want)
			}
		})
	}
}

func BenchmarkCollectSystem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CollectSystem()
	}
}

func BenchmarkFormatUptime(b *testing.B) {
	testValues := []uint64{0, 60, 3600, 86400, 604800}

	for _, val := range testValues {
		b.Run("", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				formatUptime(val)
			}
		})
	}
}
