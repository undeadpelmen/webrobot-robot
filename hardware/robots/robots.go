package robots

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	Commands []string = []string{
		"stop",
		"forward",
		"backward",
		"left",
		"right",
	}
)

func ValidCommand(cmd string) bool {
	if !slices.Contains(Commands, cmd) {
		return false
	}

	return true
}

type Robot interface {
	SetSpeed(int) error
	Forward() error
	Stop() error
	Backward() error
	Left() error
	Right() error
	Status() string
}

func RobotControlFunc(status *string, cmd chan string, errch chan error, logchan chan string, robot Robot) {
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

		*status = robot.Status()
	}
}
