package game

import (
	"errors"
	"math/rand"
	"slices"
)

type lobby struct {
	id      string
	players []*Player
	joining []string // A list of confirmation codes for players trying to join the lobby
}

func NewLobby(initialPlayer *Player) *lobby {
	return &lobby{
		id:      genLobbyCode(),
		players: []*Player{initialPlayer},
		joining: []string{},
	}
}

func (l *lobby) Player(name string) (player *Player, err error) {
	for _, p := range l.players {
		if p.Name == name {
			return p, nil
		}
	}
	return &Player{}, errors.New("Player not found")
}

// Dead lobbies are ready for cleanup
func (l *lobby) Dead() bool {
	// if all players are dead
	for _, p := range l.players {
		if !p.dead {
			return false
		}
	}
	return true
}

// RequestJoin is called when there is a request to join the lobby. It returns a confirmation code that the client will pass to complete the join.
func (l *lobby) RequestJoin() string {
	code := genConfirmationCode()
	l.joining = append(l.joining, code)
	return code
}

// Join takes a confirmation code and checks if it is in the join list.
func (l *lobby) Join(player *Player, code string) error {
	idx := slices.Index(l.joining, code)
	if idx == -1 {
		return errors.New("Confirmation code not found")
	}
	// open ws connection
	l.players = append(l.players, player)
	slices.Delete(l.joining, idx, idx)
	return nil
}

var chars = []rune("abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789")

func genLobbyCode() string {
	code := make([]rune, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func genConfirmationCode() string {
	code := make([]rune, 16)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}
