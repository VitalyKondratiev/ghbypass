package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var subdomain string
var proxy string
var expose string

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

func main() {
	parseFlags()

	conn := connectToWebSocket()
	defer conn.Close()

	handleRequests(conn)
}

// Парсинг флагов командной строки
func parseFlags() {
	flag.StringVar(&subdomain, "subdomain", "", "Subdomain on proxy host")
	flag.StringVar(&proxy, "proxy", "", "Proxy address")
	flag.StringVar(&expose, "expose", "localhost:8081", "URL on local machine to expose")
	flag.Parse()

	if subdomain == "" || proxy == "" {
		log.Fatal("subdomain and proxy are required")
	}

	log.Println("Subdomain set to:", subdomain)
	log.Println("Local address for expose:", expose)
}

// Подключение к WebSocket серверу
func connectToWebSocket() *websocket.Conn {
	wsURL := fmt.Sprintf("ws://%s/ws?subdomain=%s", proxy, subdomain)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WebSocket dial error:", err)
	}
	return conn
}

// Обработка сообщений от сервера
func handleRequests(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			continue
		}

		var reqData RequestData
		err = json.Unmarshal(message, &reqData)
		if err != nil {
			log.Println("Failed to unmarshal request data:", err)
			continue
		}

		response := handleLocalRequest(reqData)
		sendResponseToServer(conn, response)
	}
}

// Выполнение запроса к локальному серверу
func handleLocalRequest(reqData RequestData) ResponseData {
	localURL := fmt.Sprintf("http://%s%s", expose, reqData.URL)
	fmt.Println("Forwarding request to:", localURL)

	client := &http.Client{Timeout: 10 * time.Second}
	var bodyReader io.Reader
	if reqData.Body != "" {
		bodyReader = bytes.NewBuffer([]byte(reqData.Body))
	}

	req, err := http.NewRequest(reqData.Method, localURL, bodyReader)
	if err != nil {
		log.Println("Request creation error:", err)
		return ResponseData{}
	}

	// Добавляем заголовки
	for key, values := range reqData.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Local request error:", err)
		return ResponseData{}
	}
	defer resp.Body.Close()

	return buildResponseData(resp, reqData.RequestID)
}

func buildResponseData(resp *http.Response, requestID string) ResponseData {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Local request read error:", err)
	}

	responseData := ResponseData{
		Status:    resp.StatusCode,
		Headers:   resp.Header,
		Body:      body,
		RequestID: requestID,
	}

	return responseData
}

// Отправка ответа обратно на сервер через WebSocket
func sendResponseToServer(conn *websocket.Conn, responseData ResponseData) {
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		log.Println("Failed to marshal response data:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Failed to send response to server:", err)
	}
}
