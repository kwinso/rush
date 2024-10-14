package commands

import (
	"regexp"
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
	re := regexp.MustCompile(`["|']([^"']+)["|']|(\S+)`)
	matches := re.FindAllStringSubmatch(s, -1)

	var args []string
	for _, match := range matches {
		if match[1] != "" {
			// If the first capturing group (quoted string) is not empty, use it
			args = append(args, match[1])
		} else if match[2] != "" {
			// If the second capturing group (unquoted word) is not empty, use it
			args = append(args, match[2])
		}
	}

	return Cmd{
		Name:  args[0],
		IsStd: args[0][0] == ':',
		Argv:  args[1:],
		Argc:  len(args) - 1,
	}
}
