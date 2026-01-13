package l298n

import (
	"errors"
	"sync"

	"periph.io/x/conn/v3/gpio"
)

type Direction string

const (
	Forward  Direction = "forward"
	Backward Direction = "backward"
	Stop     Direction = "stop"
)

type Out struct {
	mu        sync.Mutex
	direction Direction
	in1Pin    gpio.PinOut
	in2Pin    gpio.PinOut
	speedPin  gpio.PinOut
	speed     int
}

func NewOut(in1Pin, in2Pin, speedPin gpio.PinOut, speed int) (*Out, error) {
	if in1Pin == nil || in2Pin == nil || speedPin == nil {
		return nil, errors.New("Pins can not be nil")
	}

	dev := &Out{
		direction: Stop,
		in1Pin:    in1Pin,
		in2Pin:    in2Pin,
		speedPin:  speedPin,
	}

	if err := dev.SetSpeed(speed); err != nil {
		return nil, err
	}

	return dev, nil
}

func (o *Out) Speed() int {
	return o.speed
}

func (o *Out) SetSpeed(speed int) error {
	if speed > 255 || speed < 0 {
		return errors.New("Speed must be between 0 and 255")
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	o.speed = speed

	o.speedPin.PWM(gpio.Duty(speed*100/255), 0)

	return nil
}

func (o *Out) Direction() Direction {
	return o.direction
}

func (o *Out) SetDirection(direction Direction) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	switch direction {
	case Forward:
		if err := o.in1Pin.Out(gpio.High); err != nil {
			return err
		}
		if err := o.in2Pin.Out(gpio.Low); err != nil {
			return err
		}

	case Stop:
		if err := o.in1Pin.Out(gpio.Low); err != nil {
			return err
		}
		if err := o.in2Pin.Out(gpio.Low); err != nil {
			return err
		}

	case Backward:
		if err := o.in1Pin.Out(gpio.Low); err != nil {
			return err
		}
		if err := o.in2Pin.Out(gpio.High); err != nil {
			return err
		}

	default:
		return errors.New("Wrong direction " + string(direction))
	}

	o.direction = direction

	return nil
}

func (o *Out) Forward(speed int) error {
	if err := o.SetDirection(Forward); err != nil {
		return err
	}

	if err := o.SetSpeed(speed); err != nil {
		return err
	}

	return nil
}

func (o *Out) Backward(speed int) error {
	if err := o.SetDirection(Backward); err != nil {
		return err
	}

	if err := o.SetSpeed(speed); err != nil {
		return err
	}

	return nil
}

func (o *Out) Stop() error {
	if err := o.SetDirection(Stop); err != nil {
		return err
	}

	return nil
}
