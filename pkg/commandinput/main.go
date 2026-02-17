package commandinput

import (
	"strings"

	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/anschnapp/pomodorofactory/pkg/slicehelper"
)

const cmdWidth = 50

type commandinput struct {
	width              int
	height             int
	asciRepresentation [][]runecolor.ColoredRune
}

func MakeCommandinput() *commandinput {
	c := &commandinput{
		height: 4,
		width:  cmdWidth,
	}
	c.SetTexts("[s]tart | [q]uit", "")
	return c
}

// SetTexts updates the command bar with a command line and an optional selector line.
// Both lines are padded/truncated to fit the fixed width.
func (c *commandinput) SetTexts(commandText, selectorText string) {
	sep := strings.Repeat("-", c.width)
	asci := make([][]runecolor.ColoredRune, 4)
	asci[0] = runecolor.ConvertSimpleRunes([]rune(sep))
	asci[1] = runecolor.ConvertSimpleRunes(padToWidth(commandText, c.width))
	asci[2] = runecolor.ConvertSimpleRunes(padToWidth(selectorText, c.width))
	asci[3] = runecolor.ConvertSimpleRunes([]rune(sep))
	c.asciRepresentation = asci
}

func padToWidth(s string, width int) []rune {
	runes := []rune(s)
	if len(runes) >= width {
		return runes[:width]
	}
	padded := make([]rune, width)
	copy(padded, runes)
	for i := len(runes); i < width; i++ {
		padded[i] = ' '
	}
	return padded
}

func (c *commandinput) Width() int {
	return c.width
}

func (c *commandinput) Height() int {
	return c.height
}

func (c *commandinput) Render(subview [][]runecolor.ColoredRune) {
	for i := range subview {
		for j := range subview[i] {
			subview[i][j] = runecolor.ColoredRune{Symbol: ' '}
		}
	}
	slicehelper.Copy2DSlice(c.asciRepresentation, subview)
}
