package std

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"uni_shell/internal/commands"
)

func HandleMemCommand(cmd commands.Cmd) commands.CmdResult {
	user, err := user.Current()
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Unable to get current user: %v", err),
		}
	}

	if user.Uid != "0" {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Rerun the rush shell as root in order to call this command"),
		}
	}

	if cmd.Argc != 2 {
		return commands.CmdResult{
			ExitCode: 1,
			Output: fmt.Sprintf(`Usage: :mem <pid> <start_addr> <end_addr>

<pid> format: pid 
<region> format: start and end address in format of :maps output (first column)
`),
		}
	}

	pid := cmd.Argv[0]

	// Convert pid to int
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Invalid PID: %v", err),
		}
	}

	region := cmd.Argv[1]
	parts := strings.Split(region, "-")
	if len(parts) != 2 {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Invalid region: %v", err),
		}
	}

	// Specify the address and size to read
	startAddr, err := strconv.ParseUint(parts[0], 16, 64)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Invalid address: %v", err),
		}
	}
	endAddr, err := strconv.ParseUint(parts[1], 16, 64)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Invalid address: %v", err),
		}
	}

	// Open the memory file of the process
	memFilePath := fmt.Sprintf("/proc/%d/mem", pidInt)
	memFile, err := os.Open(memFilePath)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Failed to open memory file: %v", err),
		}
	}
	defer memFile.Close()

	size := endAddr - startAddr
	if size < 0 {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   "Invalid size: must be positive",
		}
	}

	// Seek to the address
	_, err = memFile.Seek(int64(startAddr), 0)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Failed to seek in memory file: %v", err),
		}
	}

	// Read memory
	data := make([]byte, size)
	n, err := memFile.Read(data)
	if err != nil {
		return commands.CmdResult{
			ExitCode: 1,
			Output:   fmt.Sprintf("Failed to read memory: %v", err),
		}
	}

	return commands.CmdResult{
		ExitCode: 0,
		Output:   printHexDump(data[:n]),
	}
}

func printHexDump(data []byte) string {
	var buff bytes.Buffer
	const bytesPerLine = 16
	for i := 0; i < len(data); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(data) {
			end = len(data)
		}

		// Print the offset
		buff.WriteString(fmt.Sprintf("%08x  ", i))

		// Print the hex values
		for j := i; j < end; j++ {
			buff.WriteString(fmt.Sprintf("%02x ", data[j]))
		}

		// Fill the remaining space with spaces
		for j := end; j < i+bytesPerLine; j++ {
			buff.WriteString("   ")
		}

		// Print the ASCII representation
		buff.WriteString(" |")
		for j := i; j < end; j++ {
			if data[j] >= 32 && data[j] <= 126 {
				buff.WriteString(fmt.Sprintf("%c", data[j]))
			} else {
				buff.WriteString(".")
			}
		}
		buff.WriteString("|\n")
	}

	return buff.String()
}
