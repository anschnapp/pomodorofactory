package timer

import "time"

type Timer struct {
	duration  time.Duration
	startTime time.Time
	running   bool
	finished  bool
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		duration: duration,
	}
}

func (t *Timer) Start() {
	t.startTime = time.Now()
	t.running = true
	t.finished = false
}

// Reset prepares the timer for a new countdown with the given duration.
func (t *Timer) Reset(duration time.Duration) {
	t.duration = duration
	t.running = false
	t.finished = false
}

func (t *Timer) IsRunning() bool {
	return t.running
}

func (t *Timer) IsFinished() bool {
	return t.finished
}

func (t *Timer) Elapsed() time.Duration {
	if !t.running {
		return 0
	}
	elapsed := time.Since(t.startTime)
	if elapsed > t.duration {
		return t.duration
	}
	return elapsed
}

func (t *Timer) Remaining() time.Duration {
	return t.duration - t.Elapsed()
}

// Progress returns 0.0–1.0 how far the timer has progressed.
func (t *Timer) Progress() float64 {
	if !t.running {
		return 0
	}
	elapsed := time.Since(t.startTime)
	if elapsed >= t.duration {
		t.running = false
		t.finished = true
		return 1.0
	}
	return float64(elapsed) / float64(t.duration)
}

// Percentage returns 0–100 how far the timer has progressed.
func (t *Timer) Percentage() int {
	if !t.running {
		return 0
	}
	elapsed := time.Since(t.startTime)
	if elapsed >= t.duration {
		t.running = false
		t.finished = true
		return 100
	}
	pct := int(elapsed * 100 / t.duration)
	if pct > 100 {
		pct = 100
	}
	return pct
}
