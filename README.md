# portwatch

A lightweight CLI daemon that monitors port usage and alerts on unexpected bindings or conflicts.

## Features

- Real-time monitoring of network port bindings
- Configurable alerting for unexpected port usage
- Low resource footprint
- Cross-platform support (Linux, macOS, Windows)
- Simple YAML-based configuration

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch
go build -o portwatch ./cmd/portwatch
```

## Usage

Create a configuration file `portwatch.yml`:

```yaml
watch:
  - port: 8080
    expected: "myapp"
  - port: 5432
    expected: "postgres"
  - port: 3000
    alert_on_bind: true

check_interval: 5s
```

Start the daemon:

```bash
portwatch --config portwatch.yml
```

Run as a one-time check:

```bash
portwatch --config portwatch.yml --once
```

## Configuration Options

- `watch`: List of ports to monitor
- `expected`: Process name expected to bind to the port
- `alert_on_bind`: Alert when any process binds to this port
- `check_interval`: How often to check port status (default: 10s)

## License

MIT License - see [LICENSE](LICENSE) for details