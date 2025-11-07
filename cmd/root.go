package cmd

import (
	"fmt"
	"os"

	"github.com/mayvqt/sysinfo/internal/collector"
	"github.com/mayvqt/sysinfo/internal/config"
	"github.com/mayvqt/sysinfo/internal/formatter"
	"github.com/spf13/cobra"
)

var cfg *config.Config
var configFile string

var rootCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "SysInfo - Cross-platform system information tool",
	Long: `SysInfo is a comprehensive cross-platform system information tool
that collects and displays detailed information about your computer including
CPU, memory, disk, network, processes, and SMART data.`,
	RunE: runSysInfo,
}

func init() {
	cfg = config.NewConfig()

	// Configuration file
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file (default: searches for .sysinforc, ~/.config/sysinfo/config.yaml)")

	// Output options
	rootCmd.Flags().StringVarP(&cfg.Format, "format", "f", "pretty", "Output format: json, text, pretty")
	rootCmd.Flags().StringVarP(&cfg.OutputFile, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Verbose output")

	// Full dump mode
	rootCmd.Flags().BoolVar(&cfg.FullDumpToFile, "full-dump", false, "Collect ALL system information and save to sysinfo_dump.json")

	// Module selection flags
	rootCmd.Flags().BoolVar(&cfg.Modules.All, "all", true, "Collect all information")
	rootCmd.Flags().BoolVar(&cfg.Modules.System, "system", false, "Collect system information")
	rootCmd.Flags().BoolVar(&cfg.Modules.CPU, "cpu", false, "Collect CPU information")
	rootCmd.Flags().BoolVar(&cfg.Modules.Memory, "memory", false, "Collect memory information")
	rootCmd.Flags().BoolVar(&cfg.Modules.Disk, "disk", false, "Collect disk information")
	rootCmd.Flags().BoolVar(&cfg.Modules.Network, "network", false, "Collect network information")
	rootCmd.Flags().BoolVar(&cfg.Modules.Process, "process", false, "Collect process information")
	rootCmd.Flags().BoolVar(&cfg.Modules.SMART, "smart", false, "Collect SMART disk data (may require elevated privileges)")
	rootCmd.Flags().BoolVar(&cfg.Modules.GPU, "gpu", false, "Collect GPU information")
}

func Execute() error {
	return rootCmd.Execute()
}

func runSysInfo(cmd *cobra.Command, args []string) error {
	// Load configuration file if it exists
	fileConfig, err := config.LoadConfigFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	// Merge file config with CLI flags (CLI takes precedence)
	cfg.MergeWithFileConfig(fileConfig)

	// Handle full dump mode
	if cfg.FullDumpToFile {
		return runFullDump()
	}

	// If any specific module is selected, disable --all
	if cfg.Modules.System || cfg.Modules.CPU || cfg.Modules.Memory ||
		cfg.Modules.Disk || cfg.Modules.Network || cfg.Modules.Process || cfg.Modules.SMART || cfg.Modules.GPU {
		cfg.Modules.All = false
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Collecting system information...\n")
	}

	// Collect system information
	info, err := collector.Collect(cfg)
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Formatting output...\n")
	}

	// Format output
	output, err := formatter.Format(info, cfg)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Write output
	if cfg.OutputFile != "" {
		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Writing to file: %s\n", cfg.OutputFile)
		}
		err = os.WriteFile(cfg.OutputFile, []byte(output), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Output written to: %s\n", cfg.OutputFile)
	} else {
		fmt.Print(output)
	}

	// Check if we should pause (when double-clicked, not running from terminal)
	waitForEnter()

	return nil
}

// runFullDump collects all possible system information and saves to JSON file
func runFullDump() error {
	fmt.Fprintf(os.Stderr, "Starting comprehensive system information dump...\n")
	fmt.Fprintf(os.Stderr, "This will collect ALL available data (may take a moment)...\n\n")

	// Create a config to collect everything
	dumpConfig := config.NewConfig()
	dumpConfig.Modules.All = true
	dumpConfig.Format = "json"

	fmt.Fprintf(os.Stderr, "✓ Collecting system information...\n")
	info, err := collector.Collect(dumpConfig)
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Formatting data to JSON...\n")
	output, err := formatter.Format(info, dumpConfig)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Determine output filename (next to executable)
	filename := "sysinfo_dump.json"

	fmt.Fprintf(os.Stderr, "✓ Writing to file: %s\n", filename)
	err = os.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("failed to write dump file: %w", err)
	}

	fileInfo, _ := os.Stat(filename)
	sizeKB := float64(fileInfo.Size()) / 1024.0
	sizeMB := sizeKB / 1024.0

	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "═══════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(os.Stderr, "  FULL SYSTEM DUMP COMPLETE\n")
	fmt.Fprintf(os.Stderr, "═══════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(os.Stderr, "  File: %s\n", filename)
	if sizeMB >= 1.0 {
		fmt.Fprintf(os.Stderr, "  Size: %.2f MB\n", sizeMB)
	} else {
		fmt.Fprintf(os.Stderr, "  Size: %.2f KB\n", sizeKB)
	}
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Includes:\n")
	fmt.Fprintf(os.Stderr, "    • System information\n")
	fmt.Fprintf(os.Stderr, "    • Detailed CPU data\n")
	fmt.Fprintf(os.Stderr, "    • Memory information with physical modules\n")
	fmt.Fprintf(os.Stderr, "    • Disk partitions and physical disks\n")
	fmt.Fprintf(os.Stderr, "    • Network interfaces and statistics\n")
	fmt.Fprintf(os.Stderr, "    • Process information\n")
	fmt.Fprintf(os.Stderr, "    • Comprehensive SMART data with health assessment\n")
	fmt.Fprintf(os.Stderr, "    • GPU information\n")
	fmt.Fprintf(os.Stderr, "═══════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(os.Stderr, "\n")

	// Pause for user to see results when double-clicked
	waitForEnter()

	return nil
}

// waitForEnter pauses and waits for user input when not running from a terminal
func waitForEnter() {
	// On Windows, check if we're running from explorer (no terminal attached)
	// This helps when the .exe is double-clicked
	if !isTerminal() {
		fmt.Println("\nPress Enter to exit...")
		if _, err := fmt.Scanln(); err != nil {
			fmt.Fprintf(os.Stderr, "Input error: %v\n", err)
		}
	}
}

// isTerminal checks if stdout is connected to a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	// Check if it's a character device (terminal)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
