package l298n

import (
	"errors"
	"sync"

	"periph.io/x/conn/v3/gpio"
)

type Direction string

const (
	Forward  Direction = "Forward"
	Backward Direction = "Backward"
	Stop     Direction = "Stop"
)

type Out struct {
	mu        sync.Mutex
	direction Direction
	enaPin    gpio.PinOut
	enbPin    gpio.PinOut
	speedPin  gpio.PinOut
	speed     int
}

func NewOut(enaPin, enbPin, speedPin gpio.PinOut, speed int) (*Out, error) {
	if enaPin == nil || enbPin == nil || speedPin == nil {
		return nil, errors.New("Pins can not be nil")
	}

	dev := &Out{
		direction: Stop,
		enaPin:    enaPin,
		enbPin:    enbPin,
		speedPin:  speedPin,
	}

	if err := dev.SetSpeed(speed); err != nil {
		return nil, err
	}

	return dev, nil
}

func (o *Out) update(speed int, direction Direction) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if speed > 255 || speed < 0 {
		return errors.New("Speed must be between 0 and 255")
	}
	o.speed = speed

	if direction != Forward && direction != Stop && direction != Backward {
		return errors.New("Unknown direcction")
	}
	o.direction = direction

	var errA, errB error
	switch o.direction {
	case Forward:
		errA = o.enaPin.Out(gpio.High)
		errB = o.enbPin.Out(gpio.Low)
	case Backward:
		errA = o.enaPin.Out(gpio.Low)
		errB = o.enbPin.Out(gpio.High)
	case Stop:
		errA = o.enaPin.Out(gpio.Low)
		errB = o.enbPin.Out(gpio.Low)
	}

	if errA != nil {
		return errA
	}

	if errB != nil {
		return errB
	}

	if err := o.speedPin.PWM(gpio.Duty(o.speed*100/255), 0); err != nil {
		return err
	}

	return nil
}

func (o *Out) Speed() int {
	return o.speed
}

func (o *Out) SetSpeed(speed int) error {
	if err := o.update(speed, o.direction); err != nil {
		return err
	}

	return nil
}

func (o *Out) Direction() Direction {
	return o.direction
}

func (o *Out) SetDirection(direction Direction) error {
	if err := o.update(o.speed, direction); err != nil {
		return err
	}

	return nil
}

func (o *Out) Forward(speed int) error {
	if err := o.update(speed, Forward); err != nil {
		return err
	}

	return nil
}

func (o *Out) Backward(speed int) error {
	if err := o.update(speed, Backward); err != nil {
		return err
	}

	return nil
}

func (o *Out) Stop() error {
	if err := o.update(o.speed, Stop); err != nil {
		return err
	}

	return nil
}
