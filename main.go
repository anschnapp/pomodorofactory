package main

import (
	"github.com/anschnapp/pomodorofactory/pkg/commandinput"
	"github.com/anschnapp/pomodorofactory/pkg/motivationcloud"
	"github.com/anschnapp/pomodorofactory/pkg/pomodorobuild"
	"github.com/anschnapp/pomodorofactory/pkg/render"
	"github.com/anschnapp/pomodorofactory/pkg/status"
	"github.com/anschnapp/pomodorofactory/pkg/view"
)

func main() {
	var pomodorobuild render.Renderable = pomodorobuild.MakePomodoro()
	var motivationcloud render.Renderable = motivationcloud.MakeMotivationcloud()
	var status render.Renderable = status.MakeStatus()
	var commandinput render.Renderable = commandinput.MakeCommandinput()
	var view render.Renderable = view.MakeView(pomodorobuild, motivationcloud, status, commandinput)

	blankSpace := generateBlankSpace(view.Height(), view.Width())
	view.Render(&blankSpace)
}

func generateBlankSpace(height int, width int) []string {
	blankSpace := make([]string, height)

	for i := range blankSpace {
		blankSpace[i] = ""
	}
	for i := range blankSpace {
		for j := 0; j < width; j++ {
			blankSpace[i] = blankSpace[i] + " "
		}
	}
	return blankSpace
}
