package motor

import (
	"context"
	"fmt"
	"sync"

	"github.com/undeadpelmen/webrobot-robot/hardware/interfaces"
)

// Service implements MotorController interface
type Service struct {
	driver    interfaces.MotorDriver
	direction string
	speed     int
	mu        sync.RWMutex
	logger    Logger
}

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// NewService creates a new motor service
func NewService(driver interfaces.MotorDriver, logger Logger) *Service {
	return &Service{
		driver:    driver,
		direction: "stop",
		speed:     0,
		logger:    logger,
	}
}

// Initialize initializes the motor service
func (s *Service) Initialize(ctx context.Context) error {
	s.logger.Debug("Initializing motor service")

	if err := s.driver.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize motor driver: %w", err)
	}

	s.logger.Info("Motor service initialized successfully")
	return nil
}

// Shutdown shuts down the motor service
func (s *Service) Shutdown(ctx context.Context) error {
	s.logger.Debug("Shutting down motor service")

	if err := s.driver.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown motor driver: %w", err)
	}

	s.logger.Info("Motor service shutdown complete")
	return nil
}

// Forward moves the motor forward at the specified speed
func (s *Service) Forward(speed int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug().Int("speed", speed).Msg("Moving motor forward")

	if err := s.driver.SetDirection("forward"); err != nil {
		return fmt.Errorf("failed to set forward direction: %w", err)
	}

	if err := s.driver.SetSpeed(speed); err != nil {
		return fmt.Errorf("failed to set speed: %w", err)
	}

	s.direction = "forward"
	s.speed = speed
	return nil
}

// Backward moves the motor backward at the specified speed
func (s *Service) Backward(speed int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug().Int("speed", speed).Msg("Moving motor backward")

	if err := s.driver.SetDirection("backward"); err != nil {
		return fmt.Errorf("failed to set backward direction: %w", err)
	}

	if err := s.driver.SetSpeed(speed); err != nil {
		return fmt.Errorf("failed to set speed: %w", err)
	}

	s.direction = "backward"
	s.speed = speed
	return nil
}

// Stop stops the motor
func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug("Stopping motor")

	if err := s.driver.SetDirection("stop"); err != nil {
		return fmt.Errorf("failed to set stop direction: %w", err)
	}

	s.direction = "stop"
	s.speed = 0
	return nil
}

// SetSpeed sets the motor speed (0-255)
func (s *Service) SetSpeed(speed int) error {
	if speed < 0 || speed > 255 {
		return fmt.Errorf("speed must be between 0 and 255, got: %d", speed)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug().Int("speed", speed).Msg("Setting motor speed")

	if err := s.driver.SetSpeed(speed); err != nil {
		return fmt.Errorf("failed to set speed: %w", err)
	}

	s.speed = speed
	return nil
}

// GetSpeed returns the current motor speed
func (s *Service) GetSpeed() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.speed
}

// GetDirection returns the current motor direction
func (s *Service) GetDirection() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.direction
}
