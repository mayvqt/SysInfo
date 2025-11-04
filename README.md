# SysInfo

A comprehensive cross-platform system information tool written in Go. SysInfo collects and displays detailed information about your computer including CPU, memory, disks, network interfaces, processes, and more.

## Features

- 🖥️ **System Information**: Hostname, OS, platform, kernel version, uptime
- ⚡ **CPU Details**: Model, cores, frequency, per-core usage, load averages
- 💾 **Memory Stats**: RAM and swap usage with detailed metrics
- 💿 **Disk Information**: Partitions, usage, I/O statistics, SMART data support
- 🌐 **Network Interfaces**: IP addresses, MAC addresses, traffic statistics
- 📊 **Process Information**: Running processes, top consumers by CPU and memory
- 🎨 **Multiple Output Formats**: Pretty-printed (colored), plain text, JSON
- 📝 **File Output**: Save reports to files
- 🎯 **Modular Design**: Select specific information modules to display
- 🔧 **Highly Customizable**: Extensive CLI flags for fine-grained control

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mayvqt/sysinfo.git
cd sysinfo

# Download dependencies
go mod download

# Build the application
go build -o sysinfo

# (Optional) Install to your PATH
go install
```

## Usage

### Basic Usage

Display all system information with pretty formatting:

```bash
sysinfo
```

### Output Formats

```bash
# Pretty formatted output (default, with colors)
sysinfo --format pretty

# Plain text output
sysinfo --format text

# JSON output
sysinfo --format json
```

### Save to File

```bash
# Save output to a file
sysinfo --output report.txt

# Save JSON report
sysinfo --format json --output system-info.json
```

### Select Specific Modules

By default, all information is collected. You can select specific modules:

```bash
# Only CPU and memory information
sysinfo --cpu --memory

# Only disk information
sysinfo --disk

# System and network information
sysinfo --system --network

# Include SMART disk data (may require elevated privileges)
sysinfo --disk --smart
```

### Available Modules

- `--all` - Collect all information (default)
- `--system` - System information (hostname, OS, uptime, etc.)
- `--cpu` - CPU information and usage
- `--memory` - Memory and swap information
- `--disk` - Disk partitions and usage
- `--network` - Network interfaces and statistics
- `--process` - Process information and top consumers
- `--smart` - SMART disk health data (requires elevated privileges on some systems)

### Additional Options

```bash
# Verbose output (shows collection progress)
sysinfo --verbose

# Short flags
sysinfo -f json -o output.json -v
```

## Examples

### Example 1: Quick System Overview

```bash
sysinfo --system --cpu --memory
```

### Example 2: Detailed Disk Report

```bash
sysinfo --disk --smart --format pretty --output disk-report.txt
```

### Example 3: JSON Export for Monitoring

```bash
sysinfo --format json --output /var/log/sysinfo-$(date +%Y%m%d).json
```

### Example 4: Network Diagnostics

```bash
sysinfo --network --verbose
```

## Project Structure

```
SysInfo/
├── main.go                     # Application entry point
├── go.mod                      # Go module definition
├── cmd/
│   └── root.go                 # CLI command and flag definitions
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration structures
│   ├── types/
│   │   └── types.go            # Data type definitions
│   ├── collector/
│   │   ├── collector.go        # Main collection orchestrator
│   │   ├── system.go           # System information collector
│   │   ├── cpu.go              # CPU information collector
│   │   ├── memory.go           # Memory information collector
│   │   ├── disk.go             # Disk information collector
│   │   ├── network.go          # Network information collector
│   │   └── process.go          # Process information collector
│   ├── formatter/
│   │   ├── formatter.go        # Format dispatcher
│   │   ├── text.go             # Plain text formatter
│   │   └── pretty.go           # Pretty/colored formatter
│   └── utils/
│       └── format.go           # Utility functions
└── README.md
```

## Extending SysInfo

The codebase is designed to be easily extensible. Here's how to add new functionality:

### Adding a New Collector

1. Create a new file in `internal/collector/` (e.g., `battery.go`)
2. Define the data structure in `internal/types/types.go`
3. Implement the collector function
4. Add the collector to `internal/collector/collector.go`
5. Add a CLI flag in `cmd/root.go`

Example:

```go
// internal/types/types.go
type BatteryData struct {
    Present     bool    `json:"present"`
    Percent     float64 `json:"percent"`
    PluggedIn   bool    `json:"plugged_in"`
}

// internal/collector/battery.go
package collector

func CollectBattery() (*types.BatteryData, error) {
    // Implementation here
    return &types.BatteryData{}, nil
}

// internal/collector/collector.go
if cfg.ShouldCollect("battery") {
    info.Battery, err = CollectBattery()
}
```

### Adding a New Output Format

1. Create a new formatter in `internal/formatter/` (e.g., `xml.go`)
2. Implement the format function
3. Add the format to the switch statement in `formatter.go`

## Dependencies

- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system and process utilities
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [pflag](https://github.com/spf13/pflag) - Drop-in replacement for Go's flag package
- [color](https://github.com/fatih/color) - Colored terminal output
- [tablewriter](https://github.com/olekukonko/tablewriter) - ASCII table generation

## Platform Support

SysInfo is designed to work on:

- ✅ Windows
- ✅ macOS
- ✅ Linux
- ✅ FreeBSD

Some features may have limited availability on certain platforms (e.g., load averages on Windows).

## Requirements for SMART Data

SMART disk data collection may require:

- **Linux**: Root privileges or membership in the `disk` group
- **Windows**: Administrator privileges
- **macOS**: Root privileges

Run with elevated privileges:

```bash
# Linux/macOS
sudo sysinfo --smart

# Windows (Run PowerShell as Administrator)
sysinfo --smart
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Author

**mayvqt**

## Acknowledgments

- Thanks to the gopsutil project for excellent cross-platform system utilities
- Inspired by various system information tools like neofetch, htop, and System Information Viewer
