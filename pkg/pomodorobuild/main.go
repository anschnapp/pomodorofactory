package pomodorobuild

import "github.com/anschnapp/pomodorofactory/pkg/slicehelper"

type pomodorobuild struct {
	pomodoroFullAsci [][]rune
	width            int
	height           int
	percentage       int
}

func MakePomodoro() *pomodorobuild {
	pomodoroFullAsci := pomodoroAscii
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

func (p *pomodorobuild) Render(viewArea [][]rune) {
	println("renderable width is", p.Width())
	slicehelper.Copy2DSlice(p.pomodoroFullAsci, viewArea)
}
