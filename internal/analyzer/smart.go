package analyzer

import (
	"fmt"
	"time"

	"github.com/mayvqt/sysinfo/internal/types"
)

// SMARTAnalyzer analyzes SMART data for predictive failure detection
type SMARTAnalyzer struct {
	config AnalyzerConfig
}

// AnalyzerConfig contains configuration for SMART analysis
type AnalyzerConfig struct {
	// Temperature thresholds in Celsius
	TempWarning  int
	TempCritical int

	// SSD wear thresholds (percentage)
	WearWarning  float64
	WearCritical float64

	// Enable predictive analysis
	EnablePredictive bool
}

// NewSMARTAnalyzer creates a new SMART analyzer with default config
func NewSMARTAnalyzer() *SMARTAnalyzer {
	return &SMARTAnalyzer{
		config: AnalyzerConfig{
			TempWarning:      60,
			TempCritical:     70,
			WearWarning:      80.0,
			WearCritical:     90.0,
			EnablePredictive: true,
		},
	}
}

// NewSMARTAnalyzerWithConfig creates a new SMART analyzer with custom config
func NewSMARTAnalyzerWithConfig(config AnalyzerConfig) *SMARTAnalyzer {
	return &SMARTAnalyzer{config: config}
}

// AnalysisResult contains the results of SMART analysis
type AnalysisResult struct {
	Device             string
	OverallHealth      HealthStatus
	PredictedFailure   bool
	FailureProbability float64 // 0-100%
	TimeToFailure      *time.Duration
	Issues             []Issue
	Recommendations    []string
	SSDWearAnalysis    *SSDWearInfo
}

// HealthStatus represents the health status of a drive
type HealthStatus string

const (
	HealthGood     HealthStatus = "GOOD"
	HealthWarning  HealthStatus = "WARNING"
	HealthCritical HealthStatus = "CRITICAL"
	HealthFailing  HealthStatus = "FAILING"
	HealthUnknown  HealthStatus = "UNKNOWN"
)

// Issue represents a specific SMART issue
type Issue struct {
	Severity    Severity
	Code        string
	Description string
	AttributeID uint8
	Value       string
}

// Severity levels for issues
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

// SSDWearInfo contains SSD-specific wear analysis
type SSDWearInfo struct {
	WearLevelingCount uint64
	ProgramEraseCount uint64
	PercentUsed       float64
	EstimatedLifespan time.Duration
	RemainingLife     float64 // 0-100%
	WearStatus        HealthStatus
}

// Analyze performs comprehensive SMART analysis
func (a *SMARTAnalyzer) Analyze(smart *types.SMARTInfo) *AnalysisResult {
	if smart == nil {
		return &AnalysisResult{
			OverallHealth: HealthUnknown,
			Issues:        []Issue{},
		}
	}

	result := &AnalysisResult{
		Device:          smart.Device,
		Issues:          []Issue{},
		Recommendations: []string{},
	}

	// Check temperature
	a.analyzeTemperature(smart, result)

	// Check for failing attributes
	a.analyzeAttributes(smart, result)

	// Check reallocated sectors
	a.analyzeReallocatedSectors(smart, result)

	// Analyze SSD-specific metrics if applicable
	if smart.RotationRate == 0 {
		result.SSDWearAnalysis = a.analyzeSSDWear(smart)
	}

	// Predictive failure analysis
	if a.config.EnablePredictive {
		a.predictiveAnalysis(smart, result)
	}

	// Determine overall health
	a.determineOverallHealth(result)

	// Generate recommendations
	a.generateRecommendations(result)

	return result
}

// analyzeTemperature checks drive temperature
func (a *SMARTAnalyzer) analyzeTemperature(smart *types.SMARTInfo, result *AnalysisResult) {
	if smart.Temperature <= 0 {
		return
	}

	if smart.Temperature >= a.config.TempCritical {
		result.Issues = append(result.Issues, Issue{
			Severity:    SeverityCritical,
			Code:        "HIGH_TEMP_CRITICAL",
			Description: fmt.Sprintf("Drive temperature is critically high: %d°C (threshold: %d°C)", smart.Temperature, a.config.TempCritical),
			Value:       fmt.Sprintf("%d°C", smart.Temperature),
		})
	} else if smart.Temperature >= a.config.TempWarning {
		result.Issues = append(result.Issues, Issue{
			Severity:    SeverityWarning,
			Code:        "HIGH_TEMP_WARNING",
			Description: fmt.Sprintf("Drive temperature is elevated: %d°C (threshold: %d°C)", smart.Temperature, a.config.TempWarning),
			Value:       fmt.Sprintf("%d°C", smart.Temperature),
		})
	}
}

