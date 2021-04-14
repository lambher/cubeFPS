package models

import (
	"fmt"
	"time"

	"github.com/g3n/engine/math32"
)

const (
	MoveForward  = "MoveForward"
	MoveBackward = "MoveBackward"
	MoveLeft     = "MoveLeft"
	MoveRight    = "MoveRight"
	TurnLeft     = "TurnLeft"
	TurnRight    = "TurnRight"
	TurnUp       = "TurnUp"
	TurnDown     = "TurnDown"
)

type Player struct {
	Position                *math32.Vector3
	Direction               *math32.Vector3
	Velocity                *math32.Vector3
	Up                      *math32.Vector3
	VerticalAngle           float32
	VerticalAngleAngleSpeed float32
	HorizontalAngle         float32
	HorizontalAngleSpeed    float32
	Name                    string
	hp                      int
	deleted                 bool

	moves Moves
	world *World
}

type Moves map[string]bool

func newMoves() Moves {
	return map[string]bool{
		MoveForward: false,
	}
}

func NewPlayer(world *World, name string, position math32.Vector3) *Player {
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
		Velocity:                math32.NewVec3(),
		Name:                    name,
		VerticalAngle:           0,
		HorizontalAngle:         0,
		VerticalAngleAngleSpeed: 0,
		HorizontalAngleSpeed:    0,
		hp:                      100,
		deleted:                 false,
	}
	player.moves = newMoves()
	player.world = world

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

func (p *Player) TurnLeft(value bool, verticalAngleSpeed float32) {
	p.moves[TurnLeft] = value
	p.VerticalAngleAngleSpeed = verticalAngleSpeed
}

func (p *Player) TurnRight(value bool, verticalAngleSpeed float32) {
	p.moves[TurnRight] = value
	p.VerticalAngleAngleSpeed = verticalAngleSpeed
}

func (p *Player) TurnUp(value bool, horizontalAngleSpeed float32) {
	p.moves[TurnUp] = value
	p.HorizontalAngleSpeed = horizontalAngleSpeed
}

func (p *Player) TurnDown(value bool, horizontalAngleSpeed float32) {
	p.moves[TurnDown] = value
	p.HorizontalAngleSpeed = horizontalAngleSpeed
}

func (p Player) GetLeftAxis() *math32.Vector3 {
	return p.Direction.Clone().ApplyAxisAngle(p.Up, -math32.Pi/2)
}

func (p *Player) updateMoves() {
	if p.moves[MoveForward] {
		p.Velocity = p.Direction.Clone().MultiplyScalar(0.1)
	}
	if p.moves[MoveBackward] {
		p.Velocity = p.Direction.Clone().MultiplyScalar(-0.1)
	}
	if p.moves[MoveLeft] {
		p.Velocity = p.GetLeftAxis().MultiplyScalar(-0.1)
	}
	if p.moves[MoveRight] {
		p.Velocity = p.GetLeftAxis().MultiplyScalar(0.1)
	}
	if p.moves[TurnLeft] {
		p.VerticalAngle = p.VerticalAngleAngleSpeed
	}
	if p.moves[TurnRight] {
		p.VerticalAngle = -p.VerticalAngleAngleSpeed
	}
	if p.moves[TurnUp] {
		p.HorizontalAngle = p.HorizontalAngleSpeed
	}
	if p.moves[TurnDown] {
		p.HorizontalAngle = -p.HorizontalAngleSpeed
	}
}

func (p *Player) Fire() {
	bullet := NewBullet(p.world, p, p.Direction.Clone().MultiplyScalar(.5))
	p.world.AddBullet(bullet)
}

func (p *Player) Update(deltaTime time.Duration) {
	p.updateMoves()

	p.Position.Add(p.Velocity)
	p.Direction.ApplyAxisAngle(p.Up, p.VerticalAngle)
	p.Direction.ApplyAxisAngle(p.GetLeftAxis(), p.HorizontalAngle)
	p.Up.ApplyAxisAngle(p.GetLeftAxis(), p.HorizontalAngle)
	p.Velocity.MultiplyScalar(0.8)
	p.VerticalAngle *= 0.8
	p.HorizontalAngle *= 0.8
}

func (p Player) IsDeleted() bool {
	return p.deleted
}

func (b Player) GetID() string {
	return b.Name
}

func (p Player) GetHitBox() *math32.Sphere {
	return math32.NewSphere(p.Position, 1)
}

func (p *Player) BulletHit(bullet *Bullet) {
	p.hp -= bullet.hp
	if p.hp <= 0 {
		p.deleted = true
	}
	fmt.Println(p.hp)
	if p.world.eventListener != nil {
		p.world.eventListener.OnPlayerHit(p)
	}
}
