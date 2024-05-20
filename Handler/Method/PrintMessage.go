package Method

import (
	"encoding/json"
	"fmt"
	"gport/Common"
	"os"
)

func PrintMessage(rawMessage []byte, interrupt chan os.Signal) {
	var message Common.PrintMessage

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		fmt.Println("Cannot decode message:", err)
		return
	}

	fmt.Println(message.Payload.Message)

	if message.Payload.Exit {
		close(interrupt)
	}
}
