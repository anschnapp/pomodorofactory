package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const (
	SampleRate  = 44100
	numChannels = 1
	bitDepth    = 2 // 16-bit = 2 bytes per sample
)

// Engine plays raw PCM audio via platform-native tools.
// Linux: pipes raw PCM to aplay. macOS: writes a temp WAV and plays via afplay.
type Engine struct {
	platform string // "linux" or "darwin"
}

// NewEngine creates an audio engine. Returns nil, err if no playback tool is available.
func NewEngine() (*Engine, error) {
	switch runtime.GOOS {
	case "linux":
		if _, err := exec.LookPath("aplay"); err != nil {
			return nil, fmt.Errorf("aplay not found: %w", err)
		}
		return &Engine{platform: "linux"}, nil
	case "darwin":
		if _, err := exec.LookPath("afplay"); err != nil {
			return nil, fmt.Errorf("afplay not found: %w", err)
		}
		return &Engine{platform: "darwin"}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// Play sends PCM samples (signed 16-bit LE, mono, 44100Hz) to the audio output.
// Returns a channel that closes when playback finishes.
func (e *Engine) Play(samples []byte) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		switch e.platform {
		case "linux":
			e.playAplay(samples)
		case "darwin":
			e.playAfplay(samples)
		}
	}()
	return done
}

func (e *Engine) playAplay(samples []byte) {
	cmd := exec.Command("aplay",
		"-f", "S16_LE",
		"-r", "44100",
		"-c", "1",
		"--quiet",
	)
	cmd.Stdin = bytes.NewReader(samples)
	_ = cmd.Run()
}

func (e *Engine) playAfplay(samples []byte) {
	wav := pcmToWAV(samples)
	f, err := os.CreateTemp("", "pomodorofactory-*.wav")
	if err != nil {
		return
	}
	tmpPath := f.Name()
	defer os.Remove(tmpPath)

	if _, err := f.Write(wav); err != nil {
		f.Close()
		return
	}
	f.Close()

	_ = exec.Command("afplay", tmpPath).Run()
}

// pcmToWAV wraps raw PCM data (signed 16-bit LE, mono, 44100Hz) in a WAV header.
func pcmToWAV(pcm []byte) []byte {
	dataSize := uint32(len(pcm))
	fileSize := 36 + dataSize // 44-byte header minus 8 for RIFF chunk header

	buf := new(bytes.Buffer)
	buf.Grow(44 + len(pcm))

	// RIFF header
	buf.WriteString("RIFF")
	binary.Write(buf, binary.LittleEndian, fileSize)
	buf.WriteString("WAVE")

	// fmt sub-chunk
	buf.WriteString("fmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16))    // sub-chunk size
	binary.Write(buf, binary.LittleEndian, uint16(1))     // PCM format
	binary.Write(buf, binary.LittleEndian, uint16(numChannels))
	binary.Write(buf, binary.LittleEndian, uint32(SampleRate))
	binary.Write(buf, binary.LittleEndian, uint32(SampleRate*numChannels*bitDepth)) // byte rate
	binary.Write(buf, binary.LittleEndian, uint16(numChannels*bitDepth))            // block align
	binary.Write(buf, binary.LittleEndian, uint16(bitDepth*8))                      // bits per sample

	// data sub-chunk
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, dataSize)
	buf.Write(pcm)

	return buf.Bytes()
}

// WriteSample16LE writes a float64 sample [-1.0, 1.0] as signed 16-bit LE at the given index.
func WriteSample16LE(buf []byte, idx int, sample float64) {
	if sample > 1.0 {
		sample = 1.0
	}
	if sample < -1.0 {
		sample = -1.0
	}
	val := int16(sample * 32767.0)
	off := idx * bitDepth
	if off+1 < len(buf) {
		buf[off] = byte(val)
		buf[off+1] = byte(val >> 8)
	}
}

// MakeSilence generates a silent PCM buffer of the given duration.
func MakeSilence(durationSec float64) []byte {
	return make([]byte, int(SampleRate*durationSec)*bitDepth)
}

// ConcatSamples concatenates multiple PCM buffers.
func ConcatSamples(parts [][]byte) []byte {
	total := 0
	for _, p := range parts {
		total += len(p)
	}
	result := make([]byte, 0, total)
	for _, p := range parts {
		result = append(result, p...)
	}
	return result
}

// DurationSec returns the playback duration of a PCM buffer in seconds.
func DurationSec(samples []byte) float64 {
	return float64(len(samples)) / float64(SampleRate*bitDepth)
}

// Envelope applies a quick attack/release shape to a normalized time position t in [0,1].
func Envelope(t float64) float64 {
	if t < 0.05 {
		return t / 0.05
	}
	if t > 0.9 {
		return (1.0 - t) / 0.1
	}
	return 1.0
}
