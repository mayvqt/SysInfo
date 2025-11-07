package analyzer

import (
	"testing"
	"time"

	"github.com/mayvqt/sysinfo/internal/types"
)

func TestSMARTAnalyzer_Analyze_Healthy(t *testing.T) {
	analyzer := NewSMARTAnalyzer()

	smart := &types.SMARTInfo{
		Device:      "/dev/sda",
		Healthy:     true,
		Temperature: 45,
		DetailedAttribs: []types.SMARTAttribute{
			{
				ID:         1,
				Name:       "Raw_Read_Error_Rate",
				Value:      100,
				Worst:      100,
				Threshold:  50,
				Type:       "Pre-fail",
				WhenFailed: "Never",
			},
		},
	}

	result := analyzer.Analyze(smart)

	if result == nil {
		t.Fatal("Analyze returned nil")
	}

	if result.OverallHealth != HealthGood {
		t.Errorf("Expected HealthGood, got %s", result.OverallHealth)
	}

	if result.PredictedFailure {
		t.Error("Should not predict failure for healthy drive")
	}

	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues, got %d", len(result.Issues))
	}
}

func TestSMARTAnalyzer_AnalyzeTemperature(t *testing.T) {
	analyzer := NewSMARTAnalyzer()

	tests := []struct {
		name             string
		temperature      int
		expectedIssues   int
		expectedSeverity Severity
	}{
		{"Normal temp", 40, 0, ""},
		{"Warning temp", 65, 1, SeverityWarning},
		{"Critical temp", 75, 1, SeverityCritical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smart := &types.SMARTInfo{
				Device:          "/dev/sda",
				Temperature:     tt.temperature,
				DetailedAttribs: []types.SMARTAttribute{},
			}

			result := analyzer.Analyze(smart)

			issueCount := 0
			for _, issue := range result.Issues {
				if issue.Code == "HIGH_TEMP_WARNING" || issue.Code == "HIGH_TEMP_CRITICAL" {
					issueCount++
					if tt.expectedSeverity != "" && issue.Severity != tt.expectedSeverity {
						t.Errorf("Expected severity %s, got %s", tt.expectedSeverity, issue.Severity)
					}
				}
			}

			if issueCount != tt.expectedIssues {
				t.Errorf("Expected %d temperature issues, got %d", tt.expectedIssues, issueCount)
			}
		})
	}
}

func TestSMARTAnalyzer_AnalyzeReallocatedSectors(t *testing.T) {
	analyzer := NewSMARTAnalyzer()

	smart := &types.SMARTInfo{
		Device: "/dev/sda",
		DetailedAttribs: []types.SMARTAttribute{
			{
				ID:         5,
				Name:       "Reallocated_Sector_Ct",
				Value:      100,
				RawValue:   150, // High reallocated sectors
				Threshold:  36,
				Type:       "Pre-fail",
				WhenFailed: "Never",
			},
		},
	}

	result := analyzer.Analyze(smart)

	found := false
	for _, issue := range result.Issues {
		if issue.Code == "REALLOCATED_SECTORS" {
			found = true
			if issue.Severity != SeverityCritical {
				t.Errorf("Expected critical severity for high reallocated sectors, got %s", issue.Severity)
			}
		}
	}

	if !found {
		t.Error("Expected REALLOCATED_SECTORS issue")
	}

	if result.OverallHealth == HealthGood {
		t.Error("Expected health to be degraded with reallocated sectors")
	}
}

func TestSMARTAnalyzer_AnalyzeSSDWear(t *testing.T) {
	analyzer := NewSMARTAnalyzerWithConfig(AnalyzerConfig{
		TempWarning:      60,
		TempCritical:     70,
		WearWarning:      80.0,
		WearCritical:     90.0,
		EnablePredictive: true,
	})

	smart := &types.SMARTInfo{
		Device:       "/dev/nvme0n1",
		RotationRate: 0, // SSD
		PowerOnHours: 10000,
		DetailedAttribs: []types.SMARTAttribute{
			{
				ID:        231, // SSD Life Left
				Name:      "SSD_Life_Left",
				Value:     15, // 15% life remaining
				Worst:     15,
				Threshold: 10,
				Type:      "Pre-fail",
			},
		},
	}

	result := analyzer.Analyze(smart)

	if result.SSDWearAnalysis == nil {
		t.Fatal("Expected SSD wear analysis for SSD drive")
	}

	if result.SSDWearAnalysis.RemainingLife <= 0 {
		t.Error("Expected remaining life to be calculated")
	}

	// With 15% life left (85% used), wear status should be warning (80-90%)
	if result.SSDWearAnalysis.WearStatus != HealthWarning {
		t.Errorf("Expected wear status to be WARNING for 15%% life left (85%% used), got %s", result.SSDWearAnalysis.WearStatus)
	}
}

