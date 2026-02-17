package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anschnapp/pomodorofactory/pkg/audio"
	"github.com/anschnapp/pomodorofactory/pkg/celebration"
	"github.com/anschnapp/pomodorofactory/pkg/commandinput"
	"github.com/anschnapp/pomodorofactory/pkg/factoryscene"
	"github.com/anschnapp/pomodorofactory/pkg/motivationcloud"
	"github.com/anschnapp/pomodorofactory/pkg/product"
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

func randomCongrats(productName string) string {
	return fmt.Sprintf("%s we %s %s a %s %s",
		congratsWords[rand.Intn(len(congratsWords))],
		adverbWords[rand.Intn(len(adverbWords))],
		verbWords[rand.Intn(len(verbWords))],
		adjectiveWords[rand.Intn(len(adjectiveWords))],
		strings.ToLower(productName),
	)
}

func selectorLine(products []*product.Product, idx int) string {
	return fmt.Sprintf("build next:  \u2190 [%s] \u2192", products[idx].Name)
}

type appState int

const (
	stateIdle                   appState = iota // waiting for 's' to start a pomodoro
	stateWorking                                // pomodoro timer running
	stateWaitingForCelebration                  // timer done, waiting for user to press 'c'
	stateCelebrating                            // celebration animation playing
	stateOnBreak                                // break timer running (auto-started)
)

const (
	shortBreak      = 5 * time.Minute
	longBreak       = 15 * time.Minute
	pomodorosPerSet = 4
)

// Sentinel rune values for arrow keys (not valid Unicode)
const (
	keyLeft  = rune(-1)
	keyRight = rune(-2)
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

	// Product selection state
	products := product.All
	selectedProductIdx := 0
	achievedEmojis := []string{}

	// Build components
	factory := factoryscene.MakeFactoryScene(products)
	motivationcloudComp := motivationcloud.MakeMotivationcloud()
	statusComp := status.MakeStatus()
	cmdInput := commandinput.MakeCommandinput()
	cmdInput.SetTexts("[s]tart | [q]uit", selectorLine(products, selectedProductIdx))
	v := view.MakeView(factory, motivationcloudComp, statusComp, cmdInput)

	t := timer.NewTimer(workDuration)
	celeb := celebration.New(audioEngine)
	lastShuffle := time.Now()

	state := stateIdle
	congratsMsg := ""

	// Read input in a goroutine; arrow keys are decoded as sentinel rune values
	inputCh := make(chan rune)
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				b := buf[0]
				if b == 0x1b {
					// Possible escape sequence — read up to two more bytes
					seq := make([]byte, 2)
					n2, _ := os.Stdin.Read(seq[:1])
					if n2 > 0 && seq[0] == '[' {
						n3, _ := os.Stdin.Read(seq[1:2])
						if n3 > 0 {
							switch seq[1] {
							case 'C':
								inputCh <- keyRight
							case 'D':
								inputCh <- keyLeft
							}
						}
					}
					// lone ESC or unrecognised sequence: ignored
				} else {
					inputCh <- rune(b)
				}
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
			case 'h', keyLeft:
				if state == stateIdle {
					selectedProductIdx = (selectedProductIdx - 1 + len(products)) % len(products)
					cmdInput.SetTexts("[s]tart | [q]uit", selectorLine(products, selectedProductIdx))
				}
			case 'l', keyRight:
				if state == stateIdle {
					selectedProductIdx = (selectedProductIdx + 1) % len(products)
					cmdInput.SetTexts("[s]tart | [q]uit", selectorLine(products, selectedProductIdx))
				}
			case 's':
				if state == stateIdle {
					state = stateWorking
					factory.LoadArt(products[selectedProductIdx].Art)
					t.Reset(workDuration)
					t.Start()
					factory.Reset()
					cmdInput.SetTexts("[q]uit", "")
				}
			case 'c':
				if state == stateWaitingForCelebration {
					state = stateCelebrating
					congratsMsg = randomCongrats(products[selectedProductIdx].Name)
					celeb.Start(congratsMsg)
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
			statusComp.SetAchievements(
				fmt.Sprintf("Factory running  %02d:%02d", mins, secs),
				achievedEmojis,
			)

			if t.IsFinished() {
				state = stateWaitingForCelebration
				factory.SetProgress(1.0)
				statusComp.SetAchievements("Pomodoro done!  Press [c] to celebrate", achievedEmojis)
				cmdInput.SetTexts("[c]elebrate", "")
				if audioEngine != nil {
					audioEngine.Play(audio.MakeNotificationSound())
				}
			}

		case stateWaitingForCelebration:
			// Waiting for user to press 'c' — nothing to update each tick

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
				// Celebration finished — record achievement and auto-start break
				achievedEmojis = append(achievedEmojis, products[selectedProductIdx].Emoji)
				breakDuration := shortBreak
				if len(achievedEmojis)%pomodorosPerSet == 0 {
					breakDuration = longBreak
				}
				state = stateOnBreak
				t.Reset(breakDuration)
				t.Start()
				factory.SetProgress(1.0)
				cmdInput.SetTexts("[q]uit", "")
			}

		case stateOnBreak:
			t.Progress() // drive the finished flag
			remaining := t.Remaining()
			mins := int(remaining.Minutes())
			secs := int(remaining.Seconds()) % 60
			label := "Factory needs a short cooldown"
			if len(achievedEmojis)%pomodorosPerSet == 0 {
				label = "Factory needs a longer cooldown"
			}
			statusComp.SetAchievements(
				fmt.Sprintf("%s  %02d:%02d", label, mins, secs),
				achievedEmojis,
			)

			if t.IsFinished() {
				state = stateIdle
				factory.Reset()
				statusComp.SetAchievements("Factory ready  press [s] to start", achievedEmojis)
				cmdInput.SetTexts("[s]tart | [q]uit", selectorLine(products, selectedProductIdx))
			}
		}

		// Replace one phrase every 15 seconds (with animated transition)
		if time.Since(lastShuffle) >= 15*time.Second {
			motivationcloudComp.ReplaceOne()
			lastShuffle = time.Now()
		}
		motivationcloudComp.Tick()

		v.Render()
		v.Print()
	}
}
