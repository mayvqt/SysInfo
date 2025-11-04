//go:build windows
// +build windows

package collector

import (
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/yusufpapurcu/wmi"
)

// Win32_PhysicalMemory represents WMI physical memory data
type Win32_PhysicalMemory struct {
	BankLabel            string
	Capacity             uint64
	Speed                uint32
	Manufacturer         string
	PartNumber           string
	SerialNumber         string
	MemoryType           uint16
	FormFactor           uint16
	DeviceLocator        string
	ConfiguredClockSpeed uint32
	ConfiguredVoltage    uint32
	DataWidth            uint16
	TotalWidth           uint16
	Tag                  string
	TypeDetail           uint16
	MinVoltage           uint32
	MaxVoltage           uint32
}

// collectMemoryModulesPlatform implements Windows-specific memory module collection
func collectMemoryModulesPlatform() []types.MemoryModule {
	modules := make([]types.MemoryModule, 0)

	var memoryModules []Win32_PhysicalMemory
	query := "SELECT * FROM Win32_PhysicalMemory"
	err := wmi.Query(query, &memoryModules)
	if err != nil {
		return modules
	}

	for _, mem := range memoryModules {
		module := types.MemoryModule{
			Locator:      getLocator(mem),
			Capacity:     mem.Capacity,
			Speed:        uint64(mem.Speed),
			Type:         getMemoryType(mem.MemoryType),
			Manufacturer: strings.TrimSpace(mem.Manufacturer),
			PartNumber:   strings.TrimSpace(mem.PartNumber),
			SerialNumber: strings.TrimSpace(mem.SerialNumber),
			FormFactor:   getFormFactor(mem.FormFactor),
		}
		modules = append(modules, module)
	}

	return modules
}

// getLocator determines the memory module locator
func getLocator(mem Win32_PhysicalMemory) string {
	if mem.DeviceLocator != "" {
		return mem.DeviceLocator
	}
	if mem.BankLabel != "" {
		return mem.BankLabel
	}
	if mem.Tag != "" {
		return mem.Tag
	}
	return "Unknown"
}

// getMemoryType maps memory type codes to names
func getMemoryType(typeCode uint16) string {
	types := map[uint16]string{
		0:  "Unknown",
		1:  "Other",
		2:  "DRAM",
		3:  "Synchronous DRAM",
		4:  "Cache DRAM",
		5:  "EDO",
		6:  "EDRAM",
		7:  "VRAM",
		8:  "SRAM",
		9:  "RAM",
		10: "ROM",
		11: "Flash",
		12: "EEPROM",
		13: "FEPROM",
		14: "EPROM",
		15: "CDRAM",
		16: "3DRAM",
		17: "SDRAM",
		18: "SGRAM",
		19: "RDRAM",
		20: "DDR",
		21: "DDR2",
		22: "DDR2 FB-DIMM",
		24: "DDR3",
		25: "FBD2",
		26: "DDR4",
		27: "LPDDR",
		28: "LPDDR2",
		29: "LPDDR3",
		30: "LPDDR4",
		31: "Logical non-volatile device",
		32: "HBM (High Bandwidth Memory)",
		33: "HBM2 (High Bandwidth Memory Generation 2)",
		34: "DDR5",
		35: "LPDDR5",
	}

	if name, ok := types[typeCode]; ok {
		return name
	}
	return "Unknown"
}

// getFormFactor maps form factor codes to names
func getFormFactor(formFactorCode uint16) string {
	formFactors := map[uint16]string{
		0:  "Unknown",
		1:  "Other",
		2:  "SIP",
		3:  "DIP",
		4:  "ZIP",
		5:  "SOJ",
		6:  "Proprietary",
		7:  "SIMM",
		8:  "DIMM",
		9:  "TSOP",
		10: "PGA",
		11: "RIMM",
		12: "SODIMM",
		13: "SRIMM",
		14: "SMD",
		15: "SSMP",
		16: "QFP",
		17: "TQFP",
		18: "SOIC",
		19: "LCC",
		20: "PLCC",
		21: "BGA",
		22: "FPBGA",
		23: "LGA",
		24: "FB-DIMM",
	}

	if name, ok := formFactors[formFactorCode]; ok {
		return name
	}
	return "Unknown"
}
