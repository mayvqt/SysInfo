package analyzer

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mayvqt/sysinfo/internal/types"
	_ "modernc.org/sqlite"
)

// HistoryDB manages SMART data history
type HistoryDB struct {
	db *sql.DB
}

// SMARTHistoryRecord represents a historical SMART reading
type SMARTHistoryRecord struct {
	ID                 int64
	Device             string
	Timestamp          time.Time
	Temperature        int
	PowerOnHours       int64
	HealthStatus       HealthStatus
	FailureProbability float64
	RemainingLife      float64
	PercentUsed        float64
	IssueCount         int
	CriticalIssues     int
	WarningIssues      int
}

// TrendData represents trend analysis over a time period
type TrendData struct {
	Device               string
	StartTime            time.Time
	EndTime              time.Time
	AvgTemperature       float64
	MaxTemperature       int
	MinTemperature       int
	TempTrend            string // "increasing", "stable", "decreasing"
	HealthTrend          string // "improving", "stable", "degrading"
	SSDWearRate          float64
	EstimatedFailureDate *time.Time
	RecordCount          int
}

// NewHistoryDB creates a new history database
func NewHistoryDB(dbPath string) (*HistoryDB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	h := &HistoryDB{db: db}
	if err := h.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return h, nil
}

// Close closes the database connection
func (h *HistoryDB) Close() error {
	if h.db != nil {
		return h.db.Close()
	}
	return nil
}

