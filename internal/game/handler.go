package game

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var lobbies []*lobby = []*lobby{}

func handleLobbyCreate(w http.ResponseWriter, r *http.Request) {
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
	initialPlayer, err := NewPlayer(bodyObj.Player.Name)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newLobby := NewLobby(initialPlayer)
	lobbies = append(lobbies, newLobby)
	fmt.Println("New lobby created with code: " + newLobby.id)
	type output struct {
		Code string `json:"code"`
	}
	jsonOut, err := json.Marshal(output{
		Code: newLobby.id,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Data could not be marshalled into json.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.URL.String()+"/"+newLobby.id)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOut)
}

func handleLobbyRequestJoin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit: /lobby/{lobbyCode}")
	code := r.PathValue("lobbyCode")
	fmt.Println("Finding lobby with code: " + code)
	for _, l := range lobbies {
		if l.id == code {
			fmt.Println("Found lobby with code: " + code)
			// still need the user to give their name
			confCode := l.RequestJoin()
			type output struct {
				Confirm string `json:"confirm"`
			}
			jsonOut, err := json.Marshal(output{
				Confirm: confCode,
			})
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Data could not be marshalled into json.", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonOut)
			return
		}
	}
	fmt.Println("Lobby does not exist")
	http.NotFound(w, r)
}

func openWs(w http.ResponseWriter, r *http.Request) {
	// should expect a lobby code and a confirmation code
	// completes ws connection if both are valid
}

func Attach(mux *http.ServeMux) {
	fmt.Println("POST /lobby")
	mux.HandleFunc("POST /lobby", handleLobbyCreate)

	fmt.Println("GET /lobby/{lobbyCode}")
	mux.HandleFunc("GET /lobby/{lobbyCode}", handleLobbyRequestJoin)
}
