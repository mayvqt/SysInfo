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
└── # SysInfo

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

SysInfo is a professional-grade, cross-platform command-line utility designed for comprehensive system information collection and analysis. Built with Go, it provides detailed insights into hardware and software configurations including CPU specifications, memory statistics, disk information, network interfaces, process monitoring, and SMART health data.

Engineered for reliability and flexibility, SysInfo offers multiple output formats and granular module selection capabilities. Whether performing routine system audits, gathering diagnostics for support tickets, or monitoring infrastructure health, SysInfo delivers accurate, well-structured information in your preferred format.

## Features

### Core Functionality

**System Information Collection**
- CPU details including model, architecture, core count, frequency ranges, cache sizes, microcode versions, and feature flags
- Memory statistics with total, used, available, cached, buffers, swap usage, and physical RAM module specifications
- Disk information encompassing partitions, usage metrics, filesystem types, I/O statistics, physical disk details, and serial numbers
- Network interface enumeration with IP addresses, MAC addresses, bandwidth statistics, and connection states
- Process monitoring displaying running processes with CPU and memory consumption rankings
- System overview including hostname, operating system, platform details, kernel version, and uptime
- SMART disk health monitoring for predictive failure analysis (requires elevated privileges)

**Output Format Options**
- Pretty format: Human-readable colored output with formatted tables and visual hierarchy
- JSON format: Structured data output for programmatic parsing, API integration, and automation workflows
- Text format: Plain text output suitable for logging systems, archival, and simple parsing
- File output: Persistent storage of reports for historical analysis, comparison, and documentation

**Modular Architecture**
- Selective module execution to collect only required information subsets
- Default comprehensive mode for complete system snapshots
- Clean separation of concerns with distinct collector, formatter, and configuration layers
- Extensible design supporting straightforward addition of new collectors and output formats

### Cross-Platform Support

**Supported Operating Systems**
- Windows (AMD64, ARM64)
- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)

**Platform Characteristics**
- Single static binary with zero external runtime dependencies
- Pure Go implementation leveraging gopsutil for platform abstraction
- Minimal resource footprint and memory consumption
- No installation or configuration required for basic operation

## Installation

### Download Pre-built Binary

Download the appropriate binary for your platform:

**Windows**
```powershell
# Download for your architecture
# x64: sysinfo-windows-amd64.exe
# ARM64: sysinfo-windows-arm64.exe
```

**Linux**
```bash
# Download and make executable
wget https://github.com/mayvqt/SysInfo/releases/latest/download/sysinfo-linux-amd64
chmod +x sysinfo-linux-amd64
sudo mv sysinfo-linux-amd64 /usr/local/bin/sysinfo
```

**macOS**
```bash
# Intel Macs
wget https://github.com/mayvqt/SysInfo/releases/latest/download/sysinfo-darwin-amd64
# Apple Silicon Macs
wget https://github.com/mayvqt/SysInfo/releases/latest/download/sysinfo-darwin-arm64

chmod +x sysinfo-darwin-*
sudo mv sysinfo-darwin-* /usr/local/bin/sysinfo
```

### Build from Source

Requirements: Go 1.23 or later

```bash
git clone https://github.com/mayvqt/SysInfo.git
cd SysInfo/src
go build -o sysinfo .
```

Build with embedded version metadata:

```bash
VERSION="v1.0.0"
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

go build -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" -o sysinfo .
```

Cross-compilation for target platforms:

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o sysinfo-windows-amd64.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -o sysinfo-linux-amd64 .

# macOS
GOOS=darwin GOARCH=arm64 go build -o sysinfo-darwin-arm64 .
```

## Usage

### Basic Operations

Execute with default settings to collect all system information with formatted output:

```bash
./sysinfo
```

Persist results to file for documentation or analysis:

```bash
./sysinfo -o system-report.txt
```

### Output Formats

**Pretty format** (default) - Formatted output with color-coded sections and tabular data:
```bash
./sysinfo --format pretty
```

**JSON format** - Structured data for automation, scripting, and integration:
```bash
./sysinfo --format json
```

**Text format** - Plain text output for logging and archival:
```bash
./sysinfo --format text
```

### Module Selection

Execute specific information collectors as needed:

```bash
# CPU information only
./sysinfo --cpu

# Memory and disk information
./sysinfo --memory --disk

# Network interfaces
./sysinfo --network

# Top processes
./sysinfo --process

# System overview
./sysinfo --system

# SMART disk health (may require sudo/admin)
./sysinfo --smart
```

### Advanced Examples

**Generate timestamped JSON reports for automation pipelines:**
```bash
./sysinfo --format json --output system-$(date +%Y%m%d).json
```

**Monitor specific subsystems for targeted diagnostics:**
```bash
./sysinfo --cpu --memory --format pretty
```

**Collect hardware specifications for support documentation:**
```bash
./sysinfo --system --cpu --memory --disk --output support-info.txt
```

**Execute disk health assessment:**
```bash
# Linux/macOS (requires elevated privileges)
sudo ./sysinfo --smart --disk

