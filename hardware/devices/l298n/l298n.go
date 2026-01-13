package l298n

import (
	"errors"

	"periph.io/x/conn/v3/gpio"
)

type L298N struct {
	out1 *Out
	out2 *Out
}

func New(out1, out2 *Out) (*L298N, error) {
	if out1 == nil || out2 == nil {
		return nil, errors.New("Out1 or Out2 cannot be nil")
	}

	dev := &L298N{out1, out2}

	return dev, nil
}

func NewFromPins(enaOut1Pin, enbOut1Pin, speedOut1Pin, enaOut2Pin, enbOut2Pin, speedOut2Pin gpio.PinOut, speed int) (*L298N, error) {
	if enaOut1Pin == nil || enbOut1Pin == nil || enaOut2Pin == nil || enbOut2Pin == nil || speedOut1Pin == nil || speedOut2Pin == nil {
		return nil, errors.New("Pins can not be nil")
	}

	out1 := &Out{
		in1Pin:    enaOut1Pin,
		in2Pin:    enbOut1Pin,
		speedPin:  speedOut1Pin,
		direction: Stop,
	}

	if err := out1.SetSpeed(speed); err != nil {
		return nil, err
	}

	out2 := &Out{
		in1Pin:    enaOut2Pin,
		in2Pin:    enbOut2Pin,
		speedPin:  speedOut2Pin,
		direction: Stop,
	}

	if err := out2.SetSpeed(speed); err != nil {
		return nil, err
	}

	dev := &L298N{out1, out2}

	return dev, nil
}

func (l *L298N) Out1() *Out {
	return l.out1
}

func (l *L298N) Out2() *Out {
	return l.out2
}
