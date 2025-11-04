[![CI](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml/badge.svg)](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Release](https://img.shields.io/github/v/release/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/releases)

[![Go Version](https://img.shields.io/github/go-mod/go-version/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/blob/main/src/go.mod)

SysInfo
======

Overview
--------
SysInfo is a lightweight, cross-platform command-line utility and Go library for collecting system information and metrics. It gathers CPU, memory, disk, network, process and (platform-specific) SMART data and emits results in human-friendly or machine-readable formats. SysInfo is intended to be easy to build, run and integrate into automation or monitoring workflows.

Features
--------

Core functionality
~~~~~~~~~~~~~~~~~~
- Modular collectors for CPU, memory, disk, network, processes and SMART (platform-specific)
- Multiple output formats: `pretty` (human), `text`, and `json` for automation
- Selective module collection: collect everything (`--all`) or only specific modules (`--cpu`, `--disk`, etc.)
- File output support (`-o, --output`) for saving reports

Reliability & operations
~~~~~~~~~~~~~~~~~~~~~~~~
- Single static binary, minimal runtime dependencies
- Verbose/progress mode to follow collection steps
- Designed for low resource usage so it can run in containers, VMs or edge devices

Configuration & deployment
~~~~~~~~~~~~~~~~~~~~~~~~~~
- Easy build with Go and simple cross-compile examples included
- Environment-based and CLI configuration patterns (see flags)
- Platform-specific notes for SMART collection and permissions

Quickstart (Windows)
---------------------
Open PowerShell, then:

```powershell
cd src
go build -o sysinfo.exe .
.\sysinfo.exe --help
```

Or install with `go install` (Go 1.20+):

```powershell
go install github.com/mayvqt/sysinfo@latest
```

Basic usage
-----------
Run the default (all modules, pretty output):

```powershell
.\sysinfo.exe
```

CPU + Memory only, JSON:

```powershell
.\sysinfo.exe --cpu --memory --format json
```

Save a full report to a file:

```powershell
.\sysinfo.exe --output report.json --format json
```

Common flags
------------
- `--all` (default) : collect all modules
- `--system` : host, OS, kernel, uptime
- `--cpu` : CPU model, cores, frequencies, per-core usage
- `--memory` : totals, used/available, swap
- Partitions: device, mount point, filesystem type
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
