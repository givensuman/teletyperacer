package handlers

import (
	"log"

	sockets "github.com/givensuman/go-sockets/server"
	"github.com/givensuman/teletyperacer/server/types"
)

// HandleConnections sets up WebSocket event handlers
func HandleConnections(io *sockets.Namespace) {
	log.Printf("🔧 Setting up WebSocket event handlers for namespace")

	io.On("connection", func(s *sockets.Socket) {
		log.Printf("🔌 New WebSocket connection established - Client ID: %s", s.ID)

		// Handle room creation
		s.On("createRoom", func(code string) {
			log.Printf("🏠 Client %s attempting to create room with code %s", s.ID, code)

			// Join the room (this creates it if it doesn't exist)
			s.Join(code)
			log.Printf("✅ Room %s created successfully by client %s", code, s.ID)

			// Send room created confirmation
			s.Emit("roomCreated", types.RoomCreatedResponse{Code: code})
			log.Printf("📤 Sent roomCreated confirmation to client %s for room %s", s.ID, code)

			// Send initial room state (just the host for now)
			roomState := types.RoomStateResponse{
				Code:    code,
				Players: []string{"You (Host)"},
			}
			s.Emit("roomState", roomState)
			log.Printf("📤 Sent initial roomState to client %s for room %s: %d players", s.ID, code, len(roomState.Players))
		})

		// Handle room joining
		s.On("joinRoom", func(code string) {
			log.Printf("🚪 Client %s attempting to join room %s", s.ID, code)

			// Join the room
			s.Join(code)
			log.Printf("✅ Client %s successfully joined room %s", s.ID, code)

			// Send join confirmation
			s.Emit("roomJoined", types.RoomJoinedResponse{Code: code})
			log.Printf("📤 Sent roomJoined confirmation to client %s for room %s", s.ID, code)

			// Send room state with current players
			roomState := types.RoomStateResponse{
				Code:    code,
				Players: []string{"Host", "You"},
			}
			s.Emit("roomState", roomState)
			log.Printf("📤 Sent roomState to client %s for room %s: %d players", s.ID, code, len(roomState.Players))

			// Broadcast to others in the room
			s.Broadcast().To(code).Emit("playerJoined", "You")
			log.Printf("📢 Broadcasted playerJoined event to room %s (excluding sender %s)", code, s.ID)
		})

		// Handle errors
		s.On("error", func(err interface{}) {
			log.Printf("❌ WebSocket error for client %s: %v", s.ID, err)
		})

		s.On("disconnect", func() {
			log.Printf("🔌 WebSocket connection closed - Client ID: %s disconnected", s.ID)
		})
	})
}
