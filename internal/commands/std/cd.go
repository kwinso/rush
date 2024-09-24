package std

import (
	"os"
	"strings"
	"uni_shell/internal/commands"
)

func HandleCdCommand(cmd commands.Cmd) commands.CmdResult {
	path := ""
	if cmd.Argc != 0 {
		path = cmd.Argv[0]
	}

	if path == "" {
		path = os.Getenv("HOME")
	}

	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}

	err := os.Chdir(path)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   err.Error(),
		}
	}

	return commands.CmdResult{
		ExitCode: 0,
	}
}
