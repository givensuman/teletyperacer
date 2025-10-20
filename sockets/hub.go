package sockets

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	Rooms   map[uuid.UUID]*room
	Clients map[uuid.UUID]*client
)

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

// hub maintains the set of active rooms and unassigned clients.
type hub struct {
	// rooms hold clients.
	rooms Rooms
	// clients not yet committed to any room.
	clients Clients
	// Channels for room management
	roomRegistration *RegistrationHandler[*room]
	// Channels for client management
	clientRegistration *RegistrationHandler[*client]
}

// Hub implements the RegistrationHandler interface
// for clients and rooms.
// var _ RegistrationHandler[*room] = (*hub)(nil)
// var _ RegistrationHandler[*client) = (*hub)(nil)

// CreateHub initializes a new Hub.
func CreateHub() *hub {
	return &hub{
		rooms:              make(map[uuid.UUID]*room),
		clients:            make(map[uuid.UUID]*client),
		roomRegistration:   createRegistrationHandler[*room](),
		clientRegistration: createRegistrationHandler[*client](),
	}
}

// AddRoom adds a new room to the hub.
func (h *hub) AddRoom(room *room) {
	h.rooms[room.ID] = room
	go room.Run()
}

// RemoveRoom removes a room from the hub.
func (h *hub) RemoveRoom(roomID string) error {
	id, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("unable to parse roomID %s in hub.RemoveRoom call: %s", roomID, err)
	}

	delete(h.rooms, id)
	return nil
}

// GetRoom retrieves a room by its ID.
func (h *hub) GetRoom(roomID string) (*room, error) {
	id, err := uuid.Parse(roomID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse roomID %s in hub.GetRoom call: %s", roomID, err)
	}

	return h.rooms[id], nil
}

// AddClient adds a new client to the hub.
func (h *hub) AddClient(client *client) {
	h.clients[client.ID] = client
}

// RemoveClient removes a client from the hub.
func (h *hub) RemoveClient(clientID string) error {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return fmt.Errorf("unable to parse key %s in hub.RemoveClient call: %s", clientID, err)
	}

	delete(h.clients, id)
	return nil
}

// GetClient retrieves a client by its key.
func (h *hub) GetClient(clientID string) (*client, error) {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse key %s in hub.GetClient call: %s", clientID, err)
	}

	return h.clients[id], nil
}

// MoveClientToRoom moves a client from the hub to a specified room.
func (h *hub) MoveClientToRoom(clientID string, roomID string) error {
	client, err := h.GetClient(clientID)
	if err != nil {
		return fmt.Errorf("unable to move client %s to room %s: %s", clientID, roomID, err)
	}
	room, err := h.GetRoom(roomID)
	if err != nil {
		return fmt.Errorf("unable to move client %s to room %s: %s", clientID, roomID, err)
	}

	room.addClient(client)
	h.RemoveClient(clientID)

	return nil
}

// MoveClientOutOfRoom moves a client from a specified room back to the hub.
func (h *hub) MoveClientOutOfRoom(clientID string, roomID string) error {
	room, err := h.GetRoom(roomID)
	if err != nil {
		return fmt.Errorf("unable to move client %s out of room %s: %s", clientID, roomID, err)
	}

	client, err := room.getClient(clientID)
	if err != nil {
		return fmt.Errorf("unable to move client %s out of room %s: %s", clientID, roomID, err)
	}

	room.removeClient(clientID)
	h.AddClient(client)

	return nil
}

func (h *hub) handleRegistration() {
	defer func() {
		close(h.roomRegistration.Register)
		close(h.roomRegistration.Unregister)

		close(h.clientRegistration.Register)
		close(h.clientRegistration.Unregister)
	}()

	for {
		select {
		case room := <-h.roomRegistration.Register:
			// TODO
			// When are rooms created?
			h.AddRoom(room)

		case room := <-h.roomRegistration.Unregister:
			h.RemoveRoom(room.ID.String())

		case client := <-h.clientRegistration.Register:
			h.AddClient(client)

		case client := <-h.clientRegistration.Unregister:
			h.RemoveClient(client.ID.String())
		}
	}
}

func (h *hub) Run() {
	// TODO
	// Handle services here.
	go h.handleRegistration()
}
