package api

import (
	"fmt"
	"net/http"
)

type lobbySvc struct {
}

func (svc *lobbySvc) new() *Lobby {
	return nil
}

func (svc *lobbySvc) join(id string, player *Player) (lobby *Lobby, err error) {
	return nil, nil
}

func (svc *lobbySvc) close(lobby *Lobby) error {
	return nil
}

func (svc *lobbySvc) attach(mux *http.ServeMux) {
	fmt.Println("POST /lobby")
	mux.HandleFunc("POST /lobby", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit: /lobby")
		fmt.Fprint(w, "Hello World")
	})
	fmt.Println("GET /lobby/{lobbyCode}")
	mux.HandleFunc("GET /lobby/{lobbyCode}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit: /lobby/{lobbyCode}")
		code := r.PathValue("lobbyCode")
		fmt.Println(code)
		fmt.Fprint(w, "Hello World")
	})
}

func newLobbySvc() *lobbySvc {
	fmt.Println("New Lobby Service")

	return &lobbySvc{}
}
