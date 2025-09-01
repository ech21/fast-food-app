package main

//server is created HERE
import (
	"fmt"
	"github.com/ech21/fast-food-app/internal/api"
	"github.com/ech21/fast-food-app/internal/ws"
	"golang.org/x/net/websocket"
	"net/http"
)

func main() {
	entireServer := ws.NewHub()
	wsServer := websocket.Server{
		Handshake: func(cfg *websocket.Config, r *http.Request) error {
			return nil
		},

		Handler: func(conn *websocket.Conn) {
			fmt.Println("Websocket connected: ", conn.RemoteAddr())
			client := ws.NewClient(entireServer, conn)
			client.Handle()
			fmt.Println("Websocket disconnected: ", conn.RemoteAddr())

		},
	}

	http.Handle("/", http.FileServer(http.Dir("./ui/dist")))
	http.Handle("/ws", wsServer)

	api.AttachHandlers(http.DefaultServeMux)

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
