package status

import (
	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
)

type status struct {
	width              int
	height             int
	asciRepresentation [][]runecolor.ColoredRune
}

// Fixed width so the view region is large enough for any status text
const statusWidth = 30

func MakeStatus() *status {
	s := &status{
		height: 2,
		width:  statusWidth,
	}
	s.SetText("Press [s] to start", "")
	return s
}

func (s *status) SetText(line1, line2 string) {
	asci := make([][]runecolor.ColoredRune, 2)
	asci[0] = runecolor.ConvertSimpleRunes([]rune(line1))
	asci[1] = runecolor.ConvertSimpleRunes([]rune(line2))
	s.asciRepresentation = asci
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
