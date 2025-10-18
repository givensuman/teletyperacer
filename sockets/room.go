package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (r *Room) String() string {
	return fmt.Sprintf(`{
	ID: %s,
	Name: %s
	IsPrivate: %t,
	Password: %s,
}`, r.ID, r.Name, r.IsPrivate, r.Password)
}

// Room represents a game room with its own clients and messaging.
type Room struct {
	MessageBroker[*Client]
	hub           *Hub
	ID            uuid.UUID
	Name          string
	IsPrivate     bool
	Password      string
	Clients       map[uuid.UUID]*Client
	ClientsByConn map[*websocket.Conn]*Client
	eventHandlers map[string][]func(*Client, any)
}

// Room implements the IMessageBroker interface for Clients.
// var _ IMessageBroker[*Client] = (*Room)(nil)

// CreateRoom initializes a new Room.
func (h *Hub) CreateRoom(name string) *Room {
	id := uuid.New()
	return &Room{
		hub:           h,
		ID:            id,
		Name:          name,
		IsPrivate:     false,
		Password:      "",
		Clients:       make(map[uuid.UUID]*Client),
		ClientsByConn: make(map[*websocket.Conn]*Client),
		eventHandlers: make(map[string][]func(*Client, any)),
	}
}

// MakePrivate sets the room to private with the given password.
func (r *Room) MakePrivate(password string) {
	r.IsPrivate = true
	r.Password = password
}

// MakePublic sets the room to public and clears the password.
func (r *Room) MakePublic() {
	r.IsPrivate = false
	r.Password = ""
}

// AddClient adds a new client to the Room.
func (r *Room) AddClient(client *Client) {
	client.SetRoom(r)
	r.Clients[client.ID] = client
}

// RemoveClient removes a client from the Room.
func (r *Room) RemoveClient(clientID string) error {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return fmt.Errorf("unable to parse clientID %s in hub.RemoveClient call: %s", clientID, err)
	}

	delete(r.Clients, id)
	return nil
}

// GetClient retrieves a client by its ID.
func (r *Room) GetClient(clientID string) (*Client, error) {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse clientID %s in hub.GetRoom call: %s", clientID, err)
	}

	return r.Clients[id], nil
}

// GetClientByConn retrieves a client by its connection.
func (r *Room) GetClientByConn(conn *websocket.Conn) *Client {
	return r.ClientsByConn[conn]
}

// On registers an event handler.
func (r *Room) On(event string, handler func(*Client, any)) {
	r.eventHandlers[event] = append(r.eventHandlers[event], handler)
}

// SendTo sends a message to a specific client.
func (r *Room) SendTo(recipient *Client, message *Message) error {
	err := SendMessageTo(recipient.Conn, message)
	if err != nil {
		return err
	}

	return nil
}

// SendToAll sends a message to all clients in the room.
func (r *Room) SendToAll(message *Message) error {
	var err error
	for _, client := range r.Clients {
		err = SendMessageTo(client.Conn, message)
		if err != nil {
			fmt.Printf("unable to send message to client %s in room.SendToAll call: %s", client, err)
			continue
		}
	}

	return err
}

// Run starts the room to listen for register, unregister, receive, and send requests.
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.AddClient(client)

		case client := <-r.Unregister:
			r.RemoveClient(client.ID.String())

		case message := <-r.Receive:
			client, err:= r.GetClient(message.SenderID.String())
			if err != nil {
				continue
			}

			if handlers, exists := r.eventHandlers[message.Event]; exists {
				for _, handler := range handlers {
					handler(client, message.Data)
				}
			}

			// Handle callback if present
			if message.CallbackID != nil {
				ackMsg := &Message{
					Event:      message.Event + "_ack",
					Data:       []byte("Success"),
					CallbackID: message.CallbackID,
				}

				r.Send <- ackMsg
			}

		case message := <-r.Send:
			r.SendToAll(message)
		}
	}
}
