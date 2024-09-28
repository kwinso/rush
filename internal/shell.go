package internal

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"uni_shell/internal/commands"
	"uni_shell/internal/commands/std"
	"uni_shell/internal/history"

	"github.com/charmbracelet/lipgloss"
	"github.com/fiorix/go-readline"
)

func cleanup(playSound bool) {
	fmt.Println("running cleanup...")
	if playSound {
		std.PlayShutdownSound()
	}
	err := history.SaveHistoryFile()
	if err != nil {
		fmt.Println("Unable to write to history file: ", err)
	}
}

func RunShell() {
	// play sound flag parsing
	playSound := flag.Bool("silent", false, "Play startup/shutdown sound")

	highlightStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))
	errStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9"))

	history.InitHistory()
	sighup := false
	stdPrompt := highlightStyle.Render("rush> ")

	sigs := make(chan os.Signal, 5)
	signal.Notify(sigs, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		for s := range sigs {
			if s == syscall.SIGHUP {
				sighup = true
				go std.PlayBootSound()
				return
			}

			cleanup(*playSound)

			// interrupt signal received
			os.Exit(0)
		}
	}()

	if *playSound {
		go std.PlayBootSound()
	}

	var prevResult *commands.CmdResult = nil
	for {
		wd, _ := os.Getwd()
		wd = strings.Replace(wd, os.Getenv("HOME"), "~", 1)

		finalPrompt := fmt.Sprintf("(%v) %v", wd, stdPrompt)
		if prevResult != nil && prevResult.ExitCode != 0 {
			finalPrompt = errStyle.Render(fmt.Sprintf("[%d] ", prevResult.ExitCode)) + finalPrompt
		}
		input := readline.Readline(&finalPrompt)

		switch {
		case input == nil:
			if !sighup {
				cleanup(*playSound)
				os.Exit(0)
			}
		case *input != "": // Ignore blank lines
			input := *input
			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}
			history.SaveCmdToHistory(input)
			cmd := commands.ParseCommandFromString(input)

			var result *commands.CmdResult = nil
			if !cmd.IsStd {
				result = commands.RunFromPath(cmd)
			} else {
				result = std.RunStdCmd(cmd)
			}

			if result == nil {
				result = &commands.CmdResult{
					ExitCode: 127,
					Output:   fmt.Sprintf("No command found: %v", cmd.Name),
				}
			}

			output := result.Output
			if result.Output != "" {
				fmt.Println(output)
			}

			if result.Flags.Exit {
				cleanup(*playSound)
				os.Exit(result.ExitCode)
			}

			prevResult = result
		}

		if sighup {
			sighup = false
			fmt.Println(highlightStyle.Render("rush configuration reloaded..."))
		}
	}
}
