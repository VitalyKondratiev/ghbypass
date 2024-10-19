package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func forwardLocalRequest(message []byte, expose string) *ResponseData {
	var reqData RequestData
	err := json.Unmarshal(message, &reqData)
	if err != nil {
		log.Println("Failed to unmarshal request data:", err)
		return nil
	}

	localURL := fmt.Sprintf("http://%s%s", expose, reqData.URL)
	log.Println("Forwarding request to:", localURL)

	client := &http.Client{Timeout: 10 * time.Second}
	var bodyReader io.Reader
	if reqData.Body != "" {
		bodyReader = bytes.NewBuffer([]byte(reqData.Body))
	}

	req, err := http.NewRequest(reqData.Method, localURL, bodyReader)
	if err != nil {
		log.Println("Request creation error:", err)
		return nil
	}

	for key, values := range reqData.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Local request error:", err)
		return nil
	}
	defer resp.Body.Close()

	return buildResponseData(resp, reqData.RequestID)
}

func buildResponseData(resp *http.Response, requestID string) *ResponseData {
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

	return &responseData
}
