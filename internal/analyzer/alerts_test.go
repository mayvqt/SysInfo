package analyzer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mayvqt/sysinfo/internal/types"
)

func TestNewAlertManager(t *testing.T) {
	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     "https://example.com/webhook",
		WebhookTimeout: 30,
		MinLevel:       AlertWarning,
		Cooldown:       60,
	}

	manager := NewAlertManager(config)

	if manager == nil {
		t.Fatal("NewAlertManager returned nil")
	}

	if !manager.config.Enabled {
		t.Error("Expected Enabled to be true")
	}

	if manager.config.WebhookURL != config.WebhookURL {
		t.Errorf("Expected WebhookURL %s, got %s", config.WebhookURL, manager.config.WebhookURL)
	}

	if manager.config.Cooldown != config.Cooldown {
		t.Errorf("Expected Cooldown %d, got %d", config.Cooldown, manager.config.Cooldown)
	}

	// Verify internal map is initialized (check by trying to get a time)
	_, exists := manager.GetLastAlertTime("/dev/test")
	if exists {
		t.Error("Expected no alert time for non-existent device")
	}
}

func TestAlertManager_Disabled(t *testing.T) {
	config := AlertConfig{
		Enabled: false,
	}

	manager := NewAlertManager(config)

	result := &AnalysisResult{
		Device:             "/dev/sda",
		OverallHealth:      HealthCritical,
		FailureProbability: 80.0,
		PredictedFailure:   true,
	}

	// Should not send alert when disabled
	err := manager.CheckAndAlert(result)
	if err != nil {
		t.Errorf("Expected no error when disabled, got: %v", err)
	}
}

func TestAlertManager_CheckAndAlert_Critical(t *testing.T) {
	// Setup test webhook server
	var receivedAlert Alert
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&receivedAlert); err != nil {
			t.Errorf("Failed to decode webhook payload: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     server.URL,
		WebhookTimeout: 5,
		MinLevel:       AlertWarning,
		Cooldown:       0, // No cooldown for testing
	}

	manager := NewAlertManager(config)

	result := &AnalysisResult{
		Device:             "/dev/sda",
		OverallHealth:      HealthCritical,
		FailureProbability: 85.5,
		PredictedFailure:   true,
		Issues: []Issue{
			{
				Severity:    SeverityCritical,
				Description: "Predicted drive failure imminent",
			},
		},
	}

	err := manager.CheckAndAlert(result)
	if err != nil {
		t.Errorf("CheckAndAlert failed: %v", err)
	}

	// Wait a bit for webhook
	time.Sleep(100 * time.Millisecond)

	// Verify webhook was called
	if receivedAlert.Device == "" {
		t.Error("Webhook was not called or payload was empty")
	}

	if receivedAlert.Device != result.Device {
		t.Errorf("Expected device %s, got %s", result.Device, receivedAlert.Device)
	}

	if receivedAlert.Level != AlertCritical {
		t.Errorf("Expected level CRITICAL, got %s", receivedAlert.Level)
	}
}

func TestAlertManager_CheckAndAlert_Temperature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     server.URL,
		WebhookTimeout: 5,
		MinLevel:       AlertWarning,
		Cooldown:       0,
	}

	manager := NewAlertManager(config)

	// Create SMART data with high temperature
	smart := &types.SMARTInfo{
		Device:      "/dev/sda",
		DeviceModel: "Test Drive",
		Temperature: 75,
	}

	analyzer := NewSMARTAnalyzer()
	result := analyzer.Analyze(smart)

	err := manager.CheckAndAlert(result)
	if err != nil {
		t.Errorf("CheckAndAlert failed: %v", err)
	}
}

func TestAlertManager_CheckAndAlert_SSDWear(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     server.URL,
		WebhookTimeout: 5,
		MinLevel:       AlertWarning,
		Cooldown:       0,
	}

	manager := NewAlertManager(config)

	result := &AnalysisResult{
		Device:        "/dev/sda",
		OverallHealth: HealthWarning,
		SSDWearAnalysis: &SSDWearInfo{
			WearStatus:    "CRITICAL",
			RemainingLife: 5.0,
			PercentUsed:   95.0,
		},
	}

	err := manager.CheckAndAlert(result)
	if err != nil {
		t.Errorf("CheckAndAlert failed: %v", err)
	}
}

