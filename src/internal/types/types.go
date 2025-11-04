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
	CacheSize   int32        `json:"cache_size"`
	Usage       []float64    `json:"usage_percent"`
	LoadAvg     *LoadAverage `json:"load_average,omitempty"`
}

// LoadAverage contains system load averages
type LoadAverage struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// MemoryData contains memory information
type MemoryData struct {
	Total          uint64  `json:"total_bytes"`
	Available      uint64  `json:"available_bytes"`
	Used           uint64  `json:"used_bytes"`
	UsedPercent    float64 `json:"used_percent"`
	Free           uint64  `json:"free_bytes"`
	TotalFormatted string  `json:"total_formatted"`
	UsedFormatted  string  `json:"used_formatted"`
	FreeFormatted  string  `json:"free_formatted"`
	SwapTotal      uint64  `json:"swap_total_bytes"`
	SwapUsed       uint64  `json:"swap_used_bytes"`
	SwapFree       uint64  `json:"swap_free_bytes"`
	SwapPercent    float64 `json:"swap_used_percent"`
}

// DiskData contains disk and partition information
type DiskData struct {
	Partitions []PartitionInfo `json:"partitions"`
	IOStats    []DiskIOStat    `json:"io_stats,omitempty"`
	SMARTData  []SMARTInfo     `json:"smart_data,omitempty"`
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
	Device       string            `json:"device"`
	Serial       string            `json:"serial,omitempty"`
	ModelFamily  string            `json:"model_family,omitempty"`
	DeviceModel  string            `json:"device_model,omitempty"`
	Capacity     uint64            `json:"capacity_bytes,omitempty"`
	Healthy      bool              `json:"healthy"`
	Temperature  int               `json:"temperature_celsius,omitempty"`
	PowerOnHours uint64            `json:"power_on_hours,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`
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
