// Package hub manages WebSocket clients and broadcasts messages to them.
package hub

import (
	"log"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

// NewHub initializes a new Hub.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

// Run starts the hub to listen for register, unregister, and broadcast requests.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("Client registered, total clients:", len(h.Clients))
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Close()
				log.Println("Client unregistered, total clients:", len(h.Clients))
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Broadcast write error:", err)
					client.Close()
					delete(h.Clients, client)
				}
			}
		}
	}
}
