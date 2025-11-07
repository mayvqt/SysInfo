package types

import "time"

// SystemInfo holds all collected system information
type SystemInfo struct {
	Timestamp time.Time    `json:"timestamp"`
	System    *SystemData  `json:"system,omitempty"`
	CPU       *CPUData     `json:"cpu,omitempty"`
	Memory    *MemoryData  `json:"memory,omitempty"`
	Disk      *DiskData    `json:"disk,omitempty"`
	Network   *NetworkData `json:"network,omitempty"`
	Processes *ProcessData `json:"processes,omitempty"`
	GPU       *GPUData     `json:"gpu,omitempty"`
}

// SystemData contains general system information
type SystemData struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformFamily  string `json:"platform_family"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion   string `json:"kernel_version"`
	KernelArch      string `json:"kernel_arch"`
	Uptime          uint64 `json:"uptime_seconds"`
	UptimeFormatted string `json:"uptime_formatted"`
	BootTime        uint64 `json:"boot_time"`
	Procs           uint64 `json:"processes"`
}

// CPUData contains CPU information
type CPUData struct {
	ModelName   string       `json:"model_name"`
	Cores       int32        `json:"physical_cores"`
	LogicalCPUs int32        `json:"logical_cpus"`
	Vendor      string       `json:"vendor"`
	Family      string       `json:"family"`
	Model       string       `json:"model"`
	Stepping    int32        `json:"stepping"`
	MHz         float64      `json:"mhz"`
	MinMHz      float64      `json:"min_mhz,omitempty"`
	MaxMHz      float64      `json:"max_mhz,omitempty"`
	CacheSize   int32        `json:"cache_size"`
	Usage       []float64    `json:"usage_percent"`
	LoadAvg     *LoadAverage `json:"load_average,omitempty"`
	Flags       []string     `json:"flags,omitempty"`
	Microcode   string       `json:"microcode,omitempty"`
}

// LoadAverage contains system load averages
type LoadAverage struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// MemoryData contains memory information
type MemoryData struct {
	Total          uint64         `json:"total_bytes"`
	Available      uint64         `json:"available_bytes"`
	Used           uint64         `json:"used_bytes"`
	UsedPercent    float64        `json:"used_percent"`
	Free           uint64         `json:"free_bytes"`
	TotalFormatted string         `json:"total_formatted"`
	UsedFormatted  string         `json:"used_formatted"`
	FreeFormatted  string         `json:"free_formatted"`
	SwapTotal      uint64         `json:"swap_total_bytes"`
	SwapUsed       uint64         `json:"swap_used_bytes"`
	SwapFree       uint64         `json:"swap_free_bytes"`
	SwapPercent    float64        `json:"swap_used_percent"`
	Modules        []MemoryModule `json:"memory_modules,omitempty"`
	VirtualTotal   uint64         `json:"virtual_total_bytes,omitempty"`
	VirtualUsed    uint64         `json:"virtual_used_bytes,omitempty"`
	Cached         uint64         `json:"cached_bytes,omitempty"`
	Buffers        uint64         `json:"buffers_bytes,omitempty"`
	Shared         uint64         `json:"shared_bytes,omitempty"`
}

// MemoryModule contains information about a physical memory module
type MemoryModule struct {
	Locator      string `json:"locator"`
	Capacity     uint64 `json:"capacity_bytes"`
	Speed        uint64 `json:"speed_mhz,omitempty"`
	Type         string `json:"type,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	PartNumber   string `json:"part_number,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
	FormFactor   string `json:"form_factor,omitempty"`
}

// DiskData contains disk and partition information
type DiskData struct {
	Partitions    []PartitionInfo `json:"partitions"`
	PhysicalDisks []PhysicalDisk  `json:"physical_disks,omitempty"`
	IOStats       []DiskIOStat    `json:"io_stats,omitempty"`
	SMARTData     []SMARTInfo     `json:"smart_data,omitempty"`
}

// PhysicalDisk contains information about physical disks
type PhysicalDisk struct {
	Name          string `json:"name"`
	Model         string `json:"model,omitempty"`
	SerialNumber  string `json:"serial_number,omitempty"`
	Size          uint64 `json:"size_bytes"`
	SizeFormatted string `json:"size_formatted"`
	Type          string `json:"type,omitempty"`      // HDD, SSD, NVMe, etc.
	Interface     string `json:"interface,omitempty"` // SATA, NVMe, USB, etc.
	RPM           uint32 `json:"rpm,omitempty"`       // For HDDs
	Removable     bool   `json:"removable"`
}

// PartitionInfo contains information about a disk partition
type PartitionInfo struct {
	Device         string  `json:"device"`
	MountPoint     string  `json:"mount_point"`
	FSType         string  `json:"fs_type"`
	Total          uint64  `json:"total_bytes"`
	Free           uint64  `json:"free_bytes"`
	Used           uint64  `json:"used_bytes"`
	UsedPercent    float64 `json:"used_percent"`
	TotalFormatted string  `json:"total_formatted"`
	UsedFormatted  string  `json:"used_formatted"`
	FreeFormatted  string  `json:"free_formatted"`
	InodesTotal    uint64  `json:"inodes_total,omitempty"`
	InodesUsed     uint64  `json:"inodes_used,omitempty"`
	InodesFree     uint64  `json:"inodes_free,omitempty"`
}

// DiskIOStat contains disk I/O statistics
type DiskIOStat struct {
	Name       string `json:"name"`
	ReadCount  uint64 `json:"read_count"`
	WriteCount uint64 `json:"write_count"`
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
	ReadTime   uint64 `json:"read_time_ms"`
	WriteTime  uint64 `json:"write_time_ms"`
	IoTime     uint64 `json:"io_time_ms"`
}

