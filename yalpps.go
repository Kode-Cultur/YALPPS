package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/anvie/port-scanner"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"time"
)

// Command line flags
var upgrader = websocket.Upgrader{} // use default options
var addr = flag.String("addr", "localhost:8080", "http service address")
var configpath = flag.String("serverconfig", "portlist.toml", "Path to your configuration file")
var runserver = flag.Bool("server", false, "Run YALPPS as server")

type game struct {
	Name string
	Port int
}

type Games struct {
	Game []game
}

func checkInboundCon(g game, con *websocket.Conn) {

	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: g.Port,
	})
	if err != nil {
		log.Printf("Failed opening port: %d: %v", g.Port, err)
		return
	}
	defer ln.Close()

	if err := ln.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
		log.Fatalln("Failed setting deadline:", err)
	}

	byte := []byte("Hello")
	err = con.WriteMessage(websocket.TextMessage, byte)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	t, err := ln.Accept()
	if err != nil {
		log.Printf("Failed to accept connection: %v", err)
	}
	if err = t.Close(); err != nil {
		log.Fatalln("Failed to close connection: %v", err)
	}
	fmt.Println("Inbound works")
}

func checkOutboundCon(g game, con *websocket.Conn) {

	byte := []byte("Hello")
	err := con.WriteMessage(websocket.TextMessage, byte)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	_, _, err = con.ReadMessage()
	if err != nil {
		log.Println("Error:", err)
	}

	ps := portscanner.NewPortScanner(con.RemoteAddr().String(), time.Duration(time.Second*10))
	ps.IsOpen(g.Port)
	fmt.Println("outbound works")
}

func (g *Games) yalpps(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer con.Close()

	for _, game := range g.Game {
		checkInboundCon(game, con)
		checkOutboundCon(game, con)
	}

}

func main() {

	// Parsing command line flags
	flag.Parse()

	if *runserver {
		var games Games
		_, err := toml.DecodeFile(*configpath, &games)
		if err != nil {
			log.Fatalln("", err)
		}

		http.HandleFunc("/yalpps", games.yalpps)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}

}