func TestAlertManager_Cooldown(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     server.URL,
		WebhookTimeout: 5,
		MinLevel:       AlertWarning,
		Cooldown:       5, // 5 minute cooldown
	}

	manager := NewAlertManager(config)

	result := &AnalysisResult{
		Device:             "/dev/sda",
		OverallHealth:      HealthCritical,
		FailureProbability: 80.0,
		PredictedFailure:   true,
	}

	// First alert should go through (may generate 2 alerts: health + prediction)
	err := manager.CheckAndAlert(result)
	if err != nil {
		t.Errorf("First CheckAndAlert failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	firstCallCount := callCount

	// Second alert should be blocked by cooldown
	err = manager.CheckAndAlert(result)
	if err != nil {
		t.Errorf("Second CheckAndAlert failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if callCount != firstCallCount {
		t.Errorf("Expected no additional webhook calls during cooldown, first: %d, second: %d", firstCallCount, callCount)
	}

	// Verify last alert time was recorded
	lastTime, exists := manager.GetLastAlertTime(result.Device)
	if !exists {
		t.Error("Expected last alert time to be recorded")
	}
	if lastTime.IsZero() {
		t.Error("Expected last alert time to be non-zero")
	}
}

func TestAlertManager_MinLevel(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     server.URL,
		WebhookTimeout: 5,
		MinLevel:       AlertCritical, // Only send critical alerts
		Cooldown:       0,
	}

	manager := NewAlertManager(config)

	// Warning level alert - should not send
	warningResult := &AnalysisResult{
		Device:        "/dev/sda",
		OverallHealth: HealthWarning,
		Issues: []Issue{
			{
				Severity:    SeverityWarning,
				Description: "Temperature slightly elevated",
			},
		},
	}

	err := manager.CheckAndAlert(warningResult)
	if err != nil {
		t.Errorf("CheckAndAlert failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if callCount > 0 {
		t.Error("Expected warning alert to be filtered out")
	}

	// Critical level alert - should send (may be 2 alerts: health + prediction)
	criticalResult := &AnalysisResult{
		Device:             "/dev/sdb",
		OverallHealth:      HealthCritical,
		PredictedFailure:   true,
		FailureProbability: 90.0,
	}

	err = manager.CheckAndAlert(criticalResult)
	if err != nil {
		t.Errorf("CheckAndAlert failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if callCount < 1 {
		t.Errorf("Expected at least one critical alert to be sent, webhook calls: %d", callCount)
	}
}

func TestAlertManager_WebhookTimeout(t *testing.T) {
	// Create a slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Delay longer than timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     server.URL,
		WebhookTimeout: 1, // 1 second timeout
		MinLevel:       AlertWarning,
		Cooldown:       0,
	}

	manager := NewAlertManager(config)

	result := &AnalysisResult{
		Device:             "/dev/sda",
		OverallHealth:      HealthCritical,
		PredictedFailure:   true,
		FailureProbability: 80.0,
	}

	// Should return error on timeout
	err := manager.CheckAndAlert(result)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestAlertManager_InvalidWebhookURL(t *testing.T) {
	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     "http://invalid-host-that-does-not-exist.local/webhook",
		WebhookTimeout: 1,
		MinLevel:       AlertWarning,
		Cooldown:       0,
	}

	manager := NewAlertManager(config)

	result := &AnalysisResult{
		Device:             "/dev/sda",
		OverallHealth:      HealthCritical,
		PredictedFailure:   true,
		FailureProbability: 80.0,
	}

	// Should return error on invalid URL
	err := manager.CheckAndAlert(result)
	if err == nil {
		t.Error("Expected network error, got nil")
	}
}

func TestAlertManager_ClearCooldown(t *testing.T) {
	config := AlertConfig{
		Enabled:  true,
		Cooldown: 60,
	}

	manager := NewAlertManager(config)

	// Verify ClearCooldown doesn't crash on non-existent device
	manager.ClearCooldown("/dev/sda")

	// Verify it's not in map
	_, exists := manager.GetLastAlertTime("/dev/sda")
	if exists {
		t.Error("Expected device to not be in lastAlerts map")
	}
}

func TestAlertManager_HealthyDrive(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := AlertConfig{
		Enabled:        true,
		WebhookURL:     server.URL,
		WebhookTimeout: 5,
		MinLevel:       AlertWarning,
		Cooldown:       0,
	}

	manager := NewAlertManager(config)

	result := &AnalysisResult{
		Device:             "/dev/sda",
		OverallHealth:      HealthGood,
		FailureProbability: 2.0,
		PredictedFailure:   false,
	}

	err := manager.CheckAndAlert(result)
	if err != nil {
		t.Errorf("CheckAndAlert failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Should not send alert for healthy drive
	if callCount > 0 {
		t.Error("Expected no alerts for healthy drive")
	}
}
