package ws

/*
Client struct with one designated
NewClient constructor

*/

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Client struct {
	ws   *websocket.Conn
	hub  *Hub
	send chan []byte

	ID   string
	Room *Room
}

func NewClient(h *Hub, conn *websocket.Conn) *Client {
	return &Client{
		ws:   conn,
		hub:  h,
		send: make(chan []byte, 64),
	}
}

func (c *Client) Close() { _ = c.ws.Close() }

func (c *Client) Handle() {
	done := make(chan struct{})
	go c.writeLoop(done)
	c.readLoop()
	close(done)
}

func (c *Client) writeLoop(done <-chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case b, ok := <-c.send:
			if !ok {
				return
			}
			if err := websocket.Message.Send(c.ws, string(b)); err != nil {
				return
			}
		case <-ticker.C:
			_ = websocket.Message.Send(c.ws, `{"type":"ping","payload":"keepalive"}`)
		}
	}
}

func (c *Client) readLoop() {
	for {
		var raw string
		if err := websocket.Message.Receive(c.ws, &raw); err != nil {
			// on disconnect, keep saved state; just remove from live members
			if c.Room != nil {
				c.Room.remove(c)
			}
			return
		}
		fmt.Println("recv from ", c.ID)
		s := strings.TrimSpace(raw)
		if s == "" {
			continue
		}

		// Slash commands (server-authoritative)
		if strings.HasPrefix(s, "/") {
			c.handleCommand(s)
			continue
		}

		// Regular chat text -> wrap and broadcast in current room
		if c.Room == nil {
			// must join first
			c.sendJSON(Message{Type: "error", Payload: "join a room first with `/join <room> <clientID>`"})
			continue
		}

		msg := Message{
			Type:     "chat",
			Room:     c.Room.Name,
			SenderID: c.safeID(),
			Payload:  s,
		}
		c.Room.broadcast(&msg)
	}
}

func (c *Client) handleCommand(cmd string) {
	fields := strings.Fields(cmd)

	switch {
	case strings.HasPrefix(cmd, "/join "):
		// /join roomName clientID
		if len(fields) < 3 {
			c.sendJSON(Message{Type: "error", Payload: "usage: /join <room> <clientID>"})
			return
		}
		roomName := fields[1]
		clientID := fields[2]

		// leave current room (keep saved)
		if c.Room != nil {
			c.Room.remove(c)
		}

		c.ID = clientID
		c.Room = c.hub.GetOrCreateRoom(roomName)
		c.Room.add(c)

		// If snapshot exists, notify resumption
		if snap := c.Room.getSnapshot(c.ID); snap != nil {
			c.sendJSON(Message{Type: "resume", Room: c.Room.Name, SenderID: "server", Payload: snap})
		}

		sys := Message{Type: "system", Room: c.Room.Name, SenderID: "server", Payload: fmt.Sprintf("%s joined", c.safeID())}
		c.Room.broadcast(&sys)

	case strings.HasPrefix(cmd, "/leave"):
		if c.Room == nil {
			c.sendJSON(Message{Type: "error", Payload: "not in a room"})
			return
		}
		room := c.Room
		room.remove(c)
		sys := Message{Type: "system", Room: room.Name, SenderID: "server", Payload: fmt.Sprintf("%s left", c.safeID())}
		room.broadcast(&sys)
		c.Room = nil

	case strings.HasPrefix(cmd, "/save "):
		// /save {"key":"value"}  -> store arbitrary progress for this client in the room
		if c.Room == nil || c.ID == "" {
			c.sendJSON(Message{Type: "error", Payload: "join a room first (/join <room> <clientID>)"})
			return
		}
		payload := strings.TrimSpace(strings.TrimPrefix(cmd, "/save"))
		var state map[string]interface{}
		if err := json.Unmarshal([]byte(payload), &state); err != nil {
			c.sendJSON(Message{Type: "error", Payload: "invalid JSON for /save"})
			return
		}
		c.Room.saveSnapshot(c.ID, state)
		c.sendJSON(Message{Type: "system", Room: c.Room.Name, SenderID: "server", Payload: "progress saved"})

	case strings.HasPrefix(cmd, "/reset"):
		// DEV ONLY: clear saved list for current room (does not kick current connections)
		if c.Room == nil {
			c.sendJSON(Message{Type: "error", Payload: "not in a room"})
			return
		}
		c.Room.clearSaved()
		c.Room.broadcast(&Message{Type: "system", Room: c.Room.Name, SenderID: "server", Payload: "room saved list cleared"})

	default:
		c.sendJSON(Message{Type: "error", Payload: "unknown command"})
	}
}

func (c *Client) sendJSON(m Message) {
	if m.SenderID == "" {
		m.SenderID = "server"
	}
	b, _ := json.Marshal(m)
	select {
	case c.send <- b:
	default:
		c.Close()
	}
}

func (c *Client) safeID() string {
	if c.ID == "" {
		return "anon"
	}
	return c.ID
}
