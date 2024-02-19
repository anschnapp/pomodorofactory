package status

type status struct {
	width              int
	height             int
	asciRepresentation [][]rune
}

func MakeStatus() *status {
	// for now static, later dynamic status bar with different kind of entries regarding of the state of the program
	asci := [][]rune{}
	asci = append(asci, "[s]tart")
	asci = append(asci, "[q]uit")

	height := len(asci)
	width := len(asci[0])

	return &status{
		width:              width,
		height:             height,
		asciRepresentation: asci,
	}
}

func (s *status) Width() int {
	return s.width
}

func (s *status) Height() int {
	return s.height
}

func (c *status) Render(subview *[][]rune) {
	*subview = c.asciRepresentation
}
