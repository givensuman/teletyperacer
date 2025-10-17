package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection
// and the player behind it.
type Client struct {
	Conn *websocket.Conn
	ID   uuid.UUID
}

