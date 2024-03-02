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
	var view *view.View = view.MakeView(pomodorobuild, motivationcloud, status, commandinput)

	view.Render()
	view.Print()
}

