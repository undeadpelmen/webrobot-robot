package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/rs/zerolog"
	"github.com/undeadpelmen/webrobot-robot/internal/interfaces"
)

type Service struct {
	robotService interfaces.RobotController
	logger       zerolog.Logger
}

func NewService(robotService interfaces.RobotController, logger zerolog.Logger) *Service {
	return &Service{
		robotService: robotService,
		logger:       logger,
	}
}

func (s *Service) StartInteractive() {
	s.logger.Info().Msg("Starting interactive CLI")
	fmt.Println("Robot Control CLI - Type 'help' for commands")

	p := prompt.New(
		s.executor,
		s.completer,
		prompt.OptionTitle("Robot Control"),
		prompt.OptionPrefix("robot> "),
		prompt.OptionLivePrefix(s.livePrefix),
	)

	p.Run()
}

func (s *Service) StartNonInteractive() {
	s.logger.Info().Msg("Starting non-interactive CLI")
	scan := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("robot> ")
		if !scan.Scan() {
			break
		}

		command := strings.TrimSpace(scan.Text())
		if command == "" {
			continue
		}

		if err := s.executeCommand(command); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func (s *Service) executor(in string) {
	if err := s.executeCommand(in); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (s *Service) executeCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "help", "h":
		s.showHelp()
	case "status":
		s.showStatus()
	case "move", "m":
		return s.handleMove(parts[1:])
	case "stop", "s":
		return s.handleStop()
	case "speed":
		return s.handleSpeed(parts[1:])
	case "set-speed":
		return s.handleSetSpeed(parts[1:])
	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		return fmt.Errorf("unknown command: %s. Type 'help' for available commands", cmd)
	}

	return nil
}

func (s *Service) handleMove(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: move <direction> [speed]")
	}

	direction := strings.ToLower(args[0])
	speed := s.robotService.GetSpeed()

	if len(args) > 1 {
		if newSpeed, err := strconv.Atoi(args[1]); err == nil {
			speed = newSpeed
		}
	}

	if err := s.robotService.Move(direction, speed); err != nil {
		return fmt.Errorf("failed to move: %w", err)
	}

	fmt.Printf("Moving %s at speed %d\n", direction, speed)
	return nil
}

func (s *Service) handleStop() error {
	if err := s.robotService.Stop(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	fmt.Println("Robot stopped")
	return nil
}

func (s *Service) handleSpeed(args []string) error {
	speed := s.robotService.GetSpeed()
	fmt.Printf("Current speed: %d\n", speed)
	return nil
}

func (s *Service) handleSetSpeed(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: set-speed <speed>")
	}

	speed, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid speed: %v", err)
	}

	if err := s.robotService.SetSpeed(speed); err != nil {
		return fmt.Errorf("failed to set speed: %w", err)
	}

	fmt.Printf("Speed set to %d\n", speed)
	return nil
}

func (s *Service) showHelp() {
	fmt.Println(`
Available commands:
  help, h              - Show this help message
  status               - Show robot status
  move, m <dir> [speed] - Move robot in direction (forward/backward/left/right)
  stop, s              - Stop the robot
  speed                - Show current speed
  set-speed <speed>    - Set robot speed (0-255)
  exit, quit           - Exit the CLI

Directions:
  forward, f           - Move forward
  backward, b          - Move backward
  left, l              - Turn left
  right, r             - Turn right
`)
}

func (s *Service) showStatus() {
	status := s.robotService.Status()
	speed := s.robotService.GetSpeed()

	fmt.Printf("Robot Status: %s\n", status)
	fmt.Printf("Current Speed: %d\n", speed)
}

func (s *Service) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "help", Description: "Show help"},
		{Text: "status", Description: "Show robot status"},
		{Text: "move", Description: "Move robot"},
		{Text: "stop", Description: "Stop robot"},
		{Text: "speed", Description: "Show current speed"},
		{Text: "set-speed", Description: "Set robot speed"},
		{Text: "exit", Description: "Exit CLI"},
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (s *Service) livePrefix() (string, bool) {
	status := s.robotService.Status()
	return fmt.Sprintf("[%s] robot> ", status), true
}
