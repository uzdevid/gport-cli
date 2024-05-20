package Handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gport/Common"
	"gport/Handler/Method"
	"os"
)

func NewMessage(ws *websocket.Conn, rawMessage []byte, writeChan chan []byte, interrupt chan os.Signal) {
	var message Common.Message

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		fmt.Println("Cannot decode message:", err)
		return
	}

	switch message.Method {
	case "CallLocalAddress":
		go Method.CallLocalAddress(rawMessage, writeChan)
	case "SharingResponse":
		go Method.SharingResponse(rawMessage)
	case "PrintMessage":
		go Method.PrintMessage(rawMessage, interrupt)
	}
}
