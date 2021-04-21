package models

import (
	"encoding/json"
	"sync"
	"time"
)

type World struct {
	Player        *Player
	players       map[string]*Player
	models        map[string]Model
	eventListener EventListener

	playersLock sync.RWMutex
	playerLock  sync.RWMutex
}

type EventListener interface {
	OnAddPlayer(player *Player)
	OnPlayerHit(player *Player)
	OnAddBullet(bullet *Bullet)
	OnRemoveModel(model Model)
}

func (w *World) GetPlayerData() ([]byte, error) {
	w.playerLock.Lock()
	playerData, err := json.Marshal(w.Player)
	w.playerLock.Unlock()

	return playerData, err
}

func (w *World) GetPlayer(id string) *Player {
	var player *Player

	w.playersLock.Lock()
	if p, ok := w.players[id]; ok {
		player = p
	}
	w.playersLock.Unlock()

	return player
}

func (w *World) GetPlayers() []*Player {
	players := make([]*Player, 0)

	w.playersLock.Lock()
	for _, player := range w.players {
		players = append(players, player)
	}
	w.playersLock.Unlock()

	return players
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
	w.players[player.GetID()] = player
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
	for _, player := range w.players {
		player.Update(deltaTime)
		if !player.IsDeleted() {
			players[player.GetID()] = player
		} else {
			w.removeModel(player)
		}
	}
	w.players = players

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
