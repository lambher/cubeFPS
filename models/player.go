package models

import (
	"time"

	"github.com/g3n/engine/math32"
)

const (
	MoveForward  = "MoveForward"
	MoveBackward = "MoveBackward"
	MoveLeft     = "MoveLeft"
	MoveRight    = "MoveRight"
)

type Player struct {
	Position  *math32.Vector3
	Direction *math32.Vector3
	Velocity  *math32.Vector3
	Up        *math32.Vector3
	Name      string
	moves     Moves
}

type Moves map[string]bool

func newMoves() Moves {
	return map[string]bool{
		MoveForward: false,
	}
}

func NewPlayer(name string, position math32.Vector3) *Player {
	player := &Player{
		Position: &position,
		Direction: &math32.Vector3{
			X: 0,
			Y: 0,
			Z: -1,
		},
		Up: &math32.Vector3{
			X: 0,
			Y: 1,
			Z: 0,
		},
		Velocity: math32.NewVec3(),
		Name:     name,
	}
	player.moves = newMoves()

	return player
}

func (p *Player) MoveForward(value bool) {
	p.moves[MoveForward] = value
}

func (p *Player) MoveBackward(value bool) {
	p.moves[MoveBackward] = value
}

func (p *Player) MoveLeft(value bool) {
	p.moves[MoveLeft] = value
}

func (p *Player) MoveRight(value bool) {
	p.moves[MoveRight] = value
}

func (p *Player) updateMoves() {
	if p.moves[MoveForward] {
		p.Velocity = p.Direction.Clone().MultiplyScalar(0.1)
	}
	if p.moves[MoveBackward] {
		p.Velocity = p.Direction.Clone().MultiplyScalar(-0.1)
	}
	if p.moves[MoveLeft] {
		p.Velocity = p.Direction.Clone().ApplyAxisAngle(p.Up, -math32.Pi/2).MultiplyScalar(-0.1)
	}
	if p.moves[MoveRight] {
		p.Velocity = p.Direction.Clone().ApplyAxisAngle(p.Up, math32.Pi/2).MultiplyScalar(-0.1)
	}
}

func (p *Player) Update(deltaTime time.Duration) {
	p.updateMoves()

	p.Position.Add(p.Velocity)
	p.Velocity.MultiplyScalar(0.8)
}
