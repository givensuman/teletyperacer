// Package sockets manages WebSocket clients and broadcasts messages to them.
package sockets

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active rooms.
type Hub struct {
	Rooms map[uuid.UUID]*Room
}

// CreateHub initializes a new Hub.
func CreateHub() *Hub {
	return &Hub{
		Rooms: make(map[uuid.UUID]*Room),
	}
}

func (h *Hub) GetRoom(key string) (*Room, error) {
	id, err := uuid.Parse(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse key %s in hub.GetRoom call: %s", key, err)
	}

	return h.Rooms[id], nil
}

// CreateRoom creates a new room and starts its goroutine.
func (h *Hub) CreateRoom(name string, isPrivate bool, password string) *Room {
	id := uuid.New()
	room := &Room{
		ID:         id,
		Name:       name,
		IsPrivate:  isPrivate,
		Password:   password,
		Clients:    make(map[uuid.UUID]*Client),
		Broadcast:  make(chan *Message),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}

	h.Rooms[id] = room
	go room.Run()

	return room
}
