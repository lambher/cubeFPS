package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/g3n/engine/math32"
	"github.com/rs/xid"

	"github.com/lambher/video-game/models"
)

var world models.World

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, player *models.Player) {
	addYou(conn, addr, player)
	sendList(conn, addr)

	world.AddPlayer(player)

	for range time.Tick(time.Second) {
		_, err := conn.WriteToUDP([]byte("tick\n"), addr)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}
	}
}

func addYou(conn *net.UDPConn, addr *net.UDPAddr, player *models.Player) {
	playerData, err := json.Marshal(player)

	data := make([]byte, 0)

	data = append(data, []byte("you\n")...)
	data = append(data, playerData...)

	fmt.Println(string(data))
	_, err = conn.WriteToUDP(data, addr)
	if err != nil {
		fmt.Println(err)
	}
}

func sendList(conn *net.UDPConn, addr *net.UDPAddr) {
	for _, player := range world.Players {
		playerData, err := json.Marshal(player)

		data := make([]byte, 0)

		data = append(data, []byte("add_player\n")...)
		data = append(data, playerData...)

		fmt.Println(string(data))
		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
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
			go sendResponse(ser, remoteaddr, player)
		}
	}
}
