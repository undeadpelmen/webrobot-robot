package interfaces

import "context"

type RobotController interface {
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Move(direction string, speed int) error
	Stop() error
	Status() string
	GetSpeed() int
	SetSpeed(speed int) error
}

type HardwareDriver interface {
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
	MoveForward(speed int) error
	MoveBackward(speed int) error
	TurnLeft(speed int) error
	TurnRight(speed int) error
	Stop() error
}

type MessageHandler interface {
	HandleMessage(message string) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

type ConfigManager interface {
	GetServerConfig() interface{}
	GetRobotConfig() interface{}
	GetLoggingConfig() interface{}
	GetHardwareConfig() interface{}
}
