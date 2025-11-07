package analyzer

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mayvqt/sysinfo/internal/types"
)

func setupTestDB(t *testing.T) (*HistoryDB, func()) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewHistoryDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}

	return db, cleanup
}

func TestHistoryDB_RecordAnalysis(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	smart := &types.SMARTInfo{
		Device:       "/dev/sda",
		Temperature:  45,
		PowerOnHours: 1000,
		DetailedAttribs: []types.SMARTAttribute{
			{
				ID:        5,
				Name:      "Reallocated_Sector_Ct",
				Value:     100,
				Worst:     100,
				Threshold: 36,
				RawValue:  0,
			},
		},
	}

	result := &AnalysisResult{
		Device:        "/dev/sda",
		OverallHealth: HealthGood,
		Issues:        []Issue{},
	}

	err := db.RecordAnalysis(smart, result)
	if err != nil {
		t.Fatalf("Failed to record analysis: %v", err)
	}

	// Verify record was created - use a very old "since" time to ensure we get all records
	history, err := db.GetHistory("/dev/sda", time.Unix(0, 0), 10)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	t.Logf("Got %d history records", len(history))
	if len(history) != 1 {
		// Try querying all devices to debug
		allDevices, _ := db.GetDevices()
		t.Logf("Devices in DB: %v", allDevices)
		t.Errorf("Expected 1 history record, got %d", len(history))
		return
	}

	if history[0].Device != "/dev/sda" {
		t.Errorf("Expected device /dev/sda, got %s", history[0].Device)
	}

	if history[0].Temperature != 45 {
		t.Errorf("Expected temperature 45, got %d", history[0].Temperature)
	}
}

func TestHistoryDB_RecordAnalysisWithIssues(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	smart := &types.SMARTInfo{
		Device:          "/dev/sda",
		Temperature:     75,
		PowerOnHours:    1000,
		DetailedAttribs: []types.SMARTAttribute{},
	}

	result := &AnalysisResult{
		Device:        "/dev/sda",
		OverallHealth: HealthCritical,
		Issues: []Issue{
			{
				Severity:    SeverityCritical,
				Code:        "HIGH_TEMP_CRITICAL",
				Description: "Temperature critically high",
			},
			{
				Severity:    SeverityWarning,
				Code:        "HIGH_TEMP_WARNING",
				Description: "Temperature warning",
			},
		},
	}

	err := db.RecordAnalysis(smart, result)
	if err != nil {
		t.Fatalf("Failed to record analysis: %v", err)
	}

	history, err := db.GetHistory("/dev/sda", time.Unix(0, 0), 10)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(history) != 1 {
		t.Fatalf("Expected 1 history record, got %d", len(history))
	}

	if history[0].CriticalIssues != 1 {
		t.Errorf("Expected 1 critical issue, got %d", history[0].CriticalIssues)
	}

	if history[0].WarningIssues != 1 {
		t.Errorf("Expected 1 warning issue, got %d", history[0].WarningIssues)
	}
}

