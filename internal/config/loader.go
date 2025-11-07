package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfig represents the structure of the configuration file
type FileConfig struct {
	// Default output format
	Format string `yaml:"format,omitempty"`

	// Default output file
	OutputFile string `yaml:"output_file,omitempty"`

	// Verbosity
	Verbose bool `yaml:"verbose,omitempty"`

	// Default modules to collect
	Modules struct {
		System  bool `yaml:"system,omitempty"`
		CPU     bool `yaml:"cpu,omitempty"`
		Memory  bool `yaml:"memory,omitempty"`
		Disk    bool `yaml:"disk,omitempty"`
		Network bool `yaml:"network,omitempty"`
		Process bool `yaml:"process,omitempty"`
		SMART   bool `yaml:"smart,omitempty"`
		GPU     bool `yaml:"gpu,omitempty"`
	} `yaml:"modules,omitempty"`

	// SMART monitoring configuration
	SMART struct {
		EnableAlerts    bool `yaml:"enable_alerts,omitempty"`
		AlertThresholds struct {
			TemperatureCritical int `yaml:"temperature_critical,omitempty"`
			TemperatureWarning  int `yaml:"temperature_warning,omitempty"`
		} `yaml:"alert_thresholds,omitempty"`
		WebhookURL string `yaml:"webhook_url,omitempty"`
	} `yaml:"smart,omitempty"`

	// Process monitoring configuration
	Process struct {
		TopCount int `yaml:"top_count,omitempty"` // Number of top processes to show
	} `yaml:"process,omitempty"`

	// Display preferences
	Display struct {
		UseASCII bool `yaml:"use_ascii,omitempty"` // Force ASCII output instead of Unicode
	} `yaml:"display,omitempty"`
}

// LoadConfigFile attempts to load configuration from file
// Search order: ./.sysinforc, ~/.config/sysinfo/config.yaml, ~/.sysinforc
func LoadConfigFile(customPath string) (*FileConfig, error) {
	var configPath string

	if customPath != "" {
		// Use custom path if provided
		configPath = customPath
	} else {
		// Search in standard locations
		searchPaths := []string{
			".sysinforc",
			".sysinfo.yaml",
			filepath.Join(os.Getenv("HOME"), ".config", "sysinfo", "config.yaml"),
			filepath.Join(os.Getenv("HOME"), ".sysinforc"),
		}

		for _, path := range searchPaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}

		// If no config file found, return empty config (use defaults)
		if configPath == "" {
			return &FileConfig{}, nil
		}
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Parse YAML
	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	return &cfg, nil
}

// MergeWithFileConfig merges file configuration with runtime config
// CLI flags take precedence over file config
func (c *Config) MergeWithFileConfig(fileConfig *FileConfig) {
	// Only apply file config if CLI didn't override
	if c.Format == "pretty" && fileConfig.Format != "" {
		c.Format = fileConfig.Format
	}

	if c.OutputFile == "" && fileConfig.OutputFile != "" {
		c.OutputFile = fileConfig.OutputFile
	}

	if !c.Verbose && fileConfig.Verbose {
		c.Verbose = fileConfig.Verbose
	}

	// Merge module settings if --all wasn't specified
	if !c.Modules.All {
		if fileConfig.Modules.System {
			c.Modules.System = true
		}
		if fileConfig.Modules.CPU {
			c.Modules.CPU = true
		}
		if fileConfig.Modules.Memory {
			c.Modules.Memory = true
		}
		if fileConfig.Modules.Disk {
			c.Modules.Disk = true
		}
		if fileConfig.Modules.Network {
			c.Modules.Network = true
		}
		if fileConfig.Modules.Process {
			c.Modules.Process = true
		}
		if fileConfig.Modules.SMART {
			c.Modules.SMART = true
		}
		if fileConfig.Modules.GPU {
			c.Modules.GPU = true
		}
	}
}

// SaveConfigFile saves current configuration to a file
func SaveConfigFile(cfg *FileConfig, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
