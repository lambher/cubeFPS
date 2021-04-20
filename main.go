package main

import (
	"time"

	"github.com/g3n/engine/gls"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"github.com/lambher/video-game/game"
)

func main() {
	a := app.App()
	window.Get().(*window.GlfwWindow).SetFullscreen(true)

	g := game.NewGame(a)

	g.Init()
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		g.Update(deltaTime)
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(g.Scene, g.Cam)
	})
}
