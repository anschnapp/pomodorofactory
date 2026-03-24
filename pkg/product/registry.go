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

//go:embed art/orange.txt
var oragngeAsciiStr string

//go:embed art/eifeltower.txt
var eifelTowerAsciiStr string

//go:embed art/raspberry.txt
var raspberryAsciiStr string

// All is the ordered list of buildable products.
var All []*Product

func init() {
	All = []*Product{
		makeTomato(),
		makeCoffee(),
		makePenguin(),
		makeOrange(),
		makeEifenTower(),
		makeRaspberry(),
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
	colorMap['~'] = runecolor.MakeSingleColorAttributes(color.FgHiWhite)
	colorMap['#'] = runecolor.MakeSingleColorAttributes(color.FgHiBlack)
	colorMap['|'] = runecolor.MakeSingleColorAttributes(color.FgHiCyan)
	colorMap['`'] = runecolor.MakeSingleColorAttributes(color.FgHiCyan)
	colorMap['3'] = runecolor.MakeSingleColorAttributes(color.FgHiCyan)
	defaultColor := runecolor.MakeSingleColorAttributes(color.FgYellow)

	art := make([][]runecolor.ColoredRune, len(rows))
	for i, row := range rows {
		art[i] = runecolor.ConvertRunesToColoredRunes(row, colorMap, defaultColor)
	}
	return &Product{Name: "Coffee Cup", Emoji: "☕", Art: art}
}
func makeOrange() *Product {
	rows := iohelper.SplitMultilineStringToSlice(oragngeAsciiStr)
	colorMap := make(map[rune][]color.Attribute)
	colorMap['0'] = []color.Attribute{38, 2, 255, 165, 0} // RGB orange foreground
	colorMap['\\'] = runecolor.MakeSingleColorAttributes(color.FgHiGreen)
	defaultColor := []color.Attribute{38, 2, 255, 165, 0}

	art := make([][]runecolor.ColoredRune, len(rows))
	for i, row := range rows {
		art[i] = runecolor.ConvertRunesToColoredRunes(row, colorMap, defaultColor)
	}
	return &Product{Name: "Orange", Emoji: "🍊", Art: art}
}

func makeEifenTower() *Product {
	rows := iohelper.SplitMultilineStringToSlice(eifelTowerAsciiStr)
	colorMap := make(map[rune][]color.Attribute)
	colorMap['0'] = []color.Attribute{38, 2, 220, 190, 110} // RGB light iron (left side)
	colorMap['8'] = []color.Attribute{38, 2, 155, 125, 60} // RGB medium iron (crossbeam center)
	colorMap['9'] = []color.Attribute{38, 2, 80, 60, 20}   // RGB dark iron (right side)
	defaultColor := []color.Attribute{38, 2, 220, 190, 110}

	art := make([][]runecolor.ColoredRune, len(rows))
	for i, row := range rows {
		art[i] = runecolor.ConvertRunesToColoredRunes(row, colorMap, defaultColor)
	}
	return &Product{Name: "Eifeltower", Emoji: "🗼", Art: art}
}

func makeRaspberry() *Product {
	rows := iohelper.SplitMultilineStringToSlice(raspberryAsciiStr)
	colorMap := make(map[rune][]color.Attribute)
	colorMap['\\'] = runecolor.MakeSingleColorAttributes(color.FgYellow)        // hair
	colorMap['('] = runecolor.MakeSingleColorAttributes(color.FgHiMagenta)       // raspberry body
	colorMap[')'] = runecolor.MakeSingleColorAttributes(color.FgHiMagenta)       // raspberry body
	colorMap['/'] = runecolor.MakeSingleColorAttributes(color.FgMagenta)         // sticks (darker raspberry)
	colorMap['*'] = []color.Attribute{38, 2, 100, 149, 237} // ice (cornflower blue)
	defaultColor := runecolor.MakeSingleColorAttributes(color.FgHiMagenta)

	art := make([][]runecolor.ColoredRune, len(rows))
	for i, row := range rows {
		art[i] = runecolor.ConvertRunesToColoredRunes(row, colorMap, defaultColor)
	}
	return &Product{Name: "Raspberry on Ice", Emoji: "🍧", Art: art}
}

func makePenguin() *Product {
	rows := iohelper.SplitMultilineStringToSlice(penguinAsciiStr)
	colorMap := make(map[rune][]color.Attribute)
	colorMap['@'] = runecolor.MakeSingleColorAttributes(color.FgHiBlack)
	colorMap['#'] = runecolor.MakeSingleColorAttributes(color.FgHiWhite)
	colorMap['|'] = runecolor.MakeSingleColorAttributes(color.FgHiBlack)
	colorMap['\\'] = runecolor.MakeSingleColorAttributes(color.FgHiBlack)
	colorMap['<'] = runecolor.MakeSingleColorAttributes(color.FgHiWhite)
	colorMap['*'] = runecolor.MakeSingleColorAttributes(color.FgCyan)
	colorMap['%'] = runecolor.MakeSingleColorAttributes(color.FgHiMagenta)
	defaultColor := runecolor.MakeSingleColorAttributes(color.FgHiBlack)

	art := make([][]runecolor.ColoredRune, len(rows))
	for i, row := range rows {
		art[i] = runecolor.ConvertRunesToColoredRunes(row, colorMap, defaultColor)
	}
	return &Product{Name: "Penguin", Emoji: "🐧", Art: art}
}
