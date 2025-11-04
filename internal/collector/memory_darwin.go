//go:build darwin
// +build darwin

package collector

import (
	"encoding/xml"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
)

// Plist structures for system_profiler XML output
type PlistData struct {
	Array ArrayData `xml:"array"`
}

type ArrayData struct {
	Dicts []DictData `xml:"dict"`
}

type DictData struct {
	Keys   []string    `xml:"key"`
	Values []string    `xml:"string"`
	Arrays []ArrayData `xml:"array"`
}

// MemorySlot represents a single memory slot from system_profiler
type MemorySlot struct {
	Size         string
	Type         string
	Speed        string
	Status       string
	Manufacturer string
	PartNumber   string
	SerialNumber string
}

// collectMemoryModulesPlatform implements macOS-specific memory module collection using system_profiler
func collectMemoryModulesPlatform() []types.MemoryModule {
	modules := make([]types.MemoryModule, 0)

	// Run system_profiler to get memory information in XML format
	cmd := exec.Command("system_profiler", "SPMemoryDataType", "-xml")
	output, err := cmd.Output()
	if err != nil {
		return modules
	}

	// Parse the XML output
	slots := parseSystemProfilerMemory(output)

	// Convert to MemoryModule format
	for i, slot := range slots {
		// Skip empty slots
		if slot.Size == "" || slot.Size == "empty" || slot.Status == "Empty" {
			continue
		}

		module := types.MemoryModule{
			Locator:      "DIMM" + strconv.Itoa(i),
			Capacity:     parseMemorySizeMac(slot.Size),
			Speed:        parseMemorySpeedMac(slot.Speed),
			Type:         slot.Type,
			Manufacturer: slot.Manufacturer,
			PartNumber:   slot.PartNumber,
			SerialNumber: slot.SerialNumber,
			FormFactor:   "DIMM", // macOS typically uses DIMM or SODIMM
		}

		// Determine form factor from context if possible
		if strings.Contains(strings.ToLower(slot.Type), "sodimm") {
			module.FormFactor = "SODIMM"
		}

		modules = append(modules, module)
	}

	return modules
}

// parseSystemProfilerMemory parses system_profiler XML output
func parseSystemProfilerMemory(xmlData []byte) []MemorySlot {
	slots := make([]MemorySlot, 0)

	var plist PlistData
	err := xml.Unmarshal(xmlData, &plist)
	if err != nil {
		return slots
	}

	// Navigate the plist structure to find memory slots
	if len(plist.Array.Dicts) == 0 {
		return slots
	}

	// Look for the main memory dict
	for _, dict := range plist.Array.Dicts {
		// Find memory items
		for i, key := range dict.Keys {
			if key == "_items" && i < len(dict.Arrays) {
				// Found the items array containing memory banks
				for _, bankDict := range dict.Arrays[i].Dicts {
					// Each bank may have multiple slots
					for j, bankKey := range bankDict.Keys {
						if bankKey == "_items" && j < len(bankDict.Arrays) {
							// Found memory slots
							for _, slotDict := range bankDict.Arrays[j].Dicts {
								slot := parseMemorySlot(slotDict)
								slots = append(slots, slot)
							}
						}
					}
				}
			}
		}
	}

	return slots
}

// parseMemorySlot extracts memory slot information from a dict
func parseMemorySlot(dict DictData) MemorySlot {
	slot := MemorySlot{}

	for i, key := range dict.Keys {
		if i >= len(dict.Values) {
			continue
		}

		value := dict.Values[i]

		switch key {
		case "dimm_size":
			slot.Size = value
		case "dimm_type":
			slot.Type = value
		case "dimm_speed":
			slot.Speed = value
		case "dimm_status":
			slot.Status = value
		case "dimm_manufacturer":
			slot.Manufacturer = value
		case "dimm_part_number":
			slot.PartNumber = value
		case "dimm_serial_number":
			slot.SerialNumber = value
		}
	}

	return slot
}

// parseMemorySizeMac converts macOS size strings like "8 GB" to bytes
func parseMemorySizeMac(sizeStr string) uint64 {
	if sizeStr == "" || sizeStr == "empty" {
		return 0
	}

	// Remove any non-numeric prefix (like "Size: ")
	sizeStr = strings.TrimSpace(sizeStr)
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

// parseMemorySpeedMac converts speed strings like "2400 MHz" to numeric value
func parseMemorySpeedMac(speedStr string) uint64 {
	if speedStr == "" {
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
