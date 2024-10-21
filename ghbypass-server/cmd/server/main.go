package main

import (
	"flag"
	"fmt"
	"ghbypass-server/pkg/websocket"
	"log"
	"net/http"
	"sync"
)

var baseDomain string
var port int
var clients = make(map[string]*websocket.Client)
var mu sync.Mutex

func main() {
	flag.StringVar(&baseDomain, "host", "localhost", "Host domain")
	flag.IntVar(&port, "port", 80, "HTTP port")
	flag.Parse()

	log.Printf("Host set to %s:%d", baseDomain, port)

	http.HandleFunc("/ws", websocket.HandleWebSocket(clients, &mu))
	http.HandleFunc("/", websocket.HandleRequest(baseDomain, clients, &mu))

	log.Printf(`Server started on :%d`, port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
