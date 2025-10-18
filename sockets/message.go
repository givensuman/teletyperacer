package main

import (
	"github.com/google/uuid"
)

// Message represents a message sent within the room.
type Message struct {
	SenderID   uuid.UUID  `json:"senderId"`
	Event      string     `json:"event"`
	Data       []byte     `json:"data,omitempty"`
	CallbackID *uuid.UUID `json:"callbackId,omitempty"`
}
