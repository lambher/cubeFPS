package models

import "time"

type World struct {
	Player  *Player
	players map[string]*Player
}

func (w *World) AddPlayer(player *Player) {
	if w.players == nil {
		w.players = make(map[string]*Player)
	}
	if w.Player == nil {
		w.Player = player
	}
	w.players[player.Name] = player
}

func (w *World) Update(deltaTime time.Duration) {
	for _, player := range w.players {
		player.Update(deltaTime)
	}
}
