package hardware

import (
	"context"
	"fmt"

	"github.com/undeadpelmen/webrobot-robot/hardware/interfaces"
	"github.com/undeadpelmen/webrobot-robot/hardware/robot"
)

// Service implements HardwareDriver interface
type Service struct {
	robot interfaces.RobotController
}

// Config holds the configuration for hardware service
type Config struct {
	Robot robot.Config `yaml:"robot"`
}

// NewService creates a new hardware service
func NewService(robot interfaces.RobotController) *Service {
	return &Service{
		robot: robot,
	}
}

// Initialize initializes the hardware service
func (s *Service) Initialize(ctx context.Context) error {
	return s.robot.Initialize(ctx)
}

// Shutdown shuts down the hardware service
func (s *Service) Shutdown(ctx context.Context) error {
	return s.robot.Shutdown(ctx)
}

// MoveForward moves the robot forward
func (s *Service) MoveForward(speed int) error {
	return s.robot.Move("forward", speed)
}

// MoveBackward moves the robot backward
func (s *Service) MoveBackward(speed int) error {
	return s.robot.Move("backward", speed)
}

// TurnLeft turns the robot left
func (s *Service) TurnLeft(speed int) error {
	return s.robot.Move("left", speed)
}

// TurnRight turns the robot right
func (s *Service) TurnRight(speed int) error {
	return s.robot.Move("right", speed)
}

// Stop stops the robot
func (s *Service) Stop() error {
	return s.robot.Stop()
}

// Factory creates hardware services
type Factory struct{}

// NewFactory creates a new hardware factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateHardwareService creates a hardware service based on configuration
func (f *Factory) CreateHardwareService(config Config, logger Logger, useMock bool) (interfaces.HardwareDriver, error) {
	robotFactory := robot.NewFactory()

	robotService, err := robotFactory.CreateRobotFromConfig(config.Robot, logger, useMock)
	if err != nil {
		return nil, fmt.Errorf("failed to create robot: %w", err)
	}

	return NewService(robotService), nil
}

// Logger interface for hardware logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}