// SMARTInfo contains SMART data for a drive
type SMARTInfo struct {
	Device           string             `json:"device"`
	Serial           string             `json:"serial,omitempty"`
	ModelFamily      string             `json:"model_family,omitempty"`
	DeviceModel      string             `json:"device_model,omitempty"`
	FirmwareVersion  string             `json:"firmware_version,omitempty"`
	Capacity         uint64             `json:"capacity_bytes,omitempty"`
	Healthy          bool               `json:"healthy"`
	Temperature      int                `json:"temperature_celsius,omitempty"`
	PowerOnHours     uint64             `json:"power_on_hours,omitempty"`
	PowerCycleCount  uint64             `json:"power_cycle_count,omitempty"`
	RotationRate     uint32             `json:"rotation_rate,omitempty"` // 0 for SSDs, RPM for HDDs
	FormFactor       string             `json:"form_factor,omitempty"`   // 2.5", 3.5", M.2, etc.
	ATAVersion       string             `json:"ata_version,omitempty"`
	SATAVersion      string             `json:"sata_version,omitempty"`
	LocalTime        string             `json:"local_time,omitempty"`
	Attributes       map[string]string  `json:"attributes,omitempty"`
	DetailedAttribs  []SMARTAttribute   `json:"detailed_attributes,omitempty"`
	ErrorLog         *SMARTErrorLog     `json:"error_log,omitempty"`
	SelfTestLog      *SMARTSelfTestLog  `json:"self_test_log,omitempty"`
	HealthAssessment *SMARTHealthStatus `json:"health_assessment,omitempty"`
}

// SMARTAttribute contains detailed information about a SMART attribute
type SMARTAttribute struct {
	ID         uint8  `json:"id"`
	Name       string `json:"name"`
	Flag       uint16 `json:"flag"`
	Value      uint8  `json:"value"`       // Current normalized value
	Worst      uint8  `json:"worst"`       // Worst value seen
	Threshold  uint8  `json:"threshold"`   // Failure threshold
	RawValue   uint64 `json:"raw_value"`   // Raw value
	Type       string `json:"type"`        // Pre-fail or Old_age
	Updated    string `json:"updated"`     // Always or Offline
	WhenFailed string `json:"when_failed"` // Never, FAILING_NOW, In_the_past
	RawString  string `json:"raw_string"`  // Human-readable raw value
}

// SMARTErrorLog contains SMART error log information
type SMARTErrorLog struct {
	ErrorCount uint64       `json:"error_count"`
	Errors     []SMARTError `json:"errors,omitempty"`
}

// SMARTError represents a single SMART error entry
type SMARTError struct {
	ErrorNumber   uint64 `json:"error_number"`
	LifetimeHours uint64 `json:"lifetime_hours"`
	State         string `json:"state"`
	Type          string `json:"type"`
	Details       string `json:"details"`
}

// SMARTSelfTestLog contains SMART self-test log information
type SMARTSelfTestLog struct {
	TestCount uint64          `json:"test_count"`
	Tests     []SMARTSelfTest `json:"tests,omitempty"`
}

// SMARTSelfTest represents a single self-test entry
type SMARTSelfTest struct {
	TestNumber      uint64 `json:"test_number"`
	TestDescription string `json:"test_description"`
	Status          string `json:"status"`
	PercentComplete uint8  `json:"percent_complete"`
	LifetimeHours   uint64 `json:"lifetime_hours"`
	LBA             uint64 `json:"lba,omitempty"`
}

// SMARTHealthStatus contains health assessment information
type SMARTHealthStatus struct {
	Passed            bool     `json:"passed"`
	OverallAssessment string   `json:"overall_assessment"` // PASS, FAIL, WARN
	FailingAttributes []string `json:"failing_attributes,omitempty"`
	WarningAttributes []string `json:"warning_attributes,omitempty"`
	CriticalWarning   string   `json:"critical_warning,omitempty"`
	PercentUsed       float64  `json:"percent_used,omitempty"`    // For SSDs
	AvailableSpare    float64  `json:"available_spare,omitempty"` // For NVMe
	TemperatureStatus string   `json:"temperature_status,omitempty"`
}

// NetworkData contains network information
type NetworkData struct {
	Interfaces  []NetworkInterface `json:"interfaces"`
	Connections int                `json:"connection_count,omitempty"`
}

// NetworkInterface contains information about a network interface
type NetworkInterface struct {
	Name         string   `json:"name"`
	HardwareAddr string   `json:"hardware_addr"`
	Addresses    []string `json:"addresses"`
	Flags        []string `json:"flags"`
	MTU          int      `json:"mtu"`
	BytesSent    uint64   `json:"bytes_sent"`
	BytesRecv    uint64   `json:"bytes_recv"`
	PacketsSent  uint64   `json:"packets_sent"`
	PacketsRecv  uint64   `json:"packets_recv"`
	ErrorsIn     uint64   `json:"errors_in"`
	ErrorsOut    uint64   `json:"errors_out"`
	DropsIn      uint64   `json:"drops_in"`
	DropsOut     uint64   `json:"drops_out"`
}

// ProcessData contains process information
type ProcessData struct {
	TotalCount  int           `json:"total_count"`
	Running     int           `json:"running"`
	Sleeping    int           `json:"sleeping"`
	TopByMemory []ProcessInfo `json:"top_by_memory,omitempty"`
	TopByCPU    []ProcessInfo `json:"top_by_cpu,omitempty"`
}

// ProcessInfo contains information about a single process
type ProcessInfo struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	Username      string  `json:"username,omitempty"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float32 `json:"memory_percent"`
	MemoryMB      uint64  `json:"memory_mb"`
	Status        string  `json:"status"`
	CreateTime    int64   `json:"create_time,omitempty"`
}
