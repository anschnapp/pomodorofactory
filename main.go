package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/anschnapp/pomodorofactory/pkg/commandinput"
	"github.com/anschnapp/pomodorofactory/pkg/motivationcloud"
	"github.com/anschnapp/pomodorofactory/pkg/pomodorobuild"
	"github.com/anschnapp/pomodorofactory/pkg/render"
	"github.com/anschnapp/pomodorofactory/pkg/status"
	"github.com/anschnapp/pomodorofactory/pkg/timer"
	"github.com/anschnapp/pomodorofactory/pkg/view"
	"golang.org/x/term"
)

func main() {
	// Parse optional duration argument (in minutes, decimal allowed)
	duration := 25 * time.Minute
	if len(os.Args) > 1 {
		minutes, err := strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid duration: %s (expected minutes, e.g. 25 or 0.2)\n", os.Args[1])
			os.Exit(1)
		}
		duration = time.Duration(minutes * float64(time.Minute))
	}

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

	// Build components
	pomodoro := pomodorobuild.MakePomodoro()
	var motivationcloudComp render.Renderable = motivationcloud.MakeMotivationcloud()
	statusComp := status.MakeStatus()
	var commandinputComp render.Renderable = commandinput.MakeCommandinput()
	v := view.MakeView(pomodoro, motivationcloudComp, statusComp, commandinputComp)

	t := timer.NewTimer(duration)

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

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Event loop
	for {
		select {
		case b, ok := <-inputCh:
			if !ok {
				return
			}
			switch b {
			case 'q', 0x03: // 'q' or Ctrl+C
				return
			case 's':
				if !t.IsRunning() && !t.IsFinished() {
					t.Start()
				}
			}

		case <-ticker.C:
			// Only update on tick if timer is running
			if !t.IsRunning() {
				continue
			}
		}

		// Update components based on timer state
		if t.IsRunning() {
			pct := t.Percentage()
			pomodoro.SetPercentage(pct)
			remaining := t.Remaining()
			mins := int(remaining.Minutes())
			secs := int(remaining.Seconds()) % 60
			statusComp.SetText(
				fmt.Sprintf("Pomodoro running  %02d:%02d", mins, secs),
				"",
			)
		}
		if t.IsFinished() {
			pomodoro.SetPercentage(100)
			statusComp.SetText("Pomodoro complete!", "")
		}

		v.Render()
		v.Print()
	}
}
