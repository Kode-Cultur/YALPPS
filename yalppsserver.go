package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/anvie/port-scanner"
	"github.com/gorilla/websocket"
)

type YalppsServer struct {
	con  *websocket.Conn
	game *Game
}

// NewServer will return a new YalppsServer instance
func NewServer(con *websocket.Conn, game *Game) *YalppsServer {
	return &YalppsServer{
		con:  con,
		game: game,
	}
}

// CloseConnection will close current connection
func (server *YalppsServer) CloseConnection() {
	if err := server.con.Close(); err != nil {
		log.Fatalln("Failed to close connections:", err)
	}
}

// WriteMessage will write a message into the current connection
func (server *YalppsServer) WriteMessage(m []byte) {
	err := server.con.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (server *YalppsServer) checkInboundCon() {

	// Setting up TCP listener
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: server.game.Port,
	})
	if err != nil {
		log.Printf("TCP Error on ports: %d: %v", server.game.Port, err)
		return
	}
	defer ln.Close()

	// Setting deadline
	if err := ln.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
		log.Fatalln("Failed setting deadline:", err)
	}

	// Writing a message
	server.WriteMessage([]byte("Hello"))

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

func (server *YalppsServer) checkOutboundCon() {
	// Writing a message
	server.WriteMessage([]byte("Hello"))

	// Reading answer
	if _, _, err := server.con.ReadMessage(); err != nil {
		log.Println("Error:", err)
	}

	// Scanning ports
	server.ScanPort(server.con.RemoteAddr().String(), server.game.Port)

	fmt.Println("outbound works")
}

// Scan port will scan the given port on the given host
func (server *YalppsServer) ScanPort(host string, port int) bool {
	ps := portscanner.NewPortScanner(host, time.Duration(time.Second*10))
	return ps.IsOpen(port)
}
