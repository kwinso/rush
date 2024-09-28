package std

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"time"
	"uni_shell/internal/commands"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

//go:embed assets/windows_shutdown_sound.mp3
var shutdownSound []byte

//go:embed assets/windows_boot_sound.mp3
var bootSound []byte

func PlayShutdownSound() {
	// create io.Reader for omgAudio
	reader := bytes.NewReader(shutdownSound)
	d, _ := mp3.NewDecoder(reader)

	c, ready, _ := oto.NewContext(d.SampleRate(), 2, 2)
	<-ready

	p := c.NewPlayer(d)
	p.Play()
	for {
		time.Sleep(time.Second)
		if !p.IsPlaying() {
			break
		}
	}
}

func PlayBootSound() {
	// create io.Reader for omgAudio
	reader := bytes.NewReader(bootSound)
	d, _ := mp3.NewDecoder(reader)

	c, ready, _ := oto.NewContext(d.SampleRate(), 2, 2)
	<-ready

	p := c.NewPlayer(d)
	p.Play()
	for {
		time.Sleep(time.Second)
		if !p.IsPlaying() {
			break
		}
	}
}

func HandleExitCommand(cmd commands.Cmd) commands.CmdResult {
	code := 0
	var err error
	if cmd.Argc != 0 {
		code, err = strconv.Atoi(cmd.Argv[0])
		if err != nil {
			return commands.CmdResult{
				ExitCode: 1,
				Output:   fmt.Sprintf("Unable to parse exit code: %v", err),
			}
		}
	}

	return commands.CmdResult{
		ExitCode: code,
		Flags:    commands.CmdResultFlags{Exit: true},
	}
}
