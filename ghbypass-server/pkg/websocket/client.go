package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(clients map[string]*Client, mu *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subdomain := r.URL.Query().Get("subdomain")
		if subdomain == "" {
			http.Error(w, "Client not connected", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}

		client := &Client{Conn: conn}
		mu.Lock()
		clients[subdomain] = client
		mu.Unlock()

		fmt.Println("Client connected", subdomain)

		defer func() {
			mu.Lock()
			delete(clients, subdomain)
			mu.Unlock()
			conn.Close()
			fmt.Println("Client disconnected", subdomain)
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				break
			}
			HandleResponse(message)
		}
	}
}
