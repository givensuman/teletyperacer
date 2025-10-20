package sockets

import (
	"log"

	"github.com/gorilla/websocket"
)

// Player is the exposed manner in which to interact with the WebSocket server.
type Player struct {
	client *client
}

// Connect establishes a socket connection.
func Connect(url string) *Player {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server")
	}
	defer conn.Close()

	client := createClient(conn)
	go client.Run()

	return &Player{client}
}

type onEvent struct {
	p     *Player
	event string
}

type onEventDone struct {
	onEvent
}

type emitEvent struct {
	p     *Player
	event string
}

type emitEventDone struct {
	emitEvent
}

// On declares what to do when a given event occurs.
func (p *Player) On(event string) onEvent {
	return onEvent{p, event}
}

// Emit sends out a message for a given event.
func (p *Player) Emit(message *Message) emitEvent {
	p.client.sendMessage(message)

	return emitEvent{p, message.Event}
}

// Do adds a callback to be run when an event occurs.
func (e onEvent) Do(callback func(message *Message)) onEventDone {
	e.p.client.addOnCallback(e.event, callback)

	return onEventDone{e}
}

// And adds a callback to be run when an event occurs.
func (d onEventDone) And(callback func(message *Message)) onEventDone {
	d.p.client.addOnCallback(d.event, callback)
	return d
}

// Then adds a callback to be run after a successful emit.
func (e emitEvent) Then(callback func(messageResponse *MessageResponse)) emitEventDone {
	e.p.client.addEmitCallback(e.event, callback)
	return emitEventDone{e}
}

// And adds a callback to be run after a successful emit.
func (d emitEventDone) And(callback func(messageResponse *MessageResponse)) emitEventDone {
	d.p.client.addEmitCallback(d.event, callback)
	return d
}
