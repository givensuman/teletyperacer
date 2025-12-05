package screens

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/teletyperacer/client/internal/tui"
)

type LobbyModel struct {
	joinCode string
	players  []string
	isHost   bool
}

// generateJoinCode creates a random 6-character join code
func generateJoinCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func NewLobby() LobbyModel {
	return LobbyModel{
		joinCode: "",
		players:  []string{},
		isHost:   false,
	}
}

func (m LobbyModel) Init() tea.Cmd {
	return nil
}

func (m LobbyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.HomeScreen} }
		case "c":
			if m.joinCode != "" {
				clipboard.WriteAll(m.joinCode)
			}
		}
	case tui.StartHostingMsg:
		// Start hosting a room
		m.isHost = true
		m.joinCode = generateJoinCode()
		return m, func() tea.Msg { return tui.CreateRoomMsg{} }
	case tui.AppRoomHostedMsg:
		// Room hosted successfully
	case tui.AppRoomStateUpdatedMsg:
		// Update room code and players list
		m.joinCode = msg.Code
		m.players = msg.Players
	case tui.AppPlayerJoinedMsg:
		m.players = append(m.players, msg.PlayerName)
	}

	return m, nil
}

func (m LobbyModel) GetJoinCode() string {
	return m.joinCode
}

func (m LobbyModel) View() string {
	var content strings.Builder

	if m.isHost {
		content.WriteString("🎯 Host Lobby\n\n")
	} else {
		content.WriteString("🎯 Game Lobby\n\n")
	}

	if m.joinCode != "" {
		content.WriteString(fmt.Sprintf("Join Code: %s\n\n", m.joinCode))
		content.WriteString("Share this code with friends to join!\n\n")
	} else {
		content.WriteString("Creating room...\n\n")
	}

	content.WriteString("Players:\n")
	for i, player := range m.players {
		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, player))
	}

	content.WriteString("\nWaiting for players to join...\n\n")
	if m.joinCode != "" {
		content.WriteString("Press C to copy code • ")
	}
	content.WriteString("Press ESC to go back to Home • Press Q to quit")

	return lipgloss.NewStyle().
		Padding(1).
		Align(lipgloss.Center).
		Render(content.String())
}
