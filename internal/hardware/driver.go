package hardware

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/undeadpelmen/webrobot-robot/hardware"
	"github.com/undeadpelmen/webrobot-robot/hardware/interfaces"
)

type L298NDriver struct {
	service interfaces.HardwareDriver
	logger  zerolog.Logger
	config  hardware.Config
}

func NewL298NDriver(cfg hardware.Config, logger zerolog.Logger) (*L298NDriver, error) {
	factory := hardware.NewFactory()

	hardwareLogger := &hardwareLoggerAdapter{logger: logger}
	service, err := factory.CreateHardwareService(cfg, hardwareLogger, cfg.TestPins)
	if err != nil {
		return nil, fmt.Errorf("failed to create L298N hardware service: %w", err)
	}

	return &L298NDriver{
		service: service,
		logger:  logger,
		config:  cfg,
	}, nil
}

func (l *L298NDriver) Initialize(ctx context.Context) error {
	l.logger.Info().Msg("Initializing L298N hardware driver")

	if err := l.service.Initialize(ctx); err != nil {
		l.logger.Error().Err(err).Msg("Failed to initialize hardware")
		return err
	}

	l.logger.Info().Msg("L298N driver initialized successfully")
	return nil
}

func (l *L298NDriver) Shutdown(ctx context.Context) error {
	l.logger.Info().Msg("Shutting down L298N driver")

	if err := l.service.Shutdown(ctx); err != nil {
		l.logger.Error().Err(err).Msg("Failed to shutdown hardware")
		return err
	}

	l.logger.Info().Msg("L298N driver shutdown complete")
	return nil
}

func (l *L298NDriver) MoveForward(speed int) error {
	l.logger.Debug().Int("speed", speed).Msg("Moving forward")
	return l.service.MoveForward(speed)
}

func (l *L298NDriver) MoveBackward(speed int) error {
	l.logger.Debug().Int("speed", speed).Msg("Moving backward")
	return l.service.MoveBackward(speed)
}

func (l *L298NDriver) TurnLeft(speed int) error {
	l.logger.Debug().Int("speed", speed).Msg("Turning left")
	return l.service.TurnLeft(speed)
}

func (l *L298NDriver) TurnRight(speed int) error {
	l.logger.Debug().Int("speed", speed).Msg("Turning right")
	return l.service.TurnRight(speed)
}

func (l *L298NDriver) Stop() error {
	l.logger.Debug().Msg("Stopping")
	return l.service.Stop()
}

// hardwareLoggerAdapter adapts zerolog.Logger to hardware.Logger interface
type hardwareLoggerAdapter struct {
	logger zerolog.Logger
}

func (h *hardwareLoggerAdapter) Debug(msg string, fields ...interface{}) {
	h.logger.Debug().Msgf(msg, fields...)
}

func (h *hardwareLoggerAdapter) Info(msg string, fields ...interface{}) {
	h.logger.Info().Msgf(msg, fields...)
}

func (h *hardwareLoggerAdapter) Warn(msg string, fields ...interface{}) {
	h.logger.Warn().Msgf(msg, fields...)
}

func (h *hardwareLoggerAdapter) Error(msg string, fields ...interface{}) {
	h.logger.Error().Msgf(msg, fields...)
}
