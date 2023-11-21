package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/sread/app"
	scene_main "github.com/t0l1k/sread/app/scenes/start"
)

func main() {
	eui.Init(app.NewApp())
	eui.Run(scene_main.NewRRStartScene())
	eui.Quit()
}
