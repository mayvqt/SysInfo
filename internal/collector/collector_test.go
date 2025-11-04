package collector

import (
	"testing"

	"github.com/mayvqt/sysinfo/internal/config"
)

// TestCollect is an integration test that verifies Collect can run without panicking
func TestCollect(t *testing.T) {
	cfg := config.NewConfig()

	info, err := Collect(cfg)
	if err != nil {
		t.Logf("Collect returned error (may be expected on some systems): %v", err)
	}

	if info == nil {
		t.Fatal("Collect returned nil info")
	}

	// Verify timestamp is set
	if info.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}

	// With default config (All: true), we should have attempted to collect everything
	// Note: Some fields may still be nil if collection failed
}

func TestCollectWithSelectiveModules(t *testing.T) {
	cfg := &config.Config{
		Modules: config.ModuleConfig{
			All:    false,
			System: true,
			CPU:    true,
		},
	}

	info, err := Collect(cfg)
	if err != nil {
		t.Logf("Collect returned error: %v", err)
	}

	if info == nil {
		t.Fatal("Collect returned nil info")
	}

	// System and CPU should be attempted (may still be nil if collection failed)
	// Memory, Disk, Network, Processes should definitely be nil since not requested
	if cfg.ShouldCollect("memory") {
		t.Error("Memory should not be collected")
	}
	if cfg.ShouldCollect("disk") {
		t.Error("Disk should not be collected")
	}
}

func TestCollectWithVerbose(t *testing.T) {
	cfg := &config.Config{
		Verbose: true,
		Modules: config.ModuleConfig{All: true},
	}

	// This test just ensures verbose mode doesn't cause panics
	info, err := Collect(cfg)
	if err != nil {
		t.Logf("Collect with verbose returned error: %v", err)
	}

	if info == nil {
		t.Fatal("Collect returned nil info")
	}
}

func TestCollectAllModulesIndividually(t *testing.T) {
	modules := []string{"system", "cpu", "memory", "disk", "network", "process"}

	for _, module := range modules {
		t.Run(module, func(t *testing.T) {
			cfg := &config.Config{
				Modules: config.ModuleConfig{All: false},
			}

			// Enable just this module
			switch module {
			case "system":
				cfg.Modules.System = true
			case "cpu":
				cfg.Modules.CPU = true
			case "memory":
				cfg.Modules.Memory = true
			case "disk":
				cfg.Modules.Disk = true
			case "network":
				cfg.Modules.Network = true
			case "process":
				cfg.Modules.Process = true
			}

			info, err := Collect(cfg)
			if err != nil {
				t.Logf("Collect for %s returned error: %v", module, err)
			}

			if info == nil {
				t.Fatalf("Collect returned nil info for module %s", module)
			}

			if !cfg.ShouldCollect(module) {
				t.Errorf("Module %s should be collected", module)
			}
		})
	}
}

func BenchmarkCollect(b *testing.B) {
	cfg := config.NewConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Collect(cfg)
	}
}

func BenchmarkCollectSystemOnly(b *testing.B) {
	cfg := &config.Config{
		Modules: config.ModuleConfig{
			All:    false,
			System: true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Collect(cfg)
	}
}
