# [![CI](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml/badge.svg)](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Release](https://img.shields.io/github/v/release/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/releases) [![Go Version](https://img.shields.io/github/go-mod/go-version/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/blob/main/go.mod)

# SysInfo

SysInfo is a comprehensive, cross-platform command-line tool and Go library for collecting detailed system information and health data. It gathers CPU, memory (with physical module details), disk (including comprehensive SMART data with health assessment), network, process, and system metadata and outputs results in human-friendly and machine-readable formats.

## Highlights

- **Comprehensive Data Collection**: CPU, memory modules, disk (with 70+ SMART attributes), network, processes, GPU, and system metadata
- **GPU Monitoring**: Detailed GPU information including temperature, utilization, memory usage, and power draw (NVIDIA, AMD, Intel)
- **SMART Health Monitoring**: Professional-grade disk health assessment with failure prediction and SSD wear tracking
- **Multiple Output Formats**: `pretty`, `text`, and `json`
- **Live Monitoring Mode**: Real-time system stats that update in place (like htop/top)
- **Full System Dump**: Single command to capture everything to JSON for analysis
- **Single Binary**: Easy deployment and automation
- **Cross-platform**: Linux, macOS, and Windows (with platform-optimized collectors)

## Quickstart

Prerequisites
- Go 1.24 or later (building from source)
- For SMART data on Linux/macOS: `smartmontools` (`apt install smartmontools` or `brew install smartmontools`)

Build from repository root:

```powershell
go build -o sysinfo.exe .
.\\sysinfo.exe --help
```

Install with `go install`:

```powershell
go install github.com/mayvqt/sysinfo@latest
```

Examples

```powershell
# Pretty output (default)
.\\sysinfo.exe

# CPU + memory only, JSON format
.\\sysinfo.exe --cpu --memory --format json

# Live monitoring mode - updates every 2 seconds (like task manager)
.\\sysinfo.exe --cpu --memory --monitor

# Live monitoring with custom interval
.\\sysinfo.exe --cpu --memory --process --monitor --interval 5

# Comprehensive SMART data with health assessment
.\\sysinfo.exe --smart --format json

# Full system dump - captures EVERYTHING to JSON file
.\\sysinfo.exe --full-dump

# Save custom report
.\\sysinfo.exe --cpu --memory --disk --smart --format json --output report.json
```

## Flags

### Module Selection
- `--all` (default): collect all modules
- `--system`: host/OS/kernel/uptime/process count
- `--cpu`: CPU info, per-core usage, flags, microcode
- `--memory`: memory/swap info + physical RAM module details (type, speed, manufacturer)
- `--disk`: partitions, physical disks, and I/O stats
- `--network`: interface statistics and connection counts
- `--process`: process summaries (top by CPU and memory)
- `--smart`: comprehensive SMART disk data with health assessment (requires elevation)
- `--gpu`: GPU information including temperature, utilization, memory, and power draw

### Output Options
- `--format`, `-f`: output format: `pretty|text|json` (default: pretty)
- `--output`, `-o`: write output to file instead of stdout
- `--verbose`, `-v`: enable verbose logging
- `--full-dump`: collect ALL system info and save to `sysinfo_dump.json` (includes everything)

### Monitor Mode
- `--monitor`, `-m`: enable live monitoring mode (continuous updates)
- `--interval`, `-i`: update interval in seconds for monitor mode (default: 2)

Run `--help` for the complete flag list and examples.

## SMART Data Features

SysInfo provides professional-grade SMART monitoring with:

**Comprehensive Attributes** (70+ tracked):
- All standard HDD attributes (Read Error Rate, Reallocated Sectors, Power On Hours, etc.)
- SSD-specific attributes (Wear Leveling, Program/Erase Fail Counts, Life Remaining, etc.)
- Vendor-specific attributes (WD, Seagate, Samsung, Intel, Crucial, Micron, SandForce, etc.)

**Detailed Information Per Drive**:
- Firmware version, rotation rate (RPM for HDD, 0 for SSD)
- Disk geometry and form factor
- Per-attribute details: ID, current value, worst value, threshold, raw value
- Human-readable value formatting (hours/days, GB/TB, temperature, percentages)

**Health Assessment**:
- Overall status (PASS/WARN/FAIL)
- Failing and warning attribute lists
- SSD wear level and life remaining
- Temperature status (NORMAL/WARM/HIGH/CRITICAL)
- Critical attribute monitoring (reallocated sectors, pending sectors, uncorrectable errors, etc.)

## GPU Information Features

SysInfo provides comprehensive GPU monitoring with cross-platform support:

**Supported GPU Vendors**:
- NVIDIA (via nvidia-smi on Linux/Windows)
- AMD (via rocm-smi on Linux, WMI on Windows)
- Intel (via lspci on Linux, WMI on Windows)
- Apple Silicon (via system_profiler on macOS)

**Information Collected**:
- GPU model, vendor, and driver version
- Memory total, used, and utilization percentage
- Temperature monitoring
- GPU and memory utilization
- Power draw and power limit
- Clock speeds (GPU and memory)
- Fan speed percentage
- PCI bus information

**Platform Notes**:
- **Linux**: Best support with nvidia-smi (NVIDIA) or rocm-smi (AMD), falls back to lspci for basic info
- **macOS**: Uses system_profiler, full support for Apple Silicon and discrete GPUs
- **Windows**: Uses WMI for all vendors, automatically enhanced with nvidia-smi for detailed NVIDIA stats (temperature, utilization, power, clocks, fan speed)

## Platform Notes

**Windows**:
- SMART data via WMI (requires Administrator)
- Physical memory module info via WMI
- Full support for all features

**Linux**:
- SMART data requires `smartmontools` and root/sudo
- Memory module info requires dmidecode (future enhancement)
- Load averages fully supported

**macOS**:
- SMART data requires `smartmontools` (brew install) and sudo
- Memory module info via system_profiler (future enhancement)
- NVMe/Apple Silicon SSD support included

## Development

- Module path: `module github.com/mayvqt/sysinfo` (module at repository root)
- Use `go mod tidy` to reconcile dependencies and `go mod download` to prefetch.

Dev commands

```powershell
# Fetch dependencies
go mod download

# Run tests with coverage
go test ./... -v -cover

# Run linters
go vet ./...

# Build
go build -o sysinfo.exe .

# Run comprehensive dump
.\\sysinfo.exe --full-dump
```

## CI & Releases

Workflows are in `.github/workflows/`:
- **CI**: Automated testing on push/PR with coverage reporting (Codecov)
- **Release**: Cross-platform binaries (Windows, Linux, macOS) on version tags

## Contributing

- Run `go fmt` before submitting PRs
- Keep changes small and cross-platform where possible
- Include tests for new behavior
- Update README for new features

## License

MIT
