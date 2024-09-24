package std

import (
	"uni_shell/internal/commands"
)

func HandleAboutCommand(cmd commands.Cmd) commands.CmdResult {
	return commands.CmdResult{
		ExitCode: 0,
		Output: `Ruslan's Useless SHell (rush), 2024
Please never use it if you're not insane.
And don't type :omg, never.`,
	}
}
