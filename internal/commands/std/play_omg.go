package std

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"time"
	"uni_shell/internal/commands"
)

//go:embed assets/omg.mp3
var omgAudio []byte

func PlayOmg(cmd commands.Cmd) commands.CmdResult {
	// run external command paplay with os.exec

	// create io.Reader for omgAudio
	reader := bytes.NewReader(omgAudio)
	d, _ := mp3.NewDecoder(reader)

	c, ready, _ := oto.NewContext(d.SampleRate(), 2, 2)
	<-ready

	p := c.NewPlayer(d)
	defer p.Close()
	p.Play()

	fmt.Print("What")
	time.Sleep(time.Millisecond * 300)
	fmt.Print(" the")
	time.Sleep(time.Millisecond * 1000)
	fmt.Print(" h")
	for i := 0; i <= 2000; i += 50 {
		// random lipgloss color
		var style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("#%x", i+600)))

		fmt.Printf(style.Render("e"))
		time.Sleep(time.Millisecond * 50)
	}
	fmt.Println("ll")

	fmt.Println("Oooh my god")
	time.Sleep(time.Millisecond * 1700)
	fmt.Println("Nooo waaaYaYyayaaa")
	time.Sleep(time.Millisecond * 1900)

	return commands.CmdResult{
		ExitCode: 0,
	}
}
