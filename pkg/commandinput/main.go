package commandinput

import "github.com/anschnapp/pomodorofactory/pkg/slicehelper"

type commandinput struct {
	width              int
	height             int
	asciRepresentation [][]rune
}

func MakeCommandinput() *commandinput {
	// for now static, later dynamic status bar with different kind of entries regarding of the state of the program
	asci := [][]rune{}
	asci[0] = []rune("[s]tart")
	asci[1] = []rune("[q]uit")

	height := len(asci)
	width := len(asci[0])

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
