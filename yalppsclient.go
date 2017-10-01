package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/anvie/port-scanner"
	"github.com/gorilla/websocket"
)

type YalppsClient struct {
	con  *websocket.Conn
	game *Game
}

// NewClient will return a new YalppsClient instance
func NewClient(con *websocket.Conn, game *Game) *YalppsClient {
	return &YalppsClient{
		con:  con,
		game: game,
	}
}

// CloseConnection will close current connection
func (client *YalppsClient) CloseConnection() {
	if err := client.con.Close(); err != nil {
		log.Fatalln("Failed to close connections:", err)
	}
}

// WriteMessage will write a message into the current connection
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
		log.Printf("TCP Error on ports: %d: %v", client.game.Port, err)
		return
	}
	defer ln.Close()

	// Setting deadline
	if err := ln.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
		log.Fatalln("Failed setting deadline:", err)
	}

	// Writing a message
	client.con.WriteMessage(websocket.TextMessage, []byte("Hello"))

	// Accepting the next connection to the listener
	t, err := ln.Accept()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	if err := t.Close(); err != nil {
		log.Fatalln("Failed to close connection")
	}

	fmt.Println("Inbound works")
}

func (client *YalppsClient) checkOutboundCon() {
	// Writing a message
	client.con.WriteMessage(websocket.TextMessage, []byte("Hello"))

	// Reading answer
	if _, _, err := client.con.ReadMessage(); err != nil {
		log.Println("Error:", err)
	}

	// Scanning ports
	client.ScanPort(client.con.RemoteAddr().String(), client.game.Port)

	fmt.Println("outbound works")
}

// Scan port will scan the given port on the given host
func (client *YalppsClient) ScanPort(host string, port int) bool {
	ps := portscanner.NewPortScanner(host, time.Duration(time.Second*10))
	return ps.IsOpen(port)
}
