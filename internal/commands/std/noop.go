package std

import (
	"uni_shell/internal/commands"
)

func HandleNoopCommand(cmd commands.Cmd) commands.CmdResult {
	return commands.CmdResult{
		ExitCode: 0,
	}
}
