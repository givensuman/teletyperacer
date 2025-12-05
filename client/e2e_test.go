package main

import (
	"encoding/json"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestE2ERoomJoining tests end-to-end room creation and joining between two clients
func TestE2ERoomJoining(t *testing.T) {
	// Assume server is running on localhost:3000
	serverURL := "ws://localhost:3000"
	serverPath := "/ws/"

	// First, check if server is available
	testConn, _, err := websocket.DefaultDialer.Dial(serverURL+serverPath, nil)
	if err != nil {
		t.Logf("Server not available at %s%s: %v. Will attempt test anyway.", serverURL, serverPath, err)
		// Don't skip, try the test
	} else {
		testConn.Close()
	}

	// Channel to collect received messages
	type receivedMsg struct {
		clientID int
		event    string
		data     interface{}
	}
	msgChan := make(chan receivedMsg, 10)

	var wg sync.WaitGroup
	wg.Add(2)

	// Client 1: Create room
	go func() {
		defer wg.Done()
		log.Printf("Client 1 connecting to %s%s", serverURL, serverPath)
		socket, err := sockets.Connect(serverURL, serverPath, nil)
		if err != nil {
			t.Errorf("Client 1 failed to connect: %v", err)
			return
		}
		defer socket.Close()
		log.Printf("Client 1 connected")

		// Listen for messages
		socket.On("connected", func(data interface{}) {
			log.Printf("Client 1 received connected: %v", data)
			msgChan <- receivedMsg{clientID: 1, event: "connected", data: data}
		})
		socket.On("roomCreated", func(data interface{}) {
			log.Printf("Client 1 received roomCreated: %v", data)
			msgChan <- receivedMsg{clientID: 1, event: "roomCreated", data: data}
		})
		socket.On("playerJoined", func(data interface{}) {
			log.Printf("Client 1 received playerJoined: %v", data)
			msgChan <- receivedMsg{clientID: 1, event: "playerJoined", data: data}
		})

		// Create room with code "TEST123"
		log.Printf("Client 1 emitting createRoom TEST123")
		socket.Emit("createRoom", "TEST123")

		// Wait a bit for messages
		time.Sleep(2 * time.Second)
		log.Printf("Client 1 done")
	}()

	// Client 2: Join room
	go func() {
		defer wg.Done()
		// Wait a moment for client 1 to create room
		time.Sleep(500 * time.Millisecond)

		log.Printf("Client 2 connecting to %s%s", serverURL, serverPath)
		socket, err := sockets.Connect(serverURL, serverPath, nil)
		if err != nil {
			t.Errorf("Client 2 failed to connect: %v", err)
			return
		}
		defer socket.Close()
		log.Printf("Client 2 connected")

		// Listen for messages
		socket.On("connected", func(data interface{}) {
			log.Printf("Client 2 received connected: %v", data)
			msgChan <- receivedMsg{clientID: 2, event: "connected", data: data}
		})
		socket.On("roomJoined", func(data interface{}) {
			log.Printf("Client 2 received roomJoined: %v", data)
			msgChan <- receivedMsg{clientID: 2, event: "roomJoined", data: data}
		})

		// Join room "TEST123"
		log.Printf("Client 2 emitting joinRoom TEST123")
		socket.Emit("joinRoom", "TEST123")

		// Wait a bit for messages
		time.Sleep(2 * time.Second)
		log.Printf("Client 2 done")
	}()

	// Wait for both clients to finish
	wg.Wait()
	close(msgChan)

	// Collect received messages
	var messages []receivedMsg
	for msg := range msgChan {
		messages = append(messages, msg)
	}

	// Verify messages
	client1ReceivedConnected := false
	client1ReceivedRoomCreated := false
	client1ReceivedPlayerJoined := false
	client2ReceivedConnected := false
	client2ReceivedRoomJoined := false

	for _, msg := range messages {
		log.Printf("Client %d received: %s", msg.clientID, msg.event)
		switch msg.clientID {
		case 1:
			if msg.event == "connected" {
				client1ReceivedConnected = true
			} else if msg.event == "roomCreated" {
				client1ReceivedRoomCreated = true
			} else if msg.event == "playerJoined" {
				client1ReceivedPlayerJoined = true
			}
		case 2:
			if msg.event == "connected" {
				client2ReceivedConnected = true
			} else if msg.event == "roomJoined" {
				client2ReceivedRoomJoined = true
			}
		}
	}

	// Assertions
	if !client1ReceivedConnected {
		t.Error("Client 1 did not receive connected message")
	}
	if !client2ReceivedConnected {
		t.Error("Client 2 did not receive connected message")
	}
	if !client1ReceivedRoomCreated {
		t.Error("Client 1 did not receive roomCreated message")
	}
	if !client2ReceivedRoomJoined {
		t.Error("Client 2 did not receive roomJoined message")
	}
	if !client1ReceivedPlayerJoined {
		t.Error("Client 1 did not receive playerJoined message when Client 2 joined")
	}

	if client1ReceivedRoomCreated && client2ReceivedRoomJoined && client1ReceivedPlayerJoined {
		t.Log("E2E test passed: Room creation and joining works correctly")
	}
}

// TestE2EPingPong tests simple message sending between two clients in the same room
func TestE2EPingPong(t *testing.T) {
	// For this test, we'd need the server to support general message broadcasting
	// Since the current server only handles room events, this is a placeholder
	// In a real implementation, add message sending to the server

	t.Skip("Ping/pong test requires message broadcasting feature in server")
}
