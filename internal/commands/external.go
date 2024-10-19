package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func RunFromPath(cmd Cmd) *CmdResult {
	pathDirs := strings.Split(os.Getenv("PATH"), ":")

	relativePath := strings.Contains(cmd.Name, "/")
	if relativePath {
		wd, _ := os.Getwd()
		if !path.IsAbs(cmd.Name) {
			cmd.Name = path.Join(wd, cmd.Name)
		}
		res := runCommand(
			cmd.Name,
			cmd.Argv,
		)
		return &res
	}

	for _, dir := range pathDirs {
		binPath := path.Join(dir, cmd.Name)
		if _, err := os.Stat(binPath); err == nil {
			res := runCommand(binPath, cmd.Argv)
			return &res
		}
	}

	return nil
}

func runCommand(path string, args []string) CmdResult {
	cmd := exec.Command(path, args...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return CmdResult{
			ExitCode: 0,
			Output:   fmt.Sprintf("Unable to run the command: %v", err),
		}
	}

	code := 0
	if err := cmd.Wait(); err != nil {
		var exiterr *exec.ExitError

		if errors.As(err, &exiterr) {
			code = exiterr.ExitCode()
		}

		code = 1
	}

	return CmdResult{
		ExitCode: code,
		Output:   "",
	}
}
