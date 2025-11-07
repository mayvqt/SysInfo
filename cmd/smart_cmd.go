package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mayvqt/sysinfo/internal/analyzer"
	"github.com/mayvqt/sysinfo/internal/collector"
	"github.com/mayvqt/sysinfo/internal/config"
	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/spf13/cobra"
)

var (
	smartPeriod string
	smartDBPath string
)

// smartCmd represents the smart command
var smartCmd = &cobra.Command{
	Use:   "smart",
	Short: "SMART disk health monitoring and analysis",
	Long: `Advanced SMART disk health monitoring with predictive failure detection,
historical tracking, and alerting capabilities.

Examples:
  sysinfo smart analyze              # Analyze all drives with failure prediction
  sysinfo smart history              # Show 7-day trend history
  sysinfo smart history --period 30d # Show 30-day trends
  sysinfo smart check                # Quick health check all drives`,
}

// smartAnalyzeCmd performs deep SMART analysis
var smartAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Perform deep SMART analysis with failure prediction",
	Long: `Analyzes SMART data from all drives and provides:
  - Health status classification
  - Predictive failure detection with probability
  - SSD wear analysis and lifespan estimation
  - Temperature monitoring
  - Detailed recommendations

Results are automatically stored in the history database.`,
	RunE: runSmartAnalyze,
}

// smartHistoryCmd shows historical SMART data
var smartHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show SMART history and trends",
	Long: `Display historical SMART data and trend analysis including:
  - Recent health records
  - Temperature trends (increasing/stable/decreasing)
  - Health status trends
  - SSD wear rates and estimated end-of-life dates`,
	RunE: runSmartHistory,
}

// smartCheckCmd performs quick health check
var smartCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Quick SMART health check",
	Long: `Performs a quick health check on all drives without storing to history.
Useful for scripts and monitoring systems.`,
	RunE: runSmartCheck,
}

func init() {
	// Add smart command to root
	rootCmd.AddCommand(smartCmd)

	// Add subcommands
	smartCmd.AddCommand(smartAnalyzeCmd)
	smartCmd.AddCommand(smartHistoryCmd)
	smartCmd.AddCommand(smartCheckCmd)

	// Shared flags for all smart subcommands
	smartCmd.PersistentFlags().StringVar(&smartDBPath, "db", "", "Custom database path (default: ~/.config/sysinfo/smart.db)")
	smartCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Verbose output")

	// History-specific flags
	smartHistoryCmd.Flags().StringVar(&smartPeriod, "period", "7d", "Time period (e.g., 1h, 24h, 7d, 30d)")

	// Analyze-specific flags
	smartAnalyzeCmd.Flags().BoolVar(&cfg.SMARTAlerts, "alerts", false, "Send webhook alerts for critical issues")
}

func runSmartAnalyze(cmd *cobra.Command, args []string) error {
	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Initializing SMART analysis...\n")
	}

	// Setup database
	db, fileConfig, err := initSMARTDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Setup analyzer
	smartAnalyzer := createAnalyzer(fileConfig)

	// Setup alerts if enabled
	var alertMgr *analyzer.AlertManager
	if cfg.SMARTAlerts {
		alertMgr = createAlertManager(fileConfig)
	}

	// Collect SMART data
	diskData, err := collectSMARTData()
	if err != nil {
		return err
	}

	if len(diskData.SMARTData) == 0 {
		fmt.Fprintf(os.Stderr, "No SMART data available. Try running with elevated privileges (sudo).\n")
		return nil
	}

	// Analyze each drive
	for _, smart := range diskData.SMARTData {
		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Analyzing %s...\n", smart.Device)
		}

		result := smartAnalyzer.Analyze(&smart)

		// Store to history
		if err := db.RecordAnalysis(&smart, result); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to record history for %s: %v\n", smart.Device, err)
		}

		// Send alerts
		if alertMgr != nil {
			if err := alertMgr.CheckAndAlert(result); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to send alert for %s: %v\n", smart.Device, err)
			}
		}

		// Display results
		displayAnalysisResult(result)
	}

	return nil
}

func runSmartHistory(cmd *cobra.Command, args []string) error {
	// Setup database
	db, _, err := initSMARTDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Parse time period
	period, err := parseDuration(smartPeriod)
	if err != nil {
		return fmt.Errorf("invalid period format: %w", err)
	}

	since := time.Now().Add(-period)

	// Get all devices
	devices, err := db.GetDevices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}

	if len(devices) == 0 {
		fmt.Println("No historical SMART data available.")
		fmt.Println("\nRun 'sysinfo smart analyze' to start collecting data.")
		return nil
	}

	// Display header
	fmt.Printf("SMART History (Last %s)\n", smartPeriod)
	fmt.Println(repeatString("=", 70))

	// Display history for each device
	for _, device := range devices {
		if err := displayDeviceHistory(db, device, since); err != nil {
			fmt.Fprintf(os.Stderr, "Error displaying history for %s: %v\n", device, err)
			continue
		}
	}

	return nil
}

