# [![CI](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml/badge.svg)](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Release](https://img.shields.io/github/v/release/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/releases) [![Go Version](https://img.shields.io/github/go-mod/go-version/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/blob/main/go.mod)

# SysInfo

SysInfo is a comprehensive, cross-platform command-line tool and Go library for collecting detailed system information and health data. It gathers CPU, memory (with physical module details), disk (including comprehensive SMART data with health assessment), network, process, and system metadata and outputs results in human-friendly and machine-readable formats.

## Highlights

- **Comprehensive Data Collection**: CPU, memory modules, disk (with 70+ SMART attributes), network, processes, GPU, battery, and system metadata
- **Advanced SMART Analysis**: Predictive failure detection, historical tracking with trend analysis, and webhook alerting system
- **GPU Monitoring**: Detailed GPU information including temperature, utilization, memory usage, and power draw (NVIDIA, AMD, Intel)
- **Battery Monitoring**: Comprehensive battery information including charge level, health, time remaining, cycle count, temperature, and power consumption (laptops and UPS devices)
- **Multiple Output Formats**: `pretty`, `text`, and `json`
- **Full System Dump**: Single command to capture everything to JSON for analysis
- **Configuration File Support**: YAML/TOML config with sensible defaults
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

# Battery information for laptops
.\\sysinfo.exe --battery

# Comprehensive SMART data with health assessment
.\\sysinfo.exe --smart --format json

# Enhanced SMART analysis (subcommand)
sudo sysinfo smart analyze

# View SMART history and trends
sudo sysinfo smart history --period 30d

# Quick health check (for monitoring scripts)
sudo sysinfo smart check

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
- `--battery`: battery information including charge level, health, time remaining, and cycle count

### SMART Analysis Options
Use the `smart` subcommand for advanced disk health monitoring:
- `sysinfo smart analyze`: Deep SMART analysis with failure prediction, SSD wear tracking, and history storage
- `sysinfo smart history`: View historical trends, temperature patterns, and wear rate analysis
- `sysinfo smart check`: Quick health check for all drives (no history storage, perfect for monitoring scripts)

**Flags:**
- `--db <path>`: Custom database path for history storage
- `--period <duration>`: History period for `history` command (e.g., 1h, 24h, 7d, 30d, default: 7d)
- `--alerts`: Enable webhook notifications for critical events (configure in config file)
- `--verbose`: Show detailed progress and diagnostics

### Output Options
- `--format`, `-f`: output format: `pretty|text|json` (default: pretty)
- `--output`, `-o`: write output to file instead of stdout
- `--verbose`, `-v`: enable verbose logging
- `--full-dump`: collect ALL system info and save to `sysinfo_dump.json` (includes everything)
- `--config`: specify custom config file path (default: auto-detect)

Run `--help` for the complete flag list and examples.

## Configuration File

SysInfo supports configuration files for persistent settings. Create a configuration file in one of these locations:
- `./.sysinforc` (current directory)
- `./.sysinfo.yaml` (current directory)
- `~/.config/sysinfo/config.yaml` (user config directory)
- `~/.sysinforc` (home directory)

**Example Configuration** (see `.sysinforc.example`):
```yaml
# Default output format: json, text, or pretty
format: pretty

# Enable verbose output
verbose: false

# Default modules to collect (when no flags are specified)
modules:
  system: true
  cpu: true
  memory: true
  disk: true
  network: true
  process: true
  smart: false  # Requires root/admin
  gpu: true
  battery: true

# SMART monitoring configuration
smart:
  enable_alerts: false
  database_path: "~/.config/sysinfo/smart.db"
  alert_thresholds:
    temperature_critical: 70  # Celsius
    temperature_warning: 60   # Celsius
    wear_critical: 90.0       # SSD wear percentage
    wear_warning: 80.0        # SSD wear percentage

# SMART webhook alerts (requires enable_alerts: true)
smart_alerts:
  enabled: false
  webhook_url: "https://your-webhook.example.com/alerts"
  min_level: "WARNING"  # INFO, WARNING, or CRITICAL
  cooldown: 60          # minutes between alerts for same device/issue
  timeout: 10           # seconds (webhook request timeout)

# Process monitoring
process:
  top_count: 10  # Number of top processes to show

# Display preferences
display:
  use_ascii: false  # Force ASCII instead of Unicode
```

