// Package views contains the logic for rendering the various
// screens rendered for the TUI.
package views

import (
	"github.com/charmbracelet/bubbletea"

	"github.com/givensuman/go-sockets/client"
)

// Model represents application state for the TUI.
type Model struct {
	CurrentView View
	Width       int
	Height      int
	RoomID      string
	Socket      *client.Socket
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	switch m.CurrentView {
	case Home:
		return HomeView(m)
	case Practice:
		return "Practice View\nPress 'h' for Home, 'q' to Quit"
	case Host:
		return "Host View\nPress 'h' for Home, 'j' to Join, 'q' to Quit"
	case Join:
		if m.Socket != nil {
			return "Join View - Connected to room " + m.RoomID + "\nPress 'h' for Home, 'l' for Lobby, 'q' to Quit"
		}
		return "Join View - Connecting...\nPress 'h' for Home, 'q' to Quit"
	case Lobby:
		return "Lobby View - Room: " + m.RoomID + "\nPress 'h' for Home, 'y' to Play, 'q' to Quit"
	case Play:
		return "Play View\nPress 'h' for Home, 'l' for Lobby, 'q' to Quit"
	default:
		return "Unknown View\nPress 'h' for Home, 'q' to Quit"
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case ConnectedMsg:
		m.Socket = msg.Socket
		m.RoomID = msg.RoomID
		m.Socket.On("game_started", func(args ...any) {
			// Handle game start
		})
		m.CurrentView = Lobby
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			m.CurrentView = Home
		case "p":
			m.CurrentView = Practice
		case "o":
			m.CurrentView = Host
		case "j":
			m.CurrentView = Join
			// For demo, connect to a test room
			return m, connectToRoom("test-room-id")
		case "l":
			m.CurrentView = Lobby
		case "y":
			m.CurrentView = Play
			if m.Socket != nil {
				m.Socket.Emit("start_game", "hello")
			}
		}
	}
	return m, nil
}
