package main

import (
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/websocket"
)

type Config struct {
	list []Game
}

type Game struct {
	Name string
	Port int
}

// NewConfig will return a new Config object
func NewConfig(configpath string) *Config {
	var c Config
	_, err := toml.DecodeFile(configpath, &c)
	if err != nil {
		log.Fatalln("Failed to decode config:", err)
	}
	return &c
}

// Serve will hande Inbound and Outbound connections and 'serve' the config
func (conf *Config) Serve(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{} // use default options
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer con.Close()

	for _, game := range conf.list {
		server := NewServer(con, &game)
		server.checkInboundCon()
		server.checkOutboundCon()
	}

}