# Windows (requires administrator privileges)
./sysinfo --smart --disk
```

### Command-Line Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--format` | `-f` | Output format: json, text, pretty | `pretty` |
| `--output` | `-o` | Output file path (stdout if omitted) | - |
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--all` | - | Collect all available information | `true` |
| `--system` | - | Collect system information | `false` |
| `--cpu` | - | Collect CPU information | `false` |
| `--memory` | - | Collect memory information | `false` |
| `--disk` | - | Collect disk information | `false` |
| `--network` | - | Collect network information | `false` |
| `--process` | - | Collect process information | `false` |
| `--smart` | - | Collect SMART disk health data | `false` |

Note: Specifying any individual module flag automatically disables the `--all` flag. To collect all information, either omit module flags entirely or explicitly specify `--all`.

## Output Examples

### Pretty Format

```
╔══════════════════════════════════════════════════════════════╗
║                      SYSTEM INFORMATION                      ║
╚══════════════════════════════════════════════════════════════╝

Hostname:           DESKTOP-ABC123
Operating System:   windows
Platform:           Microsoft Windows 11 Pro
Kernel Version:     10.0.22631.4460
Uptime:             2h 34m 12s

╔══════════════════════════════════════════════════════════════╗
║                      CPU INFORMATION                         ║
╚══════════════════════════════════════════════════════════════╝

Model:              Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
Physical Cores:     8
Logical Cores:      8
Cache Size:         12.0 MB
Current Speed:      3600 MHz
...
```

### JSON Format

```json
{
  "system": {
    "hostname": "DESKTOP-ABC123",
    "os": "windows",
    "platform": "Microsoft Windows 11 Pro",
    "platformFamily": "Standalone Workstation",
    "platformVersion": "10.0.22631.4460",
    "kernelVersion": "10.0.22631.4460",
    "kernelArch": "x86_64",
    "uptime": 9252
  },
  "cpu": {
    "modelName": "Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz",
    "vendor": "GenuineIntel",
    "physicalCores": 8,
    "logicalCores": 8,
    "cacheSize": 12582912,
    ...
  }
}
```

## Platform-Specific Features

### Windows

- RAM module specifications (speed, manufacturer, part numbers) via Windows Management Instrumentation (WMI)
- Physical disk enumeration (model, serial number, interface type) via WMI queries
- SMART data collection through WMI (requires administrator privileges)
- Automatic console pause behavior when executed via Windows Explorer

### Linux

- RAM module specifications via dmidecode utility (requires root privileges)
- Block device information sourced from /sys filesystem
- SMART data collection via smartctl (requires root privileges and smartmontools package installation)
- Systemd-aware uptime calculation for accurate boot time reporting

### macOS

- RAM module specifications via system_profiler utility
- Disk information retrieval via diskutil
- SMART data collection via smartctl (requires root privileges and smartmontools installation)
- Native support for Apple Silicon (M1/M2/M3/M4) processors

## Development

### Project Structure

```
SysInfo/
├── src/
│   ├── main.go                         # Entry point
│   ├── go.mod                          # Go dependencies
│   ├── go.sum                          # Dependency checksums
│   ├── cmd/
│   │   └── root.go                     # CLI command definitions, flags
│   └── internal/
│       ├── config/
│       │   └── config.go               # Configuration structures
│       ├── types/
│       │   └── types.go                # Data structure definitions
│       ├── collector/
│       │   ├── system.go               # System information collector
│       │   ├── cpu.go                  # CPU information collector
│       │   ├── memory.go               # Memory information collector
│       │   ├── disk.go                 # Disk information collector
│       │   ├── network.go              # Network information collector
│       │   └── process.go              # Process information collector
│       ├── formatter/
│       │   ├── json.go                 # JSON output formatter
│       │   ├── text.go                 # Text output formatter
│       │   └── pretty.go               # Pretty/colored output formatter
│       └── utils/
│           └── format.go               # Utility functions (byte formatting)
├── .gitignore                          # Git ignore rules
├── LICENSE                             # MIT license
└── README.md                           # This file
```

### Architecture

SysInfo implements a modular collector-formatter architectural pattern:

1. **Configuration Layer** (`config/`): Command-line flag parsing and runtime configuration management
2. **Data Types Layer** (`types/`): Structured definitions for all collected system information
3. **Collector Layer** (`collector/`): Independent, specialized modules for gathering specific system metrics
4. **Formatter Layer** (`formatter/`): Output transformation engines supporting multiple presentation formats
5. **Utility Layer** (`utils/`): Shared helper functions for data formatting and conversion operations

Architectural benefits:
- Addition of new collectors without modification to existing codebase
- Support for additional output formats through isolated formatter implementation
- Component-level isolation enabling comprehensive unit testing
- Platform-specific implementation extension without core logic changes

### Adding New Collectors

Implementation workflow:

1. Define data structures in `internal/types/types.go`
2. Create collector module in `internal/collector/your_collector.go`
3. Implement collection logic using gopsutil or platform-specific APIs
4. Register module flag in `cmd/root.go`
5. Update formatters to render newly collected data
6. Document functionality and usage patterns

Example:

```go
// types/types.go
type GPUData struct {
    Model       string
    Memory      uint64
    Driver      string
    Temperature float64
}

