package audio

import (
	"math"
	"math/rand"
)

// GenerateRisingTone creates a sine sweep from startHz to endHz.
func GenerateRisingTone(startHz, endHz, durationSec float64) []byte {
	numSamples := int(SampleRate * durationSec)
	buf := make([]byte, numSamples*bitDepth)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		freq := startHz + (endHz-startHz)*t
		phase := 2.0 * math.Pi * freq * float64(i) / float64(SampleRate)
		sample := math.Sin(phase) * Envelope(t) * 0.3
		WriteSample16LE(buf, i, sample)
	}
	return buf
}

// GeneratePop creates a short noise burst with exponential decay.
func GeneratePop(durationSec float64) []byte {
	numSamples := int(SampleRate * durationSec)
	buf := make([]byte, numSamples*bitDepth)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		noise := rand.Float64()*2.0 - 1.0
		decay := math.Exp(-8.0 * t)
		sample := noise * decay * 0.4
		WriteSample16LE(buf, i, sample)
	}
	return buf
}

// GenerateFanfare creates a major chord arpeggio (C5-E5-G5-C6) with square waves.
func GenerateFanfare() []byte {
	notes := []float64{523.25, 659.25, 783.99, 1046.50}
	noteDur := 0.15
	gap := 0.03

	var parts [][]byte
	for _, freq := range notes {
		numSamples := int(SampleRate * noteDur)
		note := make([]byte, numSamples*bitDepth)
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(numSamples)
			phase := 2.0 * math.Pi * freq * float64(i) / float64(SampleRate)
			sample := math.Copysign(0.25, math.Sin(phase)) * Envelope(t)
			WriteSample16LE(note, i, sample)
		}
		parts = append(parts, note)
		parts = append(parts, MakeSilence(gap))
	}
	return ConcatSamples(parts)
}

// GeneratePartySequence builds the full party sound: rising tones + pops + fanfare.
// Returns the PCM buffer and its duration in seconds.
func GeneratePartySequence() ([]byte, float64) {
	var parts [][]byte

	// Three quick rising tones
	for i := 0; i < 3; i++ {
		startHz := 300.0 + float64(i)*200.0
		endHz := startHz + 400.0
		parts = append(parts, GenerateRisingTone(startHz, endHz, 0.2))
		parts = append(parts, MakeSilence(0.05))
	}

	// Two pops
	for i := 0; i < 2; i++ {
		parts = append(parts, GeneratePop(0.1))
		parts = append(parts, MakeSilence(0.05))
	}

	// Fanfare
	parts = append(parts, GenerateFanfare())

	combined := ConcatSamples(parts)
	return combined, DurationSec(combined)
}
