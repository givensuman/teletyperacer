package home

import (
	"github.com/charmbracelet/bubbletea"

	"github.com/givensuman/go-sockets/client"
)

// ConnectedMsg is sent when a socket connection is established.
type ConnectedMsg struct {
	Socket *client.Socket
	RoomID string
}

// connectToRoom establishes a WebSocket connection to a room.
func connectToRoom(roomID string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		socket, err := client.Connect("ws://localhost:3000/ws/", "/", func(s *client.Socket) {
			s.Emit("join", roomID)
		})
		if err != nil {
			// For simplicity, ignore error or handle later
			return nil
		}
		return ConnectedMsg{Socket: socket, RoomID: roomID}
	})
}
