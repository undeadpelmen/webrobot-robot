package robot

import (
	"context"
	"fmt"

	"github.com/undeadpelmen/webrobot-robot/hardware/interfaces"
	"github.com/undeadpelmen/webrobot-robot/hardware/motor"
)

// Factory creates robot services with different configurations
type Factory struct{}

// NewFactory creates a new robot factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateMockRobot creates a robot with mock motors for testing
func (f *Factory) CreateMockRobot(logger Logger) (interfaces.RobotController, error) {
	leftMotor := motor.NewService(motor.NewMockDriver(), logger)
	rightMotor := motor.NewService(motor.NewMockDriver(), logger)

	return NewService(leftMotor, rightMotor, logger), nil
}

// CreateL298NRobot creates a robot with L298N motor drivers
func (f *Factory) CreateL298NRobot(config Config, logger Logger) (interfaces.RobotController, error) {
	leftMotorDriver, err := motor.NewL298NDriver(config.LeftMotorConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create left motor driver: %w", err)
	}

	rightMotorDriver, err := motor.NewL298NDriver(config.RightMotorConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create right motor driver: %w", err)
	}

	leftMotor := motor.NewService(leftMotorDriver, logger)
	rightMotor := motor.NewService(rightMotorDriver, logger)

	return NewService(leftMotor, rightMotor, logger), nil
}

// CreateRobotFromConfig creates a robot based on configuration
func (f *Factory) CreateRobotFromConfig(config Config, logger Logger, useMock bool) (interfaces.RobotController, error) {
	if useMock {
		return f.CreateMockRobot(logger)
	}
	return f.CreateL298NRobot(config, logger)
}
