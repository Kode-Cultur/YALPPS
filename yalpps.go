package main

import (
	"flag"
	"net/http"

	"github.com/anvie/port-scanner"
	"log"
	"time"
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

	client, err := NewClient(*addr)
	if err != nil {
		log.Fatalln("client:", err)
	}

	client.Handler()
}

func ScanPort(host string, port int) bool {
	ps := portscanner.NewPortScanner(host, time.Duration(time.Second*10))
	return ps.IsOpen(port)
}
