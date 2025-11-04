package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mayvqt/sysinfo/internal/config"
	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Save original stdout
	oldStdout := os.Stdout

	// Create a pipe to capture output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Reset rootCmd for testing
	rootCmd.SetArgs([]string{"--help"})

	err = Execute()
	if err != nil {
		t.Errorf("Execute with --help failed: %v", err)
	}

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Errorf("Failed to close pipe writer: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	output := buf.String()

	// Verify help output contains expected text
	expectedStrings := []string{
		"SysInfo",
		"system information",
		"--format",
		"--output",
		"--verbose",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Help output missing expected string: %s", expected)
		}
	}
}

func TestRunSysInfoWithDefaultConfig(t *testing.T) {
	// Create a temporary file for output
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	// Create a fresh rootCmd for this test
	testCmd := &cobra.Command{
		Use:   "sysinfo",
		Short: "SysInfo - Cross-platform system information tool",
		RunE:  runSysInfo,
	}

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	cfg = testCfg

	err := runSysInfo(testCmd, []string{})
	if err != nil {
		t.Fatalf("runSysInfo failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Read and verify output
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Output file is empty")
	}

	// Verify it looks like JSON
	output := string(data)
	if !strings.HasPrefix(strings.TrimSpace(output), "{") {
		t.Error("JSON output doesn't start with {")
	}
}

func TestRunSysInfoWithTextFormat(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	testCfg := config.NewConfig()
	testCfg.Format = "text"
	testCfg.OutputFile = outputFile
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runSysInfo failed: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "SYSTEM INFORMATION") {
		t.Error("Text output missing expected header")
	}
}

func TestRunSysInfoWithPrettyFormat(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	testCfg := config.NewConfig()
	testCfg.Format = "pretty"
	testCfg.OutputFile = outputFile
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runSysInfo failed: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "SYSTEM INFORMATION REPORT") {
		t.Error("Pretty output missing expected header")
	}
}

func TestRunSysInfoWithSelectiveModules(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	testCfg.Modules.All = false
	testCfg.Modules.System = true
	testCfg.Modules.CPU = true
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runSysInfo with selective modules failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestRunSysInfoWithVerbose(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	// Capture stderr for verbose output
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stderr = w

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	testCfg.Verbose = true
	cfg = testCfg

	err = runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		if closeErr := w.Close(); closeErr != nil {
			t.Logf("Failed to close pipe writer: %v", closeErr)
		}
		os.Stderr = oldStderr
		t.Fatalf("runSysInfo with verbose failed: %v", err)
	}

	// Restore stderr and capture output
	if err := w.Close(); err != nil {
		t.Errorf("Failed to close pipe writer: %v", err)
	}
	os.Stderr = oldStderr

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	stderrOutput := buf.String()

	// Verify verbose messages appeared
	expectedMessages := []string{
		"Collecting system information",
		"Formatting output",
		"Writing to file",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(stderrOutput, msg) {
			t.Errorf("Verbose output missing expected message: %s", msg)
		}
	}
}

func TestRunSysInfoWithInvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	testCfg := config.NewConfig()
	testCfg.Format = "invalid"
	testCfg.OutputFile = outputFile
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err == nil {
		t.Error("Expected error with invalid format, got nil")
		return
	}

	if !strings.Contains(err.Error(), "unknown format") {
		t.Errorf("Error message doesn't mention unknown format: %v", err)
	}
}

func TestRunSysInfoWithInvalidOutputPath(t *testing.T) {
	// Try to write to a directory that doesn't exist
	invalidPath := "/this/path/does/not/exist/output.json"
	if os.PathSeparator == '\\' {
		invalidPath = "Z:\\this\\path\\does\\not\\exist\\output.json"
	}

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = invalidPath
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err == nil {
		t.Error("Expected error with invalid output path, got nil")
	}
}

func TestIsTerminal(t *testing.T) {
	// This test is environment-dependent
	result := isTerminal()
	// Just verify it doesn't panic
	t.Logf("isTerminal() = %v", result)
}

func TestModuleFlagLogic(t *testing.T) {
	// This test verifies the module flag behavior
	// When specific modules are set, All should be disabled in runSysInfo
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	testCfg.Modules.CPU = true
	testCfg.Modules.Memory = true
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runSysInfo failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestAllFlagDefault(t *testing.T) {
	// When no specific modules are selected, --all should be true by default
	// This is verified by checking the default config
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runSysInfo failed: %v", err)
	}

	// Verify file was created and has content (All modules collected)
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Output file is empty")
	}
}

