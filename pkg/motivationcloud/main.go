package motivationcloud

import "github.com/anschnapp/pomodorofactory/pkg/slicehelper"

type motivationcloud struct {
	width              int
	height             int
	asciRepresentation [][]rune
}

func MakeMotivationcloud() *motivationcloud {
	// for now static, later dynamic with wort lists and random selection
	// also different lists regarding of the state of the program
	asci := make([][]rune, 3)
	asci[0] = []rune("let's do it")
	asci[1] = []rune("           ")
	asci[2] = []rune("this will be awesome")

	height := len(asci)
	width := slicehelper.MaxWidth(asci)

	return &motivationcloud{
		width:              width,
		height:             height,
		asciRepresentation: asci,
	}
}

func (c *motivationcloud) Width() int {
	return c.width
}

func (c *motivationcloud) Height() int {
	return c.height
}

func (c *motivationcloud) Render(subview [][]rune) {
	slicehelper.Copy2DSlice(c.asciRepresentation, subview)
}
