package main

import (
	"flag"
	"log"
	"net/http"
)

// Command line flags
var addr = flag.String("addr", "localhost:8080", "http service address")
var configpath = flag.String("serverconfig", "portlist.toml", "Path to your config file")
var runserver = flag.Bool("server", false, "Run YALPPS as server")

func main() {
	// Parsing command line flags
	flag.Parse()

	// Running as server if desired
	if *runserver == true {
		// Decoding the config file
		config := NewConfig(*configpath)

		http.HandleFunc("/yalpps", config.Serve)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}

}
