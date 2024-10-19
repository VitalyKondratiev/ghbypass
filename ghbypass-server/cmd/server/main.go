package main

import (
	"flag"
	"fmt"
	"ghbypass-server/pkg/websocket"
	"net/http"
	"sync"
)

var baseDomain string
var clients = make(map[string]*websocket.Client)
var mu sync.Mutex

func main() {
	flag.StringVar(&baseDomain, "host", "localhost", "Base host domain")
	flag.Parse()

	fmt.Println("Base domain set to", baseDomain)

	http.HandleFunc("/ws", websocket.HandleWebSocket(clients, &mu))
	http.HandleFunc("/", websocket.HandleRequest(baseDomain, clients, &mu))

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