// analyzeAttributes checks SMART attributes for failures
func (a *SMARTAnalyzer) analyzeAttributes(smart *types.SMARTInfo, result *AnalysisResult) {
	for _, attr := range smart.DetailedAttribs {
		// Check if attribute has failed
		switch attr.WhenFailed {
		case "FAILING_NOW":
			result.Issues = append(result.Issues, Issue{
				Severity:    SeverityCritical,
				Code:        "ATTRIBUTE_FAILING",
				Description: fmt.Sprintf("SMART attribute %d (%s) is failing NOW", attr.ID, attr.Name),
				AttributeID: attr.ID,
				Value:       fmt.Sprintf("%d (threshold: %d)", attr.Value, attr.Threshold),
			})
		case "In_the_past":
			result.Issues = append(result.Issues, Issue{
				Severity:    SeverityWarning,
				Code:        "ATTRIBUTE_FAILED_PAST",
				Description: fmt.Sprintf("SMART attribute %d (%s) failed in the past", attr.ID, attr.Name),
				AttributeID: attr.ID,
				Value:       fmt.Sprintf("%d (threshold: %d)", attr.Value, attr.Threshold),
			})
		}

		// Check if value is close to threshold (Pre-fail attributes only)
		if attr.Type == "Pre-fail" && attr.Threshold > 0 {
			margin := int(attr.Value) - int(attr.Threshold)
			if margin <= 10 && margin > 0 {
				result.Issues = append(result.Issues, Issue{
					Severity:    SeverityWarning,
					Code:        "ATTRIBUTE_NEAR_THRESHOLD",
					Description: fmt.Sprintf("SMART attribute %d (%s) is approaching failure threshold", attr.ID, attr.Name),
					AttributeID: attr.ID,
					Value:       fmt.Sprintf("%d (threshold: %d, margin: %d)", attr.Value, attr.Threshold, margin),
				})
			}
		}
	}
}

// analyzeReallocatedSectors checks for reallocated sectors
func (a *SMARTAnalyzer) analyzeReallocatedSectors(smart *types.SMARTInfo, result *AnalysisResult) {
	// Check for reallocated sector count (ID 5)
	for _, attr := range smart.DetailedAttribs {
		if attr.ID == 5 && attr.RawValue > 0 {
			severity := SeverityWarning
			if attr.RawValue > 100 {
				severity = SeverityCritical
			}
			result.Issues = append(result.Issues, Issue{
				Severity:    severity,
				Code:        "REALLOCATED_SECTORS",
				Description: fmt.Sprintf("Drive has %d reallocated sectors", attr.RawValue),
				AttributeID: 5,
				Value:       fmt.Sprintf("%d", attr.RawValue),
			})
		}

		// Pending sectors (ID 197)
		if attr.ID == 197 && attr.RawValue > 0 {
			result.Issues = append(result.Issues, Issue{
				Severity:    SeverityCritical,
				Code:        "PENDING_SECTORS",
				Description: fmt.Sprintf("Drive has %d pending sectors (unstable)", attr.RawValue),
				AttributeID: 197,
				Value:       fmt.Sprintf("%d", attr.RawValue),
			})
		}

		// Uncorrectable sectors (ID 198)
		if attr.ID == 198 && attr.RawValue > 0 {
			result.Issues = append(result.Issues, Issue{
				Severity:    SeverityCritical,
				Code:        "UNCORRECTABLE_SECTORS",
				Description: fmt.Sprintf("Drive has %d uncorrectable sectors", attr.RawValue),
				AttributeID: 198,
				Value:       fmt.Sprintf("%d", attr.RawValue),
			})
		}
	}
}

// analyzeSSDWear analyzes SSD-specific wear metrics
func (a *SMARTAnalyzer) analyzeSSDWear(smart *types.SMARTInfo) *SSDWearInfo {
	wear := &SSDWearInfo{
		RemainingLife: 100.0,
		WearStatus:    HealthGood,
	}

	for _, attr := range smart.DetailedAttribs {
		switch attr.ID {
		case 177: // Wear Leveling Count
			wear.WearLevelingCount = attr.RawValue
		case 231, 233: // SSD Life Left / Media Wearout Indicator
			if attr.Value > 0 {
				wear.RemainingLife = float64(attr.Value)
				wear.PercentUsed = 100.0 - wear.RemainingLife
			}
		case 202, 226: // Percent Lifetime Used / Workload Media Wear Indicator
			wear.PercentUsed = float64(100 - attr.Value)
			wear.RemainingLife = 100.0 - wear.PercentUsed
		case 12: // Power Cycle Count (used for lifespan estimation)
			wear.ProgramEraseCount = attr.RawValue
		}
	}

	// Use health assessment if available
	if smart.HealthAssessment != nil && smart.HealthAssessment.PercentUsed > 0 {
		wear.PercentUsed = smart.HealthAssessment.PercentUsed
		wear.RemainingLife = 100.0 - wear.PercentUsed
	}

	// Estimate lifespan based on wear
	if wear.PercentUsed > 0 && smart.PowerOnHours > 0 {
		hoursPerPercent := float64(smart.PowerOnHours) / wear.PercentUsed
		remainingHours := hoursPerPercent * wear.RemainingLife
		wear.EstimatedLifespan = time.Duration(remainingHours) * time.Hour
	}

	// Determine wear status
	if wear.PercentUsed >= a.config.WearCritical {
		wear.WearStatus = HealthCritical
	} else if wear.PercentUsed >= a.config.WearWarning {
		wear.WearStatus = HealthWarning
	}

	return wear
}

