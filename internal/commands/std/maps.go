package std

import (
	"fmt"
	"os"
	"strconv"
	"uni_shell/internal/commands"
)

func HandleMapsCommand(cmd commands.Cmd) commands.CmdResult {
	if cmd.Argc != 1 {
		return commands.CmdResult{
			ExitCode: 1,
			Output: `Usage: :maps <pid>

Prints the memory maps of the process.`,
		}
	}

	pid, err := strconv.ParseInt(cmd.Argv[0], 10, 64)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   "Invalid PID",
		}
	}

	// Open the memory file of the process
	mapsFilePath := fmt.Sprintf("/proc/%d/maps", pid)
	data, err := os.ReadFile(mapsFilePath)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Failed to open maps file: %v", err),
		}
	}

	return commands.CmdResult{
		ExitCode: 0,
		Output:   string(data),
	}
}
