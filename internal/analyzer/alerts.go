package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AlertLevel defines the severity level for alerts
type AlertLevel string

const (
	AlertInfo     AlertLevel = "INFO"
	AlertWarning  AlertLevel = "WARNING"
	AlertCritical AlertLevel = "CRITICAL"
)

// Alert represents a disk health alert
type Alert struct {
	Level       AlertLevel             `json:"level"`
	Device      string                 `json:"device"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// AlertConfig configures the alert system
type AlertConfig struct {
	Enabled        bool       `json:"enabled"`
	WebhookURL     string     `json:"webhook_url,omitempty"`
	WebhookTimeout int        `json:"webhook_timeout"` // seconds
	MinLevel       AlertLevel `json:"min_level"`
	Cooldown       int        `json:"cooldown"` // minutes between alerts for same device
}

// AlertManager manages disk health alerts
type AlertManager struct {
	config     AlertConfig
	lastAlerts map[string]time.Time // device -> last alert time
	client     *http.Client
}

// NewAlertManager creates a new alert manager
func NewAlertManager(config AlertConfig) *AlertManager {
	if config.WebhookTimeout == 0 {
		config.WebhookTimeout = 30
	}
	if config.Cooldown == 0 {
		config.Cooldown = 60 // Default 1 hour cooldown
	}
	if config.MinLevel == "" {
		config.MinLevel = AlertWarning
	}

	return &AlertManager{
		config:     config,
		lastAlerts: make(map[string]time.Time),
		client: &http.Client{
			Timeout: time.Duration(config.WebhookTimeout) * time.Second,
		},
	}
}

// CheckAndAlert analyzes a SMART result and sends alerts if necessary
func (am *AlertManager) CheckAndAlert(result *AnalysisResult) error {
	if !am.config.Enabled {
		return nil
	}

	alerts := am.generateAlerts(result)
	for _, alert := range alerts {
		if err := am.sendAlert(alert); err != nil {
			return fmt.Errorf("failed to send alert: %w", err)
		}
	}

	return nil
}

// generateAlerts creates alerts based on analysis results
func (am *AlertManager) generateAlerts(result *AnalysisResult) []Alert {
	var alerts []Alert

	// Check if we're in cooldown period
	if lastAlert, exists := am.lastAlerts[result.Device]; exists {
		if time.Since(lastAlert) < time.Duration(am.config.Cooldown)*time.Minute {
			return alerts // Still in cooldown, skip alerts
		}
	}

	// Overall health status alert
	if result.OverallHealth == HealthCritical || result.OverallHealth == HealthFailing {
		alerts = append(alerts, Alert{
			Level:       AlertCritical,
			Device:      result.Device,
			Title:       fmt.Sprintf("Critical Disk Health: %s", result.Device),
			Description: fmt.Sprintf("Disk health is %s", result.OverallHealth),
			Timestamp:   time.Now(),
			Data: map[string]interface{}{
				"health_status":       result.OverallHealth,
				"failure_probability": result.FailureProbability,
				"predicted_failure":   result.PredictedFailure,
				"issue_count":         len(result.Issues),
			},
		})
	} else if result.OverallHealth == HealthWarning {
		if am.shouldSendAlert(AlertWarning) {
			alerts = append(alerts, Alert{
				Level:       AlertWarning,
				Device:      result.Device,
				Title:       fmt.Sprintf("Disk Health Warning: %s", result.Device),
				Description: fmt.Sprintf("Disk health degraded to %s", result.OverallHealth),
				Timestamp:   time.Now(),
				Data: map[string]interface{}{
					"health_status": result.OverallHealth,
					"issue_count":   len(result.Issues),
				},
			})
		}
	}

	// Predictive failure alert
	if result.PredictedFailure {
		alerts = append(alerts, Alert{
			Level:  AlertCritical,
			Device: result.Device,
			Title:  fmt.Sprintf("Predicted Disk Failure: %s", result.Device),
			Description: fmt.Sprintf("Drive is predicted to fail with %.1f%% probability",
				result.FailureProbability),
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"failure_probability": result.FailureProbability,
				"recommendations":     result.Recommendations,
			},
		})
	}

	// SSD wear alert
	if result.SSDWearAnalysis != nil {
		if result.SSDWearAnalysis.WearStatus == HealthCritical {
			alerts = append(alerts, Alert{
				Level:  AlertCritical,
				Device: result.Device,
				Title:  fmt.Sprintf("Critical SSD Wear: %s", result.Device),
				Description: fmt.Sprintf("SSD has %.1f%% life remaining",
					result.SSDWearAnalysis.RemainingLife),
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"remaining_life":     result.SSDWearAnalysis.RemainingLife,
					"percent_used":       result.SSDWearAnalysis.PercentUsed,
					"estimated_lifespan": result.SSDWearAnalysis.EstimatedLifespan.String(),
				},
			})
		} else if result.SSDWearAnalysis.WearStatus == HealthWarning && am.shouldSendAlert(AlertWarning) {
			alerts = append(alerts, Alert{
				Level:  AlertWarning,
				Device: result.Device,
				Title:  fmt.Sprintf("SSD Wear Warning: %s", result.Device),
				Description: fmt.Sprintf("SSD has %.1f%% life remaining",
					result.SSDWearAnalysis.RemainingLife),
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"remaining_life": result.SSDWearAnalysis.RemainingLife,
					"percent_used":   result.SSDWearAnalysis.PercentUsed,
				},
			})
		}
	}

	// Individual issue alerts (only for critical issues)
	for _, issue := range result.Issues {
		if issue.Severity == SeverityCritical {
			alerts = append(alerts, Alert{
				Level:       AlertCritical,
				Device:      result.Device,
				Title:       fmt.Sprintf("Disk Issue: %s", issue.Code),
				Description: issue.Description,
				Timestamp:   time.Now(),
				Data: map[string]interface{}{
					"issue_code":   issue.Code,
					"severity":     issue.Severity,
					"attribute_id": issue.AttributeID,
					"value":        issue.Value,
				},
			})
		}
	}

	// Update last alert time if we're sending any alerts
	if len(alerts) > 0 {
		am.lastAlerts[result.Device] = time.Now()
	}

	return alerts
}

// shouldSendAlert checks if an alert level should be sent
func (am *AlertManager) shouldSendAlert(level AlertLevel) bool {
	minLevel := am.config.MinLevel

	// Define level hierarchy
	levels := map[AlertLevel]int{
		AlertInfo:     1,
		AlertWarning:  2,
		AlertCritical: 3,
	}

	return levels[level] >= levels[minLevel]
}

// sendAlert sends an alert via configured channels
func (am *AlertManager) sendAlert(alert Alert) error {
	// Send to webhook if configured
	if am.config.WebhookURL != "" {
		if err := am.sendWebhook(alert); err != nil {
			return err
		}
	}

	// Could add other notification methods here (email, Slack, PagerDuty, etc.)

	return nil
}

// sendWebhook sends an alert to a webhook URL
func (am *AlertManager) sendWebhook(alert Alert) error {
	payload, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}

	resp, err := am.client.Post(am.config.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// ClearCooldown clears the cooldown for a specific device
func (am *AlertManager) ClearCooldown(device string) {
	delete(am.lastAlerts, device)
}

// GetLastAlertTime returns the last alert time for a device
func (am *AlertManager) GetLastAlertTime(device string) (time.Time, bool) {
	t, exists := am.lastAlerts[device]
	return t, exists
}
