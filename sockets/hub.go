package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Rooms map[uuid.UUID]*Room
type Clients map[uuid.UUID]*Client

func (r *Rooms) String() string {
	var str string
	for _, room := range *r {
		str += fmt.Sprintf("%s,\n", room)
	}

	return str
}

func (c *Clients) String() string {
	var str string
	for _, client := range *c {
		str += fmt.Sprintf("%s,\n", client.ID)
	}

	return str
}

// Hub maintains the set of active rooms.
type Hub struct {
	MessageBroker[*Room]
	// Rooms hold Clients.
	Rooms Rooms
	// Clients not yet committed to any room.
	Clients Clients
}

// Hub implements the IMessageBroker interface for Rooms.
// var _ IMessageBroker[*Room] = (*Hub)(nil)

// CreateHub initializes a new Hub.
func CreateHub() *Hub {
	return &Hub{
		Rooms: make(map[uuid.UUID]*Room),
		Clients: make(map[uuid.UUID]*Client),
	}
}

// AddRoom adds a new room to the hub.
func (h *Hub) AddRoom(room *Room) {
	h.Rooms[room.ID] = room
	go room.Run()
}

// RemoveRoom removes a room from the hub.
func (h *Hub) RemoveRoom(roomID string) error {
	id, err := uuid.Parse(roomID); if err != nil {
		return fmt.Errorf("unable to parse roomID %s in hub.RemoveRoom call: %s", roomID, err)
	}

	delete(h.Rooms, id)
	return nil
}

// GetRoom retrieves a room by its ID.
func (h *Hub) GetRoom(roomID string) (*Room, error) {
	id, err := uuid.Parse(roomID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse roomID %s in hub.GetRoom call: %s", roomID, err)
	}

	return h.Rooms[id], nil
}

// AddClient adds a new client to the hub.
func (h *Hub) AddClient(client *Client) {
	h.Clients[client.ID] = client
}

// RemoveClient removes a client from the hub.
func (h *Hub) RemoveClient(clientID string) error {
	id, err := uuid.Parse(clientID); if err != nil {
		return fmt.Errorf("unable to parse key %s in hub.RemoveClient call: %s", clientID, err)
	}

	delete(h.Clients, id)
	return nil
}

// GetClient retrieves a client by its key.
func (h *Hub) GetClient(clientID string) (*Client, error) {
	id, err := uuid.Parse(clientID); if err != nil {
		return nil, fmt.Errorf("unable to parse key %s in hub.GetClient call: %s", clientID, err)
	}

	return h.Clients[id], nil
}

// MoveClientToRoom moves a client from the hub to a specified room.
func (h *Hub) MoveClientToRoom(clientID string, roomID string) error {
	client, err := h.GetClient(clientID); if err != nil {
		return fmt.Errorf("unable to move client %s to room %s: %s", clientID, roomID, err)
	}
	room, err := h.GetRoom(roomID); if err != nil {
		return fmt.Errorf("unable to move client %s to room %s: %s", clientID, roomID, err)
	}

	room.AddClient(client)
	h.RemoveClient(clientID)

	return nil
}

// MoveClientOutOfRoom moves a client from a specified room back to the hub.
func (h *Hub) MoveClientOutOfRoom(clientID string, roomID string) error {
	client, err := h.GetClient(clientID); if err != nil {
		return fmt.Errorf("unable to move client %s out of room %s: %s", clientID, roomID, err)
	}
	room, err := h.GetRoom(roomID); if err != nil {
		return fmt.Errorf("unable to move client %s out of room %s: %s", clientID, roomID, err)
	}

	room.RemoveClient(clientID)
	h.AddClient(client)

	return nil
}

// SendTo sends a message to all clients in a specified room.
func (h *Hub) SendTo(recipient *Room, message *Message) error {
	err := recipient.SendToAll(message)
	if err != nil {
		return err
	}

	return nil
}

// SendToAll sends a message to all clients.
func (h *Hub) SendToAll(message *Message) error {
	for _, room := range h.Rooms {
		err := h.SendTo(room, message); if err != nil {
			return err
		}
	}

	for _, client := range h.Clients {
		err := SendMessageTo(client.Conn, message); if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) Run() {
	// TODO
}
