// Package root describes the root of
// the TUI application
package root

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	zone "github.com/lrstanley/bubblezone"

	"github.com/givensuman/teletyperacer/client/internal/tui/components/roominput"
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
	Code    string   `json:"code"`
	Players []string `json:"players"`
}

type PlayerJoinedData struct {
	PlayerName string `json:"playerName"`
}

type Model struct {
	// Currently rendered screen
	screen types.Screen
	// Child models
	home,
	host,
	practice tea.Model
	// Room input component (used as screen)
	roomInput tea.Model
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
	case types.HostScreen:
		content = b.root.host.View()
	case types.PracticeScreen:
		content = b.root.practice.View()
	case types.RoomInputScreen:
		content = b.root.roomInput.View()
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
			if name, ok := d["playerName"].(string); ok {
				playerData.PlayerName = name
			}
		} else if d, ok := data.(PlayerJoinedData); ok {
			playerData = d
		}
		return types.PlayerJoinedMsg{PlayerName: playerData.PlayerName}

	case "roomState":
		var stateData RoomStateData
		if d, ok := data.(map[string]interface{}); ok {
			if code, ok := d["code"].(string); ok {
				stateData.Code = code
			}
			if players, ok := d["players"].([]interface{}); ok {
				for _, p := range players {
					if name, ok := p.(string); ok {
						stateData.Players = append(stateData.Players, name)
					}
				}
			}
		} else if d, ok := data.(RoomStateData); ok {
			stateData = d
		}
		return types.RoomStateMsg{Code: stateData.Code, Players: stateData.Players}

	case "error":
		// For now, ignore errors. Could return an error message to display
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
		host:             screens.NewHost(),
		practice:         screens.NewPractice(),
		roomInput:        roominput.NewRoomInput(),
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
		if msg.Screen == types.RoomInputScreen {
			m.roomInput = roominput.NewRoomInput()
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
		// Get the join code from the host screen
		if hostModel, ok := m.host.(screens.HostModel); ok {
			m.sendWSMessage("createRoom", CreateRoomData{Code: hostModel.GetJoinCode()})
		}
		return m, nil

	case types.JoinRoomMsg:
		m.sendWSMessage("joinRoom", JoinRoomData{Code: msg.Code})
		return m, nil

	case roominput.JoinRoomMsg:
		// Send join room message to server
		return m, func() tea.Msg {
			return types.JoinRoomMsg{Code: strings.ToUpper(msg.Code)}
		}

	case roominput.HideMsg:
		return m, func() tea.Msg { return types.ScreenChangeMsg{Screen: types.HomeScreen} }

	case types.RoomJoinedMsg:
		// Successfully joined room, go back to home
		return m, func() tea.Msg { return types.ScreenChangeMsg{Screen: types.HomeScreen} }

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
	case types.HostScreen:
		m.host, cmd = m.host.Update(msg)
	case types.PracticeScreen:
		m.practice, cmd = m.practice.Update(msg)
	case types.RoomInputScreen:
		m.roomInput, cmd = m.roomInput.Update(msg)
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
	case types.HostScreen:
		content = m.host.View()
	case types.PracticeScreen:
		content = m.practice.View()
	case types.RoomInputScreen:
		content = m.roomInput.View()
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
