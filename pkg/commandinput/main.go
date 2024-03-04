package commandinput

import "github.com/anschnapp/pomodorofactory/pkg/slicehelper"

type commandinput struct {
	width              int
	height             int
	asciRepresentation [][]rune
}

func MakeCommandinput() *commandinput {
	// for now static, later dynamic status bar with different kind of entries regarding of the state of the program
	asci := make([][]rune, 3)
	asci[0] = []rune("----------------")
	asci[1] = []rune("[s]tart | [q]uit")
	asci[2] = []rune("----------------")

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

func (c *commandinput) Render(subview [][]rune) {
	slicehelper.Copy2DSlice(c.asciRepresentation, subview)
}
