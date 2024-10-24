package websocket

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"ghbypass-server/internal/utils"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//go:embed "content/free_subdomain.html"
var htmlContent string

type RequestData struct {
	Method    string              `json:"method"`
	URL       string              `json:"url"`
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body,omitempty"`
	RequestID string              `json:"request_id"`
}

func HandleRequest(baseDomain string, clients map[string]*Client, mu *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hostWithoutPort := strings.Split(r.Host, ":")[0]
		if !strings.HasSuffix(hostWithoutPort, baseDomain) {
			http.Error(w, "Invalid domain", http.StatusForbidden)
			return
		}

		subdomain := utils.GetSubdomain(r.Host, baseDomain)
		mu.Lock()
		client, ok := clients[subdomain]
		mu.Unlock()

		if !ok {
			if strings.HasPrefix(r.URL.Path, "/download/") {
				handleDownload(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, htmlContent)
			return
		}
		requestID := uuid.New().String()
		reqData, err := serializeRequest(r, requestID)
		if err != nil {
			http.Error(w, "Failed to serialize request", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(reqData)
		if err != nil {
			http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
			return
		}

		responseChan := make(chan ResponseData)
		SaveResponseWriter(requestID, w, responseChan)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp := <-responseChan
			for key, values := range resp.Headers {
				for _, value := range values {
					w.Header().Set(key, value)
				}
			}
			w.WriteHeader(resp.Status)
			w.Write(resp.Body)
		}()

		client.Conn.WriteMessage(websocket.TextMessage, jsonData)
		wg.Wait()
	}
}

func serializeRequest(r *http.Request, requestID string) (*RequestData, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	requestData := &RequestData{
		Method:    r.Method,
		URL:       r.URL.String(),
		Headers:   r.Header,
		Body:      string(bodyBytes),
		RequestID: requestID,
	}

	return requestData, nil
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)                  // Извлекаем имя файла из URL
	filepath := fmt.Sprintf("/usr/local/bin/%s", filename) // Путь к бинарному файлу

	http.ServeFile(w, r, filepath) // Отправляем файл клиенту
}
