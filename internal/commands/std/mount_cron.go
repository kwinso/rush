package std

import (
	"uni_shell/internal/commands"
	"uni_shell/internal/cronfs"
)

func HandleMountCronCommand(cmd commands.Cmd) commands.CmdResult {
	manager, err := cronfs.GetCronFSManager()
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   err.Error(),
		}
	}

	path := "./cronfs"
	err = manager.Mount(path)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   err.Error(),
		}
	}

	return commands.CmdResult{
		ExitCode: 0,
		Output:   "Mounted cronfs to " + path,
	}
}
