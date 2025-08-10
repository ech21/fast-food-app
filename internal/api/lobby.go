package api

type Settings struct {
	// in miles, the maximum distance to search for restaurants
	MaxRadius  float32
	MaxPlayers int
}

type Player struct {
	Id   string
	Name string
}

type Lobby struct {
	Id      string
	Players []Player
}

type LobbyService interface {
	// New creates a new lobby.
	New() *Lobby
	// Join takes a lobby id and a Player object and tries to add the player to
	// the lobby, or have the player rejoin if the player was already in the lobby.
	// It returns the lobby joined and any error while trying to join.
	Join(id string, player *Player) (lobby *Lobby, err error)
	// Close shuts down a lobby and removes all players. It may return an error.
	Close(lobby *Lobby) error
}
