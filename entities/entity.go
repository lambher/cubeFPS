package entities

import "github.com/g3n/engine/graphic"

type Entity interface {
	Update()
	GetMesh() *graphic.Mesh
}
