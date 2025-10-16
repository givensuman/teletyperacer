package sockets

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)


func TestRegister(t *testing.T) {
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
