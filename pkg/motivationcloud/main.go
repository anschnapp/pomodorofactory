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
	row    int
	indent int
	text   string
	color  []color.Attribute
}

type Motivationcloud struct {
	width   int
	height  int
	phrases []placedPhrase
}

func MakeMotivationcloud() *Motivationcloud {
	m := &Motivationcloud{
		width:  cloudWidth,
		height: cloudHeight,
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

	// Draw each phrase
	for _, p := range m.phrases {
		if p.row >= len(subview) {
			continue
		}
		row := subview[p.row]
		for ci, ch := range p.text {
			col := p.indent + ci
			if col < len(row) {
				row[col] = runecolor.ColoredRune{
					Symbol:          ch,
					ColorAttributes: p.color,
				}
			}
		}
	}
}
