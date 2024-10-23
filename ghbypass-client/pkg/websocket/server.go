package websocket

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type RequestData struct {
	Method    string              `json:"method"`
	URL       string              `json:"url"`
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body,omitempty"`
	RequestID string              `json:"request_id"`
}

type ResponseData struct {
	Status    int                 `json:"status"`
	Headers   map[string][]string `json:"headers"`
	Body      []byte              `json:"body,omitempty"`
	RequestID string              `json:"request_id"`
}

func ConnectToWsServer(proxy string, subdomain string) *websocket.Conn {
	log.Println("Trying SSL (wss)...")
	wsURL := fmt.Sprintf("wss://%s/ws?subdomain=%s", proxy, subdomain)
	conn, err := tryWebSocket(wsURL, true)
	if err != nil {
		log.Println("SSL connection error, trying non-SSL (ws)...")
		wsURL = fmt.Sprintf("ws://%s/ws?subdomain=%s", proxy, subdomain)
		conn, err = tryWebSocket(wsURL, false)
		if err != nil {
			log.Fatal("WebSocket connection failed on both SSL and non-SSL:", err)
		}
	}
	log.Println("Websocket connection sucess")
	return conn
}

func HandleRequests(conn *websocket.Conn, expose string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			continue
		}

		response := forwardLocalRequest(message, expose)
		if response != nil {
			sendResponseToServer(conn, response)
		}
	}
}

func sendResponseToServer(conn *websocket.Conn, responseData *ResponseData) {
	jsonData, err := json.Marshal(*responseData)
	if err != nil {
		log.Println("Failed to marshal response data:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Failed to send response to server:", err)
	}
}

func tryWebSocket(wsURL string, useSSL bool) (*websocket.Conn, error) {
	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid WebSocket URL: %v", err)
	}

	dialer := websocket.Dialer{}
	if useSSL {
		dialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	log.Printf("Trying WebSocket connection to %s...", wsURL)

	conn, _, err := dialer.Dial(parsedURL.String(), http.Header{})
	if err != nil {
		return nil, err
	}
	return conn, nil
}
