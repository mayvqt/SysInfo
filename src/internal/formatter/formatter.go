package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/mayvqt/sysinfo/src/internal/config"
	"github.com/mayvqt/sysinfo/src/internal/types"
)

// Format formats the system information according to the specified format
func Format(info *types.SystemInfo, cfg *config.Config) (string, error) {
	switch cfg.Format {
	case "json":
		return FormatJSON(info)
	case "text":
		return FormatText(info), nil
	case "pretty":
		return FormatPretty(info), nil
	default:
		return "", fmt.Errorf("unknown format: %s", cfg.Format)
	}
}

// FormatJSON formats the information as JSON
func FormatJSON(info *types.SystemInfo) (string, error) {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}
