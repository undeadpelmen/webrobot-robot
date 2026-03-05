package motor

import (
	"context"
	"fmt"
	"sync"

	"github.com/undeadpelmen/webrobot-robot/hardware/interfaces"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

// L298NDriver implements MotorDriver interface for L298N motor controller
type L298NDriver struct {
	in1Pin    gpio.PinOut
	in2Pin    gpio.PinOut
	speedPin  gpio.PinOut
	direction string
	speed     int
	mu        sync.RWMutex
}

// Config holds the configuration for L298N driver
type Config struct {
	In1Pin   string
	In2Pin   string
	SpeedPin string
}

// NewL298NDriver creates a new L298N motor driver
func NewL298NDriver(config Config) (*L298NDriver, error) {
	in1Pin := gpioreg.ByName(config.In1Pin)
	in2Pin := gpioreg.ByName(config.In2Pin)
	speedPin := gpioreg.ByName(config.SpeedPin)

	if in1Pin == nil || in2Pin == nil || speedPin == nil {
		return nil, fmt.Errorf("failed to initialize GPIO pins: in1=%s, in2=%s, speed=%s",
			config.In1Pin, config.In2Pin, config.SpeedPin)
	}

	driver := &L298NDriver{
		in1Pin:    in1Pin,
		in2Pin:    in2Pin,
		speedPin:  speedPin,
		direction: "stop",
		speed:     0,
	}

	return driver, nil
}

// Initialize initializes the L298N driver
func (l *L298NDriver) Initialize(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Initialize pins to safe state
	if err := l.in1Pin.Out(gpio.Low); err != nil {
		return fmt.Errorf("failed to initialize IN1 pin: %w", err)
	}

	if err := l.in2Pin.Out(gpio.Low); err != nil {
		return fmt.Errorf("failed to initialize IN2 pin: %w", err)
	}

	if err := l.speedPin.Out(gpio.Low); err != nil {
		return fmt.Errorf("failed to initialize speed pin: %w", err)
	}

	l.direction = "stop"
	l.speed = 0

	fmt.Println("L298N motor driver initialized")
	return nil
}

// Shutdown shuts down the L298N driver
func (l *L298NDriver) Shutdown(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Set all pins to low for safety
	if err := l.in1Pin.Out(gpio.Low); err != nil {
		return fmt.Errorf("failed to shutdown IN1 pin: %w", err)
	}

	if err := l.in2Pin.Out(gpio.Low); err != nil {
		return fmt.Errorf("failed to shutdown IN2 pin: %w", err)
	}

	if err := l.speedPin.Out(gpio.Low); err != nil {
		return fmt.Errorf("failed to shutdown speed pin: %w", err)
	}

	l.direction = "stop"
	l.speed = 0

	fmt.Println("L298N motor driver shutdown")
	return nil
}

// SetDirection sets the motor direction
func (l *L298NDriver) SetDirection(direction string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch direction {
	case "forward":
		if err := l.in1Pin.Out(gpio.High); err != nil {
			return fmt.Errorf("failed to set IN1 high: %w", err)
		}
		if err := l.in2Pin.Out(gpio.Low); err != nil {
			return fmt.Errorf("failed to set IN2 low: %w", err)
		}

	case "backward":
		if err := l.in1Pin.Out(gpio.Low); err != nil {
			return fmt.Errorf("failed to set IN1 low: %w", err)
		}
		if err := l.in2Pin.Out(gpio.High); err != nil {
			return fmt.Errorf("failed to set IN2 high: %w", err)
		}

	case "stop":
		if err := l.in1Pin.Out(gpio.Low); err != nil {
			return fmt.Errorf("failed to set IN1 low: %w", err)
		}
		if err := l.in2Pin.Out(gpio.Low); err != nil {
			return fmt.Errorf("failed to set IN2 low: %w", err)
		}

	default:
		return fmt.Errorf("invalid direction: %s", direction)
	}

	l.direction = direction
	return nil
}

// SetSpeed sets the motor speed using PWM
func (l *L298NDriver) SetSpeed(speed int) error {
	if speed < 0 || speed > 255 {
		return fmt.Errorf("speed must be between 0 and 255, got: %d", speed)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Convert speed to PWM duty cycle (0-100%)
	duty := float64(speed) / 255.0 * 100.0
	if err := l.speedPin.PWM(gpio.Duty(duty), 0); err != nil {
		return fmt.Errorf("failed to set PWM duty: %w", err)
	}

	l.speed = speed
	return nil
}

// GetSpeed returns the current speed
func (l *L298NDriver) GetSpeed() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.speed
}

// GetDirection returns the current direction
func (l *L298NDriver) GetDirection() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.direction
}
