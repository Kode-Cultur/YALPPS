package main

import (
	"fmt"
	"github.com/anvie/port-scanner"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"time"
)

type YalppsClient struct {
	con  *websocket.Conn
	game *Game
}

type Game struct {
	Name string
	Port int
}

func (client *YalppsClient) WriteMessage(m []byte) {
	err := client.con.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (client *YalppsClient) checkInboundCon() {

	// Setting up TCP listener
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: client.game.Port,
	})
	if err != nil {
		log.Printf("Failed opening port: %d: %v", client.game.Port, err)
		return
	}
	defer ln.Close()

	// Setting deadline
	if err := ln.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
		log.Fatalln("Failed setting deadline:", err)
	}

	// Writing a message
	mes := []byte("Hello")
	client.con.WriteMessage(websocket.TextMessage, mes)

	// Accepting the next connection to the listener
	t, err := ln.Accept()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	// And closing the connection now
	if err = t.Close(); err != nil {
		log.Fatalln("Failed to close connection: %v", err)
	}

	fmt.Println("Inbound works")
}

func (client *YalppsClient) checkOutboundCon() {

	mes := []byte("Hello")
	client.con.WriteMessage(websocket.TextMessage, mes)

	_, _, err := client.con.ReadMessage()
	if err != nil {
		log.Println("Error:", err)
	}

	ps := portscanner.NewPortScanner(client.con.RemoteAddr().String(),
		time.Duration(time.Second*10))
	ps.IsOpen(client.game.Port)
	fmt.Println("outbound works")
}

func (g *Games) yalpps(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer con.Close()

	for _, game := range g.games {
		client := &YalppsClient{
			con:  con,
			game: &game,
		}
		client.checkInboundCon()
		client.checkOutboundCon()
	}

}
