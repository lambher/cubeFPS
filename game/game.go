package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/lambher/video-game/conf"

	"github.com/go-gl/glfw/v3.3/glfw"
	gui2 "github.com/lambher/video-game/gui"

	"github.com/lambher/video-game/entities"

	"github.com/lambher/video-game/models"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

type Game struct {
	world         *models.World
	app           *app.Application
	Scene         *core.Node
	Cam           *camera.Camera
	gui           *gui2.GUI
	started       bool
	menu          *core.Node
	mousePosition *math32.Vector2

	entities map[string]entities.Entity

	conn net.Conn
}

func (g *Game) OnAddPlayer(player *models.Player) {
	if player == g.world.Player {
		g.Scene.Add(g.menu)
		return
	}
	if g.entities == nil {
		g.entities = make(map[string]entities.Entity)
	}

	entity := entities.NewPlayer(player)
	g.entities[player.GetID()] = entity

	g.Scene.Add(entity.Mesh)
}

func (g *Game) OnPlayerHit(player *models.Player) {
	if g.entities == nil {
		g.entities = make(map[string]entities.Entity)
	}

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
	if g.entities == nil {
		g.entities = make(map[string]entities.Entity)
	}

	entity := entities.NewBullet(bullet)
	g.entities[bullet.GetID()] = entity

	g.Scene.Add(entity.Mesh)
}

func (g *Game) OnRemoveModel(model models.Model) {
	if entity, ok := g.entities[model.GetID()]; ok {
		//if player, ok := entity.(*entities.Player); ok {
		//
		//}
		g.Scene.Remove(entity.GetMesh())
		delete(g.entities, model.GetID())
	}
}

func (g *Game) AddPlayer(player *models.Player) {
	g.world.AddPlayer(player)
}

func NewGame(app *app.Application) *Game {
	return &Game{
		app: app,
	}
}

