//go:build linux
// +build linux

package collector

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
)

// DmidecodeMemory represents dmidecode memory device output
type DmidecodeMemory struct {
	Handle          string `json:"handle"`
	Type            string `json:"type"`
	Size            string `json:"size"`
	FormFactor      string `json:"form_factor"`
	Locator         string `json:"locator"`
	BankLocator     string `json:"bank_locator"`
	MemoryType      string `json:"type_detail"`
	Speed           string `json:"speed"`
	Manufacturer    string `json:"manufacturer"`
	SerialNumber    string `json:"serial_number"`
	PartNumber      string `json:"part_number"`
	ConfiguredSpeed string `json:"configured_memory_speed"`
}

// collectMemoryModulesPlatform implements Linux-specific memory module collection using dmidecode
func collectMemoryModulesPlatform() []types.MemoryModule {
	modules := make([]types.MemoryModule, 0)

	// Try dmidecode first (requires root)
	cmd := exec.Command("dmidecode", "-t", "memory", "-q")
	output, err := cmd.Output()
	if err != nil {
		// Fall back to trying without -q flag
		cmd = exec.Command("dmidecode", "-t", "17")
		output, err = cmd.Output()
		if err != nil {
			return modules
		}
	}

	// Parse dmidecode output
	return parseDmidecodeOutput(string(output))
}

// parseDmidecodeOutput parses dmidecode text output into MemoryModule structs
func parseDmidecodeOutput(output string) []types.MemoryModule {
	modules := make([]types.MemoryModule, 0)

	lines := strings.Split(output, "\n")
	var currentModule *types.MemoryModule

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// New memory device section
		if strings.HasPrefix(line, "Memory Device") {
			if currentModule != nil && currentModule.Capacity > 0 {
				modules = append(modules, *currentModule)
			}
			currentModule = &types.MemoryModule{}
			continue
		}

		if currentModule == nil {
			continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Size":
			currentModule.Capacity = parseMemorySize(value)
		case "Locator":
			currentModule.Locator = value
		case "Bank Locator":
			if currentModule.Locator == "" {
				currentModule.Locator = value
			}
		case "Type":
			if value != "Unknown" && value != "<OUT OF SPEC>" {
				currentModule.Type = value
			}
		case "Speed":
			currentModule.Speed = parseMemorySpeed(value)
		case "Configured Memory Speed", "Configured Clock Speed":
			if speed := parseMemorySpeed(value); speed > 0 {
				currentModule.Speed = speed
			}
		case "Manufacturer":
			if value != "Unknown" && value != "NO DIMM" && value != "" {
				currentModule.Manufacturer = value
			}
		case "Serial Number":
			if value != "Unknown" && value != "NO DIMM" && value != "" {
				currentModule.SerialNumber = value
			}
		case "Part Number":
			if value != "Unknown" && value != "NO DIMM" && value != "" {
				currentModule.PartNumber = value
			}
		case "Form Factor":
			if value != "Unknown" {
				currentModule.FormFactor = value
			}
		}
	}

	// Add last module if valid
	if currentModule != nil && currentModule.Capacity > 0 {
		modules = append(modules, *currentModule)
	}

	return modules
}

// parseMemorySize converts size strings like "8192 MB" or "8 GB" to bytes
func parseMemorySize(sizeStr string) uint64 {
	if sizeStr == "No Module Installed" || sizeStr == "Unknown" {
		return 0
	}

	parts := strings.Fields(sizeStr)
	if len(parts) < 2 {
		return 0
	}

	size, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	unit := strings.ToUpper(parts[1])
	switch unit {
	case "KB":
		return uint64(size * 1024)
	case "MB":
		return uint64(size * 1024 * 1024)
	case "GB":
		return uint64(size * 1024 * 1024 * 1024)
	case "TB":
		return uint64(size * 1024 * 1024 * 1024 * 1024)
	default:
		return uint64(size)
	}
}

// parseMemorySpeed converts speed strings like "2400 MT/s" or "2400 MHz" to numeric value
func parseMemorySpeed(speedStr string) uint64 {
	if speedStr == "Unknown" || speedStr == "" {
		return 0
	}

	parts := strings.Fields(speedStr)
	if len(parts) < 1 {
		return 0
	}

	speed, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0
	}

	return speed
}
