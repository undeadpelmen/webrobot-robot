package crawler

import (
	"errors"
	"sync"

	"github.com/undeadpelmen/webrobot-robot/hardware/devices/l298n"
)

type Crawler struct {
	left  Front
	right Front
	speed int
	mu    sync.Mutex
}

type Front interface {
	Forward(int) error
	Backward(int) error
	Stop() error
}

func New(front, right Front) (*Crawler, error) {
	if front == nil {
		return nil, errors.New("Front can not be nil")
	}

	dev := &Crawler{
		left:  front,
		right: right,
	}

	return dev, nil
}

func NewFromDriver(driver *l298n.L298N) (*Crawler, error) {
	if driver == nil {
		return nil, errors.New("Driver can not be nil")
	}

	if driver.Out1() == nil || driver.Out2() == nil {
		return nil, errors.New("Fronts can not be nil")
	}

	dev := &Crawler{
		left:  driver.Out1(),
		right: driver.Out2(),
	}

	return dev, nil
}

func (c *Crawler) SetSpeed(speed int) error {
	if speed > 255 || speed < 0 {
		return errors.New("Speed must be between 0 and 255")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.speed = speed

	return nil
}

func (c *Crawler) Forward() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.left.Forward(c.speed); err != nil {
		return err
	}
	if err := c.right.Forward(c.speed); err != nil {
		return err
	}

	return nil
}

func (c *Crawler) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.left.Stop(); err != nil {
		return err
	}
	if err := c.right.Stop(); err != nil {
		return err
	}

	return nil
}

func (c *Crawler) Backward() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.left.Backward(c.speed); err != nil {
		return err
	}

	if err := c.right.Backward(c.speed); err != nil {
		return err
	}

	return nil
}

func (c *Crawler) Left() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.left.Forward(c.speed); err != nil {
		return err
	}

	if err := c.right.Stop(); err != nil {
		return err
	}

	return nil
}

func (c *Crawler) Right() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.left.Stop(); err != nil {
		return nil
	}

	if err := c.right.Forward(c.speed); err != nil {
		return err
	}

	return nil
}
