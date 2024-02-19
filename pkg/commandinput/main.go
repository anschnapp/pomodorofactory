package commandinput

type commandinput struct {
	width              int
	height             int
	asciRepresentation []string
}

func MakeCommandinput() *commandinput {
	// for now static, later dynamic status bar with different kind of entries regarding of the state of the program
	asci := []string{}
	asci = append(asci, "[s]tart")
	asci = append(asci, "[q]uit")

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

func (c *commandinput) Render(subview *[]string) {
	*subview = c.asciRepresentation
}
