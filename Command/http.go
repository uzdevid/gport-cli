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
	"sync"
	"syscall"
	"time"

	"gport/Handler"
)

var (
	mu       sync.Mutex
	isClosed bool
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
	Use:   "http",
	Short: "For share local url address",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		server, _ := cmd.Flags().GetString("server")
		domain, _ := cmd.Flags().GetString("domain")

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		ws, _, err := websocket.DefaultDialer.Dial(server, http.Header{"Cli-Version": []string{Common.CliVersion}})

		if err != nil {
			fmt.Println("Error connecting to proxy server:", err)
			return
		}

		defer func(ws *websocket.Conn) {
			_ = ws.Close()
		}(ws)

		sharingMessage := Common.Sharing{}
		sharingMessage.Method = "sharing:share"
		sharingMessage.Payload.LocalAddress = address
		sharingMessage.Payload.Domain = domain

		sharingMessageMarshal, _ := json.Marshal(sharingMessage)

		done := make(chan struct{})
		writeChan := make(chan []byte, 256) // Buffered channel to hold messages

		go writePump(ws, writeChan)

		writeChan <- sharingMessageMarshal

		go func() {
			defer close(done)

			for {
				_, message, err := ws.ReadMessage()

				mu.Lock()
				if isClosed {
					mu.Unlock()
					return
				}
				mu.Unlock()

				if err != nil {
					fmt.Println("Error reading message:", err)
					return
				}

				Handler.NewMessage(ws, message, writeChan, interrupt)
			}
		}()

		for {
			select {
			case <-done:
				return
			case <-interrupt:
				close(writeChan)

				mu.Lock()
				isClosed = true

				if err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
					fmt.Println("Error closing connection:", err)
					return
				}

				_ = ws.Close()

				mu.Unlock()

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
	shareCmd.Flags().StringP("server", "s", "wss://gport.uz/wss", "Server address")

	shareCmd.Flags().StringP("address", "a", "", "Local address")
	_ = shareCmd.MarkFlagRequired("address")

	shareCmd.Flags().String("domain", "", "Remote domain (if exists)")
}
