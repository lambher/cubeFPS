package models

import "github.com/g3n/engine/math32"

type Bullet struct {
	Player   *Player
	Position *math32.Vector3
	Velocity *math32.Vector3
}

func NewBullet(player *Player, velocity *math32.Vector3) *Bullet {
	return &Bullet{
		Player:   player,
		Position: player.Position,
		Velocity: velocity,
	}
}

func (b *Bullet) Update() {
	b.Position.Add(b.Velocity)
}
