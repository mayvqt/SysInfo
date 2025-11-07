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

func TestFullDumpMode(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Create temp directory and change to it so dump file is created there
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	testCfg := config.NewConfig()
	testCfg.FullDumpToFile = true
	cfg = testCfg

	// Run the full dump
	err = runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		t.Fatalf("runFullDump failed: %v", err)
	}

	// Verify the dump file was created
	dumpFile := filepath.Join(tmpDir, "sysinfo_dump.json")
	if _, err := os.Stat(dumpFile); os.IsNotExist(err) {
		t.Error("Full dump file was not created")
		return
	}

	// Verify file has content
	data, err := os.ReadFile(dumpFile)
	if err != nil {
		t.Fatalf("Failed to read dump file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Dump file is empty")
	}

	// Verify it's valid JSON
	output := string(data)
	if !strings.HasPrefix(strings.TrimSpace(output), "{") {
		t.Error("Dump file doesn't contain valid JSON")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{
			name:    "valid_json_format",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "valid_text_format",
			format:  "text",
			wantErr: false,
		},
		{
			name:    "valid_pretty_format",
			format:  "pretty",
			wantErr: false,
		},
		{
			name:    "invalid_format",
			format:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "output.txt")

			testCfg := config.NewConfig()
			testCfg.Format = tt.format
			testCfg.OutputFile = outputFile
			cfg = testCfg

			err := runSysInfo(&cobra.Command{}, []string{})
			if (err != nil) != tt.wantErr {
				t.Errorf("runSysInfo() error = %v, wantErr %v", err, tt.wantErr)
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
