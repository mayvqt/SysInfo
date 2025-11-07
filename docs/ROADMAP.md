# SysInfo Roadmap

## Completed ✅

- **GPU Information** - Temperature, utilization, memory (NVIDIA, AMD, Intel)
- **Physical Disk Details** - Type detection (HDD/SSD/NVMe), interface, RPM
- **Battery Monitoring** - Charge level, health, cycle count, temperature
- **Configuration File** - YAML/TOML support with `.sysinforc`
- **Enhanced SMART Analysis** - Predictive failure detection, historical tracking, webhook alerts

## Planned Features

### High Priority
- **Temperature Sensors** - CPU, motherboard, fan speeds
- **Historical Data & Trends** - SQLite storage, ASCII charts, trend analysis
- **Alerting Rules** - YAML-based conditions for CPU/disk/memory thresholds

### Medium Priority
- **Additional Export Formats** - CSV, Prometheus, HTML
- **Remote Monitoring** - HTTP API mode, REST endpoints
- **Interactive TUI** - Live updates, keyboard navigation

### Future
- **Plugins System** - Custom collectors via plugin API
- **Structured Logging** - JSON logs, log levels
- **Diff Mode** - Compare snapshots, show deltas
