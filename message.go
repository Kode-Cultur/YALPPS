package main

import (
	// "fmt"
	"github.com/gorilla/websocket"
)

const (
	SCAN_PORT   = "scan_port"
	LISTEN_PORT = "listen_port"
)

type Message struct {
	Command string
	Port    int
}

func (m *Message) Read(c *websocket.Conn) error {
	return c.ReadJSON(m)
}

// WriteMessage will write a message into the current connection
func (m *Message) Write(c *websocket.Conn) error {
	return c.WriteJSON(m)
}