**Note**: Command-line flags take precedence over configuration file settings.

## SMART Data Features

### Requirements

**For basic `--smart` flag** (raw SMART data collection):
- Linux/macOS: `smartmontools` package required
  ```bash
  # Debian/Ubuntu
  sudo apt install smartmontools
  
  # macOS
  brew install smartmontools
  
  # RHEL/CentOS/Fedora
  sudo dnf install smartmontools
  ```
- Windows: No external dependencies (uses built-in WMI)
- Elevated privileges required on all platforms (`sudo` on Linux/macOS, Administrator on Windows)

**If `smartmontools` is not installed**: 
- The `--smart` flag will silently skip SMART collection (no error)
- Only disk/partition information will be shown
- Enhanced analysis commands (`sysinfo smart analyze/history/check`) will show: "No SMART data available"

**For enhanced SMART analysis** (`sysinfo smart analyze/history/check`):
- Same `smartmontools` requirement as above
- SQLite database for history (automatically created at `~/.config/sysinfo/smart.db`)
- Go-native SQLite driver (no external dependencies)

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

### Enhanced SMART Analysis

SysInfo includes advanced SMART analysis features for predictive maintenance:

**Predictive Failure Detection**:
```bash
# Perform deep SMART analysis with failure prediction
sudo sysinfo smart analyze

# Enable webhook alerts for critical issues
sudo sysinfo smart analyze --alerts

# Use custom database location
sudo sysinfo smart analyze --db /var/lib/sysinfo/smart.db

# Example output:
# /dev/sda
# ======================================================================
# Overall Health: ✓ GOOD
# Failure Risk: 5.2%
#
# SSD Wear Analysis:
#   Status: GOOD
#   Remaining Life: 87.3%
#   Percent Used: 12.7%
#   Estimated Remaining: 1,247 days (3.4 years)
#
# Issues Detected: 0
#
# Recommendations:
#   • Drive health is good - continue monitoring
#   • Schedule backup within 90 days
```

**Historical Tracking**:
```bash
# View SMART history and trends (default: 7 days)
sudo sysinfo smart history

# View last 30 days
sudo sysinfo smart history --period 30d

# Use custom database location
sudo sysinfo smart history --db /var/lib/sysinfo/smart.db

# Example output:
# SMART History (Last 7d)
# ======================================================================
#
# Device: /dev/sda
# ----------------------------------------------------------------------
#   Recent Records: 42
#     2025-11-07 14:30 | Health: GOOD     | Temp:  42°C | Issues: 0 (Critical: 0)
#     2025-11-07 08:15 | Health: GOOD     | Temp:  38°C | Issues: 0 (Critical: 0)
#     ...
#
#   Trends:
#     Temperature: stable (Avg: 40.2°C, Min: 35°C, Max: 48°C)
#     Health Status: stable
#     SSD Wear Rate: 0.0023% per day
#     Estimated End of Life: 2028-04-15
```

**Alert System** (Webhook Notifications):
```bash
# Enable alerts for critical disk events
sudo sysinfo smart analyze --alerts

# Quick health check (useful for monitoring scripts)
sudo sysinfo smart check

# Configure webhook in config file (~/.config/sysinfo/config.yaml):
# smart_alerts:
#   enabled: true
#   webhook_url: "https://your-webhook.example.com/alerts"
#   min_level: "WARNING"  # INFO, WARNING, or CRITICAL
#   cooldown: 60  # minutes between alerts for same device
#   timeout: 10   # seconds
```

