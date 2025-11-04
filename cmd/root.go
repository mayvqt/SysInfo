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

	// Output options
	rootCmd.Flags().StringVarP(&cfg.Format, "format", "f", "pretty", "Output format: json, text, pretty")
	rootCmd.Flags().StringVarP(&cfg.OutputFile, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Verbose output")

	// Module selection flags
	rootCmd.Flags().BoolVar(&cfg.Modules.All, "all", true, "Collect all information")
	rootCmd.Flags().BoolVar(&cfg.Modules.System, "system", false, "Collect system information")
	rootCmd.Flags().BoolVar(&cfg.Modules.CPU, "cpu", false, "Collect CPU information")
	rootCmd.Flags().BoolVar(&cfg.Modules.Memory, "memory", false, "Collect memory information")
	rootCmd.Flags().BoolVar(&cfg.Modules.Disk, "disk", false, "Collect disk information")
	rootCmd.Flags().BoolVar(&cfg.Modules.Network, "network", false, "Collect network information")
	rootCmd.Flags().BoolVar(&cfg.Modules.Process, "process", false, "Collect process information")
	rootCmd.Flags().BoolVar(&cfg.Modules.SMART, "smart", false, "Collect SMART disk data (may require elevated privileges)")
}

func Execute() error {
	return rootCmd.Execute()
}

func runSysInfo(cmd *cobra.Command, args []string) error {
	// If any specific module is selected, disable --all
	if cfg.Modules.System || cfg.Modules.CPU || cfg.Modules.Memory ||
		cfg.Modules.Disk || cfg.Modules.Network || cfg.Modules.Process || cfg.Modules.SMART {
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
