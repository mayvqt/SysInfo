# SysInfo# SysInfo# SysInfo# SysInfogit clone https://github.com/mayvqt/sysinfo.git



[![CI](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml/badge.svg)](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Release](https://img.shields.io/github/v/release/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/releases)[![Go Report Card](https://goreportcard.com/badge/github.com/mayvqt/sysinfo)](https://goreportcard.com/report/github.com/mayvqt/sysinfo)

[![Go Version](https://img.shields.io/github/go-mod/go-version/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/blob/main/src/go.mod)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

SysInfo is a lightweight, cross-platform system information tool written in Go. It collects and displays detailed hardware and software information with support for JSON, text, and color-formatted output.

[![Release](https://img.shields.io/github/v/release/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/releases)SysInfo is a lightweight, cross-platform system information tool written in Go. It collects and displays detailed hardware and software information with support for JSON, text, and color-formatted output.go mod download

## Features

[![Go Version](https://img.shields.io/github/go-mod/go-version/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/blob/main/src/go.mod)

- **Cross-platform** — Windows, Linux, and macOS support via gopsutil

- **Modular collection** — Choose specific modules or collect everything

- **Multiple formats** — JSON for automation, plain text for scripting, pretty-printed for readability

- **SMART disk monitoring** — Read disk health metrics (platform-specific, requires elevation)SysInfo is a lightweight, cross-platform system information tool written in Go. It collects and displays detailed hardware and software information with support for JSON, text, and color-formatted output.

- **Process insights** — Top CPU and memory consumers

- **Network stats** — Interface details with I/O counters## FeaturesSysInfo is a cross-platform system information tool written in Go. It collects and displays comprehensive hardware and software information including CPU, memory, disk, network, processes, and SMART data.go build -o sysinfo

- **Lightweight** — Single binary, minimal dependencies

## Features

## Quickstart



```powershell

cd src- **Cross-platform** — Windows, Linux, and macOS support via gopsutil

go build -o sysinfo.exe .

.\sysinfo.exe- **Modular collection** — Choose specific modules or collect everything- **Cross-platform** — Windows, Linux, and macOS support via gopsutilgo install

```

- **Multiple formats** — JSON for automation, plain text for scripting, pretty-printed for readability

The tool runs with `--all` by default, displaying a color-formatted report of all system metrics.

- **SMART disk monitoring** — Read disk health metrics (platform-specific, requires elevation)- **Modular collection** — Choose specific modules or collect everything

## Usage

- **Process insights** — Top CPU and memory consumers

```

sysinfo [flags]- **Network stats** — Interface details with I/O counters- **Multiple formats** — JSON for automation, plain text for scripting, pretty-printed for readability## Featuressysinfo --format pretty

```

- **Lightweight** — Single binary, minimal dependencies

### Output Options

- **SMART disk monitoring** — Read disk health metrics (platform-specific, requires elevation)

| Flag | Values | Default | Description |

|------|--------|---------|-------------|## Quickstart

| `-f, --format` | `json`, `text`, `pretty` | `pretty` | Output format |

| `-o, --output` | path | - | Write to file instead of stdout |- **Process insights** — Top CPU and memory consumerssysinfo --disk --smart

| `-v, --verbose` | - | - | Show collection progress |

```powershell

### Module Selection

cd src- **Network stats** — Interface details with I/O counters

| Flag | Description |

|------|-------------|go build -o sysinfo.exe .

| `--all` | All modules (default unless specific flags used) |

| `--system` | Hostname, OS, platform, kernel, uptime |.\sysinfo.exe- **Lightweight** — Single binary, minimal dependencies- **Cross-platform support** — Windows, Linux, and macOSsysinfo --verbose

| `--cpu` | Model, cores, frequency, per-core usage, load average |

| `--memory` | Total, used, free, cached, swap |```

| `--disk` | Partitions, usage, I/O stats |

| `--network` | Interfaces, addresses, flags, I/O counters |

| `--process` | Total count, top 10 by CPU and memory |

| `--smart` | Disk health data (requires admin/root) |The tool runs with `--all` by default, displaying a color-formatted report of all system metrics.



**Note:** Specifying any module flag (e.g., `--cpu`) disables `--all`. Combine multiple flags to collect specific modules.## Quickstart- **Comprehensive data collection** — CPU, memory, disk, network, processes, and SMART disk health



## Examples## Usage



**Default output (all modules, pretty format):**

```powershell

.\sysinfo.exe```

```

sysinfo [flags]```powershell- **Multiple output formats** — JSON, text, and pretty-printed tables# SysInfo

**CPU and memory only, JSON format:**

```powershell```

.\sysinfo.exe --cpu --memory --format json

```cd src



**Save full report to file:**### Output Options

```powershell

.\sysinfo.exe --output report.json --format jsongo build -o sysinfo.exe .- **Modular design** — Collect only the information you need

```

| Flag | Values | Default | Description |

**SMART disk health (requires elevation):**

```powershell|------|--------|---------|-------------|.\sysinfo.exe

# Windows (run as Administrator)

.\sysinfo.exe --smart --disk| `-f, --format` | `json`, `text`, `pretty` | `pretty` | Output format |



# Linux/macOS| `-o, --output` | path | - | Write to file instead of stdout |```- **File output support** — Save reports to diskSysInfo is a small, cross-platform Go utility and library for collecting basic system information and metrics. It provides lightweight collectors for CPU, memory, disk, network, and processes and can be used as a CLI tool or embedded in other Go programs.

sudo ./sysinfo --smart --disk

```| `-v, --verbose` | - | - | Show collection progress |



**Network interfaces with verbose output:**

```powershell

.\sysinfo.exe --network --verbose### Module Selection

```

The tool runs with `--all` by default, displaying a color-formatted report of all system metrics.- **Lightweight** — Single binary with no dependencies

**Top processes only:**

```powershell| Flag | Description |

.\sysinfo.exe --process --format text

```|------|-------------|



## Output Format Details| `--all` | All modules (default unless specific flags used) |



### Pretty Format| `--system` | Hostname, OS, platform, kernel, uptime |## Usage- **SMART monitoring** — Disk health metrics (requires elevated privileges)## Features



- Color-coded sections with Unicode box-drawing characters| `--cpu` | Model, cores, frequency, per-core usage, load average |

- Progress bars for CPU, memory, disk usage

- Organized by module with clear headers| `--memory` | Total, used, free, cached, swap |

- Temperature warnings in SMART data (yellow >50°C, red >60°C)

- Top 5 processes by memory/CPU displayed| `--disk` | Partitions, usage, I/O stats |



### Text Format| `--network` | Interfaces, addresses, flags, I/O counters |```



- Plain text, no colors or special characters| `--process` | Total count, top 10 by CPU and memory |

- Suitable for logging or scripting

- Same structure as pretty format but simplified| `--smart` | Disk health data (requires admin/root) |sysinfo [flags]



### JSON Format



- Complete structured data**Note:** Specifying any module flag (e.g., `--cpu`) disables `--all`. Combine multiple flags to collect specific modules.```## Quickstart (works out of the box)}

- Omits null/empty fields (`omitempty`)

- Includes timestamp in ISO 8601 format

- All byte values include human-readable formatted versions

## Examples

## Platform-Specific Notes



### Windows

- SMART collection uses WMI (`Win32_DiskDrive`, `MSStorageDriver_*` classes)**Default output (all modules, pretty format):**### Output Optionsif cfg.ShouldCollect("battery") {

- Requires Administrator privileges for SMART data

- Executable pauses on exit when double-clicked (not from terminal)```powershell



### Linux.\sysinfo.exe

- SMART collection uses `smartctl` (requires `smartmontools` package)

- Supports both ATA and NVMe drives via `smartctl --json````

- Requires root for SMART data

- Load average always available (1, 5, 15 minute intervals)| Flag | Values | Default | Description |Build and run locally:



### macOS**CPU and memory only, JSON format:**

- SMART collection uses `smartctl` (install via `brew install smartmontools`)

- Requires root for SMART data```powershell|------|--------|---------|-------------|

- Some disk I/O counters may be limited compared to Linux

.\sysinfo.exe --cpu --memory --format json

## Data Collected

```| `-f, --format` | `json`, `text`, `pretty` | `pretty` | Output format |# SysInfo

### System Module

- Hostname, OS, platform family/version

- Kernel version and architecture

- Uptime (formatted and in seconds)**Save full report to file:**| `-o, --output` | path | - | Write to file instead of stdout |

- Boot time, process count

```powershell

### CPU Module

- Model name, vendor, family, model, stepping.\sysinfo.exe --output report.json --format json| `-v, --verbose` | - | - | Show collection progress |```powershell

- Physical cores and logical CPUs

- Current, min, max frequency```

- Cache size, microcode version

- Per-core usage percentages (sampled over 1 second)

- Load average (Linux/macOS)

- CPU flags**SMART disk health (requires elevation):**



### Memory Module```powershell### Module Selection# Navigate to src directorySysInfo is a small, cross-platform Go utility and library for collecting basic system information and metrics. It provides lightweight collectors for CPU, memory, disk, network, and processes and can be used as a CLI tool or embedded in other Go programs.

- Total, used, free, available (bytes and formatted)

- Usage percentage# Windows (run as Administrator)

- Cached, buffers, shared memory

- Swap total, used, free, percentage.\sysinfo.exe --smart --disk

- Physical RAM modules (placeholder for future WMI/dmidecode integration)



### Disk Module

- Partitions: device, mount point, filesystem type# Linux/macOS| Flag | Description |cd src

- Total, used, free space (bytes, formatted, percentage)

- Inode counts (Linux)sudo ./sysinfo --smart --disk

- I/O statistics per disk: read/write counts, bytes, time

- Physical disk information (placeholder)```|------|-------------|

- SMART data: device, model, serial, capacity, health status, temperature, power-on hours, critical attributes



### Network Module

- Interface name, MAC address, MTU**Network interfaces with verbose output:**| `--all` | All modules (default unless specific flags used) |## Features

- IP addresses (all assigned)

- Flags (UP, BROADCAST, LOOPBACK, MULTICAST)```powershell

- Bytes/packets sent and received

- Error and drop counts.\sysinfo.exe --network --verbose| `--system` | Hostname, OS, platform, kernel, uptime |

- Total connection count

```

### Process Module

- Total count, running, sleeping| `--cpu` | Model, cores, frequency, per-core usage, load average |# Build the binary

- Top 10 by memory: PID, name, username, memory MB, percentage

- Top 10 by CPU: PID, name, CPU percentage**Top processes only:**

- Process status, create time

```powershell| `--memory` | Total, used, free, cached, swap |

## Building

.\sysinfo.exe --process --format text

**Prerequisites:**

- Go 1.21 or later```| `--disk` | Partitions, usage, I/O stats |go build -o sysinfo.exe .- Lightweight and minimal dependencies



**Build:**

```powershell

cd src## Output Format Details| `--network` | Interfaces, addresses, flags, I/O counters |

go build -o sysinfo.exe .

```



**Cross-compile:**### Pretty Format| `--process` | Total count, top 10 by CPU and memory |- Cross-platform collectors (Windows, Linux, macOS)

```powershell

# Linux

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o sysinfo .

- Color-coded sections with Unicode box-drawing characters| `--smart` | Disk health data (requires admin/root) |

# macOS (Intel)

$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o sysinfo .- Progress bars for CPU, memory, disk usage



# macOS (Apple Silicon)- Organized by module with clear headers# Run with default settings (all modules, pretty output)- Human-friendly and machine-readable output

$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o sysinfo .

```- Temperature warnings in SMART data (yellow >50°C, red >60°C)



## Dependencies- Top 5 processes by memory/CPU displayed**Note:** Specifying any module flag (e.g., `--cpu`) disables `--all`. Combine multiple flags to collect specific modules.



Dependencies are managed via `go.mod`:



- `github.com/shirou/gopsutil/v3` — System metrics collection### Text Format.\sysinfo.exe- Easy to build and run with Go tooling

- `github.com/spf13/cobra` — Command-line interface

- `github.com/fatih/color` — Terminal colors

- `github.com/olekukonko/tablewriter` — Table rendering (used in code, minimal in current output)

- `github.com/yusufpapurcu/wmi` — Windows WMI queries (Windows only)- Plain text, no colors or special characters## Examples



Install:- Suitable for logging or scripting

```powershell

cd src- Same structure as pretty format but simplified```

go mod download

```



## Troubleshooting### JSON Format**Default output (all modules, pretty format):**



**SMART data returns empty:**

- Windows: Run as Administrator

- Linux/macOS: Run with `sudo` and ensure `smartmontools` is installed- Complete structured data```powershell## Quick start

- Not all drives/controllers support SMART

- Omits null/empty fields (`omitempty`)

**Build fails with missing dependencies:**

```powershell- Includes timestamp in ISO 8601 format.\sysinfo.exe

go mod tidy

go mod download- All byte values include human-readable formatted versions

```

```Or run directly with Go:

**Colored output not working:**

- Some terminals don't support ANSI colors## Platform-Specific Notes

- Use `--format text` or `--format json` instead



**"Press Enter to exit" appears when running from terminal:**

- This is intentional when the binary is double-clicked (no terminal attached)### Windows

- Does not occur when run from PowerShell/CMD/bash

- SMART collection uses WMI (`Win32_DiskDrive`, `MSStorageDriver_*` classes)**CPU and memory only, JSON format:**Prerequisites: Go 1.20+ installed.

**No load average on Windows:**

- Load average is a Unix/Linux concept, not available on Windows- Requires Administrator privileges for SMART data



## License- Executable pauses on exit when double-clicked (not from terminal)```powershell



MIT


### Linux.\sysinfo.exe --cpu --memory --format json```powershell

- SMART collection uses `smartctl` (requires `smartmontools` package)

- Supports both ATA and NVMe drives via `smartctl --json````

- Requires root for SMART data

- Load average always available (1, 5, 15 minute intervals)cd srcBuild and run from the repository root:



### macOS**Save full report to file:**

- SMART collection uses `smartctl` (install via `brew install smartmontools`)

- Requires root for SMART data```powershellgo run . --all

- Some disk I/O counters may be limited compared to Linux

.\sysinfo.exe --output report.json --format json

## Data Collected

`````````powershell

### System Module

- Hostname, OS, platform family/version

- Kernel version and architecture

- Uptime (formatted and in seconds)**SMART disk health (requires elevation):**cd src

- Boot time, process count

```powershell

### CPU Module

- Model name, vendor, family, model, stepping# Windows (run as Administrator)## Installationgo build -o sysinfo .

- Physical cores and logical CPUs

- Current, min, max frequency.\sysinfo.exe --smart --disk

- Cache size, microcode version

- Per-core usage percentages (sampled over 1 second)./sysinfo --help

- Load average (Linux/macOS)

- CPU flags# Linux/macOS



### Memory Modulesudo ./sysinfo --smart --disk**From source:**```

- Total, used, free, available (bytes and formatted)

- Usage percentage```

- Cached, buffers, shared memory

- Swap total, used, free, percentage

- Physical RAM modules (placeholder for future WMI/dmidecode integration)

**Network interfaces with verbose output:**

### Disk Module

- Partitions: device, mount point, filesystem type```powershell```powershellOr run directly with the Go tool for development:

- Total, used, free space (bytes, formatted, percentage)

- Inode counts (Linux).\sysinfo.exe --network --verbose

- I/O statistics per disk: read/write counts, bytes, time

- Physical disk information (placeholder)```# Clone the repository

- SMART data: device, model, serial, capacity, health status, temperature, power-on hours, critical attributes



### Network Module

- Interface name, MAC address, MTU**Top processes only:**git clone https://github.com/mayvqt/SysInfo.git```powershell

- IP addresses (all assigned)

- Flags (UP, BROADCAST, LOOPBACK, MULTICAST)```powershell

- Bytes/packets sent and received

- Error and drop counts.\sysinfo.exe --process --format textcd SysInfo\srccd src

- Total connection count

```

### Process Module

- Total count, running, sleepinggo run .

- Top 10 by memory: PID, name, username, memory MB, percentage

- Top 10 by CPU: PID, name, CPU percentage## Output Format Details

- Process status, create time

# Build```

## Building

### Pretty Format

**Prerequisites:**

- Go 1.21 or latergo build -o sysinfo.exe .



**Build:**- Color-coded sections with Unicode box-drawing characters

```powershell

cd src- Progress bars for CPU, memory, disk usage## Usage

go build -o sysinfo.exe .

```- Organized by module with clear headers



**Cross-compile:**- Temperature warnings in SMART data (yellow >50°C, red >60°C)# Install (optional - places binary in GOPATH/bin)

```powershell

# Linux- Top 5 processes by memory/CPU displayed

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o sysinfo .

go install .The binary includes command-line flags and subcommands (see `--help`). Typical usage is to run the collector and print results to stdout. The project is intentionally small so it can be integrated into other tools or pipelines.

# macOS (Intel)

$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o sysinfo .### Text Format



# macOS (Apple Silicon)```

$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o sysinfo .

```- Plain text, no colors or special characters



## Dependencies- Suitable for logging or scripting## Project layout (high level)



Dependencies are managed via `go.mod`:- Same structure as pretty format but simplified



- `github.com/shirou/gopsutil/v3` — System metrics collection**Build for different platforms:**

- `github.com/spf13/cobra` — Command-line interface

- `github.com/fatih/color` — Terminal colors### JSON Format

- `github.com/olekukonko/tablewriter` — Table rendering (used in code, minimal in current output)

- `github.com/yusufpapurcu/wmi` — Windows WMI queries (Windows only)- `src/` — Go module and application source



Install:- Complete structured data

```powershell

cd src- Omits null/empty fields (`omitempty`)```powershell- `internal/collector/` — platform-specific collectors and system probes

go mod download

```- Includes timestamp in ISO 8601 format



## Troubleshooting- All byte values include human-readable formatted versions# Windows- `cmd/` — CLI entrypoint and command wiring



**SMART data returns empty:**

- Windows: Run as Administrator

- Linux/macOS: Run with `sudo` and ensure `smartmontools` is installed## Platform-Specific Notesgo build -o sysinfo.exe .

- Not all drives/controllers support SMART



**Build fails with missing dependencies:**

```powershell### Windows## Contributing

go mod tidy

go mod download- SMART collection uses WMI (`Win32_DiskDrive`, `MSStorageDriver_*` classes)

```

- Requires Administrator privileges for SMART data# Linux

**Colored output not working:**

- Some terminals don't support ANSI colors- Executable pauses on exit when double-clicked (not from terminal)

- Use `--format text` or `--format json` instead

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o sysinfo .Contributions are welcome. Keep changes focused and include tests where appropriate. Follow existing code style and run `go vet`/`go test` before submitting PRs.

**"Press Enter to exit" appears when running from terminal:**

- This is intentional when the binary is double-clicked (no terminal attached)### Linux

- Does not occur when run from PowerShell/CMD/bash

- SMART collection uses `smartctl` (requires `smartmontools` package)

**No load average on Windows:**

- Load average is a Unix/Linux concept, not available on Windows- Supports both ATA and NVMe drives via `smartctl --json`



## License- Requires root for SMART data# macOS## License



MIT- Load average always available (1, 5, 15 minute intervals)


$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o sysinfo .

### macOS

- SMART collection uses `smartctl` (install via `brew install smartmontools`)```This project is licensed under the MIT License. See the `LICENSE` file for details.

- Requires root for SMART data

- Some disk I/O counters may be limited compared to Linux

## Command-Line Options

## Data Collected

### Output Formats

### System Module

- Hostname, OS, platform family/version| Flag | Options | Default | Description |

- Kernel version and architecture|------|---------|---------|-------------|

- Uptime (formatted and in seconds)| `-f, --format` | `json`, `text`, `pretty` | `pretty` | Output format |

- Boot time, process count| `-o, --output` | file path | stdout | Write output to file |

| `-v, --verbose` | - | false | Enable verbose logging |

### CPU Module

- Model name, vendor, family, model, stepping### Module Selection

- Physical cores and logical CPUs

- Current, min, max frequencyBy default, `--all` collects everything. Use specific flags to collect only what you need:

- Cache size, microcode version

- Per-core usage percentages (sampled over 1 second)| Flag | Description |

- Load average (Linux/macOS)|------|-------------|

- CPU flags| `--all` | Collect all information (default) |

| `--system` | System information (OS, hostname, uptime) |

### Memory Module| `--cpu` | CPU model, cores, usage, load average |

- Total, used, free, available (bytes and formatted)| `--memory` | RAM usage, swap, available memory |

- Usage percentage| `--disk` | Disk partitions, usage, I/O stats |

- Cached, buffers, shared memory| `--network` | Network interfaces, IP addresses, stats |

- Swap total, used, free, percentage| `--process` | Running processes and resource usage |

- Physical RAM modules (placeholder for future WMI/dmidecode integration)| `--smart` | SMART disk health data (requires admin/root) |



### Disk Module**Note:** Specifying any individual module flag disables `--all`. Combine flags to select multiple modules.

- Partitions: device, mount point, filesystem type

- Total, used, free space (bytes, formatted, percentage)## Usage Examples

- Inode counts (Linux)

- I/O statistics per disk: read/write counts, bytes, time### 1. Display All Information (Default)

- Physical disk information (placeholder)

- SMART data: device, model, serial, capacity, health status, temperature, power-on hours, critical attributes```powershell

.\sysinfo.exe

### Network Module```

- Interface name, MAC address, MTU

- IP addresses (all assigned)Pretty-printed output to console with all available system data.

- Flags (UP, BROADCAST, LOOPBACK, MULTICAST)

- Bytes/packets sent and received---

- Error and drop counts

- Total connection count### 2. JSON Output



### Process Module```powershell

- Total count, running, sleeping.\sysinfo.exe --format json

- Top 10 by memory: PID, name, username, memory MB, percentage```

- Top 10 by CPU: PID, name, CPU percentage

- Process status, create timeOutputs structured JSON suitable for parsing or integration with monitoring tools.



## Building---



**Prerequisites:**### 3. Save Report to File

- Go 1.21 or later

```powershell

**Build:**.\sysinfo.exe --output system-report.txt

```powershell```

cd src

go build -o sysinfo.exe .Writes pretty-printed report to `system-report.txt`.

```

---

**Cross-compile:**

```powershell### 4. Collect Specific Modules

# Linux

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o sysinfo .```powershell

# CPU and Memory only

# macOS (Intel).\sysinfo.exe --cpu --memory

$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o sysinfo .

# System info and disk with JSON output

# macOS (Apple Silicon).\sysinfo.exe --system --disk --format json

$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o sysinfo .

```# Network information saved to file

.\sysinfo.exe --network --output network-info.json --format json

## Dependencies```



Dependencies are managed via `go.mod`:---



- `github.com/shirou/gopsutil/v3` — System metrics collection### 5. SMART Disk Health (Requires Elevation)

- `github.com/spf13/cobra` — Command-line interface

- `github.com/fatih/color` — Terminal colors**Windows (Run as Administrator):**

- `github.com/olekukonko/tablewriter` — Table rendering (used in code, minimal in current output)

- `github.com/yusufpapurcu/wmi` — Windows WMI queries (Windows only)```powershell

.\sysinfo.exe --smart

Install:```

```powershell

cd src**Linux/macOS (Run with sudo):**

go mod download

``````bash

sudo ./sysinfo --smart

## Troubleshooting```



**SMART data returns empty:**Returns detailed SMART attributes including temperature, reallocated sectors, power-on hours, and health status.

- Windows: Run as Administrator

- Linux/macOS: Run with `sudo` and ensure `smartmontools` is installed---

- Not all drives/controllers support SMART

### 6. Verbose Mode

**Build fails with missing dependencies:**

```powershell```powershell

go mod tidy.\sysinfo.exe --verbose --all

go mod download```

```

Shows progress messages during collection and formatting.

**Colored output not working:**

- Some terminals don't support ANSI colors---

- Use `--format text` or `--format json` instead

### 7. Text Format for Scripting

**"Press Enter to exit" appears when running from terminal:**

- This is intentional when the binary is double-clicked (no terminal attached)```powershell

- Does not occur when run from PowerShell/CMD/bash.\sysinfo.exe --format text --cpu --memory | Select-String "Usage"

```

**No load average on Windows:**

- Load average is a Unix/Linux concept, not available on WindowsPlain text output is easily parseable with grep, findstr, or PowerShell pipelines.



## License---



MIT## Complete Example Workflow


```powershell
# Build the tool
cd src
go build -o sysinfo.exe .

# Quick system overview
.\sysinfo.exe

# Generate JSON report with all data
.\sysinfo.exe --format json --output full-report.json

# Check CPU and memory usage
.\sysinfo.exe --cpu --memory --format pretty

# Monitor disk health (as Administrator)
.\sysinfo.exe --smart --disk --output disk-health.txt

# Network diagnostics
.\sysinfo.exe --network --format json | ConvertFrom-Json | Select -ExpandProperty network
```

## Output Examples

**Pretty Format (default):**
```
╔══════════════════════════════════════════════════════════╗
║                    SYSTEM INFORMATION                    ║
╚══════════════════════════════════════════════════════════╝

Hostname:         DESKTOP-ABC123
OS:               windows
Platform:         Microsoft Windows 11 Pro
Uptime:           2 days, 5 hours, 32 minutes
Processes:        287

╔══════════════════════════════════════════════════════════╗
║                         CPU INFO                         ║
╚══════════════════════════════════════════════════════════╝

Model:            Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
Cores:            8
Logical CPUs:     8
Current Speed:    3600 MHz
Usage:            [12.5% 8.3% 15.2% 9.1% ...]
```

**JSON Format:**
```json
{
  "timestamp": "2025-11-04T10:30:00Z",
  "system": {
    "hostname": "DESKTOP-ABC123",
    "os": "windows",
    "platform": "Microsoft Windows 11 Pro",
    "uptime_seconds": 187920,
    "processes": 287
  },
  "cpu": {
    "model_name": "Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz",
    "physical_cores": 8,
    "logical_cpus": 8,
    "mhz": 3600,
    "usage_percent": [12.5, 8.3, 15.2, 9.1, ...]
  }
}
```

## Project Structure

```
SysInfo/
├── LICENSE
├── README.md
└── src/
    ├── go.mod
    ├── main.go              # Entry point
    ├── cmd/
    │   └── root.go          # CLI command definitions
    └── internal/
        ├── collector/       # Data collection modules
        │   ├── collector.go
        │   ├── cpu.go
        │   ├── disk.go
        │   ├── memory.go
        │   ├── network.go
        │   ├── process.go
        │   ├── smart_*.go   # Platform-specific SMART
        │   └── system.go
        ├── config/          # Configuration management
        │   └── config.go
        ├── formatter/       # Output formatting
        │   ├── formatter.go
        │   ├── pretty.go
        │   └── text.go
        ├── types/           # Data structures
        │   └── types.go
        └── utils/           # Utility functions
            └── format.go
```

## Dependencies

- **[gopsutil](https://github.com/shirou/gopsutil)** — Cross-platform system metrics
- **[cobra](https://github.com/spf13/cobra)** — CLI framework
- **[tablewriter](https://github.com/olekukonko/tablewriter)** — ASCII table rendering
- **[color](https://github.com/fatih/color)** — Terminal color output

Install dependencies:

```powershell
cd src
go mod download
```

## Platform Notes

### Windows
- SMART data requires Administrator privileges
- Run from PowerShell or Command Prompt
- Double-clicking the executable will pause for input before closing

### Linux
- SMART data requires root (`sudo`)
- Some metrics may require `/proc` and `/sys` filesystem access
- Install `smartmontools` for complete SMART support

### macOS
- SMART data requires root (`sudo`)
- Some disk metrics may require additional permissions
- CoreStorage/APFS volumes fully supported

## Troubleshooting

**"Permission denied" when collecting SMART data**
- Run as Administrator (Windows) or with `sudo` (Linux/macOS)

**"No data collected"**
- Ensure you're using `--all` or specific module flags
- Check verbose mode (`--verbose`) for detailed error messages

**Build errors**
- Verify Go 1.21 or later is installed: `go version`
- Run `go mod download` to fetch dependencies

**Output is empty or incomplete**
- Some virtual machines may not expose all hardware details
- Windows Subsystem for Linux (WSL) has limited hardware access

## Contributing

Contributions are welcome! Please ensure:
- Code is formatted with `go fmt`
- All modules are cross-platform compatible
- New features include appropriate error handling

## License

MIT
