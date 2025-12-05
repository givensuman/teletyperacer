package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/givensuman/teletyperacer/server/handlers"
)

func main() {
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	// Mount websocket handler on /ws/
	mux.HandleFunc("/ws/", handlers.HandleWebSocket)

	// Add REST endpoints
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Gracefully shutting down...")
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
