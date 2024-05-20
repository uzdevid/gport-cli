package Method

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gport/Common"
)

func SharingResponse(ws *websocket.Conn, rawMessage []byte, writeChan chan []byte) {
	var message Common.SharingResponse

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		fmt.Println("Cannot decode message:", err)
		return
	}

	fmt.Println("-----------------------------------------------------------------")

	for remote, local := range message.Payload.Proxies {
		fmt.Println(fmt.Sprintf("Proxy installed %s => %s", remote, local))
	}

	fmt.Println("-----------------------------------------------------------------")
}
