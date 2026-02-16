package celebration

import (
	"time"

	"github.com/anschnapp/pomodorofactory/pkg/audio"
)

// Phase represents the current stage of the celebration sequence.
type Phase int

const (
	PhaseNone   Phase = iota
	PhaseParty        // visual sparks + party sounds
	PhaseSpeech       // gibberish voice reads congratulatory message
	PhaseDone         // celebration finished
)

// Celebration coordinates the two-phase celebration after a pomodoro completes.
type Celebration struct {
	phase     Phase
	startTime time.Time
	engine    *audio.Engine // nil if audio unavailable
	message   string        // the congratulatory message for this run

	// Party phase
	partyDuration time.Duration
	partyTick     int

	// Speech phase
	charTimings []audio.CharTiming
	speechStart time.Time
	currentChar int
	speechDone  <-chan struct{}
}

// New creates a celebration coordinator. engine may be nil for visual-only mode.
func New(engine *audio.Engine) *Celebration {
	return &Celebration{
		phase:  PhaseNone,
		engine: engine,
	}
}

// Message returns the congratulatory message for the current celebration.
func (c *Celebration) Message() string { return c.message }

// Start kicks off the party phase. Call once when the timer finishes.
// message is the text that will be spoken in the speech phase.
func (c *Celebration) Start(message string) {
	c.message = message
	c.phase = PhaseParty
	c.startTime = time.Now()
	c.partyTick = 0

	if c.engine != nil {
		samples, dur := audio.GeneratePartySequence()
		c.partyDuration = time.Duration(dur * float64(time.Second))
		c.engine.Play(samples)
	} else {
		c.partyDuration = 3 * time.Second
	}
}

// Tick advances the celebration state. Call every 50ms from the event loop.
func (c *Celebration) Tick() Phase {
	switch c.phase {
	case PhaseParty:
		c.partyTick++
		if time.Since(c.startTime) >= c.partyDuration {
			c.startSpeechPhase()
		}
	case PhaseSpeech:
		elapsed := time.Since(c.speechStart)
		// Advance currentChar based on elapsed time vs charTimings
		for c.currentChar < len(c.charTimings)-1 {
			nextOffset := c.charTimings[c.currentChar+1].SampleOffset
			nextTime := time.Duration(float64(nextOffset) / float64(audio.SampleRate) * float64(time.Second))
			if elapsed >= nextTime {
				c.currentChar++
			} else {
				break
			}
		}
		// Check if speech audio is done
		select {
		case <-c.speechDone:
			c.phase = PhaseDone
		default:
		}
	}
	return c.phase
}

func (c *Celebration) startSpeechPhase() {
	c.phase = PhaseSpeech
	c.speechStart = time.Now()
	c.currentChar = 0

	msg := c.message
	if c.engine != nil {
		samples, timings := audio.GenerateAnimalese(msg)
		c.charTimings = timings
		c.speechDone = c.engine.Play(samples)
	} else {
		// Visual-only: fake timings at ~80ms per char
		c.charTimings = make([]audio.CharTiming, len([]rune(msg)))
		for i := range []rune(msg) {
			c.charTimings[i] = audio.CharTiming{
				SampleOffset: i * int(0.08*float64(audio.SampleRate)),
				CharIndex:    i,
				DurationMs:   80,
			}
		}
		totalDur := time.Duration(len([]rune(msg))) * 80 * time.Millisecond
		done := make(chan struct{})
		go func() {
			time.Sleep(totalDur)
			close(done)
		}()
		c.speechDone = done
	}
}

// IsActive returns true if the celebration is in progress.
func (c *Celebration) IsActive() bool {
	return c.phase == PhaseParty || c.phase == PhaseSpeech
}

// CurrentPhase returns the current phase.
func (c *Celebration) CurrentPhase() Phase { return c.phase }

// PartyTick returns the tick counter for animation randomization.
func (c *Celebration) PartyTick() int { return c.partyTick }

// CurrentCharIndex returns the index of the character being spoken.
func (c *Celebration) CurrentCharIndex() int { return c.currentChar }