// predictiveAnalysis performs predictive failure analysis
func (a *SMARTAnalyzer) predictiveAnalysis(smart *types.SMARTInfo, result *AnalysisResult) {
	failureScore := 0.0

	// Count critical issues
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

	// Calculate failure probability based on issues
	failureScore += float64(criticalCount) * 30.0
	failureScore += float64(warningCount) * 10.0

	// Check SSD wear if available
	if result.SSDWearAnalysis != nil {
		if result.SSDWearAnalysis.PercentUsed >= 95 {
			failureScore += 40.0
		} else if result.SSDWearAnalysis.PercentUsed >= 90 {
			failureScore += 25.0
		} else if result.SSDWearAnalysis.PercentUsed >= 80 {
			failureScore += 15.0
		}
	}

	// Check for high reallocated sectors
	for _, attr := range smart.DetailedAttribs {
		if attr.ID == 5 && attr.RawValue > 50 {
			failureScore += 20.0
		}
	}

	// Cap at 100%
	if failureScore > 100 {
		failureScore = 100
	}

	result.FailureProbability = failureScore
	result.PredictedFailure = failureScore >= 50.0

	// Estimate time to failure based on wear rate
	if result.SSDWearAnalysis != nil && result.SSDWearAnalysis.EstimatedLifespan > 0 {
		result.TimeToFailure = &result.SSDWearAnalysis.EstimatedLifespan
	}
}

// determineOverallHealth determines the overall health status
func (a *SMARTAnalyzer) determineOverallHealth(result *AnalysisResult) {
	// Default to good
	result.OverallHealth = HealthGood

	// Check for any critical issues
	hasCritical := false
	hasWarning := false
	for _, issue := range result.Issues {
		switch issue.Severity {
		case SeverityCritical:
			hasCritical = true
		case SeverityWarning:
			hasWarning = true
		}
	}

	if hasCritical || result.PredictedFailure {
		result.OverallHealth = HealthCritical
	} else if hasWarning {
		result.OverallHealth = HealthWarning
	}

	// Override if SSD wear is critical
	if result.SSDWearAnalysis != nil {
		if result.SSDWearAnalysis.WearStatus == HealthCritical {
			result.OverallHealth = HealthCritical
		} else if result.SSDWearAnalysis.WearStatus == HealthWarning && result.OverallHealth == HealthGood {
			result.OverallHealth = HealthWarning
		}
	}
}

// generateRecommendations generates actionable recommendations
func (a *SMARTAnalyzer) generateRecommendations(result *AnalysisResult) {
	if result.OverallHealth == HealthCritical {
		result.Recommendations = append(result.Recommendations, "URGENT: Back up all data immediately")
		result.Recommendations = append(result.Recommendations, "Schedule drive replacement as soon as possible")
	}

	if result.PredictedFailure {
		result.Recommendations = append(result.Recommendations, "Drive failure is predicted - plan for replacement")
	}

	// Temperature recommendations
	for _, issue := range result.Issues {
		if issue.Code == "HIGH_TEMP_CRITICAL" || issue.Code == "HIGH_TEMP_WARNING" {
			result.Recommendations = append(result.Recommendations, "Improve cooling/ventilation around the drive")
			break
		}
	}

	// SSD-specific recommendations
	if result.SSDWearAnalysis != nil && result.SSDWearAnalysis.PercentUsed >= 80 {
		result.Recommendations = append(result.Recommendations, fmt.Sprintf("SSD is %.1f%% worn - consider replacement soon", result.SSDWearAnalysis.PercentUsed))
	}

	// Reallocated sector recommendations
	for _, issue := range result.Issues {
		if issue.Code == "REALLOCATED_SECTORS" || issue.Code == "PENDING_SECTORS" {
			result.Recommendations = append(result.Recommendations, "Run a full surface scan and consider drive replacement")
			break
		}
	}

	if len(result.Recommendations) == 0 && result.OverallHealth == HealthGood {
		result.Recommendations = append(result.Recommendations, "Drive health is good - continue monitoring")
	}
}
