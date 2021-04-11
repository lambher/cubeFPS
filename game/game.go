package game

import (
	"time"

	"github.com/lambher/video-game/models"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

type Game struct {
	world models.World
	app   *app.Application
	scene *core.Node
	cam   *camera.Camera
}

func (g *Game) Init() {
	newPlayer := models.NewPlayer("Lambert", math32.Vector3{
		X: 0,
		Y: 0,
		Z: 3,
	})
	g.world.AddPlayer(newPlayer)
	g.app = app.App()
	g.scene = core.NewNode()
	gui.Manager().Set(g.scene)

	g.cam = camera.New(1)
	g.cam.SetPositionVec(newPlayer.Position)
	g.cam.SetDirectionVec(newPlayer.Direction)
	g.scene.Add(g.cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := g.app.GetSize()
		g.app.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		g.cam.SetAspect(float32(width) / float32(height))
	}
	g.app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create a blue torus and add it to the scene
	geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	g.scene.Add(mesh)

	g.listenEvent()

	g.scene.Add(graphic.NewMesh(geometry.NewCube(1), mat))

	// Create and add a button to the scene
	btn := gui.NewButton("Make Red")
	btn.SetPosition(100, 40)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		mat.SetColor(math32.NewColor("DarkGreen"))
	})
	g.scene.Add(btn)

	// Create and add lights to the scene
	g.scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	g.scene.Add(pointLight)

	// Create and add an axis helper to the scene
	g.scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	g.app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

}

func (g Game) Run() {
	// Run the application
	g.app.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		g.update(deltaTime)
		g.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(g.scene, g.cam)
	})
}

func (g *Game) listenEvent() {
	g.app.Subscribe(window.OnKeyDown, func(evname string, ev interface{}) {
		if keyEvent, ok := ev.(*window.KeyEvent); ok {
			if keyEvent.Key == window.KeyW {
				g.world.Player.MoveForward(true)
			}
			if keyEvent.Key == window.KeyS {
				g.world.Player.MoveBackward(true)
			}
			if keyEvent.Key == window.KeyD {
				g.world.Player.MoveRight(true)
			}
			if keyEvent.Key == window.KeyA {
				g.world.Player.MoveLeft(true)
			}

		}
	})
	g.app.Subscribe(window.OnKeyUp, func(evname string, ev interface{}) {
		if keyEvent, ok := ev.(*window.KeyEvent); ok {
			if keyEvent.Key == window.KeyW {
				g.world.Player.MoveForward(false)
			}
			if keyEvent.Key == window.KeyS {
				g.world.Player.MoveBackward(false)
			}
			if keyEvent.Key == window.KeyD {
				g.world.Player.MoveRight(false)
			}
			if keyEvent.Key == window.KeyA {
				g.world.Player.MoveLeft(false)
			}
		}
	})
}

func (g *Game) update(deltaTime time.Duration) {
	g.world.Update(deltaTime)
	g.cam.SetPositionVec(g.world.Player.Position)
	g.cam.SetDirectionVec(g.world.Player.Direction)
}
