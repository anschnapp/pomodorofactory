package audio

import (
	"math"
	"math/rand"
	"strings"
)

var vowels = "aeiouAEIOU"

// CharTiming maps a character index to its position in the PCM buffer.
type CharTiming struct {
	SampleOffset int
	CharIndex    int
	DurationMs   int
}

// charBlipParams returns the base frequency and duration for a character.
func charBlipParams(ch rune) (freq float64, dur float64) {
	if ch == ' ' {
		return 0, 0.06 // silence
	}
	if strings.ContainsRune(vowels, ch) {
		base := 200.0 + float64(ch%10)*20.0
		return base, 0.08
	}
	// Consonants
	base := 400.0 + float64(ch%15)*25.0
	return base, 0.06
}

// GenerateAnimalese generates gibberish speech for a message.
// Returns the PCM buffer and timing data for syncing a visual highlight.
func GenerateAnimalese(message string) ([]byte, []CharTiming) {
	var parts [][]byte
	var timings []CharTiming
	sampleOffset := 0
	gapSec := 0.02

	for i, ch := range message {
		freq, dur := charBlipParams(ch)
		// Random pitch variation +/- 15%
		variation := 1.0 + (rand.Float64()*0.3 - 0.15)
		freq *= variation

		numSamples := int(SampleRate * dur)
		blip := make([]byte, numSamples*bitDepth)

		if freq > 0 {
			for s := 0; s < numSamples; s++ {
				t := float64(s) / float64(numSamples)
				phase := 2.0 * math.Pi * freq * float64(s) / float64(SampleRate)
				// Sawtooth + sine mix for retro character
				saw := 2.0*(math.Mod(float64(s)*freq/float64(SampleRate), 1.0)) - 1.0
				sample := (0.7*saw + 0.3*math.Sin(phase)) * Envelope(t) * 0.3
				WriteSample16LE(blip, s, sample)
			}
		}

		timings = append(timings, CharTiming{
			SampleOffset: sampleOffset,
			CharIndex:    i,
			DurationMs:   int(dur * 1000),
		})
		parts = append(parts, blip)
		sampleOffset += numSamples

		gap := MakeSilence(gapSec)
		parts = append(parts, gap)
		sampleOffset += len(gap) / bitDepth
	}

	return ConcatSamples(parts), timings
}