// collector/gpu.go
package collector

import "github.com/mayvqt/sysinfo/internal/types"

func CollectGPU() (*types.GPUData, error) {
    // Implementation
    return &types.GPUData{}, nil
}
```

### Adding New Formatters

Implementation workflow:

1. Create formatter module in `internal/formatter/your_formatter.go`
2. Implement format-specific rendering logic
3. Register format option in `cmd/root.go`
4. Document format specification and use cases

Example:

```go
// formatter/xml.go
package formatter

import "encoding/xml"

func FormatXML(data *types.SystemData) (string, error) {
    output, err := xml.MarshalIndent(data, "", "  ")
    return string(output), err
}
```

### Dependencies

- **[gopsutil/v3](https://github.com/shirou/gopsutil)** - Cross-platform system and process utilities
- **[cobra](https://github.com/spf13/cobra)** - CLI framework for command structure and flags
- **[color](https://github.com/fatih/color)** - Colored terminal output
- **[tablewriter](https://github.com/olekukonko/tablewriter)** - ASCII table generation

### Testing

Execute test suite:

```bash
cd src
go test ./... -v
```

Generate coverage analysis:

```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Execute with race condition detection:

```bash
go test ./... -v -race
```

### Contributing

Contributions are welcome. Please adhere to the following guidelines:

1. Fork the repository and create a feature branch: `git checkout -b feature/descriptive-name`
2. Implement comprehensive tests for new functionality
3. Verify all tests pass: `go test ./... -v`
4. Apply standard formatting: `go fmt ./...`
5. Commit changes with clear, descriptive messages
6. Submit pull request with detailed description of changes

**Code Standards:**
- Adhere to standard Go formatting conventions (`go fmt`)
- Use descriptive, intention-revealing names for variables and functions
- Document all exported functions and complex logic with comments
- Maintain single-responsibility principle for functions
- Handle errors explicitly with appropriate context

## Troubleshooting

### SMART Data Collection Failures

SMART data collection requires elevated privileges and platform-specific tools:

**Windows:**
- Execute with administrator privileges
- Verify Windows Management Instrumentation service is operational

**Linux:**
- Execute with root privileges: `sudo ./sysinfo --smart`
- Install smartmontools package: `sudo apt-get install smartmontools` (Debian/Ubuntu) or `sudo yum install smartmontools` (RHEL/CentOS)

**macOS:**
- Execute with root privileges: `sudo ./sysinfo --smart`
- Install smartmontools via Homebrew: `brew install smartmontools`

### Missing RAM Module Specifications

Detailed RAM module information requires platform-specific utilities:

**Windows:**
- Typically functions without additional configuration via WMI
- Verify Windows Management Instrumentation service is running

**Linux:**
- Install dmidecode utility: `sudo apt-get install dmidecode`
- Execute with root privileges: `sudo ./sysinfo --memory`

**macOS:**
- Utilizes built-in system_profiler utility
- No additional configuration required

### Permission Denied Errors

Certain collectors require elevated execution privileges:

**Linux/macOS:**
```bash
sudo ./sysinfo --all
```

**Windows:**
- Right-click executable and select "Run as Administrator"
- Execute from elevated PowerShell or Command Prompt session

### Output Encoding Issues

Character encoding problems in terminal output:

**Windows:**
- Configure console for UTF-8 encoding: `chcp 65001`
- Use Windows Terminal for enhanced Unicode support

**Linux/macOS:**
- Verify terminal supports UTF-8 encoding
- Confirm LANG environment variable: `echo $LANG`

## Roadmap

Future enhancements planned:

- GPU information collection (NVIDIA, AMD, Intel)
- Battery and power information for laptops
- Temperature sensors (CPU, GPU, motherboard)
- Fan speed monitoring
- Detailed USB device enumeration
- PCI device information
- BIOS/UEFI information
- Virtualization detection and details
- Container environment detection
- HTML output format with charts
- CSV export for spreadsheet analysis
- Comparison mode (diff between two reports)
- Watch mode (continuous monitoring)
- Web interface for remote monitoring

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/mayvqt/SysInfo/issues)
- **Discussions**: Ask questions or share ideas in [GitHub Discussions](https://github.com/mayvqt/SysInfo/discussions)

## Acknowledgments

Developed with Go and leveraging the cross-platform capabilities of the [gopsutil](https://github.com/shirou/gopsutil) library for system information collection.
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
