package main

import (
	"github.com/google/uuid"
)

// Message represents a message sent within the room.
type Message struct {
	SenderID uuid.UUID
	Event    string
	Data     []byte
}