func TestHistoryDB_GetHistory(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Record multiple entries
	for i := 0; i < 5; i++ {
		smart := &types.SMARTInfo{
			Device:          "/dev/sda",
			Temperature:     40 + i,
			PowerOnHours:    uint64(1000 + i*100),
			DetailedAttribs: []types.SMARTAttribute{},
		}

		result := &AnalysisResult{
			Device:        "/dev/sda",
			OverallHealth: HealthGood,
			Issues:        []Issue{},
		}

		if err := db.RecordAnalysis(smart, result); err != nil {
			t.Fatalf("Failed to record analysis: %v", err)
		}

		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// Get all records
	history, err := db.GetHistory("/dev/sda", time.Unix(0, 0), 100)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(history) != 5 {
		t.Errorf("Expected 5 history records, got %d", len(history))
	}

	// Verify ordering (should be DESC)
	for i := 0; i < len(history)-1; i++ {
		if history[i].Timestamp.Before(history[i+1].Timestamp) {
			t.Error("History not ordered by timestamp DESC")
		}
	}
}

func TestHistoryDB_GetTrend(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Record entries with increasing temperature
	for i := 0; i < 5; i++ {
		smart := &types.SMARTInfo{
			Device:          "/dev/sda",
			Temperature:     40 + i*5,
			PowerOnHours:    1000,
			DetailedAttribs: []types.SMARTAttribute{},
		}

		result := &AnalysisResult{
			Device:        "/dev/sda",
			OverallHealth: HealthGood,
			Issues:        []Issue{},
		}

		if err := db.RecordAnalysis(smart, result); err != nil {
			t.Fatalf("Failed to record analysis: %v", err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	trend, err := db.GetTrend("/dev/sda", time.Unix(0, 0))
	if err != nil {
		t.Fatalf("Failed to get trend: %v", err)
	}

	if trend.Device != "/dev/sda" {
		t.Errorf("Expected device /dev/sda, got %s", trend.Device)
	}

	if trend.RecordCount != 5 {
		t.Errorf("Expected 5 records, got %d", trend.RecordCount)
	}

	if trend.MinTemperature != 40 {
		t.Errorf("Expected min temp 40, got %d", trend.MinTemperature)
	}

	if trend.MaxTemperature != 60 {
		t.Errorf("Expected max temp 60, got %d", trend.MaxTemperature)
	}

	if trend.TempTrend != "increasing" {
		t.Errorf("Expected increasing temp trend, got %s", trend.TempTrend)
	}
}

func TestHistoryDB_SSDWearTracking(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Simulate SSD wear over time
	for i := 0; i < 5; i++ {
		smart := &types.SMARTInfo{
			Device:          "/dev/nvme0n1",
			Temperature:     45,
			PowerOnHours:    uint64(1000 + i*200),
			RotationRate:    0, // SSD
			DetailedAttribs: []types.SMARTAttribute{},
		}

		result := &AnalysisResult{
			Device:        "/dev/nvme0n1",
			OverallHealth: HealthGood,
			SSDWearAnalysis: &SSDWearInfo{
				PercentUsed:   10.0 + float64(i)*2.0, // Increasing wear
				RemainingLife: 90.0 - float64(i)*2.0,
				WearStatus:    HealthGood,
			},
		}

		if err := db.RecordAnalysis(smart, result); err != nil {
			t.Fatalf("Failed to record analysis: %v", err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	trend, err := db.GetTrend("/dev/nvme0n1", time.Unix(0, 0))
	if err != nil {
		t.Fatalf("Failed to get trend: %v", err)
	}

	// Wear rate may be 0 or very small with such short time intervals in tests
	// Just check that the function runs without error
	if trend.Device != "/dev/nvme0n1" {
		t.Errorf("Expected device /dev/nvme0n1, got %s", trend.Device)
	}

	if trend.RecordCount != 5 {
		t.Errorf("Expected 5 records, got %d", trend.RecordCount)
	}
}

func TestHistoryDB_CleanOldRecords(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Record some entries
	smart := &types.SMARTInfo{
		Device:          "/dev/sda",
		Temperature:     45,
		PowerOnHours:    1000,
		DetailedAttribs: []types.SMARTAttribute{},
	}

	result := &AnalysisResult{
		Device:        "/dev/sda",
		OverallHealth: HealthGood,
		Issues:        []Issue{},
	}

	for i := 0; i < 3; i++ {
		if err := db.RecordAnalysis(smart, result); err != nil {
			t.Fatalf("Failed to record analysis: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Clean records older than 1 second (should remove all since we just created them)
	time.Sleep(1100 * time.Millisecond)
	err := db.CleanOldRecords(1 * time.Second)
	if err != nil {
		t.Fatalf("Failed to clean old records: %v", err)
	}

	// Verify records were deleted
	history, err := db.GetHistory("/dev/sda", time.Now().Add(-1*time.Hour), 100)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(history) != 0 {
		t.Errorf("Expected 0 records after cleanup, got %d", len(history))
	}
}

func TestHistoryDB_GetDevices(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Record entries for multiple devices
	devices := []string{"/dev/sda", "/dev/sdb", "/dev/nvme0n1"}

	for _, device := range devices {
		smart := &types.SMARTInfo{
			Device:          device,
			Temperature:     45,
			PowerOnHours:    1000,
			DetailedAttribs: []types.SMARTAttribute{},
		}

		result := &AnalysisResult{
			Device:        device,
			OverallHealth: HealthGood,
			Issues:        []Issue{},
		}

		if err := db.RecordAnalysis(smart, result); err != nil {
			t.Fatalf("Failed to record analysis: %v", err)
		}
	}

	found, err := db.GetDevices()
	if err != nil {
		t.Fatalf("Failed to get devices: %v", err)
	}

	if len(found) != 3 {
		t.Errorf("Expected 3 devices, got %d", len(found))
	}

	// Verify all devices are present
	deviceMap := make(map[string]bool)
	for _, d := range found {
		deviceMap[d] = true
	}

	for _, expected := range devices {
		if !deviceMap[expected] {
			t.Errorf("Device %s not found in results", expected)
		}
	}
}

func TestCalculateLinearTrend(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected string // "positive", "negative", "zero"
	}{
		{"Increasing", []float64{1, 2, 3, 4, 5}, "positive"},
		{"Decreasing", []float64{5, 4, 3, 2, 1}, "negative"},
		{"Stable", []float64{3, 3, 3, 3, 3}, "zero"},
		{"Too few values", []float64{1}, "zero"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trend := calculateLinearTrend(tt.values)

			switch tt.expected {
			case "positive":
				if trend <= 0 {
					t.Errorf("Expected positive trend, got %f", trend)
				}
			case "negative":
				if trend >= 0 {
					t.Errorf("Expected negative trend, got %f", trend)
				}
			case "zero":
				if trend != 0 {
					// Allow small floating point errors for stable values
					if trend < -0.01 || trend > 0.01 {
						t.Errorf("Expected zero trend, got %f", trend)
					}
				}
			}
		})
	}
}
