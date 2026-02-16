package commandinput

import (
	"strings"

	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/anschnapp/pomodorofactory/pkg/slicehelper"
)

type commandinput struct {
	width              int
	height             int
	asciRepresentation [][]runecolor.ColoredRune
}

func MakeCommandinput() *commandinput {
	c := &commandinput{
		height: 3,
		width:  20,
	}
	c.SetText("[s]tart | [q]uit")
	return c
}

func (c *commandinput) SetText(text string) {
	asci := make([][]runecolor.ColoredRune, 3)
	separator := strings.Repeat("-", len(text))
	asci[0] = runecolor.ConvertSimpleRunes([]rune(separator))
	asci[1] = runecolor.ConvertSimpleRunes([]rune(text))
	asci[2] = runecolor.ConvertSimpleRunes([]rune(separator))
	c.asciRepresentation = asci
	c.width = slicehelper.MaxWidth(asci)
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
