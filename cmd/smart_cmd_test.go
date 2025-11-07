package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mayvqt/sysinfo/internal/analyzer"
	"github.com/mayvqt/sysinfo/internal/config"
)

func TestSmartCommands_Help(t *testing.T) {
	// Test that smart command is registered
	if smartCmd == nil {
		t.Fatal("smartCmd is nil")
	}

	if smartCmd.Use != "smart" {
		t.Errorf("Expected Use='smart', got %s", smartCmd.Use)
	}

	// Test subcommands are registered
	if len(smartCmd.Commands()) != 3 {
		t.Errorf("Expected 3 subcommands, got %d", len(smartCmd.Commands()))
	}

	subcommands := make(map[string]bool)
	for _, cmd := range smartCmd.Commands() {
		subcommands[cmd.Use] = true
	}

	if !subcommands["analyze"] {
		t.Error("Expected 'analyze' subcommand to be registered")
	}
	if !subcommands["history"] {
		t.Error("Expected 'history' subcommand to be registered")
	}
	if !subcommands["check"] {
		t.Error("Expected 'check' subcommand to be registered")
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"1h", "1h0m0s", false},
		{"24h", "24h0m0s", false},
		{"1d", "24h0m0s", false},
		{"7d", "168h0m0s", false},
		{"30d", "720h0m0s", false},
		{"1w", "168h0m0s", false},
		{"1m", "720h0m0s", false},
		{"invalid", "", true},
		{"", "", true},
		{"x", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			duration, err := parseDuration(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				return
			}

			if duration.String() != tt.expected {
				t.Errorf("For input %q: expected %s, got %s", tt.input, tt.expected, duration.String())
			}
		})
	}
}

func TestRepeatString(t *testing.T) {
	tests := []struct {
		str      string
		count    int
		expected string
	}{
		{"=", 5, "====="},
		{"-", 3, "---"},
		{"*", 0, ""},
		{"ab", 2, "abab"},
	}

	for _, tt := range tests {
		result := repeatString(tt.str, tt.count)
		if result != tt.expected {
			t.Errorf("repeatString(%q, %d) = %q, want %q", tt.str, tt.count, result, tt.expected)
		}
	}
}

func TestGetHealthSymbol(t *testing.T) {
	tests := []struct {
		health   string
		expected string
	}{
		{"GOOD", "✓"},
		{"WARNING", "⚠"},
		{"CRITICAL", "✗"},
		{"FAILING", "✗"},
		{"UNKNOWN", "?"},
		{"", "?"},
	}

	for _, tt := range tests {
		symbol := getHealthSymbol(analyzer.HealthStatus(tt.health))
		if symbol != tt.expected {
			t.Errorf("getHealthSymbol(%q) = %q, want %q", tt.health, symbol, tt.expected)
		}
	}
}

func TestSmartAnalyze_NoSMARTData(t *testing.T) {
	// This test verifies the command handles missing SMART data gracefully
	// We can't easily test the full command without sudo, but we can verify
	// it doesn't panic

	// Reset flags
	smartDBPath = ""
	cfg.SMARTAlerts = false

	// Note: This will fail on most systems without smartctl, which is expected
	// The test verifies it fails gracefully, not with a panic
	err := runSmartAnalyze(smartAnalyzeCmd, []string{})

	// We expect an error (no smart data or database issues)
	// but not a panic
	if err == nil {
		t.Log("Unexpected success - system may have SMART data available")
	}
}

func TestSmartCheck_NoSMARTData(t *testing.T) {
	// Similar to analyze test - verify graceful handling
	err := runSmartCheck(smartCheckCmd, []string{})

	// We expect an error or nil (if no smart data, it prints a message and returns nil)
	// The key is it doesn't panic
	t.Logf("runSmartCheck result: %v (expected nil or error)", err)
}

func TestSmartHistory_NoDB(t *testing.T) {
	// Test with a non-existent database path
	tmpDir := t.TempDir()
	nonExistentDB := filepath.Join(tmpDir, "nonexistent.db")

	smartDBPath = nonExistentDB
	smartPeriod = "7d"

	// This should handle the missing database gracefully
	err := runSmartHistory(smartHistoryCmd, []string{})

	// May succeed (empty history) or fail (no db), but shouldn't panic
	t.Logf("runSmartHistory result: %v", err)
}

func TestInitSMARTDatabase_CustomPath(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_smart.db")

	smartDBPath = dbPath

	db, fileConfig, err := initSMARTDatabase()
	if err != nil {
		t.Fatalf("initSMARTDatabase failed: %v", err)
	}
	defer db.Close()

	// Verify database was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database file was not created at %s", dbPath)
	}

	// fileConfig may be nil if no config file exists
	t.Logf("File config: %v", fileConfig)
}

func TestCreateAnalyzer(t *testing.T) {
	// Test with nil config
	analyzer := createAnalyzer(nil)
	if analyzer == nil {
		t.Error("Expected analyzer to be created with nil config")
	}

	// Test with empty config
	analyzer = createAnalyzer(&config.FileConfig{})
	if analyzer == nil {
		t.Error("Expected analyzer to be created with empty config")
	}
}

func TestCreateAlertManager(t *testing.T) {
	// Test with nil config
	manager := createAlertManager(nil)
	if manager == nil {
		t.Error("Expected alert manager to be created with nil config")
	}

	// Test with config but no webhook URL
	cfg := &config.FileConfig{}
	manager = createAlertManager(cfg)
	if manager == nil {
		t.Error("Expected alert manager to be created with config but no webhook")
	}
}

func TestCollectSMARTData(t *testing.T) {
	// This will likely return empty data on systems without smartctl
	// But should not panic
	data, err := collectSMARTData()

	if err != nil {
		t.Logf("collectSMARTData returned error (may be expected): %v", err)
	}

	if data == nil {
		t.Error("Expected non-nil data even on error")
	}

	t.Logf("Collected SMART data for %d devices", len(data.SMARTData))
}
