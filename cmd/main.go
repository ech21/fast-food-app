package main

//server is created HERE
import (
	"fmt"
	"net/http"

	"github.com/ech21/fast-food-app/internal/ws"
	"golang.org/x/net/websocket"
)

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Pong!")
}

func main() {

	server := ws.NewServer()

	http.Handle("/", http.FileServer(http.Dir("./ui/dist")))
	http.HandleFunc("/ping", ping)
	http.Handle("/ws", websocket.Handler(server.HandleWs))

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
