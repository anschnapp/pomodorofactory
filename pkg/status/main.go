package status

import (
	"strings"

	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/fatih/color"
)

type status struct {
	width              int
	height             int
	asciRepresentation [][]runecolor.ColoredRune
}

// Fixed width so the view region is large enough for any status text
const statusWidth = 50

func MakeStatus() *status {
	s := &status{
		height: 2,
		width:  statusWidth,
	}
	s.SetText("Factory ready  press [s] to start", "")
	return s
}

func (s *status) SetText(line1, line2 string) {
	asci := make([][]runecolor.ColoredRune, 2)
	asci[0] = runecolor.ConvertSimpleRunes([]rune(line1))
	asci[1] = runecolor.ConvertSimpleRunes([]rune(line2))
	s.asciRepresentation = asci
}

// SetAchievements sets line 1 text and shows achievement emojis on line 2.
// Each emoji is typically double-width in terminal; a trailing space is added for alignment.
func (s *status) SetAchievements(line1 string, emojis []string) {
	asci := make([][]runecolor.ColoredRune, 2)
	asci[0] = runecolor.ConvertSimpleRunes([]rune(line1))
	if len(emojis) > 0 {
		asci[1] = runecolor.ConvertSimpleRunes([]rune(strings.Join(emojis, " ") + " "))
	} else {
		asci[1] = nil
	}
	s.asciRepresentation = asci
}

var celebColors = []color.Attribute{
	color.FgHiYellow, color.FgHiGreen, color.FgHiMagenta,
	color.FgHiCyan, color.FgHiRed,
}

// SetCelebrationText sets status text with a color that cycles each tick.
func (s *status) SetCelebrationText(text string, tick int) {
	clr := celebColors[tick%len(celebColors)]
	runes := []rune(text)
	colored := make([]runecolor.ColoredRune, len(runes))
	for i, r := range runes {
		colored[i] = runecolor.ColoredRune{
			Symbol:          r,
			ColorAttributes: []color.Attribute{clr},
		}
	}
	s.asciRepresentation = [][]runecolor.ColoredRune{colored, {}}
}

// SetSpeechText shows a message with the current character highlighted.
// Already-spoken characters are white, current is bold yellow, upcoming are dim.
func (s *status) SetSpeechText(message string, highlightIdx int) {
	runes := []rune(message)
	colored := make([]runecolor.ColoredRune, len(runes))
	for i, r := range runes {
		switch {
		case i == highlightIdx:
			colored[i] = runecolor.ColoredRune{
				Symbol:          r,
				ColorAttributes: []color.Attribute{color.FgHiYellow, color.Bold},
			}
		case i < highlightIdx:
			colored[i] = runecolor.ColoredRune{
				Symbol:          r,
				ColorAttributes: []color.Attribute{color.FgWhite},
			}
		default:
			colored[i] = runecolor.ColoredRune{
				Symbol:          r,
				ColorAttributes: []color.Attribute{color.FgHiBlack},
			}
		}
	}
	s.asciRepresentation = [][]runecolor.ColoredRune{colored, {}}
}

func (s *status) Width() int {
	return s.width
}

func (s *status) Height() int {
	return s.height
}

func (c *status) Render(subview [][]runecolor.ColoredRune) {
	// Clear the subview first (status text length may vary)
	for i := range subview {
		for j := range subview[i] {
			subview[i][j] = runecolor.ColoredRune{Symbol: ' '}
		}
	}
	// Copy what fits
	for i := range c.asciRepresentation {
		if i >= len(subview) {
			break
		}
		for j := range c.asciRepresentation[i] {
			if j >= len(subview[i]) {
				break
			}
			subview[i][j] = c.asciRepresentation[i][j]
		}
	}
}
