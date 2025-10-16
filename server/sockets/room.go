package sockets

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
type Client struct {
	Conn *websocket.Conn
	ID   uuid.UUID
}

// Message represents a message sent within the room.
type Message struct {
	SenderID   uuid.UUID
	Event      string
	Data       []byte
}

// Room represents a game room with its own clients and messaging.
type Room struct {
	ID         uuid.UUID
	Name       string
	IsPrivate  bool
	Password   string
	Clients    map[uuid.UUID]*Client
	Broadcast  chan *Message
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

func (r *Room) String() string {
	return r.ID.String()
}

// RegisterClient registers a new client to the room.
func (r *Room) RegisterClient(conn *websocket.Conn) {
	id := uuid.New()
	client := &Client{conn, id}

	r.Register <- conn
	r.Clients[client.ID] = client
}

// UnregisterClient unregisters an existing Client from a room.
func (r *Room) UnregisterClient(conn *websocket.Conn) {
	r.Unregister <- conn

	client := r.GetClient(conn)
	delete(r.Clients, client.ID)
}

// GetClient retrieves a client by its WebSocket connection.
func (r *Room) GetClient(conn *websocket.Conn) *Client {
	for _, client := range r.Clients {
		if client.Conn == conn {
			return client
		}
	}

	return nil
}

// GetClientByID retrieves a client by its UUID.
func (r *Room) GetClientByID(id uuid.UUID) *Client {
	return r.Clients[id]
}

// SendToClient sends a message to a specific client.
func (r *Room) SendToClient(clientID uuid.UUID, msg *Message) {
	client, exists := r.Clients[clientID]
	if !exists {
		return
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	err = client.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
	if err != nil {
		log.Println("Send to client error:", err)
		return
	}
}

// SendToAll sends a message to all clients in the room.
func (r *Room) SendToAll(msg *Message) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	for _, client := range r.Clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
		if err != nil {
			log.Println("Send to all error:", err)
			return
		}
}
}

// Run starts the room to listen for register, unregister, and broadcast requests.
func (r *Room) Run() {
	for {
		select {
		case conn := <-r.Register:
			r.RegisterClient(conn)

		case conn := <-r.Unregister:
			r.UnregisterClient(conn)

		case msg := <-r.Broadcast:
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Println("JSON marshal error:", err)
				continue
			}
			for _, client := range r.Clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
				if err != nil {
					log.Println("Broadcast write error:", err)
					r.UnregisterClient(client.Conn)
					client.Conn.Close()
				}
			}
		}
	}
}
