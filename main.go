package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/undeadpelmen/webrobot-robot/cli"
	"github.com/undeadpelmen/webrobot-robot/hardware/devices/l298n"
	"github.com/undeadpelmen/webrobot-robot/hardware/robots"
	"github.com/undeadpelmen/webrobot-robot/hardware/robots/crawler"
	"github.com/undeadpelmen/webrobot-robot/web/http"
	"periph.io/x/conn/v3/gpio/gpiotest"
)

var (
	logger zerolog.Logger
	status string
)

func main() {
	useCli := flag.Bool("c", false, "use cli instead of hardware")
	useHttp := flag.Bool("h", false, "use http instead of hardware")
	debug := flag.Bool("d", false, "enable debug logging")

	flag.Parse()

	logFile, err := os.OpenFile(filepath.Join(os.TempDir(), "/webrobot/robot.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if _, ok := err.(*os.PathError); ok {
		fmt.Println("Failed to open log file")

		if err := os.Mkdir(filepath.Join(os.TempDir(), "/webrobot"), 0775); err != nil {
			panic(err)
		}

		file, err := os.Create(filepath.Join(os.TempDir(), "/webrobot/robot.log"))
		if err != nil {
			panic(err)
		}

		logFile = file
	} else if err != nil {
		panic(err)
	}

	fileWriter := zerolog.ConsoleWriter{Out: logFile, NoColor: true}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	writers := io.MultiWriter(consoleWriter, fileWriter)

	logger = zerolog.New(writers).With().Timestamp().Logger()

	logger.Info().Msg("Start web-based robot control system")

	logger.Debug().Msg("Create driver")
	driver, err := l298n.NewFromPins(&gpiotest.Pin{}, &gpiotest.Pin{}, &gpiotest.Pin{}, &gpiotest.Pin{}, &gpiotest.Pin{}, &gpiotest.Pin{}, 255)
	if err != nil {
		logger.Panic().Err(err).Msg("")
	}

	logger.Debug().Msg("Create crawler")
	robot, err := crawler.NewFromDriver(driver)
	if err != nil {
		logger.Panic().Err(err).Msg("")
	}

	status = robot.Status()

	logger.Debug().Msg("Creating chanels")

	// Creating channels to transporting data in system
	cmdchan := make(chan string)     // Chan for robot command
	errchan := make(chan error, 10)  // Chan for error handling
	logchan := make(chan string, 10) // Chan for logging

	logger.Debug().Msg("Start robot gorutine")
	go robots.RobotControlFunc(&status, cmdchan, errchan, logchan, robot)

	if *useCli {
		// Start Cli goroutine

		logger.Debug().Msg("Start cli")
		go cli.RobotCliFunc(&status, cmdchan, errchan)
	}

	if *useHttp {
		logger.Debug().Msg("Start http")
		go http.RobotHttpFunc(&status, cmdchan, errchan)
	}

	// Create CTRL + C shortcut to close the program
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c

		os.Exit(0)
	}()

	logger.Debug().Msg("Start loop chan reading")
	for {
		select {
		case err := <-errchan:
			logger.Panic().Err(err).Msg("")

		case msg := <-logchan:
			logger.Info().Msg(msg)
		}
	}
}
