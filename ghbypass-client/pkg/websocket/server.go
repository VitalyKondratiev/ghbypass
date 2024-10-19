package websocket

import (
	"encoding/json"
	"log"

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