// initSchema creates the database schema
func (h *HistoryDB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS smart_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		temperature INTEGER,
		power_on_hours INTEGER,
		health_status TEXT,
		failure_probability REAL,
		remaining_life REAL,
		percent_used REAL,
		issue_count INTEGER,
		critical_issues INTEGER,
		warning_issues INTEGER,
		raw_data TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_device_timestamp ON smart_history(device, timestamp);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON smart_history(timestamp);

	CREATE TABLE IF NOT EXISTS smart_attributes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		history_id INTEGER NOT NULL,
		attribute_id INTEGER,
		attribute_name TEXT,
		value INTEGER,
		worst INTEGER,
		threshold INTEGER,
		raw_value INTEGER,
		when_failed TEXT,
		FOREIGN KEY(history_id) REFERENCES smart_history(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_history_attr ON smart_attributes(history_id, attribute_id);

	CREATE TABLE IF NOT EXISTS smart_issues (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		history_id INTEGER NOT NULL,
		severity TEXT,
		code TEXT,
		description TEXT,
		attribute_id INTEGER,
		FOREIGN KEY(history_id) REFERENCES smart_history(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_history_issues ON smart_issues(history_id);
	`

	_, err := h.db.Exec(schema)
	return err
}

// RecordAnalysis stores a SMART analysis result
func (h *HistoryDB) RecordAnalysis(smart *types.SMARTInfo, result *AnalysisResult) error {
	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Count issues by severity
	criticalCount := 0
	warningCount := 0
	for _, issue := range result.Issues {
		switch issue.Severity {
		case SeverityCritical:
			criticalCount++
		case SeverityWarning:
			warningCount++
		}
	}

	// Extract SSD data if available
	remainingLife := 0.0
	percentUsed := 0.0
	if result.SSDWearAnalysis != nil {
		remainingLife = result.SSDWearAnalysis.RemainingLife
		percentUsed = result.SSDWearAnalysis.PercentUsed
	}

	// Insert main record
	res, err := tx.Exec(`
		INSERT INTO smart_history (
			device, temperature, power_on_hours, health_status,
			failure_probability, remaining_life, percent_used,
			issue_count, critical_issues, warning_issues
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		smart.Device,
		smart.Temperature,
		smart.PowerOnHours,
		result.OverallHealth,
		result.FailureProbability,
		remainingLife,
		percentUsed,
		len(result.Issues),
		criticalCount,
		warningCount,
	)
	if err != nil {
		return fmt.Errorf("failed to insert history record: %w", err)
	}

	historyID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Insert attributes
	for _, attr := range smart.DetailedAttribs {
		_, err := tx.Exec(`
			INSERT INTO smart_attributes (
				history_id, attribute_id, attribute_name,
				value, worst, threshold, raw_value, when_failed
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			historyID,
			attr.ID,
			attr.Name,
			attr.Value,
			attr.Worst,
			attr.Threshold,
			attr.RawValue,
			attr.WhenFailed,
		)
		if err != nil {
			return fmt.Errorf("failed to insert attribute: %w", err)
		}
	}

	// Insert issues
	for _, issue := range result.Issues {
		_, err := tx.Exec(`
			INSERT INTO smart_issues (
				history_id, severity, code, description, attribute_id
			) VALUES (?, ?, ?, ?, ?)`,
			historyID,
			issue.Severity,
			issue.Code,
			issue.Description,
			issue.AttributeID,
		)
		if err != nil {
			return fmt.Errorf("failed to insert issue: %w", err)
		}
	}

	return tx.Commit()
}

// GetHistory retrieves historical records for a device
func (h *HistoryDB) GetHistory(device string, since time.Time, limit int) ([]SMARTHistoryRecord, error) {
	query := `
		SELECT id, device, timestamp, temperature, power_on_hours,
		       health_status, failure_probability, remaining_life,
		       percent_used, issue_count, critical_issues, warning_issues
		FROM smart_history
		WHERE device = ? AND timestamp >= datetime(?)
		ORDER BY timestamp DESC
		LIMIT ?`

	rows, err := h.db.Query(query, device, since.Format("2006-01-02 15:04:05"), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []SMARTHistoryRecord
	for rows.Next() {
		var r SMARTHistoryRecord
		var timestamp string
		err := rows.Scan(
			&r.ID, &r.Device, &timestamp, &r.Temperature, &r.PowerOnHours,
			&r.HealthStatus, &r.FailureProbability, &r.RemainingLife,
			&r.PercentUsed, &r.IssueCount, &r.CriticalIssues, &r.WarningIssues,
		)
		if err != nil {
			return nil, err
		}
		r.Timestamp, _ = time.Parse("2006-01-02 15:04:05", timestamp)
		records = append(records, r)
	}

	return records, rows.Err()
}

// GetTrend analyzes trends for a device over a time period
func (h *HistoryDB) GetTrend(device string, since time.Time) (*TrendData, error) {
	// Get aggregate stats
	query := `
		SELECT
			MIN(timestamp) as start_time,
			MAX(timestamp) as end_time,
			AVG(temperature) as avg_temp,
			MAX(temperature) as max_temp,
			MIN(temperature) as min_temp,
			COUNT(*) as record_count
		FROM smart_history
		WHERE device = ? AND timestamp >= ?`

	var trend TrendData
	trend.Device = device

	var startTime, endTime string
	err := h.db.QueryRow(query, device, since).Scan(
		&startTime, &endTime, &trend.AvgTemperature,
		&trend.MaxTemperature, &trend.MinTemperature,
		&trend.RecordCount,
	)
	if err != nil {
		return nil, err
	}

	trend.StartTime, _ = time.Parse("2006-01-02 15:04:05", startTime)
	trend.EndTime, _ = time.Parse("2006-01-02 15:04:05", endTime)

	// Analyze temperature trend
	tempTrend, err := h.calculateTrend(device, since, "temperature")
	if err == nil {
		trend.TempTrend = tempTrend
	}

	// Analyze health trend
	healthTrend, err := h.analyzeHealthTrend(device, since)
	if err == nil {
		trend.HealthTrend = healthTrend
	}

	// Calculate SSD wear rate if applicable
	wearRate, err := h.calculateWearRate(device, since)
	if err == nil && wearRate > 0 {
		trend.SSDWearRate = wearRate
		// Estimate failure date based on wear rate
		if wearRate > 0 {
			daysToFailure := (100.0 - trend.AvgTemperature) / (wearRate * 365.0)
			if daysToFailure > 0 && daysToFailure < 3650 { // Within 10 years
				failureDate := time.Now().Add(time.Duration(daysToFailure*24) * time.Hour)
				trend.EstimatedFailureDate = &failureDate
			}
		}
	}

	return &trend, nil
}

// calculateTrend determines if a metric is increasing, stable, or decreasing
func (h *HistoryDB) calculateTrend(device string, since time.Time, column string) (string, error) {
	query := fmt.Sprintf(`
		SELECT %s, timestamp
		FROM smart_history
		WHERE device = ? AND timestamp >= ?
		ORDER BY timestamp ASC`, column)

	rows, err := h.db.Query(query, device, since)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var values []float64
	for rows.Next() {
		var value float64
		var timestamp string
		if err := rows.Scan(&value, &timestamp); err != nil {
			continue
		}
		values = append(values, value)
	}

	if len(values) < 3 {
		return "stable", nil
	}

	// Simple linear regression to determine trend
	trend := calculateLinearTrend(values)
	if trend > 0.5 {
		return "increasing", nil
	} else if trend < -0.5 {
		return "decreasing", nil
	}
	return "stable", nil
}

// analyzeHealthTrend determines if health is improving, stable, or degrading
func (h *HistoryDB) analyzeHealthTrend(device string, since time.Time) (string, error) {
	query := `
		SELECT critical_issues, warning_issues, failure_probability
		FROM smart_history
		WHERE device = ? AND timestamp >= ?
		ORDER BY timestamp ASC`

	rows, err := h.db.Query(query, device, since)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var scores []float64
	for rows.Next() {
		var critical, warning int
		var failProb float64
		if err := rows.Scan(&critical, &warning, &failProb); err != nil {
			continue
		}
		// Health score (lower is better)
		score := float64(critical*10+warning*3) + failProb
		scores = append(scores, score)
	}

	if len(scores) < 3 {
		return "stable", nil
	}

	trend := calculateLinearTrend(scores)
	if trend > 1.0 {
		return "degrading", nil
	} else if trend < -1.0 {
		return "improving", nil
	}
	return "stable", nil
}

// calculateWearRate calculates the rate of SSD wear per day
func (h *HistoryDB) calculateWearRate(device string, since time.Time) (float64, error) {
	query := `
		SELECT percent_used, timestamp
		FROM smart_history
		WHERE device = ? AND timestamp >= ? AND percent_used > 0
		ORDER BY timestamp ASC`

	rows, err := h.db.Query(query, device, since)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var firstUsed, lastUsed float64
	var firstTime, lastTime time.Time
	count := 0

	for rows.Next() {
		var used float64
		var timestamp string
		if err := rows.Scan(&used, &timestamp); err != nil {
			continue
		}

		t, _ := time.Parse("2006-01-02 15:04:05", timestamp)
		if count == 0 {
			firstUsed = used
			firstTime = t
		}
		lastUsed = used
		lastTime = t
		count++
	}

	if count < 2 {
		return 0, nil
	}

	days := lastTime.Sub(firstTime).Hours() / 24.0
	if days == 0 {
		return 0, nil
	}

	wearPerDay := (lastUsed - firstUsed) / days
	return wearPerDay, nil
}

// calculateLinearTrend calculates a simple linear trend (positive = increasing)
func calculateLinearTrend(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}

	n := float64(len(values))
	var sumX, sumY, sumXY, sumX2 float64

	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate slope
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	return slope
}

// CleanOldRecords removes records older than the specified duration
func (h *HistoryDB) CleanOldRecords(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	_, err := h.db.Exec("DELETE FROM smart_history WHERE timestamp < ?", cutoff)
	return err
}

// GetDevices returns all devices with recorded history
func (h *HistoryDB) GetDevices() ([]string, error) {
	rows, err := h.db.Query("SELECT DISTINCT device FROM smart_history ORDER BY device")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []string
	for rows.Next() {
		var device string
		if err := rows.Scan(&device); err != nil {
			continue
		}
		devices = append(devices, device)
	}

	return devices, rows.Err()
}
