package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/websocket"
)

// Command line flags
var upgrader = websocket.Upgrader{} // use default options
var addr = flag.String("addr", "localhost:8080", "http service address")
var configpath = flag.String("serverconfig", "portlist.toml", "Path to your config file")
var runserver = flag.Bool("server", false, "Run YALPPS as server")

type Games struct {
	Game []Game
}

func main() {

	// Parsing command line flags
	flag.Parse()

	if *runserver {
		// Decoding the config file
		var games Games
		_, err := toml.DecodeFile(*configpath, &games)
		if err != nil {
			log.Fatalln("", err)
		}

		http.HandleFunc("/yalpps", games.yalpps)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}

}
