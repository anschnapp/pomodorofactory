package pomodorobuild

import (
	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/anschnapp/pomodorofactory/pkg/slicehelper"
	"github.com/fatih/color"
)

type pomodorobuild struct {
	// Full tomato with all x chars colored red
	pomodoroFullAsci [][]runecolor.ColoredRune
	// Which rows contain 'x' fill chars (indices into pomodoroFullAsci)
	bodyRows []int
	// Current render state (rebuilt on SetPercentage)
	currentFrame [][]runecolor.ColoredRune
	width        int
	height       int
	percentage   int
}

func MakePomodoro() *pomodorobuild {
	pomodoroFullAsci := make([][]runecolor.ColoredRune, len(pomodoroAscii))
	var bodyRows []int

	for i, v := range pomodoroAscii {
		colorMap := make(map[rune][]color.Attribute, 3)
		colorMap['|'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
		colorMap['/'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
		colorMap['\\'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
		defaultColor := runecolor.MakeSingleColorAttributes(color.FgRed)
		pomodoroFullAsci[i] = runecolor.ConvertRunesToColoredRunes(v, colorMap, defaultColor)

		// Detect body rows: any row with non-space content
		for _, r := range v {
			if r != ' ' {
				bodyRows = append(bodyRows, i)
				break
			}
		}
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

	p := &pomodorobuild{
		pomodoroFullAsci: pomodoroFullAsci,
		bodyRows:         bodyRows,
		width:            maxWidth,
		height:           height,
		percentage:       0,
	}
	p.rebuildFrame()
	return p
}

func (p *pomodorobuild) SetPercentage(pct int) {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	p.percentage = pct
	p.rebuildFrame()
}

// rebuildFrame builds currentFrame from pomodoroFullAsci,
// blanking out 'x' chars on unfilled rows.
func (p *pomodorobuild) rebuildFrame() {
	totalBody := len(p.bodyRows)
	// Number of body rows filled (from the bottom)
	filledCount := p.percentage * totalBody / 100

	// Build a set of unfilled row indices
	unfilled := make(map[int]bool)
	for i := 0; i < totalBody-filledCount; i++ {
		unfilled[p.bodyRows[i]] = true
	}

	p.currentFrame = make([][]runecolor.ColoredRune, p.height)
	emptyColor := make([]color.Attribute, 0)

	for row := range p.pomodoroFullAsci {
		src := p.pomodoroFullAsci[row]
		dst := make([]runecolor.ColoredRune, len(src))
		copy(dst, src)

		if unfilled[row] {
			// Replace all visible chars with spaces on unfilled rows
			for col := range dst {
				if dst[col].Symbol != ' ' {
					dst[col] = runecolor.ColoredRune{
						Symbol:          ' ',
						ColorAttributes: emptyColor,
					}
				}
			}
		}
		p.currentFrame[row] = dst
	}
}

func (p *pomodorobuild) Width() int {
	return p.width
}

func (p *pomodorobuild) Height() int {
	return p.height
}

func (p *pomodorobuild) Render(viewArea [][]runecolor.ColoredRune) {
	slicehelper.Copy2DSlice(p.currentFrame, viewArea)
}
