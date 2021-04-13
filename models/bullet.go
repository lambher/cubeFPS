package models

import (
	"time"

	"github.com/g3n/engine/math32"
	"github.com/rs/xid"
)

type Bullet struct {
	ID       string
	Player   *Player
	Position *math32.Vector3
	Velocity *math32.Vector3

	hp      int
	deleted bool
	world   *World
}

func NewBullet(world *World, player *Player, velocity *math32.Vector3) *Bullet {
	return &Bullet{
		ID:       xid.New().String(),
		Player:   player,
		Position: player.Position.Clone().Add(player.Velocity),
		Velocity: velocity,
		world:    world,
		hp:       10,
	}
}

func (b *Bullet) Update(deltaTime time.Duration) {
	b.Position.Add(b.Velocity)
	for _, player := range b.world.players {
		if player != b.Player {
			if player.GetHitBox().ContainsPoint(b.Position) {
				b.world.BulletHitsPlayer(b, player)
				b.deleted = true
			}
		}
	}
	if b.Position.Length() > 10 {
		b.deleted = true
	}
}

func (b Bullet) IsDeleted() bool {
	return b.deleted
}

func (b Bullet) GetID() string {
	return b.ID
}
