package pomodoro

import (
	"fmt"
	"github.com/anschnapp/pomodorofactory/pkg/iohelper"
)

// todo how to put this in the binary itself
const pomodoroAsciFile = "pomodoro-asci"

var newPomodoroPlainAsci []string

func init() {
	pomodoroPlainAsci, error := iohelper.ReadFileInArray(pomodoroAsciFile)
	if error == nil {
		panic("static pomodoro file cannot be loaded (should never happen)")
	}
	newPomodoroPlainAsci = pomodoroPlainAsci
	fmt.Println(len(newPomodoroPlainAsci))
}

type pomodoro struct {
	pomodoroAsciRepresentation []string
	width                      int
	height                     int
	percentage                 int
}

func (p *pomodoro) Width() int {
	return p.width
}

func (p *pomodoro) Height() int {
	return p.height
}

func MakePomodoro() *pomodoro {
	height := len(newPomodoroPlainAsci)
	fmt.Println(height)
	if height < 1 {
		panic("pomodoro file must have at least a length of 1")
	}
	firstEntry := newPomodoroPlainAsci[0]
	width := len(firstEntry)

	return &pomodoro{
		pomodoroAsciRepresentation: newPomodoroPlainAsci,
		width:                      width,
		height:                     height,
		percentage:                 0,
	}
}
