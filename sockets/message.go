package sockets

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Message represents a message sent between clients.
type Message struct {
	SenderID uuid.UUID `json:"senderId"`
	Event    string    `json:"event"`
	Data     []byte    `json:"data,omitempty"`
	Respond  bool      `json:"respond"`
}

// MessageResponse represents a response sent back after mesasge handling.
type MessageResponse struct {
	Event string `json:"event"`
}

// sendMessageTo sends a message to a specific WebSocket connection.
func (c *connection) sendMessageTo(conn *connection, message any) error {
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

// parseMessage parses a WebSocket message into a typed message.
func parseMessage(data []byte) (*Message, error) {
	var message *Message
	err := json.Unmarshal(data, message)
	if err != nil {
		return nil, fmt.Errorf("error parsing message: %s", err)
	}

	return message, nil
}

// CreateMessage creates a new message for a client.
func (c *client) CreateMessage(event string, data []byte) *Message {
	return &Message{
		SenderID: c.ID,
		Event:    event,
		Data:     data,
		Respond:  false,
	}
}

// Pretty please?
func (m *Message) PleaseRespond() {
	m.Respond = true
}
