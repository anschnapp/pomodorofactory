package render

import "github.com/anschnapp/pomodorofactory/pkg/runecolor"

type Renderable interface {
	Render([][]runecolor.ColoredRune)
	Width() int
	Height() int
}
