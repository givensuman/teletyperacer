package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/givensuman/teletyperacer/sockets"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

var roomCounter int
var roomMutex sync.Mutex

type CreateRoomRequest struct {
	Name      string `json:"name"`
	IsPrivate bool   `json:"isPrivate"`
	Password  string `json:"password,omitempty"`
}

type CreateRoomResponse struct {
	ID string `json:"id"`
}

type IncomingMessage struct {
	Event      string      `json:"event"`
	Data       interface{} `json:"data,omitempty"`
	CallbackID *string     `json:"callbackId,omitempty"`
}

func handleCreateRoom(hub *sockets.Hub, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	roomMutex.Lock()
	roomCounter++
	roomMutex.Unlock()

	room := hub.CreateRoom(req.Name, req.IsPrivate, req.Password)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateRoomResponse{ID: room.String()})
}

func handleWebSocket(hub *sockets.Hub, w http.ResponseWriter, r *http.Request) {
	// Parse roomID from URL path, assuming /ws/{roomID}
	path := r.URL.Path
	if len(path) <= 4 || path[:4] != "/ws/" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	roomID := path[4:] // Remove "/ws/"

	room, err := hub.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if room.IsPrivate {
		password := r.URL.Query().Get("password")
		if password != room.Password {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	client := sockets.CreateClient(conn)
	room.Register <- client

	// Start a goroutine to read messages from the client
	go func() {
		defer func() {
			client := room.GetClientByConn(conn)
			if client != nil {
				room.Unregister <- client
			}
		}()
		for {
			_, rawMsg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received in room %s: %s", roomID, rawMsg)

			var incoming IncomingMessage
			if err := json.Unmarshal(rawMsg, &incoming); err != nil {
				log.Println("JSON unmarshal error:", err)
				continue
			}

			client := room.GetClientByConn(conn)
			if client == nil {
				log.Println("Client not found")
				continue
			}

			var callbackID *uuid.UUID
			if incoming.CallbackID != nil {
				if id, err := uuid.Parse(*incoming.CallbackID); err == nil {
					callbackID = &id
				}
			}

			msg := &sockets.Message{
				SenderID:   client.ID,
				Event:      incoming.Event,
				Data:       incoming.Data,
				CallbackID: callbackID,
			}

			room.Receive <- msg
		}
	}()
}

func main() {
	h := sockets.CreateHub()

	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		handleCreateRoom(h, w, r)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(h, w, r)
	})

	log.Println("WebSocket server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
