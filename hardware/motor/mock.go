package motor

import (
	"context"
	"fmt"

	"github.com/undeadpelmen/webrobot-robot/hardware/interfaces"
)

// MockDriver implements MotorDriver interface for testing
type MockDriver struct {
	direction string
	speed     int
}

// NewMockDriver creates a new mock motor driver
func NewMockDriver() *MockDriver {
	return &MockDriver{
		direction: "stop",
		speed:     0,
	}
}

// Initialize initializes the mock driver
func (m *MockDriver) Initialize(ctx context.Context) error {
	fmt.Println("[MOCK] Motor driver initialized")
	return nil
}

// Shutdown shuts down the mock driver
func (m *MockDriver) Shutdown(ctx context.Context) error {
	fmt.Println("[MOCK] Motor driver shutdown")
	return nil
}

// SetDirection sets the motor direction
func (m *MockDriver) SetDirection(direction string) error {
	switch direction {
	case "forward", "backward", "stop":
		m.direction = direction
		fmt.Printf("[MOCK] Motor direction set to: %s\n", direction)
		return nil
	default:
		return fmt.Errorf("invalid direction: %s", direction)
	}
}

// SetSpeed sets the motor speed
func (m *MockDriver) SetSpeed(speed int) error {
	if speed < 0 || speed > 255 {
		return fmt.Errorf("speed must be between 0 and 255, got: %d", speed)
	}
	m.speed = speed
	fmt.Printf("[MOCK] Motor speed set to: %d\n", speed)
	return nil
}

// GetSpeed returns the current speed
func (m *MockDriver) GetSpeed() int {
	return m.speed
}

// GetDirection returns the current direction
func (m *MockDriver) GetDirection() string {
	return m.direction
}