func TestSMARTAnalyzer_PredictiveAnalysis(t *testing.T) {
	analyzer := NewSMARTAnalyzer()

	// Create a drive with multiple critical issues
	smart := &types.SMARTInfo{
		Device:       "/dev/sda",
		Temperature:  75, // Critical temp
		RotationRate: 0,  // SSD
		PowerOnHours: 20000,
		DetailedAttribs: []types.SMARTAttribute{
			{
				ID:         5,
				Name:       "Reallocated_Sector_Ct",
				Value:      100,
				RawValue:   200, // Many reallocated sectors
				Threshold:  36,
				Type:       "Pre-fail",
				WhenFailed: "Never",
			},
			{
				ID:         197,
				Name:       "Current_Pending_Sector",
				Value:      100,
				RawValue:   10, // Pending sectors
				Threshold:  0,
				Type:       "Old_age",
				WhenFailed: "Never",
			},
			{
				ID:        231,
				Name:      "SSD_Life_Left",
				Value:     5, // Only 5% life left
				Worst:     5,
				Threshold: 10,
				Type:      "Pre-fail",
			},
		},
	}

	result := analyzer.Analyze(smart)

	if !result.PredictedFailure {
		t.Error("Expected failure prediction for drive with multiple critical issues")
	}

	if result.FailureProbability < 50.0 {
		t.Errorf("Expected high failure probability, got %.1f%%", result.FailureProbability)
	}

	if result.OverallHealth != HealthCritical {
		t.Errorf("Expected CRITICAL health status, got %s", result.OverallHealth)
	}

	if len(result.Recommendations) == 0 {
		t.Error("Expected recommendations for failing drive")
	}
}

func TestSMARTAnalyzer_FailingAttribute(t *testing.T) {
	analyzer := NewSMARTAnalyzer()

	smart := &types.SMARTInfo{
		Device: "/dev/sda",
		DetailedAttribs: []types.SMARTAttribute{
			{
				ID:         184,
				Name:       "End-to-End_Error",
				Value:      1,
				Worst:      1,
				Threshold:  97,
				Type:       "Pre-fail",
				WhenFailed: "FAILING_NOW",
			},
		},
	}

	result := analyzer.Analyze(smart)

	found := false
	for _, issue := range result.Issues {
		if issue.Code == "ATTRIBUTE_FAILING" {
			found = true
			if issue.Severity != SeverityCritical {
				t.Error("Expected critical severity for failing attribute")
			}
		}
	}

	if !found {
		t.Error("Expected ATTRIBUTE_FAILING issue")
	}

	if result.OverallHealth != HealthCritical {
		t.Error("Expected critical health for failing attribute")
	}
}

func TestSMARTAnalyzer_NilInput(t *testing.T) {
	analyzer := NewSMARTAnalyzer()

	result := analyzer.Analyze(nil)

	if result == nil {
		t.Fatal("Expected result even with nil input")
	}

	if result.OverallHealth != HealthUnknown {
		t.Errorf("Expected HealthUnknown for nil input, got %s", result.OverallHealth)
	}
}

func TestSMARTAnalyzer_SSDLifespanEstimation(t *testing.T) {
	analyzer := NewSMARTAnalyzer()

	smart := &types.SMARTInfo{
		Device:       "/dev/nvme0n1",
		RotationRate: 0,
		PowerOnHours: 10000, // 10k hours used
		DetailedAttribs: []types.SMARTAttribute{
			{
				ID:    202,
				Name:  "Percent_Lifetime_Remain",
				Value: 80, // 20% used
			},
		},
	}

	result := analyzer.Analyze(smart)

	if result.SSDWearAnalysis == nil {
		t.Fatal("Expected SSD wear analysis")
	}

	if result.SSDWearAnalysis.EstimatedLifespan == 0 {
		t.Error("Expected lifespan estimation")
	}

	// With 20% used in 10k hours, total lifespan would be 50k hours total
	// Remaining is 80% of 50k = 40k hours
	expectedHours := 40000 * time.Hour
	tolerance := 1 * time.Hour

	if result.SSDWearAnalysis.EstimatedLifespan < expectedHours-tolerance ||
		result.SSDWearAnalysis.EstimatedLifespan > expectedHours+tolerance {
		t.Errorf("Expected lifespan around %v, got %v",
			expectedHours, result.SSDWearAnalysis.EstimatedLifespan)
	}
}
