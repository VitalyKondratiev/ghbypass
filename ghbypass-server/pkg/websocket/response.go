package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type ResponseData struct {
	Status    int                 `json:"status"`
	Headers   map[string][]string `json:"headers"`
	Body      []byte              `json:"body,omitempty"`
	RequestID string              `json:"request_id"`
}

var responses = make(map[string]chan ResponseData)
var mu sync.Mutex

func SaveResponseWriter(requestID string, w http.ResponseWriter, responseChan chan ResponseData) {
	mu.Lock()
	responses[requestID] = responseChan
	mu.Unlock()
}

func HandleResponse(message []byte) {
	var respData ResponseData
	err := json.Unmarshal(message, &respData)
	if err != nil {
		log.Println("Failed to unmarshal response data:", err)
		return
	}

	mu.Lock()
	responseChan, ok := responses[respData.RequestID]
	mu.Unlock()

	if ok {
		responseChan <- respData
		mu.Lock()
		delete(responses, respData.RequestID)
		mu.Unlock()
	}
}
