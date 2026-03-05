package robot

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"github.com/undeadpelmen/webrobot-robot/internal/interfaces"
)

type Service struct {
	driver interfaces.HardwareDriver
	logger zerolog.Logger
	status string
	speed  int
	mu     sync.RWMutex
}

func NewService(driver interfaces.HardwareDriver, logger zerolog.Logger) *Service {
	return &Service{
		driver: driver,
		logger: logger,
		status: "stopped",
		speed:  255,
	}
}

func (s *Service) Initialize(ctx context.Context) error {
	s.logger.Info().Msg("Initializing robot service")

	if err := s.driver.Initialize(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to initialize hardware driver")
		return fmt.Errorf("failed to initialize robot service: %w", err)
	}

	s.setStatus("ready")
	s.logger.Info().Msg("Robot service initialized successfully")
	return nil
}

func (s *Service) Move(direction string, speed int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != "ready" {
		return fmt.Errorf("robot is not ready, current status: %s", s.status)
	}

	s.logger.Info().Str("direction", direction).Int("speed", speed).Msg("Moving robot")

	var err error
	switch direction {
	case "forward", "f":
		err = s.driver.MoveForward(speed)
	case "backward", "b":
		err = s.driver.MoveBackward(speed)
	case "left", "l":
		err = s.driver.TurnLeft(speed)
	case "right", "r":
		err = s.driver.TurnRight(speed)
	default:
		return fmt.Errorf("invalid direction: %s", direction)
	}

	if err != nil {
		s.logger.Error().Err(err).Str("direction", direction).Msg("Failed to move robot")
		return err
	}

	s.speed = speed
	s.setStatus("moving")
	return nil
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info().Msg("Stopping robot")

	if err := s.driver.Stop(); err != nil {
		s.logger.Error().Err(err).Msg("Failed to stop robot")
		return err
	}

	s.setStatus("stopped")
	return nil
}

func (s *Service) Status() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

func (s *Service) GetSpeed() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.speed
}

func (s *Service) SetSpeed(speed int) error {
	if speed < 0 || speed > 255 {
		return fmt.Errorf("speed must be between 0 and 255, got: %d", speed)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.speed = speed
	s.logger.Debug().Int("speed", speed).Msg("Speed updated")
	return nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("Shutting down robot service")

	if err := s.Stop(); err != nil {
		s.logger.Error().Err(err).Msg("Failed to stop robot during shutdown")
	}

	if err := s.driver.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to shutdown hardware driver")
		return err
	}

	s.setStatus("shutdown")
	s.logger.Info().Msg("Robot service shutdown complete")
	return nil
}

func (s *Service) setStatus(status string) {
	s.status = status
	s.logger.Debug().Str("status", status).Msg("Robot status updated")
}
