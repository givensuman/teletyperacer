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
	hub := NewHub()
	if hub.Clients == nil {
		t.Error("Clients map not initialized")
	}
	if hub.Broadcast == nil {
		t.Error("Broadcast channel not initialized")
	}
	if hub.Register == nil {
		t.Error("Register channel not initialized")
	}
	if hub.Unregister == nil {
		t.Error("Unregister channel not initialized")
	}
}

func TestRegister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		hub.Register <- conn
		go func() {
			defer func() { hub.Unregister <- conn }()
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
	if len(hub.Clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(hub.Clients))
	}
}

func TestUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		hub.Register <- conn
		go func() {
			defer func() { hub.Unregister <- conn }()
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

	if len(hub.Clients) != 0 {
		t.Errorf("Expected 0 clients, got %d", len(hub.Clients))
	}
}

func TestBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		hub.Register <- conn
		go func() {
			defer func() { hub.Unregister <- conn }()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}
				hub.Broadcast <- message
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
