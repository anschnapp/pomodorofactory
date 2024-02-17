package pomodorobuild

import (
	_ "embed"
	"github.com/anschnapp/pomodorofactory/pkg/iohelper"
)

//go:embed pomodoro-asci
var pomodoroAsciiSingleString string

var pomodoroAscii []string

func init() {
	pomodoroAscii = iohelper.SplitMultilineStringToArray(pomodoroAsciiSingleString)
}
