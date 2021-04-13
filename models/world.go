package models

import (
	"fmt"
	"time"
)

type World struct {
	Player        *Player
	players       map[string]*Player
	models        map[string]Model
	eventListener EventListener
}

type EventListener interface {
	OnAddPlayer(player *Player)
	OnAddBullet(bullet *Bullet)
	OnRemoveModel(model Model)
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

func (w *World) BulletHitsPlayer(bullet *Bullet, player *Player) {
	player.hp -= bullet.hp
	fmt.Println("hit player hp :", player.hp)
}

func (w *World) AddBullet(bullet *Bullet) {
	if w.models == nil {
		w.models = make(map[string]Model)
	}
	w.models[bullet.ID] = bullet
	if w.eventListener != nil {
		w.eventListener.OnAddBullet(bullet)
	}
}

func (w *World) Update(deltaTime time.Duration) {
	for _, player := range w.players {
		player.Update(deltaTime)
	}

	models := make(map[string]Model)

	for _, model := range w.models {
		model.Update(deltaTime)
		if !model.IsDeleted() {
			models[model.GetID()] = model
		} else {
			if w.eventListener != nil {
				w.eventListener.OnRemoveModel(model)
			}
		}
	}
	w.models = models
}
