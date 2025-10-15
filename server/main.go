package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"server/lib"
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

func handleCreateRoom(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
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
	id := fmt.Sprintf("room%d", roomCounter)
	roomMutex.Unlock()

	room := hub.CreateRoom(id, req.Name, req.IsPrivate, req.Password)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateRoomResponse{ID: room.ID})
}

func handleWebSocket(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	// Parse roomID from URL path, assuming /ws/{roomID}
	path := r.URL.Path
	if len(path) <= 4 || path[:4] != "/ws/" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	roomID := path[4:] // Remove "/ws/"

	room, exists := hub.Rooms[roomID]
	if !exists {
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

	room.RegisterClient(conn)

	// Start a goroutine to read messages from the client
	go func() {
		defer func() { room.UnregisterClient(conn) }()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received in room %s: %s", roomID, message)
			room.Broadcast <- message // Broadcast the message to room clients
		}
	}()
}

func main() {
	h := hub.CreateHub()

	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		handleCreateRoom(h, w, r)
	})

	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(h, w, r)
	})

	log.Println("WebSocket server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
