package cli

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/c-bata/go-prompt"
)

var (
	errch    chan error
	outch    chan string
	commands []string = []string{
		"stop",
		"left",
		"right",
		"forward",
		"backward",
	}
)

func validCommand(cmd string) bool {
	if !slices.Contains(commands, cmd) {
		return false
	}

	return true
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Close the program"},
		{Text: "command", Description: "Send command t robot"},
		{Text: "help", Description: "Show help page"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(str string) {
	in := strings.Split(str, " ")

	switch in[0] {
	case "exit":
		os.Exit(0)
	case "command":
		if len(in) < 2 {
			fmt.Println("Usage: command [command]")
			return
		}
		cmd := strings.ToLower(in[1])

		if !validCommand(cmd) {
			fmt.Println("Unknown command")
			return
		}

		outch <- cmd
	default:
		fmt.Println("Unknown option: " + in[0])
		fmt.Println("Use help to learn more")
	}
}

func RobotCliFunc(out chan string, err chan error) {
	outch = out
	errch = err

	fmt.Println("Welcome to robot control system")
	defer fmt.Println("Bye!")
	prom := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionSwitchKeyBindMode(prompt.CommonKeyBind),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn:  func(*prompt.Buffer) { os.Exit(0) },
		}),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.Escape,
			Fn:  func(*prompt.Buffer) { os.Exit(0) },
		}),
	)
	prom.Run()
}
