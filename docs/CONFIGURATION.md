# SysInfo Configuration Guide

This guide covers everything you need to know about configuring SysInfo using configuration files.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Configuration File Locations](#configuration-file-locations)
- [Configuration Options](#configuration-options)
- [Use Cases & Examples](#use-cases--examples)
- [Best Practices](#best-practices)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)

## Overview

SysInfo supports YAML-based configuration files that allow you to:

- **Set default behavior** - Avoid typing the same flags repeatedly
- **Standardize monitoring** - Deploy consistent configs across multiple systems
- **Simplify automation** - Make scripts cleaner by moving settings to config
- **Customize per environment** - Different configs for dev/staging/production
- **Enable advanced features** - Configure alerting thresholds, webhooks, and more

**Priority**: Command-line flags always override config file settings.

## Quick Start

### 1. Create a basic config file

```bash
# Copy the example config
cp .sysinforc.example ~/.sysinforc

# Or create your own
cat > ~/.sysinforc << 'EOF'
format: pretty
modules:
  cpu: true
  memory: true
  disk: true
EOF
```

### 2. Run with config

```bash
# Uses ~/.sysinforc automatically
sysinfo

# Use a specific config file
sysinfo --config /path/to/custom.yaml

# Override config with CLI flags
sysinfo --format json  # Overrides format from config
```

### 3. Verify it's working

```bash
# With verbose flag, you'll see config being loaded
sysinfo -v
```

## Configuration File Locations

SysInfo searches for configuration files in this order (first found wins):

1. **Current directory**: `./.sysinforc` or `./.sysinfo.yaml`
2. **User config directory**: `~/.config/sysinfo/config.yaml`
3. **Home directory**: `~/.sysinforc`
4. **Custom path**: Via `--config /path/to/file.yaml` flag

### Platform-Specific Paths

**Linux/macOS**:
- `~/.sysinforc`
- `~/.config/sysinfo/config.yaml`
- `./sysinforc`

**Windows**:
- `%USERPROFILE%\.sysinforc`
- `%APPDATA%\sysinfo\config.yaml`
- `.\.sysinforc`

## Configuration Options

### Complete Configuration Reference

```yaml
# Output format: json, text, or pretty
format: pretty

# Output file path (leave empty for stdout)
output_file: /var/log/sysinfo.json

# Enable verbose output
verbose: false

# Default modules to collect (when no flags specified)
modules:
  system: true
  cpu: true
  memory: true
  disk: true
  network: true
  process: true
  smart: false   # Requires root/admin
  gpu: true

# SMART monitoring configuration
smart:
  # Enable alerts for SMART issues
  enable_alerts: false
  
  # Alert thresholds (Celsius)
  alert_thresholds:
    temperature_critical: 70
    temperature_warning: 60
  
  # Webhook URL for alerts (optional)
  webhook_url: https://monitoring.example.com/webhook

# Process monitoring configuration
process:
  # Number of top processes to show
  top_count: 10

# Display preferences
display:
  # Force ASCII output instead of Unicode box drawing
  use_ascii: false
```

### Option Details

#### `format`
- **Type**: String
- **Values**: `json`, `text`, `pretty`
- **Default**: `pretty`
- **Description**: Default output format. CLI `-f/--format` flag overrides.

#### `output_file`
- **Type**: String
- **Default**: empty (stdout)
- **Description**: Default output file path. CLI `-o/--output` flag overrides.
- **Supports**: Absolute and relative paths

#### `verbose`
- **Type**: Boolean
- **Default**: `false`
- **Description**: Enable verbose logging. CLI `-v/--verbose` flag overrides.

#### `modules.*`
- **Type**: Boolean
- **Default**: All `true` except `smart: false`
- **Description**: Which modules to collect by default.
- **Note**: CLI module flags (e.g., `--cpu`) override these settings.

#### `smart.enable_alerts`
- **Type**: Boolean
- **Default**: `false`
- **Description**: Enable SMART health alerting (future feature).

#### `smart.alert_thresholds.temperature_critical`
- **Type**: Integer (Celsius)
- **Default**: `70`
- **Description**: Temperature threshold for critical alerts.

#### `smart.alert_thresholds.temperature_warning`
- **Type**: Integer (Celsius)
- **Default**: `60`
- **Description**: Temperature threshold for warning alerts.

#### `smart.webhook_url`
- **Type**: String
- **Default**: empty
- **Description**: HTTP endpoint to POST alerts (future feature).

#### `process.top_count`
- **Type**: Integer
- **Default**: `10`
- **Description**: Number of top processes to display (future feature).

#### `display.use_ascii`
- **Type**: Boolean
- **Default**: `false`
- **Description**: Force ASCII output for limited terminals (future feature).

## Use Cases & Examples

### 1. System Administrator - Daily Health Checks

**Scenario**: Run daily checks on production servers without typing flags.

```yaml
# /etc/sysinfo/config.yaml
format: json
output_file: /var/log/sysinfo/health-check.json
verbose: false

modules:
  system: true
  cpu: true
  memory: true
  disk: true
  smart: true
  network: true
```

**Usage**:
```bash
# Crontab entry
0 9 * * * /usr/local/bin/sysinfo

# Manual check with different output
sysinfo --format pretty  # Override format for human viewing
```

**Benefits**:
- Consistent daily logging
- No flags needed in cron jobs
- Easy to override when needed

### 2. Developer - JSON Output for Scripts

**Scenario**: Developer frequently queries system info in shell scripts.

```yaml
# ~/.sysinforc
format: json
verbose: false

modules:
  cpu: true
  memory: true
  gpu: true
  disk: true
```

**Usage**:
```bash
# Get memory usage percentage
sysinfo | jq '.memory.used_percent'

# Check if GPU temperature is high
temp=$(sysinfo | jq '.gpu.gpus[0].temperature')
if [ $temp -gt 80 ]; then
  echo "GPU temperature critical: ${temp}°C"
fi

# Still override when needed
sysinfo --format pretty --all  # Pretty view of everything
```

**Benefits**:
- Clean scripts without flag clutter
- Consistent output format
- Easy to parse with jq

### 3. Support/Helpdesk - Standardized Reports

**Scenario**: Support team needs consistent reports from customer machines.

```yaml
# C:\SupportTools\sysinfo\.sysinforc
format: text
output_file: C:\Support\Reports\system-info.txt
verbose: true

modules:
  system: true
  cpu: true
  memory: true
  disk: true
  smart: true
  gpu: true
  network: true
  process: true
```

**Usage**:
```powershell
# Run on customer machine
sysinfo.exe

# Report automatically saved to standard location
# Support staff can attach C:\Support\Reports\system-info.txt to ticket
```

**Benefits**:
- Uniform reports across all machines
- No training needed (just run the tool)
- Automatic file naming

### 4. CI/CD Pipeline - Build Environment Logging

**Scenario**: Log build machine specs in CI pipeline without cluttering logs.

```yaml
# .sysinforc in project repository
format: json
verbose: false

modules:
  cpu: true
  memory: true
  disk: true
```

**Usage**:
```yaml
# .github/workflows/build.yml
- name: Log build environment
  run: sysinfo -o build-env.json

- name: Upload environment info
  uses: actions/upload-artifact@v3
  with:
    name: build-environment
    path: build-env.json
```

**Benefits**:
- Minimal config in repository
- No verbose output cluttering CI logs
- Consistent across all builds

### 5. Multi-Environment Setup

**Scenario**: Different configs for different purposes.

```yaml
# ~/.sysinfo-dev.yaml - Development quick check
format: pretty
modules:
  cpu: true
  gpu: true
  memory: true
```

```yaml
# ~/.sysinfo-prod.yaml - Production monitoring
format: json
output_file: /var/log/sysinfo/metrics.json
modules:
  all: true
smart:
  enable_alerts: true
  alert_thresholds:
    temperature_critical: 65
```

```yaml
# ~/.sysinfo-minimal.yaml - Quick status
format: text
modules:
  cpu: true
  memory: true
```

**Usage**:
```bash
# Development - quick pretty view
sysinfo --config ~/.sysinfo-dev.yaml

# Production check
sysinfo --config ~/.sysinfo-prod.yaml

# Quick text status
sysinfo --config ~/.sysinfo-minimal.yaml

# Default (uses ~/.sysinforc)
sysinfo
```

### 6. Server Monitoring with Alerting

**Scenario**: Automated disk health monitoring with webhook alerts.

```yaml
# /etc/sysinfo/monitoring.yaml
format: json
verbose: false

modules:
  disk: true
  smart: true
  memory: true
  cpu: true

smart:
  enable_alerts: true
  alert_thresholds:
    temperature_critical: 70
    temperature_warning: 60
  webhook_url: https://monitoring.company.internal/api/sysinfo-alerts
```

**Usage**:
```bash
# Hourly cron job
0 * * * * /usr/bin/sysinfo --config /etc/sysinfo/monitoring.yaml | \
  /usr/local/bin/process-metrics.sh
```

**Future enhancement**: Direct webhook posting from sysinfo.

### 7. Data Center - Uniform Deployment

**Scenario**: Deploy identical monitoring to all servers via Ansible/Puppet.

```yaml
# Ansible template: roles/monitoring/files/sysinfo.yaml
format: json
output_file: /var/log/sysinfo/metrics.json

modules:
  system: true
  cpu: true
  memory: true
  disk: true
  smart: true
  network: true

smart:
  enable_alerts: true
  alert_thresholds:
    temperature_critical: {{ smart_temp_critical | default(70) }}
    temperature_warning: {{ smart_temp_warning | default(60) }}
  webhook_url: {{ monitoring_webhook_url }}
```

**Ansible playbook**:
```yaml
- name: Deploy SysInfo config
  copy:
    src: sysinfo.yaml
    dest: /etc/sysinfo/config.yaml
    mode: 0644

- name: Install cron job
  cron:
    name: "SysInfo metrics collection"
    minute: "*/15"
    job: "/usr/bin/sysinfo --config /etc/sysinfo/config.yaml"
```

**Benefits**:
- Centralized configuration management
- Easy to update all servers
- Templated for environment differences

### 8. Performance Testing Lab

**Scenario**: Capture system state before/after performance tests.

```yaml
# lab-baseline.yaml
format: json
output_file: baseline-{{ test_run_id }}.json
verbose: false

modules:
  cpu: true
  memory: true
  gpu: true
  disk: true
  network: true
  process: true

process:
  top_count: 20
```

**Test script**:
```bash
#!/bin/bash
TEST_ID=$(date +%Y%m%d-%H%M%S)

# Capture baseline
sysinfo -o "baseline-${TEST_ID}.json"

# Run performance test
./run-performance-test.sh

# Capture after state
sysinfo -o "after-${TEST_ID}.json"

# Compare
./compare-metrics.py baseline-${TEST_ID}.json after-${TEST_ID}.json
```

### 9. Home Lab - Always Pretty

**Scenario**: Home lab enthusiast wants full pretty output every time.

```yaml
# ~/.sysinforc
format: pretty
verbose: false

modules:
  system: true
  cpu: true
  memory: true
  disk: true
  smart: true
  network: true
  gpu: true
  process: true

process:
  top_count: 15

display:
  use_ascii: false
```

**Usage**:
```bash
# Just run it - always get full pretty output
sysinfo

# Still can override
sysinfo --format json -o server-state.json
```

### 10. Laptop User - Battery Focused

**Scenario**: Laptop user wants to monitor battery and resource usage.

```yaml
# ~/.config/sysinfo/config.yaml
format: pretty

modules:
  battery: true     # Future feature
  cpu: true
  memory: true
  gpu: true
  process: true

process:
  top_count: 5

display:
  use_ascii: false
```

**Usage**:
```bash
# Quick battery + resource check
sysinfo

# Just battery
sysinfo --battery
```

## Best Practices

### 1. Version Control Your Configs

Store configs in your dotfiles repository:

```bash
# In your dotfiles repo
~/dotfiles/
  └── sysinfo/
      ├── default.yaml      # General purpose
      ├── monitoring.yaml   # Production monitoring
      └── dev.yaml          # Development

# Symlink to home
ln -s ~/dotfiles/sysinfo/default.yaml ~/.sysinforc
```

### 2. Use Comments Liberally

```yaml
# Production Server Monitoring
# Updated: 2025-11-07
# Owner: DevOps Team

format: json  # JSON for easy parsing

# Output to rotating log
output_file: /var/log/sysinfo/metrics.json

modules:
  # Core metrics always needed
  cpu: true
  memory: true
  disk: true
  
  # SMART requires root - enabled in production only
  smart: true
  
  # Network not needed in monitoring
  network: false
```

### 3. Test Config Changes

```bash
# Test config syntax
sysinfo --config new-config.yaml --verbose

# Dry run with different output to verify
sysinfo --config new-config.yaml -o /tmp/test-output.json

# Compare with known good output
diff <(sysinfo --config old-config.yaml) \
     <(sysinfo --config new-config.yaml)
```

### 4. Use Environment-Specific Configs

```bash
# Set up aliases
alias sysinfo-dev='sysinfo --config ~/.sysinfo-dev.yaml'
alias sysinfo-prod='sysinfo --config ~/.sysinfo-prod.yaml'
alias sysinfo-quick='sysinfo --config ~/.sysinfo-minimal.yaml'

# Or use environment variables (future enhancement)
export SYSINFO_CONFIG=~/.sysinfo-${ENVIRONMENT}.yaml
```

### 5. Keep Secrets Out of Configs

For webhook URLs with tokens:

```yaml
# DON'T DO THIS
smart:
  webhook_url: https://api.example.com/webhook?token=secret123

# DO THIS - use environment variables (future enhancement)
smart:
  webhook_url: ${WEBHOOK_URL}

# Or keep sensitive configs separately
smart:
  webhook_url: https://api.example.com/webhook
```

### 6. Document Your Configs

Add a README alongside your config:

```markdown
# SysInfo Configuration

## Purpose
Production server monitoring for Project X

## Usage
Runs hourly via cron: `/etc/cron.hourly/sysinfo-metrics`

## Alerts
SMART alerts post to Slack #infrastructure channel

## Maintenance
- Update temperature thresholds seasonally
- Review webhook URL quarterly
```

## Advanced Usage

### Dynamic Output Paths

Use date/hostname in output paths (requires shell wrapper):

```bash
#!/bin/bash
# sysinfo-logger.sh
OUTPUT_DIR="/var/log/sysinfo"
HOSTNAME=$(hostname)
DATE=$(date +%Y%m%d)

sysinfo --config /etc/sysinfo/config.yaml \
  -o "${OUTPUT_DIR}/${HOSTNAME}-${DATE}.json"
```

### Conditional Module Loading

Different modules based on system type (requires wrapper):

```bash
#!/bin/bash
# sysinfo-adaptive.sh

if lspci | grep -i nvidia > /dev/null; then
  CONFIG="~/.sysinfo-gpu.yaml"
else
  CONFIG="~/.sysinfo-nogpu.yaml"
fi

sysinfo --config "$CONFIG" "$@"
```

### Config Validation Script

```bash
#!/bin/bash
# validate-sysinfo-config.sh

CONFIG_FILE="${1:-~/.sysinforc}"

echo "Validating $CONFIG_FILE..."

# Check YAML syntax
if ! python3 -c "import yaml; yaml.safe_load(open('$CONFIG_FILE'))" 2>/dev/null; then
  echo "ERROR: Invalid YAML syntax"
  exit 1
fi

# Test run
if ! sysinfo --config "$CONFIG_FILE" -o /tmp/sysinfo-validate.json 2>&1 | grep -q "Collecting"; then
  echo "ERROR: Config failed to run"
  exit 1
fi

echo "Config is valid!"
rm -f /tmp/sysinfo-validate.json
```

### Templating Configs

Using envsubst for environment variable substitution:

```yaml
# config.template.yaml
format: ${OUTPUT_FORMAT:-pretty}
output_file: ${OUTPUT_PATH:-}

modules:
  cpu: ${COLLECT_CPU:-true}
  memory: ${COLLECT_MEMORY:-true}
  smart: ${COLLECT_SMART:-false}

smart:
  webhook_url: ${WEBHOOK_URL}
```

```bash
# Generate config
export OUTPUT_FORMAT=json
export WEBHOOK_URL=https://monitoring.example.com/webhook
envsubst < config.template.yaml > config.yaml

sysinfo --config config.yaml
```

## Troubleshooting

### Config Not Loading

**Problem**: Config file seems to be ignored.

**Debug steps**:
```bash
# 1. Check if file exists
ls -la ~/.sysinforc

# 2. Verify YAML syntax
python3 -c "import yaml; print(yaml.safe_load(open('~/.sysinforc')))"

# 3. Run with verbose to see what's loaded
sysinfo -v

# 4. Explicitly specify config
sysinfo --config ~/.sysinforc -v
```

**Common causes**:
- YAML indentation errors (use spaces, not tabs)
- File in wrong location
- Permissions issue (make sure it's readable)
- CLI flags overriding config

### YAML Syntax Errors

**Problem**: Error parsing config file.

**Common mistakes**:

```yaml
# WRONG - tabs instead of spaces
format:	json   # Uses tab

# RIGHT
format: json   # Uses spaces

# WRONG - inconsistent indentation
modules:
  cpu: true
   memory: true  # Extra space

# RIGHT
modules:
  cpu: true
  memory: true

# WRONG - missing space after colon
format:json

# RIGHT
format: json
```

**Validation**:
```bash
# Check YAML online
cat ~/.sysinforc | python3 -c "import sys, yaml; yaml.safe_load(sys.stdin)"

# Or use yamllint
yamllint ~/.sysinforc
```

### CLI Flags Not Overriding Config

**Problem**: CLI flag doesn't override config setting.

**Explanation**: This might be expected behavior:

```yaml
# config.yaml
format: json
```

```bash
# This DOES override
sysinfo --format pretty  # Uses pretty, not json

# This doesn't change anything (--all is default anyway)
sysinfo --all  # Still uses json from config
```

**Solution**: Always check what you're trying to override:
```bash
# See effective config with verbose
sysinfo -v
```

### Permission Denied

**Problem**: Can't read config file.

```bash
# Check permissions
ls -l ~/.sysinforc

# Should be readable
# -rw-r--r-- or -rw-------

# Fix permissions
chmod 644 ~/.sysinforc
```

### Config in Wrong Location

**Problem**: Config file not found.

**Debug**:
```bash
# Check search paths
sysinfo -v  # Will show where it looked

# Test each location
cat ~/.sysinforc
cat ~/.config/sysinfo/config.yaml
cat ./.sysinforc
```

**Solution**: Put config in one of the searched locations or use `--config`.

## Migration Guide

### From CLI Flags to Config

**Before**:
```bash
# Your typical usage
sysinfo --format json --cpu --memory --disk -o /tmp/metrics.json
```

**After**:
```yaml
# ~/.sysinforc
format: json
output_file: /tmp/metrics.json
modules:
  cpu: true
  memory: true
  disk: true
```

```bash
# Now just run
sysinfo
```

### Updating Existing Configs

When new options are added to SysInfo:

```yaml
# Old config (still works!)
format: json
modules:
  cpu: true
  memory: true

# Updated config (with new features)
format: json
modules:
  cpu: true
  memory: true
  battery: true  # New module

process:
  top_count: 15  # New option
```

Old configs remain compatible - new options are optional.

## FAQ

**Q: Can I have multiple config files?**  
A: Yes! Use `--config` to specify different files for different purposes.

**Q: What if I don't have a config file?**  
A: No problem - SysInfo works fine without one. All defaults are sensible.

**Q: Do CLI flags override config settings?**  
A: Yes, always. Config provides defaults, CLI flags override.

**Q: Can I use environment variables in config?**  
A: Not yet, but this is planned for a future release.

**Q: Is the config file required?**  
A: No - it's completely optional. Use it if it makes your workflow easier.

**Q: Can I share my config file?**  
A: Absolutely! Version control it, share with your team, deploy via Ansible, etc.

**Q: What happens if config has invalid YAML?**  
A: SysInfo will report an error and not run. Fix the YAML syntax.

**Q: Can I disable modules in the config?**  
A: Yes - setting a module to `false` disables it when using the config defaults.

**Q: How do I see what config was loaded?**  
A: Run with `-v/--verbose` flag to see config loading details.

## See Also

- [Main README](../README.md) - General usage and features
- [ROADMAP](../ROADMAP.md) - Upcoming configuration features
- [Example Config](.sysinforc.example) - Full example with all options
- [SMART Guide](SMART.md) - SMART monitoring configuration details (planned)

## Contributing

Have ideas for config options? Open an issue or PR!

Common requests:
- Environment variable substitution
- Config file validation command
- More granular module controls
- Additional output path variables

---

**Last Updated**: November 7, 2025  
**SysInfo Version**: 1.0.0+ (Configuration support)