func runSmartCheck(cmd *cobra.Command, args []string) error {
	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Performing SMART health check...\n")
	}

	// Collect SMART data
	diskData, err := collectSMARTData()
	if err != nil {
		return err
	}

	if len(diskData.SMARTData) == 0 {
		fmt.Fprintf(os.Stderr, "No SMART data available. Try running with elevated privileges (sudo).\n")
		return nil
	}

	// Quick analysis without storing
	smartAnalyzer := analyzer.NewSMARTAnalyzer()

	allHealthy := true
	for _, smart := range diskData.SMARTData {
		result := smartAnalyzer.Analyze(&smart)

		status := "✓"
		switch result.OverallHealth {
		case analyzer.HealthCritical, analyzer.HealthFailing:
			status = "✗"
			allHealthy = false
		case analyzer.HealthWarning:
			status = "⚠"
			allHealthy = false
		}

		fmt.Printf("%s %-20s %s", status, smart.Device, result.OverallHealth)

		if result.PredictedFailure {
			fmt.Printf("  [FAILURE PREDICTED: %.1f%%]", result.FailureProbability)
		} else if result.FailureProbability > 20 {
			fmt.Printf("  [Risk: %.1f%%]", result.FailureProbability)
		}

		if result.SSDWearAnalysis != nil && result.SSDWearAnalysis.RemainingLife < 20 {
			fmt.Printf("  [SSD Life: %.1f%%]", result.SSDWearAnalysis.RemainingLife)
		}

		fmt.Println()
	}

	if allHealthy {
		fmt.Println("\n✓ All drives healthy")
		return nil
	}

	fmt.Println("\n⚠ Issues detected - run 'sysinfo smart analyze' for details")
	return nil
}

// Helper functions

func initSMARTDatabase() (*analyzer.HistoryDB, *config.FileConfig, error) {
	// Load config file
	fileConfig, _ := config.LoadConfigFile(configFile)

	// Determine database path
	dbPath := smartDBPath
	if dbPath == "" && fileConfig != nil {
		dbPath = fileConfig.SMART.DBPath
	}
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dbPath = filepath.Join(home, ".config", "sysinfo", "smart.db")
	}

	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database
	db, err := analyzer.NewHistoryDB(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open SMART history database: %w", err)
	}

	return db, fileConfig, nil
}

func createAnalyzer(fileConfig *config.FileConfig) *analyzer.SMARTAnalyzer {
	if fileConfig != nil && fileConfig.SMART.AlertThresholds.TemperatureCritical > 0 {
		return analyzer.NewSMARTAnalyzerWithConfig(analyzer.AnalyzerConfig{
			TempWarning:      fileConfig.SMART.AlertThresholds.TemperatureWarning,
			TempCritical:     fileConfig.SMART.AlertThresholds.TemperatureCritical,
			WearWarning:      80.0,
			WearCritical:     90.0,
			EnablePredictive: true,
		})
	}
	return analyzer.NewSMARTAnalyzer()
}

func createAlertManager(fileConfig *config.FileConfig) *analyzer.AlertManager {
	webhookURL := ""
	if fileConfig != nil {
		webhookURL = fileConfig.SMART.WebhookURL
	}

	if webhookURL == "" && cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Warning: Alerts enabled but no webhook URL configured\n")
		fmt.Fprintf(os.Stderr, "Add 'webhook_url' to smart section in config file\n")
	}

	return analyzer.NewAlertManager(analyzer.AlertConfig{
		Enabled:        true,
		WebhookURL:     webhookURL,
		WebhookTimeout: 30,
		MinLevel:       analyzer.AlertWarning,
		Cooldown:       60,
	})
}

func collectSMARTData() (*types.DiskData, error) {
	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Collecting SMART data...\n")
	}

	diskData, err := collector.CollectDisk(true)
	if err != nil {
		return nil, fmt.Errorf("failed to collect disk data: %w", err)
	}

	return diskData, nil
}

