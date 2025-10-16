package views

import (
	"encoding/json"
	"log"
	"net/url"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Socket represents a WebSocket connection with event handling.
type Socket struct {
	conn          *websocket.Conn
	eventHandlers map[string][]func(interface{})
	callbacks     map[uuid.UUID]func(interface{})
	mu            sync.RWMutex
}

// Message represents a message sent over the socket.
type Message struct {
	Event      string    `json:"event"`
	Data       []byte    `json:"data,omitempty"`
	CallbackID *uuid.UUID `json:"callbackId,omitempty"`
}

// Connect establishes a WebSocket connection to the given room.
func Connect(roomID string) (*Socket, error) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws/" + roomID}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	socket := &Socket{
		conn:          conn,
		eventHandlers: make(map[string][]func(interface{})),
		callbacks:     make(map[uuid.UUID]func(interface{})),
	}

	go socket.listen()

	return socket, nil
}

// Emit sends an event with data to the server.
func (s *Socket) Emit(event string, data []byte) {
	msg := Message{
		Event: event,
		Data:  data,
	}
	s.send(msg)
}

// EmitWithCallback sends an event with data and a callback for response.
func (s *Socket) EmitWithCallback(event string, data []byte, callback func(interface{})) {
	id := uuid.New()
	msg := Message{
		Event:      event,
		Data:       data,
		CallbackID: &id,
	}
	s.mu.Lock()
	s.callbacks[id] = callback
	s.mu.Unlock()
	s.send(msg)
}

// On registers an event handler.
func (s *Socket) On(event string, handler func(interface{})) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventHandlers[event] = append(s.eventHandlers[event], handler)
}

// send marshals and sends a message.
func (s *Socket) send(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}
	err = s.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("Write error:", err)
	}
}

// listen reads incoming messages and dispatches them.
func (s *Socket) listen() {
	defer s.conn.Close()
	for {
		_, data, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		if msg.CallbackID != nil {
			s.mu.RLock()
			callback, exists := s.callbacks[*msg.CallbackID]
			s.mu.RUnlock()
			if exists {
				callback(msg.Data)
				s.mu.Lock()
				delete(s.callbacks, *msg.CallbackID)
				s.mu.Unlock()
			}
		} else {
			s.mu.RLock()
			handlers, exists := s.eventHandlers[msg.Event]
			s.mu.RUnlock()
			if exists {
				for _, handler := range handlers {
					handler(msg.Data)
				}
			}
		}
	}
}
