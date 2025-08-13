package ws

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type Server struct {
	connections map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWs(ws *websocket.Conn) {
	fmt.Println("New connection from client:", ws.RemoteAddr())
	s.connections[ws] = true
	s.readLoop(ws)

}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Read error:", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
		ws.Write([]byte("Message received."))
	}
}
