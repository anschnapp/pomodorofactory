package runecolor

import (
	"github.com/fatih/color"
)

type ColoredRune struct {
	Symbol          rune
	ColorAttributes []color.Attribute
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

		if configuredAttribtues == nil {
			configuredAttribtues = defaultColor
		}

		coloredRunes[i] = ColoredRune{
			Symbol:          runes[i],
			ColorAttributes: configuredAttribtues,
		}
	}
	return coloredRunes
}

func MakeSingleColorAttributes(attribute color.Attribute) []color.Attribute {
	result := make([]color.Attribute, 1)
	result[0] = attribute
	return result
}
