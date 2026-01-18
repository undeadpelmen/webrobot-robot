package robots

import (
	"errors"
	"fmt"
	"strings"
)

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
