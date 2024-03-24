package commandinput

import (
	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/anschnapp/pomodorofactory/pkg/slicehelper"
)

type commandinput struct {
	width              int
	height             int
	asciRepresentation [][]runecolor.ColoredRune
}

func MakeCommandinput() *commandinput {
	// for now static, later dynamic status bar with different kind of entries regarding of the state of the program
	asci := make([][]runecolor.ColoredRune, 3)
	asci[0] = runecolor.ConvertSimpleRunes([]rune("----------------"))
	asci[1] = runecolor.ConvertSimpleRunes([]rune("[s]tart | [q]uit"))
	asci[2] = runecolor.ConvertSimpleRunes([]rune("----------------"))

	height := len(asci)
	width := slicehelper.MaxWidth(asci)

	return &commandinput{
		width:              width,
		height:             height,
		asciRepresentation: asci,
	}
}

func (c *commandinput) Width() int {
	return c.width
}

func (c *commandinput) Height() int {
	return c.height
}

func (c *commandinput) Render(subview [][]runecolor.ColoredRune) {
	slicehelper.Copy2DSlice(c.asciRepresentation, subview)
}
