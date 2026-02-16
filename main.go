package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/anschnapp/pomodorofactory/pkg/audio"
	"github.com/anschnapp/pomodorofactory/pkg/celebration"
	"github.com/anschnapp/pomodorofactory/pkg/commandinput"
	"github.com/anschnapp/pomodorofactory/pkg/factoryscene"
	"github.com/anschnapp/pomodorofactory/pkg/motivationcloud"
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

	// Initialize audio (optional — celebration works visually without it)
	audioEngine, _ := audio.NewEngine()

	// Build components
	factory := factoryscene.MakeFactoryScene()
	motivationcloudComp := motivationcloud.MakeMotivationcloud()
	statusComp := status.MakeStatus()
	var commandinputComp render.Renderable = commandinput.MakeCommandinput()
	v := view.MakeView(factory, motivationcloudComp, statusComp, commandinputComp)

	t := timer.NewTimer(duration)
	celeb := celebration.New(audioEngine)
	celebrationStarted := false
	lastShuffle := time.Now()

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

	ticker := time.NewTicker(50 * time.Millisecond)
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
			// Skip tick if nothing is animating
			if !t.IsRunning() && !celeb.IsActive() {
				continue
			}
		}

		// Update components based on timer state
		if t.IsRunning() {
			factory.SetProgress(t.Progress())
			remaining := t.Remaining()
			mins := int(remaining.Minutes())
			secs := int(remaining.Seconds()) % 60
			statusComp.SetText(
				fmt.Sprintf("Pomodoro running  %02d:%02d", mins, secs),
				"",
			)
		}
		if t.IsFinished() {
			if !celebrationStarted {
				celebrationStarted = true
				celeb.Start()
			}

			if celeb.IsActive() {
				phase := celeb.Tick()
				switch phase {
				case celebration.PhaseParty:
					factory.SetCelebrating(celeb.PartyTick())
					statusComp.SetCelebrationText("POMODORO COMPLETE!", celeb.PartyTick())
				case celebration.PhaseSpeech:
					factory.SetProgress(1.0)
					statusComp.SetSpeechText(celebration.Message, celeb.CurrentCharIndex())
				}
			} else {
				factory.SetProgress(1.0)
				statusComp.SetText("Pomodoro complete!", "")
			}
		}

		// Refresh motivation cloud every 5 minutes
		if time.Since(lastShuffle) >= 5*time.Minute {
			motivationcloudComp.Shuffle()
			lastShuffle = time.Now()
		}

		v.Render()
		v.Print()
	}
}
