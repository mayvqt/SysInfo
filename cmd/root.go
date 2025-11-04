package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Monitor mode options
	rootCmd.Flags().BoolVarP(&cfg.Monitor, "monitor", "m", false, "Enable live monitoring mode (continuously update)")
	rootCmd.Flags().IntVarP(&cfg.MonitorInterval, "interval", "i", 2, "Update interval in seconds for monitor mode")

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

	// Validate monitor mode
	if cfg.Monitor {
		if cfg.OutputFile != "" {
			return fmt.Errorf("monitor mode cannot be used with file output")
		}
		if cfg.MonitorInterval < 1 {
			cfg.MonitorInterval = 1
		}
		return runMonitorMode()
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

// runMonitorMode continuously updates the output at the specified interval
func runMonitorMode() error {
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create a ticker for periodic updates
	ticker := time.NewTicker(time.Duration(cfg.MonitorInterval) * time.Second)
	defer ticker.Stop()

	// Hide cursor for cleaner display
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit

	fmt.Fprintf(os.Stderr, "Live monitoring mode - Press Ctrl+C to exit\n")
	fmt.Fprintf(os.Stderr, "Update interval: %d second(s)\n\n", cfg.MonitorInterval)
	time.Sleep(500 * time.Millisecond) // Brief pause so user can see the message

	// Display initial data immediately
	if err := displayLiveData(true); err != nil {
		fmt.Print("\033[?25h") // Show cursor on error
		return err
	}

	// Main monitor loop
	for {
		select {
		case <-ticker.C:
			if err := displayLiveData(false); err != nil {
				fmt.Print("\033[?25h") // Show cursor on error
				return err
			}
		case <-sigChan:
			fmt.Print("\033[?25h") // Show cursor
			fmt.Fprintf(os.Stderr, "\n\nMonitoring stopped.\n")
			return nil
		}
	}
}

// displayLiveData collects and displays current system information
func displayLiveData(isFirstUpdate bool) error {
	// Collect system information (this might take time)
	info, err := collector.Collect(cfg)
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	// Format output (buffer everything first)
	output, err := formatter.Format(info, cfg)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Now that we have all the data, update the display atomically
	if !isFirstUpdate {
		// Clear the screen completely and move cursor to home
		// Using the full sequence for better Windows PowerShell compatibility
		fmt.Print("\033[2J\033[H")
	}

	// Print the complete output all at once
	fmt.Print(output)

	// Ensure output is flushed immediately
	os.Stdout.Sync()

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
