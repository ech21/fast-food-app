package ws

import "time"

type Message struct {
	Type     string      `json:"type"`
	Room     string      `json:"room"`
	SenderID string      `json:"senderID"`
	Payload  interface{} `json:"payload"`
}

type ClientSnapshot struct {
	ClientID  string
	UpdatedAt time.Time
	State     map[string]interface{}
}
