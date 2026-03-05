package interfaces

import "context"

// MotorController defines the interface for controlling individual motors
type MotorController interface {
	Forward(speed int) error
	Backward(speed int) error
	Stop() error
	SetSpeed(speed int) error
	GetSpeed() int
	GetDirection() string
}

// RobotController defines the interface for robot movement control
type RobotController interface {
	Move(direction string, speed int) error
	Stop() error
	Status() string
	GetSpeed() int
	SetSpeed(speed int) error
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// HardwareDriver defines the interface for hardware abstraction
type HardwareDriver interface {
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
	MoveForward(speed int) error
	MoveBackward(speed int) error
	TurnLeft(speed int) error
	TurnRight(speed int) error
	Stop() error
}

// MotorDriver defines the interface for motor hardware control
type MotorDriver interface {
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
	SetDirection(direction string) error
	SetSpeed(speed int) error
	GetSpeed() int
	GetDirection() string
}

// GPIOController defines the interface for GPIO pin control
type GPIOController interface {
	SetPin(pin string, level string) error
	SetPWMPin(pin string, duty int) error
	ReadPin(pin string) (string, error)
	Initialize() error
	Shutdown() error
}
