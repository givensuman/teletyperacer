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

const MaxPlayers = 10

type LobbyMode int

const (
	HostMode LobbyMode = iota
	PlayerMode
)

type LobbyModel struct {
	mode     LobbyMode
	joinCode string
	players  []string
}

func generateJoinCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func NewHostLobby() LobbyModel {
	return LobbyModel{
		mode:     HostMode,
		joinCode: generateJoinCode(),
		players:  []string{}, // Host is automatically a player
	}
}

func NewPlayerLobby(code string) LobbyModel {
	return LobbyModel{
		mode:     PlayerMode,
		joinCode: code,
		players:  []string{}, // Will be populated when room state is received
	}
}

func (m LobbyModel) Init() tea.Cmd {
	if m.mode == HostMode {
		return func() tea.Msg { return types.CreateRoomMsg{} }
	}
	return nil
}

func (m LobbyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return types.ScreenChangeMsg{Screen: types.HomeScreen} }
		case "c":
			if m.joinCode != "" {
				// Try to copy to clipboard
				return m, func() tea.Msg {
					return types.CopyCodeMsg{Code: m.joinCode}
				}
			}
		}
	case types.RoomCreatedMsg:
		// Server confirmed room creation
		if m.mode == HostMode {
			m.players = []string{"You (Host)"}
		}
	case types.RoomJoinedMsg:
		// Successfully joined room as player
		if m.mode == PlayerMode {
			m.joinCode = msg.Code
			// Player name will be set when room state is received
		}
	case types.PlayerJoinedMsg:
		m.players = append(m.players, msg.PlayerName)
	case types.RoomStateMsg:
		// Update room state
		m.joinCode = msg.Code
		m.players = msg.Players
	}

	return m, nil
}

func (m LobbyModel) GetJoinCode() string {
	return m.joinCode
}

func (m LobbyModel) View() string {
	var content strings.Builder

	if m.mode == HostMode {
		content.WriteString("🎯 Host Lobby\n\n")
	} else {
		content.WriteString("🎯 Player Lobby\n\n")
	}

	if m.joinCode != "" {
		content.WriteString(fmt.Sprintf("Join Code: %s\n", m.joinCode))
		content.WriteString("(press 'c' to copy to clipboard)\n\n")
		if m.mode == HostMode {
			content.WriteString("Share this code with friends to join!\n\n")
		}
	} else {
		if m.mode == HostMode {
			content.WriteString("Creating room...\n\n")
		} else {
			content.WriteString("Joining room...\n\n")
		}
	}

	content.WriteString(fmt.Sprintf("Players (%d/%d):\n", len(m.players), MaxPlayers))
	for i, player := range m.players {
		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, player))
	}

	if m.mode == HostMode {
		if len(m.players) >= MaxPlayers {
			content.WriteString("\nRoom is full!\n\n")
		} else {
			content.WriteString("\nWaiting for players to join...\n\n")
		}
	} else {
		content.WriteString("\nWaiting for host to start...\n\n")
	}
	content.WriteString("Press ESC to go back to Home • Press Q to quit")

	return lipgloss.NewStyle().
		Padding(1).
		Align(lipgloss.Center).
		Render(content.String())
}
