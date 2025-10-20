package sockets

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (c *client) String() string {
	return c.ID.String()
}

// client represents a single WebSocket connection.
// A client can send and receive messages to/from a room.
// A client must be registered by either a room or a hub, but not both.
type client struct {
	MessageHandler
	room          *room
	ID            uuid.UUID
	Conn          *connection
	onCallbacks   map[string]func(message *Message)
	emitCallbacks map[string]func(messageResponse *MessageResponse)
}

// createClient initializes a new Client.
func createClient(conn *websocket.Conn) *client {
	return &client{
		room:          nil,
		ID:            uuid.New(),
		Conn:          createConnection(conn),
		onCallbacks:   make(map[string]func(*Message)),
		emitCallbacks: make(map[string]func(*MessageResponse)),
	}
}

// joinRoom joins a room.
func (c *client) joinRoom(room *room) {
	c.room = room
	c.room.Register <- c
}

// leaveRoom leaves a room.
func (c *client) leaveRoom(room *room) {
	c.room = nil
	c.room.Unregister <- c
}

// addOnCallback adds a callback for receiving a given event.
func (c *client) addOnCallback(event string, callback func(message *Message)) {
	fn, exists := c.onCallbacks[event]
	if !exists {
		c.onCallbacks[event] = callback
		return
	}

	c.onCallbacks[event] = func(message *Message) {
		fn(message)
		callback(message)
	}
}

func (c *client) getOnCallback(event string) func(message *Message) {
	return c.onCallbacks[event]
}

// addEmitCallback adds a callback for successful emits.
func (c *client) addEmitCallback(event string, callback func(messageResponse *MessageResponse)) {
	fn, exists := c.emitCallbacks[event]
	if !exists {
		c.emitCallbacks[event] = callback
		return
	}

	c.emitCallbacks[event] = func(messageResponse *MessageResponse) {
		fn(messageResponse)
		callback(messageResponse)
	}
}

func (c *client) getEmitCallback(event string) func(messageResponse *MessageResponse) {
	return c.emitCallbacks[event]
}

// sendMessage constructs and sends a message to the client's room.
func (c *client) sendMessage(message *Message) {
	c.Conn.sendMessageTo(c.room.Conn, message)
}

// receiveMessage handles callbacks for the message event.
func (c *client) receiveMessage(message *Message) *Message {
	callback := c.getOnCallback(message.Event)
	if callback != nil {
		callback(message)
	}

	return message
}

// receiveMessageResponse handles callbacks for successful emits.
func (c *client) receiveMessageResponse(messageResponse *MessageResponse) *MessageResponse {
	callback := c.getEmitCallback(messageResponse.Event)
	if callback != nil {
		callback(messageResponse)
	}

	return nil
}

func (c *client) Run() {
	for {
		select {
		case message := <-c.Receive:
			c.receiveMessage(message)

		case message := <-c.Send:
			c.sendMessage(message)
		}

		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error in reading message", err)
			continue
		}

		message, err := parseMessage(msg)
		if err != nil {
			log.Println("Error in parsing message", err)
			continue
		}

		c.Receive <- message
	}
}
