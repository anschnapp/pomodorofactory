package product

import (
	_ "embed"

	"github.com/anschnapp/pomodorofactory/pkg/iohelper"
	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/fatih/color"
)

//go:embed art/tomato.txt
var tomatoAsciiStr string

//go:embed art/coffee.txt
var coffeeAsciiStr string

//go:embed art/penguin.txt
var penguinAsciiStr string

// All is the ordered list of buildable products.
var All []*Product

func init() {
	All = []*Product{
		makeTomato(),
		makeCoffee(),
		makePenguin(),
	}
}

func makeTomato() *Product {
	rows := iohelper.SplitMultilineStringToSlice(tomatoAsciiStr)
	colorMap := make(map[rune][]color.Attribute)
	colorMap['|'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
	colorMap['/'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
	colorMap['\\'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
	defaultColor := runecolor.MakeSingleColorAttributes(color.FgRed)

	art := make([][]runecolor.ColoredRune, len(rows))
	for i, row := range rows {
		art[i] = runecolor.ConvertRunesToColoredRunes(row, colorMap, defaultColor)
	}
	return &Product{Name: "Tomato", Emoji: "🍅", Art: art}
}

func makeCoffee() *Product {
	rows := iohelper.SplitMultilineStringToSlice(coffeeAsciiStr)
	colorMap := make(map[rune][]color.Attribute)
	colorMap['|'] = runecolor.MakeSingleColorAttributes(color.FgHiYellow)
	colorMap['_'] = runecolor.MakeSingleColorAttributes(color.FgHiYellow)
	colorMap['-'] = runecolor.MakeSingleColorAttributes(color.FgHiYellow)
	colorMap['='] = runecolor.MakeSingleColorAttributes(color.FgHiYellow)
	colorMap['~'] = runecolor.MakeSingleColorAttributes(color.FgHiWhite)
	defaultColor := runecolor.MakeSingleColorAttributes(color.FgYellow)

	art := make([][]runecolor.ColoredRune, len(rows))
	for i, row := range rows {
		art[i] = runecolor.ConvertRunesToColoredRunes(row, colorMap, defaultColor)
	}
	return &Product{Name: "Coffee Cup", Emoji: "☕", Art: art}
}

func makePenguin() *Product {
	rows := iohelper.SplitMultilineStringToSlice(penguinAsciiStr)
	colorMap := make(map[rune][]color.Attribute)
	colorMap['|'] = runecolor.MakeSingleColorAttributes(color.FgHiCyan)
	colorMap['/'] = runecolor.MakeSingleColorAttributes(color.FgHiCyan)
	colorMap['\\'] = runecolor.MakeSingleColorAttributes(color.FgHiCyan)
	colorMap['_'] = runecolor.MakeSingleColorAttributes(color.FgHiCyan)
	colorMap['^'] = runecolor.MakeSingleColorAttributes(color.FgHiCyan)
	colorMap['o'] = runecolor.MakeSingleColorAttributes(color.FgHiWhite)
	defaultColor := runecolor.MakeSingleColorAttributes(color.FgHiBlack)

	art := make([][]runecolor.ColoredRune, len(rows))
	for i, row := range rows {
		art[i] = runecolor.ConvertRunesToColoredRunes(row, colorMap, defaultColor)
	}
	return &Product{Name: "Penguin", Emoji: "🐧", Art: art}
}
