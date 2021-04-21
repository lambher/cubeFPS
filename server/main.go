package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/lambher/video-game/conf"

	"github.com/g3n/engine/math32"
	"github.com/rs/xid"

	"github.com/lambher/video-game/models"
)

var world models.World

type Client struct {
	Addr   *net.UDPAddr
	Conn   *net.UDPConn
	Player *models.Player
}

var clients map[string]*Client

func (c *Client) sendResponse() {
	c.addYou()
	c.sendList()
	c.populatePlayer()

	world.AddPlayer(c.Player)
}

func (c *Client) parse(message string) {
	messages := strings.Split(message, "\n")
	if len(messages) <= 1 {
		return
	}
	switch messages[0] {
	case "refresh_player":
		c.handleRefreshPlayer([]byte(messages[1]))
	case "move":
		c.handleMove([]byte(messages[1]))
	}
}

func (c *Client) handleMove(data []byte) {
	var move models.Moves

	err := json.Unmarshal(data, &move)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Player.RefreshMoves(move)
}

func (c *Client) handleRefreshPlayer(data []byte) {
	var player models.Player

	err := json.Unmarshal(data, &player)
	if err != nil {
		fmt.Println(err)
		return
	}
	if player.Position == nil {
		fmt.Println("player position is null")
		return
	}

	c.Player.Refresh(player)
}

func (c *Client) populatePlayer() {
	for _, client := range clients {
		if client.Player.ID != c.Player.ID {
			client.addPlayer(c.Player)
		}
	}
}

func (c *Client) refreshPlayer(player *models.Player) {
	playerData, err := json.Marshal(player)

	data := make([]byte, 0)

	data = append(data, []byte("refresh_player\n")...)
	data = append(data, playerData...)

	_, err = c.Conn.WriteToUDP(data, c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) addYou() {
	playerData, err := json.Marshal(c.Player)

	data := make([]byte, 0)

	data = append(data, []byte("you\n")...)
	data = append(data, playerData...)

	_, err = c.Conn.WriteToUDP(data, c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) addPlayer(player *models.Player) {
	playerData, err := json.Marshal(player)

	data := make([]byte, 0)

	data = append(data, []byte("add_player\n")...)
	data = append(data, playerData...)

	_, err = c.Conn.WriteToUDP(data, c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) sendList() {
	for _, player := range world.GetPlayers() {
		c.addPlayer(player)
	}
}

func main() {
	clients = make(map[string]*Client)
	addr := net.UDPAddr{
		Port: conf.Port,
		IP:   net.ParseIP(conf.Host),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	fmt.Printf("listen on port %d\n", addr.Port)

	go gameLoop()
	go tick()

	for {
		p := make([]byte, 2048)

		n, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		value := string(p[:n])
		fmt.Printf("Read a message from %s %s \n", remoteaddr.String(), p)
		if value == "hello" {
			player := models.NewPlayer(xid.New().String(), &world, "", *math32.NewVec3())
			clients[remoteaddr.String()] = &Client{
				Addr:   remoteaddr,
				Conn:   ser,
				Player: player,
			}
			go clients[remoteaddr.String()].sendResponse()
			//go clients[remoteaddr.String()].listen()
		} else {
			if client, ok := clients[remoteaddr.String()]; ok {
				client.parse(value)
			}
		}
	}
}

func gameLoop() {
	tick := time.Tick(16 * time.Millisecond)
	t := time.Now()
	for {
		select {
		case <-tick:
			deltaTime := time.Since(t)
			t = time.Now()
			world.Update(deltaTime)
		}
	}
}

func tick() {
	for range time.Tick(conf.TickTimeServer) {
		refreshPlayers()
	}
}

func refreshPlayers() {
	for _, client := range clients {
		for _, player := range world.GetPlayers() {
			client.refreshPlayer(player)
		}
	}
}
