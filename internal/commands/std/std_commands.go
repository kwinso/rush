package std

import (
	"uni_shell/internal/commands"
)

var stdCommands = map[string]commands.CmdHandler{
	":e":    HandleEnvCommand,
	":?":    HandleAboutCommand,
	":q":    HandleExitCommand,
	":noop": HandleNoopCommand,
	":cd":   HandleCdCommand,
	":part": HandleShowPartitionsCommand,
	":omg":  PlayOmg,
	":mem":  HandleMemCommand,
	":maps": HandleMapsCommand,
}

func RunStdCmd(cmd commands.Cmd) *commands.CmdResult {
	handler, ok := stdCommands[cmd.Name]

	if !ok {
		return nil
	}

	res := handler(cmd)
	return &res
}
