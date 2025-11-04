# [![CI](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml/badge.svg)](https://github.com/mayvqt/sysinfo/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Release](https://img.shields.io/github/v/release/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/releases) [![Go Version](https://img.shields.io/github/go-mod/go-version/mayvqt/sysinfo)](https://github.com/mayvqt/sysinfo/blob/main/go.mod)

# SysInfo

SysInfo is a small, cross-platform command-line tool and Go library for collecting system information and health data. It gathers CPU, memory, disk (including SMART where supported), network, process, and system metadata and outputs results in human-friendly and machine-readable formats.

## Highlights

- Modular collectors: CPU, memory, disk (SMART optional), network, processes, and system metadata
- Output formats: `pretty`, `text`, and `json`
- Single binary for easy deployment and automation
- Cross-platform: Linux, macOS and Windows (collection varies by platform)

## Quickstart

Prerequisites
- Go 1.21 or later (building from source)

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
# Pretty (default)
.\\sysinfo.exe

# CPU + memory only, JSON
.\\sysinfo.exe --cpu --memory --format json

# Save full JSON report
.\\sysinfo.exe --format json --output full-report.json
```

## Flags

- `--all` (default): collect all modules
- `--system`: host/OS/kernel/uptime
- `--cpu`: CPU info and per-core usage
- `--memory`: memory and swap info
- `--disk`: partitions and disk info (SMART optional)
- `--network`: interface statistics
- `--process`: process summaries
- `--format`: `pretty|text|json`
- `--output`: write output to file
- `--verbose`: enable verbose logging

Run `--help` for the complete flag list and examples.

## Platform notes

- SMART: requires elevation (Administrator/root) and `smartmontools` on Linux/macOS. On Windows SMART data may be available via WMI but also requires elevation.
- Load averages: not available on Windows (Unix concept).

## Development

- Module path: `module github.com/mayvqt/sysinfo` (module at repository root)
- Use `go mod tidy` to reconcile dependencies and `go mod download` to prefetch.

Dev commands

```powershell
# fetch deps
go mod download

# vet + tests
go vet ./...
go test ./... -v

# optional linters
# golangci-lint run
```

## CI & Releases

Workflows are in `.github/workflows/`. The release workflow builds cross-platform binaries and places them in `releases/`.

## Contributing

- Run `go fmt` before submitting PRs
- Keep changes small and cross-platform where possible
- Include tests for new behavior

## License

MIT