func displayDeviceHistory(db *analyzer.HistoryDB, device string, since time.Time) error {
	fmt.Printf("\nDevice: %s\n", device)
	fmt.Println(repeatString("-", 70))

	// Get history records
	history, err := db.GetHistory(device, since, 100)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	if len(history) == 0 {
		fmt.Println("  No records in this period")
		return nil
	}

	// Display recent records
	fmt.Printf("  Recent Records: %d\n", len(history))
	maxRecords := 5
	if len(history) < maxRecords {
		maxRecords = len(history)
	}

	for i := 0; i < maxRecords; i++ {
		record := history[i]
		fmt.Printf("    %s | Health: %-8s | Temp: %3d°C | Issues: %d (Critical: %d)\n",
			record.Timestamp.Format("2006-01-02 15:04"),
			record.HealthStatus,
			record.Temperature,
			record.IssueCount,
			record.CriticalIssues,
		)
	}

	// Get trend analysis
	trend, err := db.GetTrend(device, since)
	if err != nil {
		return fmt.Errorf("failed to calculate trends: %w", err)
	}

	fmt.Println("\n  Trends:")
	fmt.Printf("    Temperature: %s (Avg: %.1f°C, Min: %d°C, Max: %d°C)\n",
		trend.TempTrend, trend.AvgTemperature, trend.MinTemperature, trend.MaxTemperature)
	fmt.Printf("    Health Status: %s\n", trend.HealthTrend)

	if trend.SSDWearRate > 0 {
		fmt.Printf("    SSD Wear Rate: %.4f%% per day\n", trend.SSDWearRate)
		if trend.EstimatedFailureDate != nil {
			fmt.Printf("    Estimated End of Life: %s\n",
				trend.EstimatedFailureDate.Format("2006-01-02"))
		}
	}

	return nil
}

func displayAnalysisResult(result *analyzer.AnalysisResult) {
	fmt.Printf("\n%s\n", result.Device)
	fmt.Println(repeatString("=", 70))

	// Overall health
	healthSymbol := getHealthSymbol(result.OverallHealth)
	fmt.Printf("Overall Health: %s %s\n", healthSymbol, result.OverallHealth)

	// Predictive analysis
	if result.PredictedFailure {
		fmt.Printf("⚠ PREDICTED FAILURE (%.1f%% probability)\n", result.FailureProbability)
	} else if result.FailureProbability > 20 {
		fmt.Printf("Failure Risk: %.1f%%\n", result.FailureProbability)
	}

	// SSD wear analysis
	if result.SSDWearAnalysis != nil {
		displaySSDWear(result.SSDWearAnalysis)
	}

	// Issues
	if len(result.Issues) > 0 {
		displayIssues(result.Issues)
	}

	// Recommendations
	if len(result.Recommendations) > 0 {
		displayRecommendations(result.Recommendations)
	}

	fmt.Println()
}

func getHealthSymbol(health analyzer.HealthStatus) string {
	switch health {
	case analyzer.HealthGood:
		return "✓"
	case analyzer.HealthWarning:
		return "⚠"
	case analyzer.HealthCritical, analyzer.HealthFailing:
		return "✗"
	default:
		return "?"
	}
}

func displaySSDWear(wear *analyzer.SSDWearInfo) {
	fmt.Println("\nSSD Wear Analysis:")
	fmt.Printf("  Status: %s\n", wear.WearStatus)
	fmt.Printf("  Remaining Life: %.1f%%\n", wear.RemainingLife)
	fmt.Printf("  Percent Used: %.1f%%\n", wear.PercentUsed)
	if wear.EstimatedLifespan > 0 {
		days := int(wear.EstimatedLifespan.Hours() / 24)
		years := float64(days) / 365.0
		fmt.Printf("  Estimated Remaining: %d days (%.1f years)\n", days, years)
	}
}

func displayIssues(issues []analyzer.Issue) {
	fmt.Printf("\nIssues Found: %d\n", len(issues))
	for _, issue := range issues {
		var severity string
		switch issue.Severity {
		case analyzer.SeverityCritical:
			severity = "CRITICAL"
		case analyzer.SeverityWarning:
			severity = "WARNING"
		default:
			severity = "INFO"
		}
		fmt.Printf("  [%s] %s\n", severity, issue.Description)
	}
}

func displayRecommendations(recommendations []string) {
	fmt.Println("\nRecommendations:")
	for _, rec := range recommendations {
		fmt.Printf("  • %s\n", rec)
	}
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	unit := s[len(s)-1]
	value := s[:len(s)-1]

	var multiplier time.Duration
	switch unit {
	case 'h':
		multiplier = time.Hour
	case 'd':
		multiplier = 24 * time.Hour
	case 'w':
		multiplier = 7 * 24 * time.Hour
	case 'm':
		multiplier = 30 * 24 * time.Hour
	default:
		return time.ParseDuration(s)
	}

	var num int
	if _, err := fmt.Sscanf(value, "%d", &num); err != nil {
		return 0, fmt.Errorf("invalid number in duration: %w", err)
	}

	return time.Duration(num) * multiplier, nil
}

func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
