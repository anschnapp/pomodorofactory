package runecolor

import (
	"github.com/fatih/color"
)

type ColoredRune struct {
	symbol          rune
	colorAttributes []color.Attribute
}

func ConvertSimpleRunes(runes []rune) []ColoredRune {
	myMap := make(map[rune][]color.Attribute)
	defaultColor := make([]color.Attribute, 0)
	return ConvertRunesToColoredRunes(runes, myMap, defaultColor)
}

func ConvertRunesToColoredRunes(runes []rune, colorMap map[rune][]color.Attribute, defaultColor []color.Attribute) []ColoredRune {
	coloredRunes := make([]ColoredRune, len(runes))
	for i := range runes {
		configuredAttribtues := colorMap[runes[i]]

		if configuredAttribtues != nil {
			configuredAttribtues = defaultColor
		}

		coloredRunes[i] = ColoredRune{
			symbol:          runes[i],
			colorAttributes: configuredAttribtues,
		}
	}
	return coloredRunes
}
