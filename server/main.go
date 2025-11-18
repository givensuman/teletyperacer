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

		s.On("disconnect", func() {
			log.Printf("Client %s disconnected", s.ID)
		})
	})

	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: server,
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

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}
