package Message

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

const (
	cliVersion = "0.0.0"
)

// Callback represents the structure of the callback message
type Callback struct {
	Method  string `json:"method"`
	Payload struct {
		RequestId string `json:"requestId"`
		Content   struct {
			Type   string `json:"type"`
			Base64 string `json:"base64"`
		} `json:"content"`
	} `json:"payload"`
}

// MessagePayload represents the structure of the incoming message payload
type MessagePayload struct {
	RequestId string `json:"requestId"`
	Address   string `json:"address"`
}

// IncomingMessage represents the structure of the incoming WebSocket message
type IncomingMessage struct {
	Payload MessagePayload `json:"payload"`
}

func compressData(data string) (string, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err := writer.Write([]byte(data))
	if err != nil {
		return "", err
	}
	writer.Close()
	compressedData := buffer.Bytes()
	return base64.StdEncoding.EncodeToString(compressedData), nil
}

func handleMessage(ws *websocket.Conn, message []byte, writeChan chan []byte) {
	var incomingMessage IncomingMessage

	if err := json.Unmarshal(message, &incomingMessage); err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return
	}

	address := incomingMessage.Payload.Address
	response, err := http.Get(address)

	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	compressedContent, err := compressData(string(body))

	if err != nil {
		fmt.Println("Error compressing data:", err)
		return
	}

	callback := Callback{
		Method: "app:local-content",
	}

	callback.Payload.RequestId = incomingMessage.Payload.RequestId
	callback.Payload.Content.Type = response.Header.Get("Content-Type")
	callback.Payload.Content.Base64 = compressedContent

	callbackMessage, err := json.Marshal(callback)

	if err != nil {
		fmt.Println("Error marshalling callback message:", err)
		return
	}

	writeChan <- callbackMessage
}

func writePump(ws *websocket.Conn, writeChan chan []byte) {
	for message := range writeChan {
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/", http.Header{"Cli-Version": []string{cliVersion}})

	if err != nil {
		fmt.Println("Error connecting to WebSocket:", err)
		return
	}

	defer ws.Close()

	done := make(chan struct{})
	writeChan := make(chan []byte, 256) // Buffered channel to hold messages

	go writePump(ws, writeChan)

	go func() {
		defer close(done)
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			go handleMessage(ws, message, writeChan)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("Interrupt received, closing connection...")
			close(writeChan) // Close the write channel to terminate the writePump
			if err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				fmt.Println("Error closing connection:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
