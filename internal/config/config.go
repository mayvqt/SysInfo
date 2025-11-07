package config

// Config holds the runtime configuration for the application
type Config struct {
	// Output format: json, text, pretty
	Format string

	// Output file path (empty means stdout)
	OutputFile string

	// Verbosity level
	Verbose bool

	// Full dump mode - collect everything and save to JSON file
	FullDumpToFile bool

	// Module selection flags
	Modules ModuleConfig
}

// ModuleConfig controls which information modules to collect
type ModuleConfig struct {
	All     bool
	System  bool
	CPU     bool
	Memory  bool
	Disk    bool
	Network bool
	Process bool
	SMART   bool
	GPU     bool
}

// NewConfig creates a default configuration
func NewConfig() *Config {
	return &Config{
		Format:     "pretty",
		OutputFile: "",
		Verbose:    false,
		Modules: ModuleConfig{
			All: true,
		},
	}
}

// ShouldCollect determines if a module should be collected
func (c *Config) ShouldCollect(module string) bool {
	if c.Modules.All {
		return true
	}

	switch module {
	case "system":
		return c.Modules.System
	case "cpu":
		return c.Modules.CPU
	case "memory":
		return c.Modules.Memory
	case "disk":
		return c.Modules.Disk
	case "network":
		return c.Modules.Network
	case "process":
		return c.Modules.Process
	case "smart":
		return c.Modules.SMART
	case "gpu":
		return c.Modules.GPU
	default:
		return false
	}
}
