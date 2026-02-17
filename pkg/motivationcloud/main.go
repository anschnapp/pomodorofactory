package motivationcloud

import (
	"math/rand"

	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/fatih/color"
)

const (
	cloudWidth  = 30
	cloudHeight = 10
	maxIndent   = 5
	phraseCount = 5
)

var phraseColors = []color.Attribute{
	color.FgHiCyan,
	color.FgWhite,
	color.FgHiMagenta,
}

type placedPhrase struct {
	row         int
	indent      int
	text        string
	color       []color.Attribute
	revealChars int // how many chars are visible (for reveal animation)
	fadingOut   bool
}

type Motivationcloud struct {
	width      int
	height     int
	phrases    []placedPhrase
	pendingNew *placedPhrase // new phrase waiting for the fading-out slot to clear
	fadeOutIdx int           // index of the phrase currently fading out (-1 = none)
}

func MakeMotivationcloud() *Motivationcloud {
	m := &Motivationcloud{
		width:      cloudWidth,
		height:     cloudHeight,
		fadeOutIdx: -1,
	}
	m.Shuffle()
	return m
}

// Shuffle picks new random phrases and scatters them across the available rows.
func (m *Motivationcloud) Shuffle() {
	// Pick phraseCount random phrases that fit within width
	picked := pickPhrases(phraseCount, cloudWidth-maxIndent)

	// Distribute across rows with spacing
	m.phrases = distribute(picked, cloudHeight, cloudWidth)

	// All phrases start fully revealed
	for i := range m.phrases {
		m.phrases[i].revealChars = len([]rune(m.phrases[i].text))
	}
	m.fadeOutIdx = -1
	m.pendingNew = nil
}

// ReplaceOne starts fading out one random phrase, which will be replaced by a
// new phrase once the fade-out completes.
func (m *Motivationcloud) ReplaceOne() {
	if len(m.phrases) == 0 || m.fadeOutIdx >= 0 {
		return // already animating
	}

	// Pick a new phrase that isn't currently displayed
	current := make(map[string]bool, len(m.phrases))
	for _, p := range m.phrases {
		current[p.text] = true
	}
	candidates := make([]string, 0, len(phrases))
	maxLen := cloudWidth - maxIndent
	for _, p := range phrases {
		if len(p) <= maxLen && !current[p] {
			candidates = append(candidates, p)
		}
	}
	if len(candidates) == 0 {
		return
	}
	newText := candidates[rand.Intn(len(candidates))]

	// Pick a random slot to replace
	idx := rand.Intn(len(m.phrases))
	old := &m.phrases[idx]

	// Prepare the new phrase on the same row
	maxInd := cloudWidth - len([]rune(newText))
	if maxInd > maxIndent {
		maxInd = maxIndent
	}
	if maxInd < 0 {
		maxInd = 0
	}
	col := phraseColors[rand.Intn(len(phraseColors))]
	pending := &placedPhrase{
		row:         old.row,
		indent:      rand.Intn(maxInd + 1),
		text:        newText,
		color:       []color.Attribute{col},
		revealChars: 0,
	}
	m.pendingNew = pending
	m.fadeOutIdx = idx
	old.fadingOut = true
}

// Tick advances the reveal/fade animation by one character. Call this on every
// render tick (e.g. 50ms). Returns true if any visual change happened.
func (m *Motivationcloud) Tick() bool {
	changed := false

	// Advance fade-out
	if m.fadeOutIdx >= 0 && m.fadeOutIdx < len(m.phrases) {
		p := &m.phrases[m.fadeOutIdx]
		if p.revealChars > 0 {
			p.revealChars--
			changed = true
		} else {
			// Fade-out done — swap in the new phrase
			if m.pendingNew != nil {
				m.phrases[m.fadeOutIdx] = *m.pendingNew
				m.pendingNew = nil
			}
			m.fadeOutIdx = -1
			changed = true
		}
	}

	// Advance reveal on any phrase that isn't fully visible yet
	for i := range m.phrases {
		p := &m.phrases[i]
		if !p.fadingOut && p.revealChars < len([]rune(p.text)) {
			p.revealChars++
			changed = true
		}
	}

	return changed
}

// pickPhrases selects n unique random phrases that fit within maxLen.
func pickPhrases(n int, maxLen int) []string {
	// Build candidate list of phrases that fit
	candidates := make([]string, 0, len(phrases))
	for _, p := range phrases {
		if len(p) <= maxLen {
			candidates = append(candidates, p)
		}
	}

	// Shuffle and pick first n
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if n > len(candidates) {
		n = len(candidates)
	}
	return candidates[:n]
}

// distribute places phrases on rows with at least 1 empty row between them,
// with random indentation and color.
func distribute(picked []string, height int, width int) []placedPhrase {
	n := len(picked)
	if n == 0 {
		return nil
	}

	// Calculate available rows: we need n rows for phrases + (n-1) gaps of at least 1
	// Spread them evenly across the height
	spacing := height / n
	if spacing < 2 {
		spacing = 2
	}

	result := make([]placedPhrase, n)
	for i, text := range picked {
		row := i * spacing
		if row >= height {
			row = height - 1
		}

		// Random indent, but ensure text fits
		maxInd := width - len(text)
		if maxInd > maxIndent {
			maxInd = maxIndent
		}
		if maxInd < 0 {
			maxInd = 0
		}
		indent := rand.Intn(maxInd + 1)

		col := phraseColors[rand.Intn(len(phraseColors))]

		result[i] = placedPhrase{
			row:    row,
			indent: indent,
			text:   text,
			color:  []color.Attribute{col},
		}
	}
	return result
}

func (m *Motivationcloud) Width() int {
	return m.width
}

func (m *Motivationcloud) Height() int {
	return m.height
}

func (m *Motivationcloud) Render(subview [][]runecolor.ColoredRune) {
	// Clear the view region
	for i := range subview {
		for j := range subview[i] {
			subview[i][j] = runecolor.ColoredRune{Symbol: ' '}
		}
	}

	// Draw each phrase (only up to revealChars)
	for _, p := range m.phrases {
		if p.row >= len(subview) {
			continue
		}
		row := subview[p.row]
		runes := []rune(p.text)
		visible := p.revealChars
		if visible > len(runes) {
			visible = len(runes)
		}
		for ci := 0; ci < visible; ci++ {
			col := p.indent + ci
			if col >= len(row) {
				break
			}
			attrs := p.color
			// Dim the leading character for a subtle reveal/fade effect
			if ci == visible-1 && visible < len(runes) {
				attrs = []color.Attribute{color.Faint}
			}
			row[col] = runecolor.ColoredRune{
				Symbol:          runes[ci],
				ColorAttributes: attrs,
			}
		}
	}
}
