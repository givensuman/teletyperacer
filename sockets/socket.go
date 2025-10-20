// Sockets provides logic and control for WebSocket relationships 
// between players and the server.
package sockets

import (
	"github.com/gorilla/websocket"
)

type WrapsConnection interface {
	*room | *client
	// Run goroutine handles concurrent logic for active
	// WebSocket connections.
	Run()
}

// connection wraps a *websocket.Conn.
type connection struct {
	*websocket.Conn
}

// createConnection creates a new connection.
func createConnection(conn *websocket.Conn) *connection {
	conn.SetReadLimit(512)

	return &connection{ conn }
}

// RegistrationHandler handles registration of a room or a client.
type RegistrationHandler[T WrapsConnection] struct {
	Register   chan T
	Unregister chan T
}

// createRegistrationHandler creates a new RegistrationHandler.
func createRegistrationHandler[T WrapsConnection]() *RegistrationHandler[T] {
	return &RegistrationHandler[T]{
		Register:   make(chan T),
		Unregister: make(chan T),
	}
}

// IClientRegistrationHandler defines the interface for client registrars.
type IRegistrationHandler[T WrapsConnection] interface {
	// handleRegistration goroutine handles incoming connections.
	handleRegistration()
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
type IMessageHandler interface {
	// handleMessage goroutine handles incoming messages to connection.
	handleMessage()
}
