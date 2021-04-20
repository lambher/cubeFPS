package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

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

	for range time.Tick(time.Second) {
		_, err := c.Conn.WriteToUDP([]byte("tick\n"), c.Addr)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}
	}
}

func (c *Client) listen() {
	for {
		p := make([]byte, 2048)

		n, _, err := c.Conn.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			return
		}
		c.parse(string(p[:n]))
	}
}

func (c *Client) parse(message string) {
	messages := strings.Split(message, "\n")
	if len(messages) <= 1 {
		return
	}
	switch messages[0] {
	case "refresh":
		c.handleRefresh([]byte(messages[1]))
	}
}

func (c *Client) handleRefresh(data []byte) {
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

func (c *Client) addYou() {
	playerData, err := json.Marshal(c.Player)

	data := make([]byte, 0)

	data = append(data, []byte("you\n")...)
	data = append(data, playerData...)

	fmt.Println(string(data))
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

	fmt.Println(string(data))
	_, err = c.Conn.WriteToUDP(data, c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) sendList() {
	for _, player := range world.Players {
		c.addPlayer(player)
	}
}

func main() {
	clients = make(map[string]*Client)
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	fmt.Printf("listen on port %d\n", addr.Port)
	for {
		n, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		value := string(p[:n])
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if value == "hello" {
			player := models.NewPlayer(xid.New().String(), &world, "", *math32.NewVec3())
			clients[player.GetID()] = &Client{
				Addr:   remoteaddr,
				Conn:   ser,
				Player: player,
			}
			go clients[player.GetID()].sendResponse()
			go clients[player.GetID()].listen()
		}
	}
}
