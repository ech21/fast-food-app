package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type lobbySvc struct {
	lobbies []Lobby
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
		var bodyObj struct {
			Player struct {
				Name string `json:"name"`
			} `json:"player"`
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error reading request body.", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &bodyObj)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Body could not be unmarshalled.", http.StatusInternalServerError)
			return
		}
		// check fields
		if len(bodyObj.Player.Name) == 0 {
			http.Error(w, "Missing in body 'player.name'", http.StatusBadRequest)
			return
		}
		code := randomCode()
		svc.lobbies = append(svc.lobbies, Lobby{
			Id: code,
			Players: []Player{{
				Name: bodyObj.Player.Name,
			}},
		})
		fmt.Println("New lobby created with code: " + code)
		jsonOut, err := json.Marshal(struct {
			Code string `json:"code"`
		}{
			Code: code,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Data could not be marshalled into json.", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.String()+"/"+code)
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonOut)
	})

	fmt.Println("GET /lobby/{lobbyCode}")
	mux.HandleFunc("GET /lobby/{lobbyCode}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit: /lobby/{lobbyCode}")
		code := r.PathValue("lobbyCode")
		fmt.Println("Finding lobby with code: " + code)
		for _, l := range svc.lobbies {
			if l.Id == code {
				fmt.Println("Found lobby with code: " + code)
				// still need the user to give their name
				return
			}
		}
		fmt.Println("Lobby not found")
		http.NotFound(w, r)
	})
}

// purge removes all empty lobbies
func (svc *lobbySvc) purge() {
	alive := []Lobby{}
	for _, lobby := range svc.lobbies {
		if len(lobby.Players) > 0 {
			alive = append(alive, lobby)
		}
	}
	svc.lobbies = alive
}

func newLobbySvc() *lobbySvc {
	fmt.Println("New Lobby Service")

	return &lobbySvc{}
}

var chars = []rune("abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789")

func randomCode() string {
	code := make([]rune, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}
