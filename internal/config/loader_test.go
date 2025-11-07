package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".sysinforc")

	configContent := `
format: json
verbose: true
modules:
  cpu: true
  memory: true
  disk: true
smart:
  enable_alerts: true
  alert_thresholds:
    temperature_critical: 70
    temperature_warning: 60
process:
  top_count: 10
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load the config
	cfg, err := LoadConfigFile(configPath)
	if err != nil {
		t.Fatalf("LoadConfigFile() error = %v", err)
	}

	// Verify loaded values
	if cfg.Format != "json" {
		t.Errorf("Format = %q; want %q", cfg.Format, "json")
	}

	if !cfg.Verbose {
		t.Error("Verbose should be true")
	}

	if !cfg.Modules.CPU {
		t.Error("Modules.CPU should be true")
	}

	if !cfg.Modules.Memory {
		t.Error("Modules.Memory should be true")
	}

	if !cfg.Modules.Disk {
		t.Error("Modules.Disk should be true")
	}

	if !cfg.SMART.EnableAlerts {
		t.Error("SMART.EnableAlerts should be true")
	}

	if cfg.SMART.AlertThresholds.TemperatureCritical != 70 {
		t.Errorf("SMART.AlertThresholds.TemperatureCritical = %d; want %d",
			cfg.SMART.AlertThresholds.TemperatureCritical, 70)
	}

	if cfg.Process.TopCount != 10 {
		t.Errorf("Process.TopCount = %d; want %d", cfg.Process.TopCount, 10)
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	// Should return empty config when file doesn't exist
	cfg, err := LoadConfigFile("")
	if err != nil {
		t.Fatalf("LoadConfigFile() error = %v; want nil", err)
	}

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}

	// Empty config should have default/zero values
	if cfg.Format != "" {
		t.Errorf("Expected empty format, got %q", cfg.Format)
	}
}

func TestLoadConfigFileInvalid(t *testing.T) {
	// Create invalid YAML file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".sysinforc")

	invalidYAML := `
format: json
  invalid: indentation
    more: bad
`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Should return error for invalid YAML
	_, err := LoadConfigFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestMergeWithFileConfig(t *testing.T) {
	tests := []struct {
		name        string
		runtime     *Config
		file        *FileConfig
		wantFormat  string
		wantVerbose bool
	}{
		{
			name: "File provides defaults",
			runtime: &Config{
				Format:  "pretty", // default
				Verbose: false,
			},
			file: &FileConfig{
				Format:  "json",
				Verbose: true,
			},
			wantFormat:  "json", // File provides non-default
			wantVerbose: true,   // File sets it
		},
		{
			name: "Empty file config",
			runtime: &Config{
				Format:  "pretty",
				Verbose: false,
			},
			file:        &FileConfig{},
			wantFormat:  "pretty",
			wantVerbose: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.runtime.MergeWithFileConfig(tt.file)

			if tt.runtime.Format != tt.wantFormat {
				t.Errorf("Format = %q; want %q", tt.runtime.Format, tt.wantFormat)
			}

			if tt.runtime.Verbose != tt.wantVerbose {
				t.Errorf("Verbose = %v; want %v", tt.runtime.Verbose, tt.wantVerbose)
			}
		})
	}
}

func TestMergeWithFileConfigModules(t *testing.T) {
	// Test that file config doesn't override when --all is set
	runtime := &Config{
		Modules: ModuleConfig{
			All: true,
		},
	}

	file := &FileConfig{}
	file.Modules.CPU = true
	file.Modules.Memory = true

	runtime.MergeWithFileConfig(file)

	// With All=true, individual flags shouldn't be overridden
	if runtime.Modules.CPU {
		t.Error("CPU should not be set when All is true")
	}

	// Test that file config sets modules when All is false
	runtime2 := &Config{
		Modules: ModuleConfig{
			All: false,
		},
	}

	runtime2.MergeWithFileConfig(file)

	if !runtime2.Modules.CPU {
		t.Error("CPU should be set from file config")
	}

	if !runtime2.Modules.Memory {
		t.Error("Memory should be set from file config")
	}
}

func TestSaveConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config", "sysinfo.yaml")

	cfg := &FileConfig{
		Format:  "json",
		Verbose: true,
	}
	cfg.Modules.CPU = true
	cfg.Modules.Memory = true
	cfg.Process.TopCount = 15

	// Save the config
	if err := SaveConfigFile(cfg, configPath); err != nil {
		t.Fatalf("SaveConfigFile() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load it back
	loaded, err := LoadConfigFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	// Verify values
	if loaded.Format != "json" {
		t.Errorf("Format = %q; want %q", loaded.Format, "json")
	}

	if !loaded.Verbose {
		t.Error("Verbose should be true")
	}

	if !loaded.Modules.CPU {
		t.Error("Modules.CPU should be true")
	}

	if loaded.Process.TopCount != 15 {
		t.Errorf("Process.TopCount = %d; want %d", loaded.Process.TopCount, 15)
	}
}
