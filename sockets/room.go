package sockets

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (r *room) String() string {
	return fmt.Sprintf(`{
	ID: %s,
	Name: %s
	IsPrivate: %t,
	Password: %s,
}`, r.ID, r.Name, r.IsPrivate, r.Password)
}

// room represents a game room with its own clients and messaging.
// A room can receive messages from a client, and distribute it to
// other clients in the same room.
// A room must be registered by a hub.
type room struct {
	*MessageHandler
	*RegistrationHandler[*client]
	hub       *hub
	ID        uuid.UUID
	Conn      *connection
	Name      string
	IsPrivate bool
	Password  string
	Clients   Clients
	// done acts as a kill switch for concurrent operations.
	done      chan any
}

// room implements the following interfaces.
// var _ IClientRegistrationHandler = (*room)(nil)
// var _ IMessageHandler[*client] = (*room)(nil)

// createRoom initializes a new Room.
func (h *hub) createRoom(name string, conn *websocket.Conn) *room {
	id := uuid.New()
	room := &room{
		MessageHandler:      createMessageHandler(),
		RegistrationHandler: createRegistrationHandler[*client](),
		hub:                 h,
		ID:                  id,
		Conn:                createConnection(conn),
		Name:                name,
		IsPrivate:           false,
		Password:            "",
		Clients:             make(map[uuid.UUID]*client),
	}

	return room
}

// makePrivate sets the room to private with the given password.
func (r *room) makePrivate(password string) {
	r.IsPrivate = true
	r.Password = password
}

// makePublic sets the room to public and clears the password.
func (r *room) makePublic() {
	r.IsPrivate = false
	r.Password = ""
}

// addClient adds a new client to the Room.
func (r *room) addClient(client *client) {
	r.Clients[client.ID] = client
}

// removeClient removes a client from the Room.
func (r *room) removeClient(clientID string) error {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return fmt.Errorf("unable to parse clientID %s in hub.RemoveClient call: %s", clientID, err)
	}

	_, exists := r.Clients[id]
	if exists {
		delete(r.Clients, id)
	}

	return nil
}

// getClient retrieves a client by its ID.
func (r *room) getClient(clientID string) (*client, error) {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse clientID %s in hub.GetRoom call: %s", clientID, err)
	}

	return r.Clients[id], nil
}

// respondTo sends a mesasge response to a specific client.
func (r *room) respondTo(sender *client, messageResponse *MessageResponse) error {
	err := r.Conn.sendMessageTo(sender.Conn, messageResponse)
	return err
}

// sendMessage sends a message to all clients in the room.
func (r *room) sendMessage(message *Message) error {
	var errs []error
	for clientID, client := range r.Clients {
		if message.SenderID == clientID {
			continue
		}

		err := r.Conn.sendMessageTo(client.Conn, message)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to send message to client %s in room.sendMessage call: %s", client, err))
			continue
		}
	}

	return errors.Join(errs...)
}

func (r *room) handleRegistration() {
	defer func() {
		close(r.Register)
		close(r.Unregister)
	}()

	for {
		select {
		case _ = <-r.done:
			return

		case client := <-r.Register:
			r.addClient(client)

		case client := <-r.Unregister:
			r.removeClient(client.ID.String())
		}
	}
}

func (r *room) handleMessage() {
	defer func() {
		close(r.Send)
		close(r.Receive)
	}()

	for {
		select {
		case _ = <-r.done:
			return

		case message := <-r.Send:
			r.sendMessage(message)

		case message := <-r.Receive:
			// On receipt of a message, forward it to all
			// other connected clients.
			client, err := r.getClient(message.SenderID.String())
			if err != nil {
				continue
			}

			if message.Respond {
				r.respondTo(client, &MessageResponse{
					Event: message.Event,
				})
			}

			r.Send <- message
		}
	}
}

// Run starts the room to listen for register, unregister, receive, and send requests.
func (r *room) Run() {
	go r.handleRegistration()
	go r.handleMessage()

	defer r.Conn.Close()
	for {
		_, rawMessage, err := r.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("Unexpected error: %s", err)
		}

		message, err := parseMessage(rawMessage)
		if err != nil {
			fmt.Printf("Parse error: %s", err)
		}

		r.Receive <- message
	}
}
