package hardware

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/undeadpelmen/webrobot-robot/hardware"
	"github.com/undeadpelmen/webrobot-robot/hardware/interfaces"
)

type MockDriver struct {
	service interfaces.HardwareDriver
	logger  zerolog.Logger
}

func NewMockDriver() *MockDriver {
	// Create mock configuration
	config := hardware.Config{
		TestPins: true,
	}

	factory := hardware.NewFactory()
	hardwareLogger := &hardwareLoggerAdapter{logger: zerolog.Nop()}
	service, _ := factory.CreateHardwareService(config, hardwareLogger, true)

	return &MockDriver{
		service: service,
		logger:  zerolog.Nop(),
	}
}

func (m *MockDriver) Initialize(ctx context.Context) error {
	return m.service.Initialize(ctx)
}

func (m *MockDriver) Shutdown(ctx context.Context) error {
	return m.service.Shutdown(ctx)
}

func (m *MockDriver) MoveForward(speed int) error {
	return m.service.MoveForward(speed)
}

func (m *MockDriver) MoveBackward(speed int) error {
	return m.service.MoveBackward(speed)
}

func (m *MockDriver) TurnLeft(speed int) error {
	return m.service.TurnLeft(speed)
}

func (m *MockDriver) TurnRight(speed int) error {
	return m.service.TurnRight(speed)
}

func (m *MockDriver) Stop() error {
	return m.service.Stop()
}
