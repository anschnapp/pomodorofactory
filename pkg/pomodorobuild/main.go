package pomodorobuild

import (
	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/anschnapp/pomodorofactory/pkg/slicehelper"
	"github.com/fatih/color"
)

type pomodorobuild struct {
	pomodoroFullAsci [][]runecolor.ColoredRune
	width            int
	height           int
	percentage       int
}

func MakePomodoro() *pomodorobuild {
	pomodoroFullAsci := make([][]runecolor.ColoredRune, len(pomodoroAscii))
	for i, v := range pomodoroAscii {
		// todo make some helper methods for this...
		colorMap := make(map[rune][]color.Attribute, 0)
		defaultColor := make([]color.Attribute, 1)
		defaultColor[0] = color.FgRed
		pomodoroFullAsci[i] = runecolor.ConvertRunesToColoredRunes(v, colorMap, defaultColor)
	}
	height := len(pomodoroFullAsci)
	if height < 1 {
		panic("pomodoro file must have at least a length of 1")
	}
	maxWidth := 0
	for i := range pomodoroFullAsci {
		if len(pomodoroFullAsci[i]) > maxWidth {
			maxWidth = len(pomodoroFullAsci[i])
		}

	}
	width := maxWidth

	return &pomodorobuild{
		pomodoroFullAsci: pomodoroFullAsci,
		width:            width,
		height:           height,
		percentage:       0,
	}
}

func (p *pomodorobuild) Width() int {
	return p.width
}

func (p *pomodorobuild) Height() int {
	return p.height
}

func (p *pomodorobuild) Render(viewArea [][]runecolor.ColoredRune) {
	println("renderable width is", p.Width())
	slicehelper.Copy2DSlice(p.pomodoroFullAsci, viewArea)
}
