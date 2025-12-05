// Package root describes the root of
// the TUI application
package root

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	zone "github.com/lrstanley/bubblezone"

	"github.com/givensuman/teletyperacer/client/internal/tui/components/input"
	"github.com/givensuman/teletyperacer/client/internal/tui/screens"
	"github.com/givensuman/teletyperacer/client/internal/types"
)

// WebSocket message types
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type CreateRoomData struct {
	Code string `json:"code"`
}

type JoinRoomData struct {
	Code string `json:"code"`
}

type RoomStateData struct {
	Code        string `json:"code"`
	PlayerCount int    `json:"playerCount"`
	YourIndex   int    `json:"yourIndex"`
}

type PlayerJoinedData struct {
	PlayerIndex int `json:"playerIndex"`
}

type Model struct {
	// Currently rendered screen
	screen types.Screen
	// Child models
	home,
	lobby,
	practice tea.Model
	// Join screen
	join tea.Model
	// WebSocket connection
	conn    *websocket.Conn
	spinner spinner.Model
	// Window dimensions
	width, height int
	// Connection status
	connectionStatus types.ConnectionStatus
	// WebSocket message channel
	wsChan chan tea.Msg
}

type backgroundModel struct {
	root *Model
}

func (b backgroundModel) Init() tea.Cmd { return nil }

func (b backgroundModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return b, nil }

func (b backgroundModel) View() string {
	var content string
	switch b.root.screen {
	case types.HomeScreen:
		content = b.root.home.View()
	case types.LobbyScreen:
		content = b.root.lobby.View()
	case types.PracticeScreen:
		content = b.root.practice.View()
	case types.JoinScreen:
		content = b.root.join.View()
	default:
		content = b.root.home.View()
	}

	if b.root.connectionStatus == types.Connecting {
		spinnerView := lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(b.root.width).
			Height(b.root.height).
			Render("Connecting to server...\n" + b.root.spinner.View())
		return zone.Scan(lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(b.root.width).
			Height(b.root.height).
			Render(content + "\n\n" + spinnerView))
	}

	return zone.Scan(lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Width(b.root.width).
		Height(b.root.height).
		Render(content))
}

// categorizeConnectionError determines the type of connection error
func categorizeConnectionError(err error) types.ConnectionStatus {
	if err == nil {
		return types.Connected
	}

	errStr := err.Error()

	// Check for client-side configuration errors
	if strings.Contains(errStr, "invalid URL") ||
		strings.Contains(errStr, "unsupported protocol") ||
		strings.Contains(errStr, "malformed") {
		return types.ClientError
	}

	// Check for network-related errors that indicate server unreachable
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			return types.ServerUnreachable
		}
		// Check for specific network error types
		if opErr, ok := netErr.(*net.OpError); ok {
			if opErr.Op == "dial" {
				return types.ServerUnreachable
			}
		}
	}

	// Check for common server unreachable indicators
	if strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "connection timed out") {
		return types.ServerUnreachable
	}

	// Default to client error for any other issues
	return types.ClientError
}

// sendWSMessage sends a message to the WebSocket server
func (m Model) sendWSMessage(msgType string, data interface{}) {
	if m.conn == nil {
		return
	}

	msg := WSMessage{Type: msgType, Data: data}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return
	}
	m.conn.WriteMessage(websocket.TextMessage, jsonData)
}

// copyToClipboard attempts to copy text to system clipboard
func (m Model) copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		if err := clipboard.WriteAll(text); err != nil {
			return types.ClipboardErrorMsg{Message: "Failed to copy to clipboard"}
		}
		return types.ClipboardSuccessMsg{Message: "Code copied to clipboard!"}
	}
}

