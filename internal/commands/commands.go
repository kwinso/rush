package commands

import (
	"strings"
)

type Cmd struct {
	Name  string
	Argv  []string
	Argc  int
	IsStd bool
}

type CmdResultFlags struct {
	Exit bool
}

type CmdResult struct {
	ExitCode int
	Output   string
	Flags    CmdResultFlags
}

type CmdHandler = func(cmd Cmd) CmdResult

func ParseCommandFromString(s string) Cmd {
	args := strings.Split(s, " ")

	return Cmd{
		Name:  args[0],
		IsStd: args[0][0] == ':',
		Argv:  args[1:],
		Argc:  len(args) - 1,
	}
}
