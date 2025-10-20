// Sockets provides logic and control for WebSocket relationships 
// between players and the server.
package sockets

import (
	"github.com/gorilla/websocket"
)

type CanBeRegistered interface {
	*room | *client
}

// RegistrationHandler handles registration of a room or a client.
type RegistrationHandler[T CanBeRegistered] struct {
	Register   chan T
	Unregister chan T
}

// createRegistrationHandler creates a new RegistrationHandler.
func createRegistrationHandler[T CanBeRegistered]() *RegistrationHandler[T] {
	return &RegistrationHandler[T]{
		Register:   make(chan T),
		Unregister: make(chan T),
	}
}

// IClientRegistrationHandler defines the interface for client registrars.
type IClientRegistrationHandler[T any] interface {
	AddClient(client *client)
	RemoveClient(client *client)
	Run()
}

// IRoomRegistrationHandler defines the interface for room registrars.
type IRoomRegistrationHandler interface {
	AddRoom(room *room)
	RemoveRoom(room *room)
	Run()
}

// MessageHandler handles sending and receiving messages to a room or client.
type MessageHandler struct {
	Send    chan *Message
	Receive chan *Message
}

// createMessageHandler creates a new MessageBroker.
func createMessageHandler() *MessageHandler {
	return &MessageHandler{
		Send:    make(chan *Message),
		Receive: make(chan *Message),
	}
}

// IMessageHandler defines the interface for passing messages.
type IMessageHandler[T any] interface {
	SendMessage(message *Message, args... []any)
	ReceiveMessage(message *Message, args... []any)
	Run()
}

type connection struct {
	*websocket.Conn
}

func createConnection(conn *websocket.Conn) *connection {
	return &connection{ conn }
}
