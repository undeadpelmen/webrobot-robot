# WebRobot Control System

A modular robot control system with support for multiple interfaces including CLI, HTTP API, and WebSocket.

## Architecture

The project has been refactored with a clean, modular architecture:

- **`internal/config/`** - Configuration management with YAML support
- **`internal/interfaces/`** - Core interfaces for dependency injection and testability
- **`internal/robot/`** - Robot control service layer
- **`internal/hardware/`** - Hardware abstraction layer with mock support
- **`internal/api/`** - HTTP REST API handlers
- **`internal/websocket/`** - WebSocket server implementation
- **`internal/cli/`** - Command-line interface service
- **`internal/app/`** - Application orchestration and lifecycle management
- **`cmd/robot/`** - Application entry point

## Features

- **Modular Design**: Clean separation of concerns with dependency injection
- **Configuration Management**: YAML-based configuration with environment support
- **Multiple Interfaces**: CLI, HTTP API, and WebSocket support
- **Hardware Abstraction**: Support for real hardware and mock/testing modes
- **Structured Logging**: Configurable logging with zerolog
- **Graceful Shutdown**: Proper cleanup and shutdown handling

## Quick Start

### Build

```bash
go build -o bin/robot ./cmd/robot
```

### Run with Default Configuration

```bash
# Start with CLI interface
./bin/robot -cli

# Start with HTTP API
./bin/robot -http

# Start with WebSocket server
./bin/robot -websocket

# Start with all interfaces
./bin/robot -cli -http -websocket
```

### Run with Configuration File

```bash
./bin/robot -config config.yaml -cli -http -websocket
```

### Test Mode

```bash
./bin/robot -test -cli -http -websocket
```

## Configuration

The system uses YAML configuration files. See `config.yaml` for an example:

```yaml
server:
  host: "0.0.0.0"
  http_port: 8080
  websocket_port: 8081
  enable_http: true
  enable_websocket: true

robot:
  speed: 255
  test_mode: false
  enable_cli: true
  auto_connect: true

logging:
  level: "info"
  file: "/tmp/webrobot/robot.log"
  console: true
  max_size: 10
  max_backups: 5

hardware:
  driver: "l298n"
  test_pins: false
  pins:
    enable_a: "11"
    input1: "13"
    input2: "15"
    enable_b: "16"
    input3: "18"
    input4: "22"
```

## API Endpoints

### HTTP API

- `GET /api/v1/status` - Get robot status
- `POST /api/v1/move` - Move robot (JSON: `{"direction": "forward", "speed": 255}`)
- `POST /api/v1/stop` - Stop robot
- `GET /api/v1/speed` - Get current speed
- `PUT /api/v1/speed` - Set speed (JSON: `{"speed": 200}`)
- `GET /health` - Health check

### WebSocket

Connect to `ws://localhost:8081/ws` and send JSON messages:

```json
{"type": "move", "payload": {"direction": "forward", "speed": 255}}
{"type": "stop", "payload": {}}
{"type": "status", "payload": {}}
{"type": "set_speed", "payload": {"speed": 200}}
```

### CLI Commands

```
robot> help                    # Show help
robot> status                  # Show robot status
robot> move forward 255        # Move forward at speed 255
robot> move left               # Turn left at current speed
robot> stop                    # Stop robot
robot> speed                   # Show current speed
robot> set-speed 200           # Set speed to 200
robot> exit                    # Exit CLI
```

## Development

### Project Structure

```
webrobot-robot/
├── cmd/robot/           # Application entry point
├── internal/
│   ├── app/            # Application orchestration
│   ├── api/            # HTTP API handlers
│   ├── cli/            # CLI service
│   ├── config/         # Configuration management
│   ├── hardware/       # Hardware abstraction
│   ├── interfaces/     # Core interfaces
│   ├── robot/          # Robot control service
│   └── websocket/      # WebSocket server
├── hardware/           # Legacy hardware implementations
├── web/               # Legacy web implementations
├── cli/               # Legacy CLI implementations
├── config.yaml        # Example configuration
└── go.mod
```

### Adding New Features

1. Define interfaces in `internal/interfaces/`
2. Implement services in appropriate `internal/` modules
3. Wire up dependencies in `internal/app/application.go`
4. Add configuration options in `internal/config/config.go`

### Testing

The architecture supports easy testing with mock implementations:

```go
// Use mock hardware for testing
driver := hardware.NewMockDriver()
robotService := robot.NewService(driver, logger)
```

## Dependencies

- **Gin** - HTTP web framework
- **Gorilla WebSocket** - WebSocket implementation
- **Zerolog** - Structured logging
- **Go-prompt** - Interactive CLI
- **Periph.io** - Hardware GPIO control
- **YAML.v3** - Configuration parsing

## Installation

### Automated Installation (Recommended)

Use the provided install script for automated setup:

```bash
# Clone the repository
git clone https://github.com/undeadpelmen/webrobot-robot.git
cd webrobot-robot

# Run the install script
./install.sh

# Optional: Install with systemd service (requires sudo)
./install.sh --systemd

# Optional: Run tests after installation
./install.sh --test

# Show install script options
./install.sh --help
```

The install script will:
- Check for required dependencies (Go, Git)
- Install Go modules
- Build the application
- Create default configuration file
- Set up log directory
- Optionally create a systemd service

### Manual Installation

```shell
git clone https://github.com/undeadpelmen/webrobot-robot.git
cd webrobot-robot
go mod tidy
go build -o bin/robot ./cmd/robot
```
