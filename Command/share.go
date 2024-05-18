package Command

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"gport/Common"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gport/Handler"
)

func writePump(ws *websocket.Conn, writeChan chan []byte) {
	for message := range writeChan {
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:     "share",
	Aliases: []string{"sh"},
	Short:   "For share local url address",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		server, _ := cmd.Flags().GetString("server")

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		url := fmt.Sprintf("wss://%s", server)

		ws, _, err := websocket.DefaultDialer.Dial(url, http.Header{"Cli-Version": []string{Common.CliVersion}})

		if err != nil {
			fmt.Println("Error connecting to proxy server:", err)
			return
		}

		defer ws.Close()

		sharingMessage := Common.Sharing{}
		sharingMessage.Method = "sharing:share"
		sharingMessage.Payload.LocalAddress = address

		sharingMessageMarshal, _ := json.Marshal(sharingMessage)

		done := make(chan struct{})
		writeChan := make(chan []byte, 256) // Buffered channel to hold messages

		go writePump(ws, writeChan)

		writeChan <- sharingMessageMarshal

		go func() {
			defer close(done)
			for {
				_, message, err := ws.ReadMessage()

				if err != nil {
					fmt.Println("Error reading message:", err)
					return
				}

				Handler.NewMessage(ws, message, writeChan)
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
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().StringP("server", "s", "gport.uzdevid.com/wss", "Server address")

	shareCmd.Flags().StringP("address", "a", "http://localhost", "Local address")
}
