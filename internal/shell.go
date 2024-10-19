package internal

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"uni_shell/internal/commands"
	"uni_shell/internal/commands/std"
	"uni_shell/internal/cronfs"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/chzyer/readline"
)

// TODO: This cleanup function does not guarantee that readlinen will exist safely
func cleanup(playSound bool) {
	log.Debug("Cleaning up...")
	if playSound {
		std.PlayShutdownSound()
	}

	manager, err := cronfs.GetCronFSManager()
	if err != nil {
		fmt.Println("Failed to cleanup cronfs:", err)
	}

	if manager.IsMounted() {
		log.Debug("Found cronfs mounted, unmounting...")
		err = manager.Unmount()
		if err != nil {
			fmt.Println("Failed to cleanup cronfs:", err)
		}
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
			}
		}
	}()
}

func RunShell() {
	isDebug := os.Getenv("DEBUG") == "1"
	if isDebug {
		log.SetLevel(log.DebugLevel)
	}

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
		HistoryFile:     filepath.Join(usr.HomeDir, ".rush_history"),
		InterruptPrompt: "^C",
		EOFPrompt:       ":q",
		VimMode:         true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	go handleConfigReload(l)

	// play sound flag parsing
	silent := flag.Bool("silent", false, "Do not play any sounds")

	fmt.Println("Rush command line arguments:")
	flag.PrintDefaults()

	flag.Parse()
	if !*silent {
		go std.PlayBootSound()
	}

	var prevResult *commands.CmdResult = nil
	log.Debug("Shell started successfully")
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

		if result.Output != "" {
			fmt.Println(result.Output)
		}

		if result.Flags.Exit {
			log.Debug("Command required exit, leaving shell")
			cleanup(!*silent)
			os.Exit(result.ExitCode)
		}

		prevResult = result
	}

	cleanup(!*silent)
}
