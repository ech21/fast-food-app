package ws

import "sync"

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewHub() *Hub { return &Hub{rooms: make(map[string]*Room)} }

func (h *Hub) GetOrCreateRoom(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	if r, ok := h.rooms[name]; ok {
		return r
	}
	r := NewRoom(name)
	h.rooms[name] = r
	return r
}

func (h *Hub) Delete(name string) {
	h.mu.Lock()
	delete(h.rooms, name)
	h.mu.Unlock()
}
