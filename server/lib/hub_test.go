package hub

import (
	"net/http"
	"net/http/httptest"
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
	room := hub.CreateRoom("testroom", "Test Room", false, "")
	if room.ID != "testroom" {
		t.Errorf("Expected room ID 'testroom', got %s", room.ID)
	}
	if room.Name != "Test Room" {
		t.Errorf("Expected room name 'Test Room', got %s", room.Name)
	}
	if room.IsPrivate {
		t.Error("Expected room to be public")
	}
	if _, exists := hub.Rooms["testroom"]; !exists {
		t.Error("Room not added to hub")
	}
}

func TestRegister(t *testing.T) {
	hub := CreateHub()
	room := hub.CreateRoom("testroom", "Test Room", false, "")

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
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		}()
	}))
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	time.Sleep(10 * time.Millisecond)
	if len(room.Clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(room.Clients))
	}
}

func TestUnregister(t *testing.T) {
	hub := CreateHub()
	room := hub.CreateRoom("testroom", "Test Room", false, "")

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
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		}()
	}))
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)
	conn.Close()
	time.Sleep(10 * time.Millisecond)

	if len(room.Clients) != 0 {
		t.Errorf("Expected 0 clients, got %d", len(room.Clients))
	}
}

func TestBroadcast(t *testing.T) {
	hub := CreateHub()
	room := hub.CreateRoom("testroom", "Test Room", false, "")

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
				room.Broadcast <- message
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

	message := []byte("test message")
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
	if string(received) != string(message) {
		t.Errorf("Expected %s, got %s", message, received)
	}
}
