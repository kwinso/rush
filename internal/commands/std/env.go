package std

import (
	"fmt"
	"os"
	"uni_shell/internal/commands"
)

func HandleEnvCommand(cmd commands.Cmd) commands.CmdResult {
	if cmd.Argc == 1 {
		env := os.Getenv(cmd.Argv[0])
		return commands.CmdResult{
			Output: env,
		}
	}

	if cmd.Argc == 2 {
		name := cmd.Argv[0]
		val := cmd.Argv[1]
		err := os.Setenv(name, val)
		if err != nil {
			return commands.CmdResult{
				ExitCode: 1,
				Output:   fmt.Sprintf("\\e: Unable to set env variable: %v\n", err),
			}
		}
		return commands.CmdResult{
			ExitCode: 0,
			Output:   fmt.Sprintf("%v = %v", name, val),
		}
	}

	return commands.CmdResult{
		ExitCode: 1,
		Output: `Invalid args. Usage: 
:e VAR - show env var VAR
:e VAR VAL - set env var VAR to value VAL`,
	}
}