**Alert Webhook Payload** (JSON):
```json
{
  "timestamp": "2025-11-07T14:30:00Z",
  "device": "/dev/sda",
  "level": "CRITICAL",
  "health_status": "CRITICAL",
  "failure_probability": 87.2,
  "issues": [
    {
      "severity": "CRITICAL",
      "code": "REALLOCATED_SECTORS",
      "description": "Drive has 150 reallocated sectors",
      "attribute_id": 5
    }
  ],
  "recommendations": [
    "URGENT: Back up all data immediately",
    "Schedule drive replacement as soon as possible"
  ]
}
```

**Features**:
- **Health Classification**: GOOD, WARNING, CRITICAL, FAILING, UNKNOWN
- **Predictive Analysis**: Calculates failure probability (0-100%) based on SMART attributes, temperature, and wear
- **Temperature Monitoring**: Configurable warning (60°C) and critical (70°C) thresholds with trend tracking
- **SSD Lifespan Estimation**: Calculates remaining lifetime based on wear metrics and usage patterns
- **Trend Analysis**: Tracks temperature changes, health degradation, and wear rate over time
- **Webhook Alerts**: Configurable JSON notifications for critical events (failure predictions, high wear, temperature issues)
- **SQLite History**: Automatic tracking of SMART metrics with configurable retention and cleanup
- **Intelligent Alerting**: Cooldown periods prevent alert fatigue, minimum severity levels filter noise
- **Zero Configuration**: Works out-of-box with sensible defaults, fully customizable via config file

**Analysis Capabilities**:
- Critical attribute detection (reallocated sectors, pending sectors, uncorrectable errors)
- Pre-fail attribute threshold monitoring with early warning
- SSD-specific wear leveling and program/erase cycle tracking
- Temperature anomaly detection and historical comparison
- Multi-vendor SMART attribute support (WD, Seagate, Samsung, Intel, Crucial, etc.)

**Custom Database Path**:
```bash
# Use custom database location
sudo sysinfo smart analyze --db /var/lib/sysinfo/smart.db
```

**Quick Health Check**:
```bash
# Fast health check for all drives (no history storage)
sudo sysinfo smart check

# Example output:
# ✓ /dev/sda            GOOD
# ⚠ /dev/nvme0n1        WARNING  [SSD Life: 15.3%]
# ✗ /dev/sdb            CRITICAL [FAILURE PREDICTED: 87.2%]
#
# ⚠ Issues detected - run 'sysinfo smart analyze' for details
```

## Integration & Automation

SysInfo is designed for easy integration into monitoring systems and automation workflows:

**Cron/Scheduled Tasks**:
```bash
# Daily SMART analysis with webhook alerts (crontab)
0 2 * * * /usr/local/bin/sysinfo smart analyze --alerts >> /var/log/sysinfo.log 2>&1

# Hourly quick health check
0 * * * * /usr/local/bin/sysinfo smart check || echo "SMART check failed" | mail -s "Disk Alert" admin@example.com
```

**Monitoring Scripts**:
```bash
#!/bin/bash
# Check disk health and exit with error code if issues found
sysinfo smart check --format json > /tmp/smart_status.json
if [ $? -ne 0 ]; then
    # Send alert, page admin, etc.
    curl -X POST https://monitoring.example.com/alert \
        -H "Content-Type: application/json" \
        -d @/tmp/smart_status.json
    exit 1
fi
```

**JSON API Integration**:
```bash
# Export full system info as JSON for ingestion
sysinfo --all --format json --output /var/www/api/sysinfo.json

# SMART analysis JSON output for monitoring systems
sysinfo smart analyze --format json | jq '.results[] | select(.overall_health != "GOOD")'
```

**Docker/Container Monitoring**:
```dockerfile
# Include in container health checks
FROM alpine:latest
COPY sysinfo /usr/local/bin/
HEALTHCHECK --interval=5m --timeout=10s \
  CMD sysinfo smart check || exit 1
```

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

## Documentation

- **[Configuration Guide](docs/CONFIGURATION.md)** - Complete configuration file reference and examples
- **[Roadmap](docs/ROADMAP.md)** - Development roadmap and completed features
- **[SMART Monitoring](docs/ROADMAP.md#phase-2-advanced-features)** - Enhanced SMART analysis, historical tracking, and alerts

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
