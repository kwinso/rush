package std

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"uni_shell/internal/commands"

	"github.com/dustin/go-humanize"
)

func HandleShowPartitionsCommand(cmd commands.Cmd) commands.CmdResult {
	if cmd.Argc != 1 {
		return commands.CmdResult{
			ExitCode: 1,
			Output: `Usage: :part <disk_name>

<disk_name> format: disk name without namespace, e.g. nvme0`,
		}
	}

	requestedDeviceName := cmd.Argv[0]
	devices, _ := filepath.Glob(path.Join("/sys/block/", requestedDeviceName) + "*")

	if len(devices) == 0 {
		avaliableDisks, _ := filepath.Glob("/sys/block/*")
		var buff bytes.Buffer
		buff.WriteString("Unable to find disk device with name " + requestedDeviceName + ".\nAvailable devices:")
		for _, d := range avaliableDisks {
			buff.WriteString(fmt.Sprintf("\n%v", path.Base(d)))
		}
		return commands.CmdResult{
			ExitCode: 127,
			Output:   buff.String(),
		}
	}

	var buff bytes.Buffer

	// Find only the first one for the sake of simplicity
	devicePath := devices[0]
	deviceName := path.Base(devicePath)
	if deviceName != requestedDeviceName {
		buff.WriteString("* Found device with name " + requestedDeviceName + ", which is not an exact match for the disk name.\n")
	}
	partitions, _ := filepath.Glob(path.Join(devicePath, deviceName) + "*")

	buff.WriteString(fmt.Sprintf("Disk %v, Size: %v\n\nPartitions:\n", requestedDeviceName, humanize.Bytes(readSize(devicePath))))
	for _, p := range partitions {
		buff.WriteString(fmt.Sprintf(
			"%v\t%v\n",
			filepath.Base(p), humanize.Bytes(readSize(p)),
		))
	}

	return commands.CmdResult{
		ExitCode: 0,
		Output:   buff.String(),
	}
}

func readSize(devicePath string) uint64 {
	diskSizeBytes, _ := os.ReadFile(filepath.Join(devicePath, "size"))

	size, err := strconv.Atoi(strings.TrimSpace(string(diskSizeBytes)))
	if err != nil {
		fmt.Println("err: ", err)
		return 0
	}

	// Size is returned in 512 byte sectors
	return uint64(512 * size)
}
