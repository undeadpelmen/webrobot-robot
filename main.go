package main

import (
	"log"
	"os"

	"github.com/undeadpelmen/webrobot-robot/cli"
	"github.com/undeadpelmen/webrobot-robot/hardware/devices/l298n"
	"github.com/undeadpelmen/webrobot-robot/hardware/robots"
	"github.com/undeadpelmen/webrobot-robot/hardware/robots/crawler"
	"periph.io/x/conn/v3/gpio/gpiotest"
)

var (
	lg *log.Logger
)

func main() {
	logout, err := os.OpenFile("log/robot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	if _, err := logout.Write([]byte("\n\n")); err != nil {
		panic(err)
	}

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

	lg.Println("Start robot gorutine")
	go robots.RobotControlFunc(cmdchan, errchan, logchan, robot)

	go cli.RobotCliFunc(cmdchan, errchan)

	lg.Println("Start for loop chan reading")
	for {
		select {
		case err := <-errchan:
			lg.Panicln(err)

		case msg := <-logchan:
			lg.Print(msg)
		}
	}
}
