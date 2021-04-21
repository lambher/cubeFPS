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
	ID                      string
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

func NewPlayer(id string, world *World, name string, position math32.Vector3) *Player {
	player := &Player{
		ID:       id,
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

func (p Player) GetHP() int {
	return p.hp
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

const maxAngleSpeed = 0.5

func (p *Player) TurnLeft(value bool, verticalAngleSpeed float32) {
	if verticalAngleSpeed <= 0 {
		return
	}
	p.moves[TurnLeft] = value
	if verticalAngleSpeed > maxAngleSpeed {
		verticalAngleSpeed = maxAngleSpeed
	}
	p.VerticalAngleAngleSpeed = verticalAngleSpeed
}

func (p *Player) TurnRight(value bool, verticalAngleSpeed float32) {
	if verticalAngleSpeed <= 0 {
		return
	}
	p.moves[TurnRight] = value
	if verticalAngleSpeed > maxAngleSpeed {
		verticalAngleSpeed = maxAngleSpeed
	}
	p.VerticalAngleAngleSpeed = verticalAngleSpeed
}

func (p *Player) TurnUp(value bool, horizontalAngleSpeed float32) {
	if horizontalAngleSpeed <= 0 {
		return
	}
	p.moves[TurnUp] = value
	if horizontalAngleSpeed > maxAngleSpeed {
		horizontalAngleSpeed = maxAngleSpeed
	}
	p.HorizontalAngleSpeed = horizontalAngleSpeed
}

func (p *Player) TurnDown(value bool, horizontalAngleSpeed float32) {
	if horizontalAngleSpeed <= 0 {
		return
	}
	p.moves[TurnDown] = value
	if horizontalAngleSpeed > maxAngleSpeed {
		horizontalAngleSpeed = maxAngleSpeed
	}
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
	bullet := NewBullet(p.world, p, p.Direction.Clone().MultiplyScalar(.5).Add(p.Velocity))
	p.world.AddBullet(bullet)
}

func (p *Player) Update(deltaTime time.Duration) {
	p.updateMoves()

	p.Position.Add(p.Velocity)

	//direction := p.Direction.Clone()
	leftAxis := p.GetLeftAxis().Clone()
	up := p.Up.Clone()

	p.Up.ApplyAxisAngle(leftAxis, p.HorizontalAngle)
	p.Direction.ApplyAxisAngle(up, p.VerticalAngle)
	p.Direction.ApplyAxisAngle(leftAxis, p.HorizontalAngle)

	if p.Up.Length() > 1 || p.Up.Length() < 0.9 {
		p.Up.Normalize()
	}
	if p.Direction.Length() > 1 || p.Direction.Length() < 0.9 {
		p.Direction.Normalize()
	}

	p.Velocity.MultiplyScalar(0.8)
	p.VerticalAngle *= 0.8
	p.HorizontalAngle *= 0.8
}

func (p Player) IsDeleted() bool {
	return p.deleted
}

func (b Player) GetID() string {
	return b.ID
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

func (p *Player) Refresh(player Player) {
	p.Name = player.Name
	p.Up = player.Up
	p.Position = player.Position
	p.Direction = player.Direction
	p.Velocity = player.Velocity
	//p.HorizontalAngle = player.HorizontalAngle
	//p.VerticalAngle = player.VerticalAngle
}

func (p *Player) RefreshMoves(moves Moves) {
	p.moves = moves
}
