package main

import (
	"fmt"
	"os"

	"github.com/anschnapp/pomodorofactory/pkg/commandinput"
	"github.com/anschnapp/pomodorofactory/pkg/motivationcloud"
	"github.com/anschnapp/pomodorofactory/pkg/pomodorobuild"
	"github.com/anschnapp/pomodorofactory/pkg/render"
	"github.com/anschnapp/pomodorofactory/pkg/status"
	"github.com/anschnapp/pomodorofactory/pkg/view"
	"golang.org/x/term"
)

func main() {
	// Put terminal in raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to set raw mode: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Enter alternate screen buffer
	fmt.Print("\033[?1049h")
	defer fmt.Print("\033[?1049l")

	// Hide cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	// Build view
	var pomodorobuild render.Renderable = pomodorobuild.MakePomodoro()
	var motivationcloud render.Renderable = motivationcloud.MakeMotivationcloud()
	var status render.Renderable = status.MakeStatus()
	var commandinput render.Renderable = commandinput.MakeCommandinput()
	var v *view.View = view.MakeView(pomodorobuild, motivationcloud, status, commandinput)

	// Read input in a goroutine
	inputCh := make(chan byte)
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				inputCh <- buf[0]
			}
			if err != nil {
				close(inputCh)
				return
			}
		}
	}()

	// Initial render
	v.Render()
	v.Print()

	// Event loop
	for b := range inputCh {
		switch b {
		case 'q', 0x03: // 'q' or Ctrl+C
			return
		}
		v.Render()
		v.Print()
	}
}
