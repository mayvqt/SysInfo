package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg == nil {
		t.Fatal("NewConfig() returned nil")
	}

	if cfg.Format != "pretty" {
		t.Errorf("Default format = %q; want %q", cfg.Format, "pretty")
	}

	if cfg.OutputFile != "" {
		t.Errorf("Default output file = %q; want empty string", cfg.OutputFile)
	}

	if cfg.Verbose {
		t.Error("Default verbose = true; want false")
	}

	if cfg.Monitor {
		t.Error("Default monitor = true; want false")
	}

	if cfg.MonitorInterval != 2 {
		t.Errorf("Default monitor interval = %d; want 2", cfg.MonitorInterval)
	}

	if !cfg.Modules.All {
		t.Error("Default modules.All = false; want true")
	}
}

func TestShouldCollect(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		module   string
		expected bool
	}{
		{
			name: "all modules enabled",
			config: &Config{
				Modules: ModuleConfig{All: true},
			},
			module:   "system",
			expected: true,
		},
		{
			name: "all modules enabled - any module",
			config: &Config{
				Modules: ModuleConfig{All: true},
			},
			module:   "cpu",
			expected: true,
		},
		{
			name: "only system module",
			config: &Config{
				Modules: ModuleConfig{
					All:    false,
					System: true,
				},
			},
			module:   "system",
			expected: true,
		},
		{
			name: "system module disabled when not all",
			config: &Config{
				Modules: ModuleConfig{
					All:    false,
					System: false,
					CPU:    true,
				},
			},
			module:   "system",
			expected: false,
		},
		{
			name: "cpu module enabled",
			config: &Config{
				Modules: ModuleConfig{
					All: false,
					CPU: true,
				},
			},
			module:   "cpu",
			expected: true,
		},
		{
			name: "memory module enabled",
			config: &Config{
				Modules: ModuleConfig{
					All:    false,
					Memory: true,
				},
			},
			module:   "memory",
			expected: true,
		},
		{
			name: "disk module enabled",
			config: &Config{
				Modules: ModuleConfig{
					All:  false,
					Disk: true,
				},
			},
			module:   "disk",
			expected: true,
		},
		{
			name: "network module enabled",
			config: &Config{
				Modules: ModuleConfig{
					All:     false,
					Network: true,
				},
			},
			module:   "network",
			expected: true,
		},
		{
			name: "process module enabled",
			config: &Config{
				Modules: ModuleConfig{
					All:     false,
					Process: true,
				},
			},
			module:   "process",
			expected: true,
		},
		{
			name: "smart module enabled",
			config: &Config{
				Modules: ModuleConfig{
					All:   false,
					SMART: true,
				},
			},
			module:   "smart",
			expected: true,
		},
		{
			name: "unknown module",
			config: &Config{
				Modules: ModuleConfig{All: false},
			},
			module:   "unknown",
			expected: false,
		},
		{
			name: "unknown module with all disabled",
			config: &Config{
				Modules: ModuleConfig{All: false},
			},
			module:   "unknown",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ShouldCollect(tt.module)
			if result != tt.expected {
				t.Errorf("ShouldCollect(%q) = %v; want %v", tt.module, result, tt.expected)
			}
		})
	}
}

func TestModuleConfigCombinations(t *testing.T) {
	// Test that when All is true, individual module settings are ignored
	cfg := &Config{
		Modules: ModuleConfig{
			All:     true,
			System:  false,
			CPU:     false,
			Memory:  false,
			Disk:    false,
			Network: false,
			Process: false,
			SMART:   false,
		},
	}

	modules := []string{"system", "cpu", "memory", "disk", "network", "process", "smart"}
	for _, module := range modules {
		if !cfg.ShouldCollect(module) {
			t.Errorf("With All=true, ShouldCollect(%q) should be true", module)
		}
	}
}

func TestConfigFields(t *testing.T) {
	cfg := &Config{
		Format:     "json",
		OutputFile: "/tmp/output.json",
		Verbose:    true,
		Modules: ModuleConfig{
			All:     false,
			System:  true,
			CPU:     true,
			Memory:  false,
			Disk:    false,
			Network: false,
			Process: false,
			SMART:   false,
		},
	}

	if cfg.Format != "json" {
		t.Errorf("Format = %q; want %q", cfg.Format, "json")
	}

	if cfg.OutputFile != "/tmp/output.json" {
		t.Errorf("OutputFile = %q; want %q", cfg.OutputFile, "/tmp/output.json")
	}

	if !cfg.Verbose {
		t.Error("Verbose = false; want true")
	}

	if cfg.ShouldCollect("system") != true {
		t.Error("ShouldCollect(system) = false; want true")
	}

	if cfg.ShouldCollect("memory") != false {
		t.Error("ShouldCollect(memory) = true; want false")
	}
}

func TestMonitorModeConfig(t *testing.T) {
	tests := []struct {
		name            string
		monitor         bool
		monitorInterval int
		wantMonitor     bool
		wantInterval    int
	}{
		{
			name:            "monitor enabled with default interval",
			monitor:         true,
			monitorInterval: 2,
			wantMonitor:     true,
			wantInterval:    2,
		},
		{
			name:            "monitor enabled with custom interval",
			monitor:         true,
			monitorInterval: 5,
			wantMonitor:     true,
			wantInterval:    5,
		},
		{
			name:            "monitor disabled",
			monitor:         false,
			monitorInterval: 2,
			wantMonitor:     false,
			wantInterval:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Monitor:         tt.monitor,
				MonitorInterval: tt.monitorInterval,
			}

			if cfg.Monitor != tt.wantMonitor {
				t.Errorf("Monitor = %v; want %v", cfg.Monitor, tt.wantMonitor)
			}

			if cfg.MonitorInterval != tt.wantInterval {
				t.Errorf("MonitorInterval = %d; want %d", cfg.MonitorInterval, tt.wantInterval)
			}
		})
	}
}
