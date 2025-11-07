# SysInfo Development Roadmap

This document tracks planned improvements and upgrades for the SysInfo project. Items are organized by priority and implementation status.

## 🚀 High-Impact Features

### Hardware Information
- [x] **GPU Information Module** ✅ Completed
  - GPU model, vendor, memory
  - Temperature and utilization
  - Platform-specific implementations (nvidia-smi, DirectX, Metal)
  - Add `--gpu` flag and include in `--all`
  - Windows NVIDIA enhancement with nvidia-smi for advanced metrics
  - Files: `gpu.go`, `gpu_linux.go`, `gpu_darwin.go`, `gpu_windows.go`
  
- [x] **Complete Physical Disk Collection** ✅ Completed
  - Linux: `lsblk -J` for comprehensive disk info, fallback to `/sys/block`
  - macOS: `diskutil list/info -plist` for disk details
  - Windows: WMI `MSFT_PhysicalDisk` (Windows 8+) and `Win32_DiskDrive` (fallback)
  - Detects disk type (HDD, SSD, NVMe), interface (SATA, NVMe, USB), size, model, serial
  - RPM detection for HDDs, removable media detection
  - Files: `disk_linux.go`, `disk_darwin.go`, `disk_windows.go`
  
- [ ] **Battery Information Module**
  - Capacity, charge level, state
  - Time remaining, cycle count
  - Health percentage
  - Laptop and UPS support

- [ ] **Temperature Sensors Module**
  - CPU package temperature
  - Motherboard sensors
  - Fan speeds (RPM)
  - Voltage rails
  - Platform: `sensors` (Linux), `smc` (macOS), WMI (Windows)

### Data Analysis & Intelligence

- [ ] **Enhanced SMART Analysis**
  - Predictive failure alerts
  - SSD wear leveling analysis with lifespan estimation
  - Historical tracking (SQLite storage)
  - Webhook alerts for critical disk health
  
- [ ] **Historical Data & Trends**
  - `--record` flag to store to local DB
  - `--history 7d` to show 7-day trends
  - `--chart cpu,mem` for ASCII charts in terminal
  - Use SQLite or BoltDB for local storage
  
- [ ] **Diff Mode**
  - `--diff baseline.json current.json`
  - Show changes in CPU/memory usage, disk space
  - SMART attribute deltas
  - New/removed processes comparison

### Configuration & Customization

- [x] **Configuration File Support** ✅ Completed
  - `.sysinforc` or `.sysinfo.yaml` support
  - Persist default output format
  - Default modules to collect
  - SMART monitoring configuration
  - Process top count configuration
  - Display preferences (ASCII mode)
  - Files: `loader.go`, `loader_test.go`, `.sysinforc.example`

- [ ] **Custom Formatters/Templates**
  - `--template custom.tmpl` using Go templates
  - `--jq '.cpu | {cores, usage}'` built-in JQ filtering
  - User-defined output formats

- [ ] **Alerting Rules**
  - YAML-based alert configuration
  - Conditions on metrics (CPU, disk, memory thresholds)
  - Actions: notify, email, webhook
  - Example: high CPU, low disk space alerts

## 🎯 Medium-Impact Features

### Output & Export

- [ ] **Additional Export Formats**
  - CSV for spreadsheet analysis
  - Prometheus metrics format
  - HTML report with Chart.js visualizations
  - XML for enterprise tools
  
- [ ] **Comparison with Baseline**
  - `--baseline > baseline.json` to create baseline
  - `--compare baseline.json` to show deltas
  - Highlight significant changes

### Monitoring & Remote Access

- [ ] **Remote Monitoring**
  - `--serve :8080` for HTTP API mode
  - REST API endpoints for programmatic access
  - `--push https://...` to push to monitoring server
  - WebSocket support for live updates
  
- [ ] **Interactive Mode**
  - TUI using bubbletea or termui
  - Navigate modules, drill down into details
  - Live updates in interactive mode
  - Keyboard shortcuts for common actions

### Developer Experience

- [ ] **Plugins System**
  - Allow custom collectors via plugins
  - `~/.sysinfo/plugins/` directory
  - `--load-plugin docker` for custom modules
  - Plugin API documentation

## 🔧 Code Quality Improvements

### Testing & Validation

- [ ] **Enhanced Unit Test Coverage**
  - Mock interfaces for `gopsutil` calls
  - Table-driven tests for parsers
  - Snapshot testing for formatters
  - Platform-specific test fixtures
  - Target 80%+ coverage

- [ ] **Configuration Validation**
  - `--validate` flag to check environment
  - Check required tools installed
  - Verify permissions (SMART needs root)
  - Ensure output directory writable
  - Validate config file syntax

### Logging & Error Handling

- [ ] **Structured Logging Framework**
  - Replace `fmt.Fprintf` with `log/slog` or `zerolog`
  - Log levels: DEBUG, INFO, WARN, ERROR
  - JSON log output option
  - Log aggregation support

- [ ] **Graceful Degradation Messages**
  - Helpful hints when tools are missing
  - Install commands for required utilities
  - Clear permission error messages
  - Better error context

- [ ] **Error Handling Improvements**
  - Fix potential panic in process collection (empty status array)
  - Handle edge cases in parsers
  - Retry logic for transient failures
  - Timeout handling for slow operations

## 🎨 User Experience

### UI/UX Enhancements

- [ ] **Progress Indicators**
  - Show progress for slow operations (SMART collection)
  - Spinner or progress bar for long-running tasks
  - Estimated time remaining
  
- [ ] **Large Output Handling**
  - Pagination for monitor mode
  - Configurable limits (`--top-processes 5`)
  - Terminal size detection
  - Responsive formatting

- [ ] **Unicode/Emoji Fallback**
  - Detect terminal capabilities
  - Fallback to ASCII for limited terminals
  - `--ascii` flag to force ASCII mode
  - Better Windows console support

### Quick Wins

- [ ] **Version Command Enhancements**
  - `--version` shows detailed build info
  - Include commit hash, build time
  - Show Go version used for build
  
- [ ] **Quiet Mode**
  - `--quiet` flag for silent operation
  - Only output data, suppress stderr
  - Useful for scripting
  
- [ ] **JSON Output Filtering**
  - `--only cpu,memory` to filter JSON output
  - Include only requested modules in output
  - Reduce output size for specific use cases
  
- [ ] **Human-Readable Timestamps**
  - `--timestamp-format relative` for "2 hours ago"
  - ISO8601, Unix, or custom formats
  - Timezone support

## 📋 Implementation Priority

### Phase 1 (Current Sprint)
1. ✅ Create roadmap document
2. ✅ GPU Information Module - **COMPLETED**
3. ✅ Complete Physical Disk Collection - **COMPLETED**

### Phase 2 (Next Sprint)
1. ✅ Configuration File Support - **COMPLETED**
2. Enhanced SMART Analysis
3. Battery Information Module

### Phase 3 (Future)
1. Logging Framework
2. Additional Export Formats
3. Historical Data & Trends

### Phase 4 (Long-term)
1. Remote Monitoring
2. Interactive Mode
3. Plugins System

---

## Contributing

When implementing features from this roadmap:

1. Create a feature branch: `git checkout -b feature/gpu-info`
2. Update the checkbox in this file: `- [x] GPU Information Module`
3. Implement with tests
4. Update documentation (README.md, copilot-instructions.md)
5. Submit PR with reference to roadmap item

## Notes

- Check items off with `[x]` when merged to main
- Add `(In Progress)` annotation for active development
- Link to issues/PRs in GitHub for tracking
- Update priority based on user feedback and community requests
