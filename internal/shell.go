package internal

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"uni_shell/internal/commands"
	"uni_shell/internal/commands/std"

	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
)

// TODO: This cleanup function does not guarantee that readlinen will exist safely
func cleanup(playSound bool) {
	if playSound {
		std.PlayShutdownSound()
	}
}

func buildPrompt(user *user.User, prevResult *commands.CmdResult, highlightStyle lipgloss.Style, errStyle lipgloss.Style) string {
	stdPrompt := highlightStyle.Render("» ")

	wd, _ := os.Getwd()
	wd = strings.Replace(wd, user.HomeDir, "~", 1)

	finalPrompt := fmt.Sprintf("(%v) %v", wd, stdPrompt)
	if prevResult != nil && prevResult.ExitCode != 0 {
		finalPrompt = errStyle.Render(fmt.Sprintf("[%d] ", prevResult.ExitCode)) + finalPrompt
	}

	return finalPrompt
}

func handleConfigReload(l *readline.Instance) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	go func() {
		for s := range sigs {
			if s == syscall.SIGHUP {
				fmt.Println("\nrush configuration reloaded...")
				l.Refresh()
				return
			}
		}
	}()
}

func RunShell() {
	highlightStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))
	errStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9"))

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	l, err := readline.NewEx(&readline.Config{
		Prompt:            "\033[31m\033[0m ",
		HistoryFile:       filepath.Join(usr.HomeDir, ".rush_history"),
		InterruptPrompt:   "^C",
		EOFPrompt:         ":q",
		HistorySearchFold: true,
		VimMode:           true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	setPasswordCfg := l.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})

	log.SetOutput(l.Stderr())

	go handleConfigReload(l)

	// play sound flag parsing
	bootSound := flag.Bool("silent", false, "Play startup/shutdown sound")
	if *bootSound {
		go std.PlayBootSound()
	}

	var prevResult *commands.CmdResult = nil
	for {
		l.SetPrompt(buildPrompt(
			usr,
			prevResult,
			highlightStyle,
			errStyle,
		))
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// I know that this is a hack, but I don't wanna refactor so that comamnds have access to the readline instance
		if line == ":exitvim" {
			l.SetVimMode(false)
			continue
		}
		if line == ":vim" {
			l.SetVimMode(true)
			continue
		}

		cmd := commands.ParseCommandFromString(line)
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
			cleanup(*bootSound)
			os.Exit(result.ExitCode)
		}

		prevResult = result
	}
}
