package pomodorobuild

type pomodorobuild struct {
	pomodoroFullAsci []string
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
	firstEntry := pomodoroFullAsci
	width := len(firstEntry)

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

func (p *pomodorobuild) Render(viewArea *[]string) {
	*viewArea = p.pomodoroFullAsci
}
