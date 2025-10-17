package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewHub(t *testing.T) {
	hub := CreateHub()
	if hub.Rooms == nil {
		t.Error("Rooms map not initialized")
	}
}

func TestCreateRoom(t *testing.T) {
	hub := CreateHub()
	room := hub.CreateRoom("Test Room", false, "")
	if room.Name != "Test Room" {
		t.Errorf("Expected room name 'Test Room', got %s", room.Name)
	}
	if room.IsPrivate {
		t.Error("Expected room to be public")
	}
	if _, exists := hub.Rooms[room.ID]; !exists {
		t.Error("Room not added to hub")
	}
}

func TestBroadcast(t *testing.T) {
	hub := CreateHub()
	room := hub.CreateRoom("Test Room", false, "")

	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		room.RegisterClient(conn)
		go func() {
			defer func() { room.UnregisterClient(conn) }()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}
				room.Broadcast <- &Message{
					Event: "test",
					Data:  message,
				}
			}
		}()
	}))
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http")
	conn1, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn2.Close()

	time.Sleep(10 * time.Millisecond)

	testMsg := Message{
		Event: "test",
		Data:  []byte("test message"),
	}
	message, _ := json.Marshal(testMsg)
	err = conn1.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	conn2.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, received, err := conn2.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	var receivedMsg Message
	json.Unmarshal(received, &receivedMsg)
	if receivedMsg.Event != "test" || reflect.DeepEqual(receivedMsg.Data, []byte("test message")) == false {
		t.Errorf("Expected event 'test' with data 'test message', got %+v", receivedMsg)
	}
}
