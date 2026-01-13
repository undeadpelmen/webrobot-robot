package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/undeadpelmen/webrobot-robot/hardware/devices/l298n"
	"github.com/undeadpelmen/webrobot-robot/hardware/robots/crawler"
	"periph.io/x/conn/v3/gpio/gpiotest"
)

var (
	lg *log.Logger
)

func main() {
	logout, err := os.OpenFile("robot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	logout.Write([]byte("\n\n"))

	lg = log.New(logout, "", log.Ltime|log.Lshortfile)
	defer logout.Close()
	lg.Println("Start web-based robot control system")

	lg.Println("Create driver")
	driver, err := l298n.NewFromPins(&gpiotest.Pin{}, &gpiotest.Pin{}, &gpiotest.Pin{}, &gpiotest.Pin{}, &gpiotest.Pin{}, &gpiotest.Pin{}, 255)
	if err != nil {
		lg.Panicln(err)
	}

	lg.Println("Create crawler")
	robot, err := crawler.NewFromDriver(driver)
	if err != nil {
		lg.Panicln(err)
	}

	lg.Println("Creating chanels")

	cmdchan := make(chan string)
	errchan := make(chan error, 10)
	logchan := make(chan string, 10)
	in := make(chan string)

	lg.Println("Start robot gorutine")
	go RobotControlFunc(cmdchan, errchan, logchan, robot)

	go func() {
		var cmd string
		for {
			fmt.Scan(&cmd)
			in <- cmd
		}
	}()

	lg.Println("Start for loop chan reading")
	for {
		select {
		case cmd := <-in:
			lg.Println("Command:", cmd)
			cmdchan <- cmd

		case err := <-errchan:
			lg.Panicln(err)

		case msg := <-logchan:
			lg.Print(msg)
		}
	}
}

type Robot interface {
	SetSpeed(int) error
	Forward() error
	Stop() error
	Backward() error
	Left() error
	Right() error
}

func RobotControlFunc(cmd chan string, errch chan error, logchan chan string, robot Robot) {
	for {
		command := strings.ToLower(<-cmd)

		switch command {
		case "forward":
			if err := robot.Forward(); err != nil {
				errch <- err
			}
			logchan <- fmt.Sprintln("Forward")

		case "stop":
			if err := robot.Stop(); err != nil {
				errch <- err
			}
			logchan <- fmt.Sprintln("Stop")

		case "backward":
			if err := robot.Backward(); err != nil {
				errch <- err
			}
			logchan <- fmt.Sprintln("Backward")

		case "left":
			if err := robot.Left(); err != nil {
				errch <- err
			}
			logchan <- fmt.Sprintln("Left")

		case "right":
			if err := robot.Right(); err != nil {
				errch <- err
			}
			logchan <- fmt.Sprintln("Right")

		default:
			errch <- errors.New("Wrong command: " + command)
		}
	}
}
