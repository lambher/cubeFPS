package gui

import (
	"fmt"
	"os"

	"github.com/g3n/engine/text"

	"github.com/g3n/engine/core"

	"github.com/g3n/engine/gui"
	"github.com/lambher/video-game/models"
)

type GUI struct {
	hpLabel *gui.Label
	world   *models.World

	*core.Node
}

func NewGUI(world *models.World, width, height int) *GUI {
	font, err := text.NewFont("./assets/fonts/joystix monospace.ttf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var GUI GUI

	GUI.world = world

	GUI.Node = core.NewNode()

	GUI.hpLabel = gui.NewLabel("HP")
	GUI.hpLabel.SetFontSize(50)
	GUI.hpLabel.SetFont(font)
	GUI.hpLabel.SetPosition(float32(width/2), float32(height)-100)

	GUI.Node.Add(GUI.hpLabel)

	return &GUI
}

func (g *GUI) Update() {
	g.hpLabel.SetText(fmt.Sprintf("HP:%d", g.world.Player.GetHP()))
}
