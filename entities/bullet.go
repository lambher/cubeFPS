package entities

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/lambher/video-game/models"
)

type Bullet struct {
	model    *models.Bullet
	material *material.Standard
	geometry *geometry.Geometry
	Mesh     *graphic.Mesh
}

func NewBullet(model *models.Bullet) *Bullet {
	var bullet Bullet

	bullet.material = material.NewStandard(math32.NewColor("DarkRed"))
	bullet.geometry = geometry.NewCube(0.1)
	bullet.model = model
	bullet.Mesh = graphic.NewMesh(bullet.geometry, bullet.material)
	bullet.Mesh.SetPositionVec(bullet.model.Position)
	return &bullet
}

func (b *Bullet) Update() {
	b.Mesh.SetPositionVec(b.model.Position)
}

func (b Bullet) GetMesh() *graphic.Mesh {
	return b.Mesh
}
