package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// MessageBroker manages WebSocket connections for a specific type T.
type MessageBroker[T any] struct {
	Send       chan *Message
	Receive    chan *Message
	Register   chan T
	Unregister chan T
}

// IMessageBroker defines the interface for socket handlers.
type IMessageBroker[T any] interface {
	SendTo(recipient T, message *Message) error
	SendToAll(message *Message) error
	Run()
}

// SendMessageTo sends a message to a specific WebSocket connection.
func SendMessageTo(conn *websocket.Conn, message *Message) error {
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSON marshal error: %s", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonMsg)
	if err != nil {
		return fmt.Errorf("send to client error: %s", err)
	}

	return nil
}
