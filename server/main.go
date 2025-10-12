package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"server/lib"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func handleWebSocket(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	hub.Register <- conn // Register the new client

	// Start a goroutine to read messages from the client
	go func() {
		defer func() { hub.Unregister <- conn }()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s", message)
			hub.Broadcast <- message // Broadcast the message to all clients
		}
	}()
}

func main() {
	h := hub.NewHub()
	go h.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(h, w, r)
	})

	log.Println("WebSocket server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