// handleWSMessageFromEvent processes incoming WebSocket messages from specific events
func (m Model) handleWSMessageFromEvent(eventType string, data interface{}) tea.Msg {
	switch eventType {
	case "roomCreated":
		var roomData CreateRoomData
		if d, ok := data.(map[string]interface{}); ok {
			if code, ok := d["code"].(string); ok {
				roomData.Code = code
			}
		} else if d, ok := data.(CreateRoomData); ok {
			roomData = d
		}
		return types.RoomCreatedMsg{Code: roomData.Code}

	case "roomJoined":
		var roomData JoinRoomData
		if d, ok := data.(map[string]interface{}); ok {
			if code, ok := d["code"].(string); ok {
				roomData.Code = code
			}
		} else if d, ok := data.(JoinRoomData); ok {
			roomData = d
		}
		return types.RoomJoinedMsg{Code: roomData.Code}

	case "playerJoined":
		var playerData PlayerJoinedData
		if d, ok := data.(map[string]interface{}); ok {
			if index, ok := d["playerIndex"].(float64); ok {
				playerData.PlayerIndex = int(index)
			}
		} else if d, ok := data.(PlayerJoinedData); ok {
			playerData = d
		}
		return types.PlayerJoinedMsg{PlayerIndex: playerData.PlayerIndex}

	case "roomState":
		var stateData RoomStateData
		if d, ok := data.(map[string]interface{}); ok {
			if code, ok := d["code"].(string); ok {
				stateData.Code = code
			}
			if playerCount, ok := d["playerCount"].(float64); ok {
				stateData.PlayerCount = int(playerCount)
			}
			if yourIndex, ok := d["yourIndex"].(float64); ok {
				stateData.YourIndex = int(yourIndex)
			}
		} else if d, ok := data.(RoomStateData); ok {
			stateData = d
		}
		return types.RoomStateMsg{Code: stateData.Code, PlayerCount: stateData.PlayerCount, YourIndex: stateData.YourIndex}

	case "error":
		// Handle specific error types
		if d, ok := data.(map[string]interface{}); ok {
			if reason, ok := d["reason"].(string); ok {
				return types.RoomJoinFailedMsg{Reason: reason}
			}
		}
		return nil
	}

	return nil
}

// handleWSMessage processes incoming WebSocket messages (legacy, kept for compatibility)
func (m Model) handleWSMessage(data interface{}) tea.Msg {
	var wsMsg WSMessage
	var jsonData []byte

	// Handle different data types from WebSocket
	switch d := data.(type) {
	case string:
		jsonData = []byte(d)
	case []byte:
		jsonData = d
	default:
		// Try to marshal the interface{} to JSON
		if marshaled, err := json.Marshal(d); err == nil {
			jsonData = marshaled
		} else {
			return nil
		}
	}

	if err := json.Unmarshal(jsonData, &wsMsg); err != nil {
		return nil
	}

	return m.handleWSMessageFromEvent(wsMsg.Type, wsMsg.Data)
}

// waitForWSMessage waits for WebSocket messages
func (m Model) waitForWSMessage() tea.Cmd {
	return func() tea.Msg {
		msg := <-m.wsChan
		return msg
	}
}

func New() Model {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/ws/", nil)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	connectionStatus := types.Connected
	if err != nil {
		connectionStatus = categorizeConnectionError(err)
	}

	return Model{
		screen:           types.HomeScreen,
		home:             screens.NewHome(),
		lobby:            screens.NewHostLobby(),
		practice:         screens.NewPractice(),
		join:             screens.NewJoin(),
		conn:             conn,
		spinner:          s,
		width:            80,
		height:           24,
		connectionStatus: connectionStatus,
		wsChan:           make(chan tea.Msg, 10),
	}
}

