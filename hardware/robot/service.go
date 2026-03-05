package robot

import (
	"context"
	"fmt"
	"sync"

	"github.com/undeadpelmen/webrobot-robot/hardware/interfaces"
	"github.com/undeadpelmen/webrobot-robot/hardware/motor"
)

// Service implements RobotController interface
type Service struct {
	leftMotor  interfaces.MotorController
	rightMotor interfaces.MotorController
	status     string
	speed      int
	mu         sync.RWMutex
	logger     Logger
}

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// Config holds the configuration for robot service
type Config struct {
	LeftMotorConfig  motor.Config `yaml:"left_motor"`
	RightMotorConfig motor.Config `yaml:"right_motor"`
	DefaultSpeed     int          `yaml:"default_speed"`
}

// NewService creates a new robot service
func NewService(leftMotor, rightMotor interfaces.MotorController, logger Logger) *Service {
	return &Service{
		leftMotor:  leftMotor,
		rightMotor: rightMotor,
		status:     "stopped",
		speed:      255,
		logger:     logger,
	}
}

// Initialize initializes the robot service
func (s *Service) Initialize(ctx context.Context) error {
	s.logger.Debug("Initializing robot service")

	if err := s.leftMotor.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize left motor: %w", err)
	}

	if err := s.rightMotor.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize right motor: %w", err)
	}

	s.setStatus("ready")
	s.logger.Info("Robot service initialized successfully")
	return nil
}

// Shutdown shuts down the robot service
func (s *Service) Shutdown(ctx context.Context) error {
	s.logger.Debug("Shutting down robot service")

	if err := s.Stop(); err != nil {
		s.logger.Error("Failed to stop robot during shutdown")
	}

	if err := s.leftMotor.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown left motor: %w", err)
	}

	if err := s.rightMotor.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown right motor: %w", err)
	}

	s.setStatus("shutdown")
	s.logger.Info("Robot service shutdown complete")
	return nil
}

// Move controls robot movement
func (s *Service) Move(direction string, speed int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != "ready" {
		return fmt.Errorf("robot is not ready, current status: %s", s.status)
	}

	s.logger.Debug().Str("direction", direction).Int("speed", speed).Msg("Moving robot")

	var err error
	switch direction {
	case "forward", "f":
		err = s.moveForward(speed)
	case "backward", "b":
		err = s.moveBackward(speed)
	case "left", "l":
		err = s.turnLeft(speed)
	case "right", "r":
		err = s.turnRight(speed)
	default:
		return fmt.Errorf("invalid direction: %s", direction)
	}

	if err != nil {
		return fmt.Errorf("failed to move robot %s: %w", direction, err)
	}

	s.speed = speed
	s.setStatus("moving")
	return nil
}

// Stop stops the robot
func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug("Stopping robot")

	if err := s.leftMotor.Stop(); err != nil {
		return fmt.Errorf("failed to stop left motor: %w", err)
	}

	if err := s.rightMotor.Stop(); err != nil {
		return fmt.Errorf("failed to stop right motor: %w", err)
	}

	s.setStatus("stopped")
	return nil
}

// Status returns the current robot status
func (s *Service) Status() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// GetSpeed returns the current robot speed
func (s *Service) GetSpeed() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.speed
}

// SetSpeed sets the robot speed
func (s *Service) SetSpeed(speed int) error {
	if speed < 0 || speed > 255 {
		return fmt.Errorf("speed must be between 0 and 255, got: %d", speed)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.leftMotor.SetSpeed(speed); err != nil {
		return fmt.Errorf("failed to set left motor speed: %w", err)
	}

	if err := s.rightMotor.SetSpeed(speed); err != nil {
		return fmt.Errorf("failed to set right motor speed: %w", err)
	}

	s.speed = speed
	s.logger.Debug().Int("speed", speed).Msg("Robot speed updated")
	return nil
}

// moveForward moves the robot forward
func (s *Service) moveForward(speed int) error {
	if err := s.leftMotor.Forward(speed); err != nil {
		return err
	}
	return s.rightMotor.Forward(speed)
}

// moveBackward moves the robot backward
func (s *Service) moveBackward(speed int) error {
	if err := s.leftMotor.Backward(speed); err != nil {
		return err
	}
	return s.rightMotor.Backward(speed)
}

// turnLeft turns the robot left
func (s *Service) turnLeft(speed int) error {
	if err := s.leftMotor.Forward(speed); err != nil {
		return err
	}
	return s.rightMotor.Stop()
}

// turnRight turns the robot right
func (s *Service) turnRight(speed int) error {
	if err := s.leftMotor.Stop(); err != nil {
		return err
	}
	return s.rightMotor.Forward(speed)
}

// setStatus updates the robot status
func (s *Service) setStatus(status string) {
	s.status = status
	s.logger.Debug().Str("status", status).Msg("Robot status updated")
}
