package entities

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/lambher/video-game/models"
)

type Player struct {
	model    *models.Player
	material *material.Standard
	geometry *geometry.Geometry
	*graphic.Mesh
}

func NewPlayer(model *models.Player) *Player {
	var player Player

	player.material = material.NewStandard(math32.NewColor("DarkBlue"))
	player.geometry = geometry.NewCube(1)
	player.model = model
	player.Mesh = graphic.NewMesh(player.geometry, player.material)
	player.Mesh.SetPositionVec(player.model.Position)
	return &player
}

func (p *Player) Update() {
	p.Mesh.SetPositionVec(p.model.Position)
}

func (p Player) GetMesh() *graphic.Mesh {
	return p.Mesh
}
