package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection
// and the player behind it.
type Client struct {
	room *Room
	Conn *websocket.Conn
	ID   uuid.UUID
}

// CreateClient initializes a new Client.
func CreateClient(conn *websocket.Conn) *Client {
	return &Client{
		room: nil,
		Conn: conn,
		ID:   uuid.New(),
	}
}

func (c *Client) SetRoom(room *Room) {
	c.room = room
}

func (c *Client) String() string {
	return c.ID.String()
}
