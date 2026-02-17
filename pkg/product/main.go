package product

import "github.com/anschnapp/pomodorofactory/pkg/runecolor"

type Product struct {
	Name  string
	Emoji string
	Art   [][]runecolor.ColoredRune // pre-colored, ready for factoryscene
}
