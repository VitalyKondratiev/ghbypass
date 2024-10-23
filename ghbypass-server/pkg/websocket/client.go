package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
			fmt.Println("Websocket upgrade error:", err)
			return
		}

		client := &Client{Conn: conn}
		mu.Lock()
		_, hasClient := clients[subdomain]
		if !hasClient {
			clients[subdomain] = client
		}
		mu.Unlock()

		if hasClient {
			closeConnection(conn, "Subdomain is taken")
			log.Println("Client trying use already taken subdomain:", subdomain)
			return
		}

		log.Println("Client connected:", subdomain)

		defer func() {
			mu.Lock()
			delete(clients, subdomain)
			mu.Unlock()
			closeConnection(conn, "Client disconnected")
			log.Println("Client disconnected:", subdomain)
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			HandleResponse(message)
		}
	}
}

func closeConnection(conn *websocket.Conn, cause string) {
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, cause))
	if err != nil {
		log.Println("Error sending close message:", err)
	}

	time.Sleep(1 * time.Second)

	err = conn.Close()
	if err != nil {
		log.Println("Error closing websocket connection:", err)
	}
}
