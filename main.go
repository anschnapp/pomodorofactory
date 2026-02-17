package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/anschnapp/pomodorofactory/pkg/audio"
	"github.com/anschnapp/pomodorofactory/pkg/celebration"
	"github.com/anschnapp/pomodorofactory/pkg/commandinput"
	"github.com/anschnapp/pomodorofactory/pkg/factoryscene"
	"github.com/anschnapp/pomodorofactory/pkg/motivationcloud"
	"github.com/anschnapp/pomodorofactory/pkg/status"
	"github.com/anschnapp/pomodorofactory/pkg/timer"
	"github.com/anschnapp/pomodorofactory/pkg/view"
	"golang.org/x/term"
)

var (
	congratsWords = []string{
		"Congratulations", "Well done", "Fantastic", "Bravo",
		"Impressive", "Outstanding", "Remarkable", "Excellent",
		"Stupendous", "Wonderful", "Sensational", "Phenomenal",
		"Incredible", "Marvelous", "Brilliant", "Spectacular",
		"Superb", "Terrific", "Magnificent", "Splendid",
	}
	adverbWords = []string{
		"successfully", "masterfully", "skillfully", "proudly",
		"brilliantly", "flawlessly", "expertly", "elegantly",
		"perfectly", "superbly", "gracefully", "precisely",
		"diligently", "gloriously", "beautifully", "boldly",
		"heroically", "effortlessly", "passionately", "triumphantly",
	}
	verbWords = []string{
		"built", "assembled", "crafted", "manufactured",
		"constructed", "forged", "produced", "engineered",
		"fabricated", "created", "welded", "shaped",
		"molded", "formed", "designed", "composed",
		"erected", "fashioned", "devised", "completed",
	}
	adjectiveWords = []string{
		"beautiful", "magnificent", "splendid", "glorious",
		"stunning", "exquisite", "pristine", "majestic",
		"radiant", "dazzling", "fabulous", "grand",
		"supreme", "flawless", "legendary", "epic",
		"divine", "stellar", "remarkable", "perfect",
	}
)

func randomCongrats() string {
	return fmt.Sprintf("%s we %s %s a %s pomodoro",
		congratsWords[rand.Intn(len(congratsWords))],
		adverbWords[rand.Intn(len(adverbWords))],
		verbWords[rand.Intn(len(verbWords))],
		adjectiveWords[rand.Intn(len(adjectiveWords))],
	)
}

type appState int

const (
	stateIdle        appState = iota // waiting for 's' to start a pomodoro
	stateWorking                     // pomodoro timer running
	stateCelebrating                 // celebration animation playing
	stateOnBreak                     // break timer running (auto-started)
)

const (
	shortBreak      = 5 * time.Minute
	longBreak       = 15 * time.Minute
	pomodorosPerSet = 4
)

func main() {
	// Parse optional duration argument (in minutes, decimal allowed)
	workDuration := 25 * time.Minute
	if len(os.Args) > 1 {
		minutes, err := strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid duration: %s (expected minutes, e.g. 25 or 0.2)\n", os.Args[1])
			os.Exit(1)
		}
		workDuration = time.Duration(minutes * float64(time.Minute))
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
	cmdInput := commandinput.MakeCommandinput()
	v := view.MakeView(factory, motivationcloudComp, statusComp, cmdInput)

	t := timer.NewTimer(workDuration)
	celeb := celebration.New(audioEngine)
	lastShuffle := time.Now()

	state := stateIdle
	completedPomodoros := 0
	congratsMsg := ""

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
				if state == stateIdle {
					state = stateWorking
					t.Reset(workDuration)
					t.Start()
					factory.Reset()
					cmdInput.SetText("[q]uit")
				}
			}

		case <-ticker.C:
			// tick proceeds — motivation cloud always animates
		}

		switch state {
		case stateWorking:
			factory.SetProgress(t.Progress())
			remaining := t.Remaining()
			mins := int(remaining.Minutes())
			secs := int(remaining.Seconds()) % 60
			statusComp.SetTextWithTomatoes(
				fmt.Sprintf("Factory running  %02d:%02d", mins, secs),
				completedPomodoros,
			)

			if t.IsFinished() {
				state = stateCelebrating
				congratsMsg = randomCongrats()
				celeb.Start(congratsMsg)
			}

		case stateCelebrating:
			if celeb.IsActive() {
				phase := celeb.Tick()
				switch phase {
				case celebration.PhaseParty:
					factory.SetCelebrating(celeb.PartyTick())
					statusComp.SetCelebrationText("POMODORO COMPLETE!", celeb.PartyTick())
				case celebration.PhaseSpeech:
					factory.SetProgress(1.0)
					statusComp.SetSpeechText(congratsMsg, celeb.CurrentCharIndex())
				}
			} else {
				// Celebration finished — count pomodoro and auto-start break
				completedPomodoros++
				breakDuration := shortBreak
				if completedPomodoros%pomodorosPerSet == 0 {
					breakDuration = longBreak
				}
				state = stateOnBreak
				t.Reset(breakDuration)
				t.Start()
				factory.SetProgress(1.0)
				cmdInput.SetText("[q]uit")
			}

		case stateOnBreak:
			t.Progress() // drive the finished flag
			remaining := t.Remaining()
			mins := int(remaining.Minutes())
			secs := int(remaining.Seconds()) % 60
			label := "Factory needs a short cooldown"
			if completedPomodoros%pomodorosPerSet == 0 {
				label = "Factory needs a longer cooldown"
			}
			statusComp.SetTextWithTomatoes(
				fmt.Sprintf("%s  %02d:%02d", label, mins, secs),
				completedPomodoros,
			)

			if t.IsFinished() {
				state = stateIdle
				factory.Reset()
				statusComp.SetTextWithTomatoes("Factory ready  press [s] to start", completedPomodoros)
				cmdInput.SetText("[s]tart | [q]uit")
			}
		}

		// Replace one phrase every 30 seconds (with animated transition)
		if time.Since(lastShuffle) >= 15*time.Second {
			motivationcloudComp.ReplaceOne()
			lastShuffle = time.Now()
		}
		motivationcloudComp.Tick()

		v.Render()
		v.Print()
	}
}