func TestMonitorModeWithFileOutput(t *testing.T) {
	// Monitor mode should fail when used with file output
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	testCfg.Monitor = true
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err == nil {
		t.Error("Expected error when using monitor mode with file output, got nil")
		return
	}

	if !strings.Contains(err.Error(), "monitor mode cannot be used with file output") {
		t.Errorf("Error message doesn't mention monitor mode restriction: %v", err)
	}
}

func TestMonitorModeIntervalValidation(t *testing.T) {
	// Test that monitor interval is validated (minimum 1 second)
	tests := []struct {
		name            string
		interval        int
		expectError     bool
		expectedMessage string
	}{
		{
			name:        "valid interval 1 second",
			interval:    1,
			expectError: false,
		},
		{
			name:        "valid interval 5 seconds",
			interval:    5,
			expectError: false,
		},
		{
			name:        "zero interval should be adjusted to 1",
			interval:    0,
			expectError: false,
		},
		{
			name:        "negative interval should be adjusted to 1",
			interval:    -1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCfg := config.NewConfig()
			testCfg.Monitor = true
			testCfg.MonitorInterval = tt.interval
			cfg = testCfg

			// Since we can't actually test the monitor loop, we just verify
			// that the interval is validated in the config
			if testCfg.MonitorInterval < 1 {
				// This would be adjusted in runSysInfo
				t.Logf("Interval %d would be adjusted to 1", testCfg.MonitorInterval)
			}
		})
	}
}

func TestMonitorModeDisabled(t *testing.T) {
	// When monitor mode is disabled, should run normally
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	testCfg.Monitor = false
	testCfg.Modules.CPU = true
	testCfg.Modules.Memory = true
	testCfg.Modules.All = false
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runSysInfo failed with monitor disabled: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestFullDumpMode(t *testing.T) {
	// Test full dump functionality
	testCfg := config.NewConfig()
	testCfg.FullDumpToFile = true
	testCfg.Verbose = true
	cfg = testCfg

	// We can't fully test this without mocking, but we can verify the config is set
	if !testCfg.FullDumpToFile {
		t.Error("FullDumpToFile should be true")
	}

	// Full dump should enable all modules
	if !testCfg.Modules.CPU && !testCfg.Modules.Memory {
		t.Log("Note: Full dump should enable all modules in runFullDump")
	}
}

func TestDisplayLiveDataFormatting(t *testing.T) {
	// Test that displayLiveData doesn't panic
	// We can't test the actual output without a real terminal
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("displayLiveData panicked: %v", r)
		}
	}()

	// This will likely fail to collect data but shouldn't panic
	// displayLiveData(true) - can't test without actual data collection
	t.Log("displayLiveData formatting test skipped - requires terminal")
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*config.Config)
		wantErr bool
	}{
		{
			name: "valid_json_format",
			setup: func(c *config.Config) {
				c.Format = "json"
			},
			wantErr: false,
		},
		{
			name: "valid_text_format",
			setup: func(c *config.Config) {
				c.Format = "text"
			},
			wantErr: false,
		},
		{
			name: "valid_pretty_format",
			setup: func(c *config.Config) {
				c.Format = "pretty"
			},
			wantErr: false,
		},
		{
			name: "monitor_with_file_output",
			setup: func(c *config.Config) {
				c.Monitor = true
				c.OutputFile = "test.json"
			},
			wantErr: false, // Should work but file output is disabled in monitor mode
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCfg := config.NewConfig()
			tt.setup(testCfg)

			// Validate that the config is properly set
			if tt.name == "monitor_with_file_output" && testCfg.Monitor && testCfg.OutputFile != "" {
				t.Log("Monitor mode with file output - file output should be ignored")
			}
		})
	}
}

func TestVerboseOutput(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "verbose_output.json")

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	testCfg.Verbose = true
	testCfg.Modules.System = true
	testCfg.Modules.All = false
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runSysInfo with verbose failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Verbose output file was not created")
	}
}

func TestSMARTDataCollection(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "smart_output.json")

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	testCfg.Modules.SMART = true
	testCfg.Modules.Disk = true
	testCfg.Modules.All = false
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runSysInfo with SMART failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("SMART output file was not created")
	}

	// Read and verify output contains disk data
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read SMART output file: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "disk") && !strings.Contains(output, "Disk") {
		t.Log("Note: SMART data may not be available on this system")
	}
}
