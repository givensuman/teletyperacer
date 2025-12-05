package screens

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/teletyperacer/client/internal/types"
)

type HostModel struct {
	joinCode string
	players  []string
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

func NewHost() HostModel {
	return HostModel{
		joinCode: generateJoinCode(),
		players:  []string{}, // Host is automatically a player
	}
}

func (m HostModel) Init() tea.Cmd {
	return func() tea.Msg { return types.CreateRoomMsg{} }
}

func (m HostModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return types.ScreenChangeMsg{Screen: types.HomeScreen} }
		}
	case types.RoomCreatedMsg:
		// Server confirmed room creation
		m.players = []string{"You (Host)"}
	case types.PlayerJoinedMsg:
		m.players = append(m.players, msg.PlayerName)
	}

	return m, nil
}

func (m HostModel) GetJoinCode() string {
	return m.joinCode
}

func (m HostModel) View() string {
	var content strings.Builder

	content.WriteString("🎯 Host Lobby\n\n")

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
	content.WriteString("Press ESC to go back to Home • Press Q to quit")

	return lipgloss.NewStyle().
		Padding(1).
		Align(lipgloss.Center).
		Render(content.String())
}
