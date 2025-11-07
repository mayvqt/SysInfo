package types

// GPUData contains GPU information
type GPUData struct {
	GPUs []GPUInfo `json:"gpus"`
}

// GPUInfo contains information about a single GPU
type GPUInfo struct {
	Index            int     `json:"index"`
	Name             string  `json:"name"`
	Vendor           string  `json:"vendor"`
	Driver           string  `json:"driver,omitempty"`
	DriverVersion    string  `json:"driver_version,omitempty"`
	MemoryTotal      uint64  `json:"memory_total_bytes,omitempty"`
	MemoryUsed       uint64  `json:"memory_used_bytes,omitempty"`
	MemoryFree       uint64  `json:"memory_free_bytes,omitempty"`
	MemoryFormatted  string  `json:"memory_total_formatted,omitempty"`
	Temperature      int     `json:"temperature_celsius,omitempty"`
	FanSpeed         int     `json:"fan_speed_percent,omitempty"`
	PowerDraw        float64 `json:"power_draw_watts,omitempty"`
	PowerLimit       float64 `json:"power_limit_watts,omitempty"`
	Utilization      int     `json:"utilization_percent,omitempty"`
	MemoryUtilization int    `json:"memory_utilization_percent,omitempty"`
	ClockSpeed       int     `json:"clock_speed_mhz,omitempty"`
	ClockSpeedMemory int     `json:"clock_speed_memory_mhz,omitempty"`
	PCIBus           string  `json:"pci_bus,omitempty"`
	UUID             string  `json:"uuid,omitempty"`
}
