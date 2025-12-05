package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sockets "github.com/givensuman/go-sockets/server"
)

func main() {
	server := sockets.NewServer()
	io := server.Of("/")

	io.On("connection", func(s *sockets.Socket) {
		log.Printf("Client %s connected", s.ID)

		// Handle room creation
		s.On("createRoom", func(code string) {
			log.Printf("Client %s creating room with code %s", s.ID, code)

			// Join the room (this creates it if it doesn't exist)
			s.Join(code)

			// Send room created confirmation
			s.Emit("roomCreated", code)

			// Send initial room state (just the host for now)
			s.Emit("roomState", map[string]interface{}{
				"code":    code,
				"players": []string{"You (Host)"},
			})
		})

		// Handle room joining
		s.On("joinRoom", func(code string) {
			log.Printf("Client %s joining room %s", s.ID, code)

			// Join the room
			s.Join(code)

			// Send join confirmation
			s.Emit("roomJoined", code)

			// For now, just send a simple room state
			// In a real implementation, you'd track room members
			s.Emit("roomState", map[string]interface{}{
				"code":    code,
				"players": []string{"Host", fmt.Sprintf("Player %s", s.ID[:6])},
			})

			// Broadcast to others in the room
			s.Broadcast().To(code).Emit("playerJoined", fmt.Sprintf("Player %s", s.ID[:6]))
		})

		s.On("disconnect", func() {
			log.Printf("Client %s disconnected", s.ID)
		})
	})

	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	// Mount websocket server on /ws/
	mux.Handle("/ws/", server)

	// Add REST endpoints
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Gracefully shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP server did not close gracefully: %v", err)
		}

		os.Exit(0)
	}()

	log.Println("Server starting on :3000")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}
