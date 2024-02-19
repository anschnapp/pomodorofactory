package motivationcloud

type motivationcloud struct {
	width              int
	height             int
	asciRepresentation []string
}

func MakeMotivationcloud() *motivationcloud {
	// for now static, later dynamic with wort lists and random selection
	// also different lists regarding of the state of the program
	asci := []string{}
	asci = append(asci, "let's do it")
	asci = append(asci, "           ")
	asci = append(asci, "this will be awesome")

	height := len(asci)
	width := len(asci[0])

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

func (c *motivationcloud) Render(subview *[]string) {
	*subview = c.asciRepresentation
}
