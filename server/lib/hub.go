// Package hub manages WebSocket clients and broadcasts messages to them.
package hub

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
type client struct {
	*websocket.Conn
}

// Room represents a game room with its own clients and messaging.
type Room struct {
	ID         string
	Name       string
	IsPrivate  bool
	Password   string
	Clients    map[*client]bool
	Broadcast  chan []byte
	Register   chan *client
	Unregister chan *client
}

// RegisterClient registers a new client to the room.
func (r *Room) RegisterClient(conn *websocket.Conn) {
	r.Register <- &client{Conn: conn}
}

// UnregisterClient unregisters an existing Client from a room.
func (r *Room) UnregisterClient(conn *websocket.Conn) {
	r.Unregister <- &client{Conn: conn}
}

// Hub maintains the set of active rooms.
type Hub struct {
	Rooms map[string]*Room
}

// CreateHub initializes a new Hub.
func CreateHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

// CreateRoom creates a new room and starts its goroutine.
func (h *Hub) CreateRoom(id, name string, isPrivate bool, password string) *Room {
	room := &Room{
		ID:         id,
		Name:       name,
		IsPrivate:  isPrivate,
		Password:   password,
		Clients:    make(map[*client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *client),
		Unregister: make(chan *client),
	}
	h.Rooms[id] = room
	go room.Run()
	return room
}

// Run starts the room to listen for register, unregister, and broadcast requests.
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
			log.Printf("Client registered in room %s, total clients: %d", r.ID, len(r.Clients))
		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				client.Close()
				log.Printf("Client unregistered from room %s, total clients: %d", r.ID, len(r.Clients))
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Broadcast write error:", err)
					client.Close()
					delete(r.Clients, client)
				}
			}
		}
	}
}
