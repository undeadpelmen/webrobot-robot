package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/undeadpelmen/webrobot-robot/hardware/robots"
)

var (
	errch chan error
	outch chan string

	status *string

	clear map[string]func()

	suggests []prompt.Suggest = []prompt.Suggest{
		{Text: "exit", Description: "Close the program"},
		{Text: "clear", Description: "Clear the screen"},
		{Text: "command", Description: fmt.Sprintf("Send command to robot\nSupported commands:\n%v", robots.Commands)},
		{Text: "status", Description: "Show status of program"},
		{Text: "help", Description: "Show help page"},
	}
)

func completer(d prompt.Document) []prompt.Suggest {
	s := suggests
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

		if !robots.ValidCommand(cmd) {
			fmt.Println("Unknown command")
			return
		}

		outch <- cmd

	case "hi":
		fmt.Println("Hello!!")

	case "help":
		fmt.Println("Helpl page")
		for _, cmd := range suggests {
			fmt.Println(cmd)
		}

	case "status":
		fmt.Println("Status")
		fmt.Println(*status)

	case "clear":
		value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
		if ok {                          //if we defined a clear func for that platform:
			value() //we execute it
		} else { //unsupported platform
			panic("Your platform is unsupported! I can't clear terminal screen :(")
		}

	default:
		fmt.Println("Unknown option: " + in[0])
		fmt.Println("Use help to learn more")
	}
}

func RobotCliFunc(stat *string, out chan string, err chan error) {
	outch = out
	errch = err

	status = stat

	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	fmt.Println("Welcome to robot control system")
	defer fmt.Println("Bye!")

	prom := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">> "),
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
