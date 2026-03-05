#!/bin/bash

# WebRobot Control System Installation Script
# This script installs the webrobot-robot project and its dependencies

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.24.4 or later."
        echo "Visit https://golang.org/dl/ for installation instructions."
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_status "Found Go version: $GO_VERSION"
}

# Check if Git is installed
check_git() {
    if ! command -v git &> /dev/null; then
        print_error "Git is not installed. Please install Git."
        exit 1
    fi
    
    print_status "Found Git: $(git --version)"
}

# Clone repository if not in project directory
setup_repo() {
    if [ ! -f "go.mod" ]; then
        print_status "Cloning webrobot-robot repository..."
        git clone https://github.com/undeadpelmen/webrobot-robot.git
        cd webrobot-robot
        print_status "Repository cloned successfully."
    else
        print_status "Already in project directory."
    fi
}

# Install Go dependencies
install_dependencies() {
    print_status "Installing Go dependencies..."
    go mod tidy
    print_status "Dependencies installed successfully."
}

# Build the application
build_application() {
    print_status "Building webrobot-robot..."
    
    # Create bin directory if it doesn't exist
    mkdir -p bin
    
    # Build the application
    go build -o bin/robot ./cmd/robot
    
    if [ $? -eq 0 ]; then
        print_status "Build successful! Binary created at bin/robot"
    else
        print_error "Build failed!"
        exit 1
    fi
}

# Create systemd service file (optional)
create_systemd_service() {
    if [ "$1" = "--systemd" ]; then
        print_status "Creating systemd service file..."
        
        SERVICE_FILE="/etc/systemd/system/webrobot.service"
        
        # Check if running as root for systemd service creation
        if [ "$EUID" -ne 0 ]; then
            print_warning "Creating systemd service requires root privileges."
            print_warning "Please run with sudo to create the service: sudo ./install.sh --systemd"
            return
        fi
        
        cat > "$SERVICE_FILE" << EOF
[Unit]
Description=WebRobot Control System
After=network.target

[Service]
Type=simple
User=$SUDO_USER
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/bin/robot -http -websocket
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
        
        systemctl daemon-reload
        systemctl enable webrobot.service
        print_status "Systemd service created and enabled."
        print_status "Start the service with: sudo systemctl start webrobot"
        print_status "Check status with: sudo systemctl status webrobot"
    fi
}

# Create configuration file if it doesn't exist
setup_config() {
    if [ ! -f "config.yaml" ]; then
        print_status "Creating default configuration file..."
        cat > config.yaml << EOF
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
EOF
        print_status "Default configuration file created."
    else
        print_status "Configuration file already exists."
    fi
}

# Create log directory
setup_log_directory() {
    LOG_DIR="/tmp/webrobot"
    if [ ! -d "$LOG_DIR" ]; then
        print_status "Creating log directory: $LOG_DIR"
        mkdir -p "$LOG_DIR"
    fi
}

# Run tests
run_tests() {
    if [ "$1" = "--test" ]; then
        print_status "Running tests..."
        go test ./...
        if [ $? -eq 0 ]; then
            print_status "All tests passed!"
        else
            print_warning "Some tests failed."
        fi
    fi
}

# Display usage information
show_usage() {
    echo "WebRobot Control System Installation Complete!"
    echo ""
    echo "Usage:"
    echo "  ./bin/robot -cli                    # Start with CLI interface"
    echo "  ./bin/robot -http                   # Start with HTTP API"
    echo "  ./bin/robot -websocket              # Start with WebSocket server"
    echo "  ./bin/robot -cli -http -websocket   # Start with all interfaces"
    echo "  ./bin/robot -config config.yaml -cli -http -websocket  # Start with config file"
    echo "  ./bin/robot -test -cli -http -websocket               # Start in test mode"
    echo ""
    echo "API Endpoints:"
    echo "  HTTP API: http://localhost:8080"
    echo "  WebSocket: ws://localhost:8081/ws"
    echo ""
    echo "For more information, see README.md"
}

# Main installation function
main() {
    print_status "Starting WebRobot Control System installation..."
    
    check_go
    check_git
    setup_repo
    install_dependencies
    build_application
    setup_config
    setup_log_directory
    run_tests "$1"
    create_systemd_service "$1"
    
    show_usage
    print_status "Installation completed successfully!"
}

# Parse command line arguments
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "WebRobot Control System Installation Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --systemd    Create and enable systemd service (requires sudo)"
    echo "  --test       Run tests after installation"
    echo "  --help, -h   Show this help message"
    echo ""
    exit 0
fi

# Run main function with all arguments
main "$@"
