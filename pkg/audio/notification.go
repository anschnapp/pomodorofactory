package audio

import "math"

// MakeNotificationSound generates a mechanical alarm ring reminiscent of an analog
// kitchen timer — a bell tone rapidly modulated by a mechanical clapper (brrrrr).
func MakeNotificationSound() []byte {
	dur := 1.2        // total ring duration in seconds
	bellFreq := 1100.0 // Hz — the resonant bell pitch
	clapRate := 24.0   // Hz — clapper strikes per second

	n := int(SampleRate * dur)
	buf := make([]byte, n*bitDepth)
	for i := 0; i < n; i++ {
		t := float64(i) / SampleRate

		// Fade in (20ms) and fade out (200ms) to avoid clicks and soften the end
		env := 1.0
		if t < 0.02 {
			env = t / 0.02
		} else if t > dur-0.2 {
			env = (dur - t) / 0.2
		}

		// Half-wave rectified sine = mechanical clapper hitting (not a smooth vibrato)
		clapper := math.Max(0, math.Sin(2*math.Pi*clapRate*t))

		// Bell tone: fundamental + inharmonic overtone for metallic texture
		bell := math.Sin(2*math.Pi*bellFreq*t)*0.7 +
			math.Sin(2*math.Pi*bellFreq*2.42*t)*0.3

		sample := bell * clapper * env * 0.5
		WriteSample16LE(buf, i, sample)
	}
	return buf
}
