package pomodorobuild

import (
	_ "embed"

	"github.com/anschnapp/pomodorofactory/pkg/iohelper"
)

//go:embed pomodoro-asci
var pomodoroAsciiSingleString string

var pomodoroAscii [][]rune

func init() {
	pomodoroAscii = iohelper.SplitMultilineStringToSlice(pomodoroAsciiSingleString)
}
