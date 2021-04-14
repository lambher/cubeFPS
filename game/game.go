package game

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/lambher/video-game/entities"

	"github.com/go-gl/glfw/v3.3/glfw"
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
	world *models.World
	app   *app.Application
	scene *core.Node
	cam   *camera.Camera

	mousePosition *math32.Vector2

	entities map[string]entities.Entity
}

func (g *Game) OnAddPlayer(player *models.Player) {
	if player == g.world.Player {
		return
	}
	if g.entities == nil {
		g.entities = make(map[string]entities.Entity)
	}

	entity := entities.NewPlayer(player)
	g.entities[player.GetID()] = entity

	g.scene.Add(entity.Mesh)
}

func (g *Game) OnPlayerHit(player *models.Player) {
	if p, ok := g.entities[player.GetID()].(*entities.Player); ok {
		p.Hit()
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (g *Game) OnAddBullet(bullet *models.Bullet) {
	entity := entities.NewBullet(bullet)
	g.entities[bullet.GetID()] = entity

	g.scene.Add(entity.Mesh)
}

func (g *Game) OnRemoveModel(model models.Model) {
	if entity, ok := g.entities[model.GetID()]; ok {
		//if player, ok := entity.(*entities.Player); ok {
		//
		//}
		g.scene.Remove(entity.GetMesh())
		delete(g.entities, model.GetID())
	}
}

func (g *Game) AddPlayer(player *models.Player) {
	g.world.AddPlayer(player)
}

func (g *Game) Init() {
	g.world = &models.World{}
	g.world.SubscribeEventListener(g)

	g.app = app.App()
	g.scene = core.NewNode()
	gui.Manager().Set(g.scene)
	window.Get().(*window.GlfwWindow).SetFullscreen(true)

	newPlayer := models.NewPlayer(g.world, "Lambert", math32.Vector3{
		X: 0,
		Y: 0,
		Z: 3,
	})

	//g.mouseVelocity = math32.NewVec2()

	g.AddPlayer(newPlayer)

	g.AddPlayer(models.NewPlayer(g.world, "Milande", math32.Vector3{
		X: 0,
		Y: 0,
		Z: 3,
	}))

	g.cam = camera.New(1)
	g.cam.SetPositionVec(newPlayer.Position)
	g.cam.SetDirectionVec(newPlayer.Direction)
	g.scene.Add(g.cam)

	skybox, err := graphic.NewSkybox(graphic.SkyboxData{
		DirAndPrefix: "./assets/textures/skyboxes/lambert/",
		Extension:    "png",
		Suffixes: [6]string{
			"right",
			"left",
			"top",
			"bottom",
			"front",
			"back",
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	g.scene.Add(skybox)

	//fmt.Println(glfw.CursorMode)

	g.app.IWindow.(*window.GlfwWindow).SetInputMode(glfw.CursorMode, glfw.CursorHidden)
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

	// Create and add a button to the scene
	btn = gui.NewButton("Reset Player")
	btn.SetPosition(100, 90)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		g.resetPlayer()
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

func (g *Game) Run() {
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
			if keyEvent.Key == window.KeyLeft {
				g.world.Player.TurnLeft(true, 0.01)
			}
			if keyEvent.Key == window.KeyRight {
				g.world.Player.TurnRight(true, 0.01)
			}
			if keyEvent.Key == window.KeyUp {
				g.world.Player.TurnUp(true, 0.01)
			}
			if keyEvent.Key == window.KeyDown {
				g.world.Player.TurnDown(true, 0.01)
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
			if keyEvent.Key == window.KeyLeft {
				g.world.Player.TurnLeft(false, 0.01)
			}
			if keyEvent.Key == window.KeyRight {
				g.world.Player.TurnRight(false, 0.01)
			}
			if keyEvent.Key == window.KeyUp {
				g.world.Player.TurnUp(false, 0.01)
			}
			if keyEvent.Key == window.KeyDown {
				g.world.Player.TurnDown(false, 0.01)
			}
		}
	})

	g.app.Subscribe(window.OnMouseDown, func(evname string, ev interface{}) {
		if mouseEvent, ok := ev.(*window.MouseEvent); ok {
			if mouseEvent.Button == window.MouseButton1 {
				g.world.Player.Fire()
			}
		}
	})

	x, y := g.app.GetSize()
	g.mousePosition = math32.NewVector2(float32(x/2), float32(y/2))

	g.app.Subscribe(window.OnCursor, func(evname string, ev interface{}) {
		if cursorEvent, ok := ev.(*window.CursorEvent); ok {
			g.world.Player.TurnLeft(false, 0)
			g.world.Player.TurnRight(false, 0)
			g.world.Player.TurnDown(false, 0)
			g.world.Player.TurnUp(false, 0)

			x := -g.mousePosition.X + cursorEvent.Xpos
			y := -g.mousePosition.Y + cursorEvent.Ypos

			x *= 0.002
			y *= 0.002

			if x < 0 {
				g.world.Player.TurnLeft(true, -x)
			}
			if x > 0 {
				g.world.Player.TurnRight(true, x)
			}
			if y < 0 {
				g.world.Player.TurnUp(true, -y)
			}
			if y > 0 {
				g.world.Player.TurnDown(true, y)
			}
		}
	})
}

func (g *Game) resetPlayer() {
	g.world.Player = models.NewPlayer(g.world, "Lambert", math32.Vector3{
		X: 0,
		Y: 0,
		Z: 3,
	})
}

func (g *Game) update(deltaTime time.Duration) {
	g.world.Update(deltaTime)
	g.cam.SetPositionVec(g.world.Player.Position)
	//g.cam.SetDirectionVec(g.world.Player.Direction)
	g.cam.LookAt(g.world.Player.Direction.Clone().Add(g.world.Player.Position), g.world.Player.Up)
	x, y := g.app.GetSize()
	g.mousePosition = math32.NewVector2(float32(x/2), float32(y/2))
	g.app.IWindow.(*window.GlfwWindow).SetCursorPos(float64(g.mousePosition.X), float64(g.mousePosition.Y))

	for _, entity := range g.entities {
		entity.Update()
	}
}
