package main

import (
	"log"
	"net/http"

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

	// Mount websocket server on /ws/
	mux.Handle("/ws/", server)

	// Add REST endpoints
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Fatal(http.ListenAndServe(":3000", mux))
}
