package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Robot    RobotConfig    `yaml:"robot"`
	Logging  LoggingConfig  `yaml:"logging"`
	Hardware HardwareConfig `yaml:"hardware"`
}

type ServerConfig struct {
	HTTPPort      int    `yaml:"http_port"`
	WebSocketPort int    `yaml:"websocket_port"`
	Host          string `yaml:"host"`
	EnableHTTP    bool   `yaml:"enable_http"`
	EnableWS      bool   `yaml:"enable_websocket"`
}

type RobotConfig struct {
	Speed       int  `yaml:"speed"`
	TestMode    bool `yaml:"test_mode"`
	EnableCLI   bool `yaml:"enable_cli"`
	AutoConnect bool `yaml:"auto_connect"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	Console    bool   `yaml:"console"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
}

type HardwareConfig struct {
	Driver   string            `yaml:"driver"`
	Pins     map[string]string `yaml:"pins"`
	TestPins bool              `yaml:"test_pins"`
}

func Load(configPath string) (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			HTTPPort:      8080,
			WebSocketPort: 8081,
			Host:          "0.0.0.0",
			EnableHTTP:    false,
			EnableWS:      false,
		},
		Robot: RobotConfig{
			Speed:       255,
			TestMode:    false,
			EnableCLI:   false,
			AutoConnect: true,
		},
		Logging: LoggingConfig{
			Level:      "info",
			File:       "/tmp/webrobot/robot.log",
			Console:    true,
			MaxSize:    10,
			MaxBackups: 5,
		},
		Hardware: HardwareConfig{
			Driver: "l298n",
			Pins: map[string]string{
				"enable_a": "11",
				"input1":   "13",
				"input2":   "15",
				"enable_b": "16",
				"input3":   "18",
				"input4":   "22",
			},
			TestPins: false,
		},
	}

	if configPath == "" {
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, fmt.Errorf("config file not found: %s", configPath)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

func (c *Config) EnsureLogDirectory() error {
	logDir := filepath.Dir(c.Logging.File)
	return os.MkdirAll(logDir, 0755)
}

func (c *Config) GetLogLevel() zerolog.Level {
	switch c.Logging.Level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