func (g *Game) connect() {
	p := make([]byte, 2048)
	var err error
	g.conn, err = net.Dial("udp", conf.Host+":"+strconv.Itoa(conf.Port))
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(g.conn, "hello")
	n, err := bufio.NewReader(g.conn).Read(p)
	if err == nil {
		g.parse(string(p[:n]))
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	defer g.conn.Close()
	g.listen()
}

func (g *Game) refreshPlayer() {
	playerData, err := g.world.GetPlayerData()

	data := make([]byte, 0)

	data = append(data, []byte("refresh_player\n")...)
	data = append(data, playerData...)

	_, err = fmt.Fprintf(g.conn, string(data))
	if err != nil {
		fmt.Println(err)
	}
}

func (g *Game) sendMove() {
	playerMoveData, err := g.world.GetPlayerMoveData()

	data := make([]byte, 0)

	data = append(data, []byte("move\n")...)
	data = append(data, playerMoveData...)

	_, err = fmt.Fprintf(g.conn, string(data))
	if err != nil {
		fmt.Println(err)
	}
}

func (g *Game) sendFire() {
	data := make([]byte, 0)

	data = append(data, []byte("fire\n")...)

	_, err := fmt.Fprintf(g.conn, string(data))
	if err != nil {
		fmt.Println(err)
	}
}

func (g *Game) listen() {
	for {
		p := make([]byte, 2048)
		n, err := bufio.NewReader(g.conn).Read(p)
		if err == nil {
			g.parse(string(p[:n]))
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	}
}

func (g *Game) parse(message string) {
	messages := strings.Split(message, "\n")
	if len(messages) <= 1 {
		return
	}
	switch messages[0] {
	case "you":
		g.handleYou([]byte(messages[1]))
	case "add_player":
		g.handleAddPlayer([]byte(messages[1]))
	case "exit":
		g.handleExit([]byte(messages[1]))
	case "refresh_player":
		g.handleRefreshPlayer([]byte(messages[1]))
	case "fire":
		g.handleFire([]byte(messages[1]))
	}
}

func (g *Game) handleRefreshPlayer(data []byte) {
	var player models.Player

	err := json.Unmarshal(data, &player)
	if err != nil {
		fmt.Println(err)
		return
	}
	if player.Position == nil {
		fmt.Println("player position is null")
		return
	}
	if p := g.world.GetPlayer(player.GetID()); p != nil {
		p.Refresh(player)
	}
}

func (g *Game) handleFire(data []byte) {
	var player models.Player

	err := json.Unmarshal(data, &player)
	if err != nil {
		fmt.Println(err)
		return
	}
	if player.Position == nil {
		fmt.Println("player position is null")
		return
	}
	if p := g.world.GetPlayer(player.GetID()); p != nil {
		p.Refresh(player)
		p.Fire()
	}
}

func (g *Game) handleAddPlayer(data []byte) {
	var player models.Player

	err := json.Unmarshal(data, &player)
	if err != nil {
		fmt.Println(err)
		return
	}
	if player.Position == nil {
		fmt.Println("player position is null")
		return
	}
	newPlayer := models.NewPlayer(player.GetID(), g.world, player.Name, *player.Position)

	g.AddPlayer(newPlayer)
}

func (g *Game) handleExit(data []byte) {
	var player models.Player

	err := json.Unmarshal(data, &player)
	if err != nil {
		fmt.Println(err)
		return
	}
	if player.Position == nil {
		fmt.Println("player position is null")
		return
	}

	g.world.RemovePlayer(&player)
}

func (g *Game) handleYou(data []byte) {
	var player models.Player

	err := json.Unmarshal(data, &player)
	if err != nil {
		fmt.Println(err)
		return
	}
	if player.Position == nil {
		fmt.Println("player position is null")
		return
	}
	newPlayer := models.NewPlayer(player.GetID(), g.world, player.Name, *player.Position)

	g.AddPlayer(newPlayer)
}

func (g *Game) Init() {
	g.world = &models.World{}
	g.world.SubscribeEventListener(g)

	go g.connect()

	g.Scene = core.NewNode()
	gui.Manager().Set(g.Scene)

	//newPlayer := models.NewPlayer(g.world, "Lambert", math32.Vector3{
	//	X: 0,
	//	Y: 0,
	//	Z: 3,
	//})
	//
	////g.mouseVelocity = math32.NewVec2()
	//
	//g.AddPlayer(newPlayer)
	//
	//g.AddPlayer(models.NewPlayer(g.world, "Milande", math32.Vector3{
	//	X: 0,
	//	Y: 0,
	//	Z: -3,
	//}))
	//g.AddPlayer(models.NewPlayer(g.world, "Etienne", math32.Vector3{
	//	X: 0,
	//	Y: 3,
	//	Z: -3,
	//}))
	//g.AddPlayer(models.NewPlayer(g.world, "Patrick", math32.Vector3{
	//	X: -3,
	//	Y: 3,
	//	Z: -3,
	//}))

	width, height := g.app.GetSize()

	g.menu = core.NewNode()
	editName := gui.NewEdit(200, "Enter your name")
	startButton := gui.NewButton("Start")
	exitButton := gui.NewButton("Exit")
	startButton.Subscribe(gui.OnClick, func(s string, i interface{}) {
		if editName.Text() == "" {
			return
		}
		g.world.Player.Name = editName.Text()
		g.start()
	})
	exitButton.Subscribe(gui.OnClick, func(s string, i interface{}) {
		g.app.Exit()
	})
	editName.SetPosition(float32(width)/2, float32(height)/2)
	startButton.SetPosition(float32(width)/2, float32(height)/2+editName.ContentHeight())
	exitButton.SetPosition(float32(width)/2, float32(height)/2+editName.ContentHeight()+startButton.ContentHeight())
	g.menu.Add(editName)
	g.menu.Add(startButton)
	g.menu.Add(exitButton)

	g.Cam = camera.New(1)
	//g.Cam.SetPositionVec(newPlayer.Position)
	//g.Cam.SetDirectionVec(newPlayer.Direction)
	g.Scene.Add(g.Cam)

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

	g.Scene.Add(skybox)

	//fmt.Println(glfw.CursorMode)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly

		g.app.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		g.Cam.SetAspect(float32(width) / float32(height))
	}
	g.app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	g.initGUI()
	//// Create and add a button to the scene
	//btn := gui.NewButton("Make Red")
	//btn.SetPosition(100, 40)
	//btn.SetSize(40, 40)
	//btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
	//})
	//g.scene.Add(btn)

	// Create and add lights to the scene
	g.Scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	g.Scene.Add(pointLight)

	// Create and add an axis helper to the scene
	g.Scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	g.app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
}

//func (g *Game) Run() {
//	// Run the application
//	g.app.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
//		g.update(deltaTime)
//		g.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
//		renderer.Render(g.scene, g.cam)
//	})
//}

func (g *Game) start() {
	g.refreshPlayer()
	g.Scene.Remove(g.menu)
	g.app.IWindow.(*window.GlfwWindow).SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	g.started = true
	g.listenEvent()
	fmt.Println("start")
}

func (g *Game) pause() {
	g.Scene.Add(g.menu)
	g.app.IWindow.(*window.GlfwWindow).SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	g.started = false
}

func (g *Game) initGUI() {
	width, height := g.app.GetSize()
	g.gui = gui2.NewGUI(g.world, width, height)
	g.Scene.Add(g.gui)
}

func (g *Game) listenEvent() {
	g.app.Subscribe(window.OnKeyDown, func(evname string, ev interface{}) {
		if !g.started {
			return
		}
		if keyEvent, ok := ev.(*window.KeyEvent); ok {
			if keyEvent.Key == window.KeyW {
				g.world.Player.MoveForward(true)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyS {
				g.world.Player.MoveBackward(true)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyD {
				g.world.Player.MoveRight(true)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyA {
				g.world.Player.MoveLeft(true)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyLeft {
				g.world.Player.TurnLeft(true, 0.5)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyRight {
				g.world.Player.TurnRight(true, 0.5)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyUp {
				g.world.Player.TurnUp(true, 0.5)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyDown {
				g.world.Player.TurnDown(true, 0.5)
				g.sendMove()
			}

			if keyEvent.Key == window.KeyEscape {
				g.pause()
			}
		}
	})

	g.app.Subscribe(window.OnKeyUp, func(evname string, ev interface{}) {
		if !g.started {
			return
		}
		if keyEvent, ok := ev.(*window.KeyEvent); ok {
			if keyEvent.Key == window.KeyW {
				g.world.Player.MoveForward(false)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyS {
				g.world.Player.MoveBackward(false)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyD {
				g.world.Player.MoveRight(false)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyA {
				g.world.Player.MoveLeft(false)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyLeft {
				g.world.Player.TurnLeft(false, 0.01)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyRight {
				g.world.Player.TurnRight(false, 0.01)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyUp {
				g.world.Player.TurnUp(false, 0.01)
				g.sendMove()
			}
			if keyEvent.Key == window.KeyDown {
				g.world.Player.TurnDown(false, 0.01)
				g.sendMove()
			}
		}
	})

	g.app.Subscribe(window.OnMouseDown, func(evname string, ev interface{}) {
		if !g.started {
			return
		}
		if mouseEvent, ok := ev.(*window.MouseEvent); ok {
			if mouseEvent.Button == window.MouseButton1 {
				g.sendFire()
			}
		}
	})

	x, y := g.app.GetSize()
	g.mousePosition = math32.NewVector2(float32(x/2), float32(y/2))

	g.app.Subscribe(window.OnCursor, func(evname string, ev interface{}) {
		if !g.started {
			return
		}
		if cursorEvent, ok := ev.(*window.CursorEvent); ok {
			g.world.Player.TurnLeft(false, 1)
			g.world.Player.TurnRight(false, 1)
			g.world.Player.TurnDown(false, 1)
			g.world.Player.TurnUp(false, 1)

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
			g.sendMove()

		}
	})
}

func (g *Game) SendExit() {
	_, err := fmt.Fprintf(g.conn, "exit")
	if err != nil {
		fmt.Println(err)
	}
}

func (g *Game) Update(deltaTime time.Duration) {
	if !g.started {
		return
	}

	//g.axes.SetDirectionVec(g.world.Player.Direction)
	g.gui.Update()
	//g.world.Player.Update(deltaTime)
	g.world.UpdatePositions(deltaTime)
	g.Cam.SetPositionVec(g.world.Player.Position)
	//g.cam.SetDirectionVec(g.world.Player.Direction)
	g.Cam.LookAt(g.world.Player.Direction.Clone().Add(g.world.Player.Position), g.world.Player.Up)
	x, y := g.app.GetSize()
	g.mousePosition = math32.NewVector2(float32(x/2), float32(y/2))
	g.app.IWindow.(*window.GlfwWindow).SetCursorPos(float64(g.mousePosition.X), float64(g.mousePosition.Y))

	for _, entity := range g.entities {
		entity.Update()
	}
}