func (m Model) Init() tea.Cmd {
	// Start WebSocket message reader
	if m.conn != nil {
		go func() {
			for {
				_, data, err := m.conn.ReadMessage()
				if err != nil {
					// Connection closed or error
					return
				}

				var wsMsg WSMessage
				if err := json.Unmarshal(data, &wsMsg); err != nil {
					continue
				}

				if msg := m.handleWSMessageFromEvent(wsMsg.Type, wsMsg.Data); msg != nil {
					select {
					case m.wsChan <- msg:
					default:
						// Channel full, drop message
					}
				}
			}
		}()
	}

	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			return types.ConnectionStatusMsg{Status: m.connectionStatus}
		},
		m.waitForWSMessage(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Forward window size to current screen
		return m.updateCurrentScreen(msg)

	case types.ScreenChangeMsg:
		m.screen = msg.Screen
		if msg.Screen == types.JoinScreen {
			m.join = screens.NewJoin()
			return m, m.join.Init()
		}
		if msg.Screen == types.LobbyScreen {
			if m.screen == types.HomeScreen {
				m.lobby = screens.NewHostLobby()
			}
			return m, m.lobby.Init()
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case types.ConnectionStatusMsg:
		m.connectionStatus = msg.Status
		// Forward connection status to current screen
		return m.updateCurrentScreen(msg)

	case types.CreateRoomMsg:
		// Get the join code from the lobby screen
		if lobbyModel, ok := m.lobby.(screens.LobbyModel); ok {
			m.sendWSMessage("createRoom", CreateRoomData{Code: lobbyModel.GetJoinCode()})
		}
		return m, nil

	case types.JoinRoomMsg:
		m.sendWSMessage("joinRoom", JoinRoomData{Code: msg.Code})
		return m, nil

	case types.GetRoomStateMsg:
		m.sendWSMessage("getRoomState", map[string]string{"code": msg.Code})
		return m, nil

	case input.SubmitMsg:
		// Send join room message to server
		return m, func() tea.Msg {
			return types.JoinRoomMsg{Code: strings.ToUpper(msg.Value)}
		}

	case input.HideMsg:
		return m, func() tea.Msg { return types.ScreenChangeMsg{Screen: types.HomeScreen} }

	case types.CopyCodeMsg:
		// Try to copy code to clipboard using common commands
		cmd := m.copyToClipboard(msg.Code)
		return m, cmd

	case types.ClipboardSuccessMsg:
		// Could show a brief success notification here
		return m, nil

	case types.ClipboardErrorMsg:
		// Could show an error notification here
		return m, nil

	case types.RoomJoinedMsg:
		// Successfully joined room, switch to player lobby
		m.lobby = screens.NewPlayerLobby(msg.Code)
		return m, tea.Batch(func() tea.Msg { return types.ScreenChangeMsg{Screen: types.LobbyScreen} }, func() tea.Msg { return types.GetRoomStateMsg{} })

	case types.RoomJoinFailedMsg:
		// Forward to current screen to handle the error
		return m.updateCurrentScreen(msg)

	default:
		// Forward all other messages to current screen
		return m.updateCurrentScreen(msg)
	}
}

func (m Model) updateCurrentScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.screen {
	case types.HomeScreen:
		m.home, cmd = m.home.Update(msg)
	case types.LobbyScreen:
		m.lobby, cmd = m.lobby.Update(msg)
	case types.PracticeScreen:
		m.practice, cmd = m.practice.Update(msg)
	case types.JoinScreen:
		m.join, cmd = m.join.Update(msg)
	default:
		cmd = nil
	}

	// Always continue waiting for WebSocket messages
	return m, tea.Batch(cmd, m.waitForWSMessage())
}

func (m Model) View() string {
	var content string
	switch m.screen {
	case types.HomeScreen:
		content = m.home.View()
	case types.LobbyScreen:
		content = m.lobby.View()
	case types.PracticeScreen:
		content = m.practice.View()
	case types.JoinScreen:
		content = m.join.View()
	default:
		content = m.home.View()
	}

	if m.connectionStatus == types.Connecting {
		spinnerView := lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(m.width).
			Height(m.height).
			Render("Connecting to server...\n" + m.spinner.View())
		content = lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(m.width).
			Height(m.height).
			Render(content + "\n\n" + spinnerView)
	}

	return zone.Scan(lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Width(m.width).
		Height(m.height).
		Render(content))
}
