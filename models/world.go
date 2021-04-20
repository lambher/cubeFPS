package models

import (
	"time"
)

type World struct {
	Player        *Player
	Players       map[string]*Player
	models        map[string]Model
	eventListener EventListener
}

type EventListener interface {
	OnAddPlayer(player *Player)
	OnPlayerHit(player *Player)
	OnAddBullet(bullet *Bullet)
	OnRemoveModel(model Model)
}

func (w *World) SubscribeEventListener(e EventListener) {
	w.eventListener = e
}

func (w *World) AddPlayer(player *Player) {
	if w.Players == nil {
		w.Players = make(map[string]*Player)
	}
	if w.Player == nil {
		w.Player = player
	}
	w.Players[player.GetID()] = player
	if w.eventListener != nil {
		w.eventListener.OnAddPlayer(player)
	}
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

func (w *World) removeModel(model Model) {
	if w.eventListener != nil {
		w.eventListener.OnRemoveModel(model)
	}
}

func (w *World) Update(deltaTime time.Duration) {
	players := make(map[string]*Player)
	for _, player := range w.Players {
		player.Update(deltaTime)
		if !player.IsDeleted() {
			players[player.GetID()] = player
		} else {
			w.removeModel(player)
		}
	}
	w.Players = players

	models := make(map[string]Model)

	for _, model := range w.models {
		model.Update(deltaTime)
		if !model.IsDeleted() {
			models[model.GetID()] = model
		} else {
			w.removeModel(model)
		}
	}
	w.models = models
}
