package main

import (
	// fmt
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

type Client struct {
	con *websocket.Conn
}

func NewClient(addr string) (*Client, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   "/yalpps",
	}
	// Running as client if desired
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return &Client{con: c}, err
}

func (client *Client) Handler() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	message := make(chan Message)

	go func() {
		defer client.con.Close()
		defer close(message)
		for {
			m := Message{}
			err := m.Read(client.con)
			if err != nil {
				log.Println("read:", err)
				return
			}
			message <- m
		}
	}()

	for {
		select {
		case m := <-message:
			client.HandleMessage(&m)
		case <-interrupt:
			log.Println("interrupt")

			err := client.con.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
		}
	}
}

func (client *Client) HandleMessage(m *Message) {
	switch m.Command {
	case SCAN_PORT:
		isOpen := ScanPort(client.con.RemoteAddr().String(), m.Port)
		fmt.Println("Host:", client.con.RemoteAddr().String(),
			"Port:", m.Port, "Is open:", isOpen,
		)

	case LISTEN_PORT:
		return

	default:
		log.Fatalln("Unkown Message:", m)
	}
}
