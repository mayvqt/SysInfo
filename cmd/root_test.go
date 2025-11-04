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
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset rootCmd for testing
	rootCmd.SetArgs([]string{"--help"})

	err := Execute()
	if err != nil {
		t.Errorf("Execute with --help failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
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
	r, w, _ := os.Pipe()
	os.Stderr = w

	testCfg := config.NewConfig()
	testCfg.Format = "json"
	testCfg.OutputFile = outputFile
	testCfg.Verbose = true
	cfg = testCfg

	err := runSysInfo(&cobra.Command{}, []string{})
	if err != nil {
		w.Close()
		os.Stderr = oldStderr
		t.Fatalf("runSysInfo with verbose failed: %v", err)
	}

	// Restore stderr and capture output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
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
