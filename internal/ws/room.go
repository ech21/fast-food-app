package ws

import (
	"encoding/json"
	"sync"
	"time"
)

type Room struct {
	Name    string
	mu      sync.RWMutex
	members map[*Client]bool
	saved   map[string]*ClientSnapshot
}

func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		members: make(map[*Client]bool),
		saved:   make(map[string]*ClientSnapshot),
	}
}

func (r *Room) add(c *Client) {
	r.mu.Lock()
	r.members[c] = true
	r.mu.Unlock()
}

func (r *Room) remove(c *Client) {
	r.mu.Lock()
	delete(r.members, c)
	r.mu.Unlock()
}

func (r *Room) saveSnapshot(clientID string, state map[string]interface{}) {
	r.mu.Lock()
	r.saved[clientID] = &ClientSnapshot{
		ClientID:  clientID,
		UpdatedAt: time.Now(),
		State:     state,
	}
	r.mu.Unlock()
}

func (r *Room) getSnapshot(clientID string) *ClientSnapshot {
	r.mu.RLock()
	s := r.saved[clientID]
	r.mu.RUnlock()
	return s
}

func (r *Room) clearSaved() {
	r.mu.Lock()
	r.saved = make(map[string]*ClientSnapshot)
	r.mu.Unlock()
}

func (r *Room) broadcast(msg *Message) {
	data, _ := json.Marshal(msg)
	r.mu.RLock()
	for c := range r.members {
		select {
		case c.send <- data:
		default:
			go c.Close()
		}
	}
	r.mu.RUnlock()
}
