package models

import (
	"time"
)

type World struct {
	Player        *Player
	players       map[string]*Player
	bullets       []*Bullet
	eventListener EventListener
}

type EventListener interface {
	OnAddPlayer(player *Player)
	OnAddBullet(bullet *Bullet)
}

func (w *World) SubscribeEventListener(e EventListener) {
	w.eventListener = e
}

func (w *World) AddPlayer(player *Player) {
	if w.players == nil {
		w.players = make(map[string]*Player)
	}
	if w.Player == nil {
		w.Player = player
	}
	w.players[player.Name] = player
	if w.eventListener != nil {
		w.eventListener.OnAddPlayer(player)
	}
}

func (w *World) AddBullet(bullet *Bullet) {
	if w.bullets == nil {
		w.bullets = make([]*Bullet, 0)
	}
	w.bullets = append(w.bullets, bullet)
	if w.eventListener != nil {
		w.eventListener.OnAddBullet(bullet)
	}
}

func (w *World) Update(deltaTime time.Duration) {
	for _, player := range w.players {
		player.Update(deltaTime)
	}
}
