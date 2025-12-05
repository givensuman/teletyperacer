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
	mode        LobbyMode
	joinCode    string
	playerCount int
	playerIndex int // 0-based index of current player
	lastVersion int // last received state version
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
		mode:        HostMode,
		joinCode:    generateJoinCode(),
		playerCount: 1, // Host is automatically a player
		playerIndex: 0, // Host is always P1
		lastVersion: -1,
	}
}

func NewPlayerLobby(code string) LobbyModel {
	return LobbyModel{
		mode:        PlayerMode,
		joinCode:    code,
		playerCount: 0,  // Will be updated by server
		playerIndex: -1, // Will be updated by server
		lastVersion: -1,
	}
}

func (m LobbyModel) Init() tea.Cmd {
	if m.mode == HostMode {
		return func() tea.Msg { return types.CreateRoomMsg{} }
	} else if m.mode == PlayerMode {
		// Request current room state when joining as player
		return func() tea.Msg { return types.GetRoomStateMsg{Code: m.joinCode} }
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
			m.playerCount = 1 // Host is the first player
			m.playerIndex = 0 // Host is always at index 0
		}
	case types.RoomJoinedMsg:
		// Successfully joined room as player
		if m.mode == PlayerMode {
			m.joinCode = msg.Code
			// Temporarily assume we're the only player until room state arrives
			m.playerCount = 1
			m.playerIndex = 0
		}

	case types.RoomStateMsg:
		// Update room state only if version is newer
		if msg.Version > m.lastVersion {
			m.lastVersion = msg.Version
			m.joinCode = msg.Code
			m.playerCount = msg.PlayerCount
			m.playerIndex = msg.YourIndex
			// Validate yourIndex
			if m.playerIndex < 0 || m.playerIndex >= m.playerCount {
				// Invalid, but for now, set to 0 or something
				m.playerIndex = 0
			}
		}
	}

	return m, nil
}

func (m LobbyModel) GetJoinCode() string {
	return m.joinCode
}

// ANSI colors for players
var playerColors = []lipgloss.Color{
	lipgloss.Color("1"),  // Red
	lipgloss.Color("2"),  // Green
	lipgloss.Color("4"),  // Blue
	lipgloss.Color("3"),  // Yellow
	lipgloss.Color("5"),  // Magenta
	lipgloss.Color("6"),  // Cyan
	lipgloss.Color("9"),  // Bright Red
	lipgloss.Color("10"), // Bright Green
	lipgloss.Color("12"), // Bright Blue
	lipgloss.Color("11"), // Bright Yellow
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

	// Create player grid (2 columns x 5 rows)
	playerSlots := make([]string, MaxPlayers)
	for i := 0; i < MaxPlayers; i++ {
		if i < m.playerCount {
			displayName := fmt.Sprintf("P%d", i+1)

			// Add special labels for current player and host
			if i == m.playerIndex {
				displayName += " (you)"
			} else if m.mode == HostMode && i == 0 {
				displayName += " (host)"
			}

			// Style the player slot
			playerStyle := lipgloss.NewStyle().
				Foreground(playerColors[i%len(playerColors)]).
				Background(lipgloss.Color("236")). // Dark gray background
				Padding(0, 1).
				Align(lipgloss.Center).
				Width(15)

			playerSlots[i] = playerStyle.Render(displayName)
		} else {
			// Empty slot
			emptyStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Background(lipgloss.Color("235")).
				Padding(0, 1).
				Align(lipgloss.Center).
				Width(15)

			playerSlots[i] = emptyStyle.Render("Empty")
		}
	}

	// Arrange in 2 columns
	var leftColumn, rightColumn []string
	for i := 0; i < MaxPlayers; i++ {
		if i < 5 {
			leftColumn = append(leftColumn, playerSlots[i])
		} else {
			rightColumn = append(rightColumn, playerSlots[i])
		}
	}

	gridStyle := lipgloss.NewStyle().
		Padding(0, 1)

	leftCol := gridStyle.Render(strings.Join(leftColumn, "\n"))
	rightCol := gridStyle.Render(strings.Join(rightColumn, "\n"))

	playerGrid := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	content.WriteString("Players:\n\n")
	content.WriteString(playerGrid)
	content.WriteString("\n\n")

	if m.mode == HostMode {
		if m.playerCount >= MaxPlayers {
			content.WriteString("Room is full!\n\n")
		} else {
			content.WriteString("Waiting for players to join...\n\n")
		}
	} else {
		content.WriteString("Waiting for host to start...\n\n")
	}
	content.WriteString("Press ESC to go back to Home • Press Q to quit")

	return lipgloss.NewStyle().
		Padding(1).
		Align(lipgloss.Center).
		Render(content.String())
}
