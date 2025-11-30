// Package root describes the root of
// the TUI application
package root

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	sockets "github.com/givensuman/go-sockets/client"
	"github.com/givensuman/teletyperacer/client/internal/components/alert"
	"github.com/givensuman/teletyperacer/client/internal/components/roominput"
	"github.com/givensuman/teletyperacer/client/internal/store"
	"github.com/givensuman/teletyperacer/client/internal/tui"
	"github.com/givensuman/teletyperacer/client/internal/tui/screens"
)

type Model struct {
	// Currently rendered screen
	screen tui.Screen
	// Child models
	home,
	host,
	practice tea.Model
	// WebSocket connection
	socket  *sockets.Socket
	spinner spinner.Model
	// Alert overlay
	alert     tea.Model
	showAlert bool
	// Room input overlay
	roomInput     tea.Model
	showRoomInput bool
}

type backgroundModel struct {
	root *Model
}

func (b backgroundModel) Init() tea.Cmd { return nil }

func (b backgroundModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return b, nil }

func (b backgroundModel) View() string {
	var content string
	switch b.root.screen {
	case tui.HomeScreen:
		content = b.root.home.View()
	case tui.HostScreen:
		content = b.root.host.View()
	case tui.PracticeScreen:
		content = b.root.practice.View()
	default:
		content = b.root.home.View()
	}

	store := store.GetStore()
	if store.ConnectionStatus == tui.Connecting {
		spinnerView := lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(store.Width).
			Height(store.Height).
			Render("Connecting to server...\n" + b.root.spinner.View())
		return zone.Scan(lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(store.Width).
			Height(store.Height).
			Render(content + "\n\n" + spinnerView))
	}

	return zone.Scan(lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Width(store.Width).
		Height(store.Height).
		Render(content))
}

func New() Model {
	socket, err := sockets.Connect("ws://localhost:3000/ws/", "/", nil)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	store := store.GetStore()
	if err != nil {
		store.ConnectionStatus = tui.Failed
	} else {
		store.ConnectionStatus = tui.Connected
	}

	return Model{
		screen:        tui.HomeScreen,
		home:          screens.NewHome(),
		host:          screens.NewHost(),
		practice:      screens.NewPractice(),
		socket:        socket,
		spinner:       s,
		showAlert:     false,
		showRoomInput: false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		store := store.GetStore()
		store.Width = msg.Width
		store.Height = msg.Height
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tui.ConnectionStatusMsg:
		store := store.GetStore()
		store.ConnectionStatus = msg.Status
		return m, nil

	case alert.ShowMsg:
		m.alert = alert.NewAlert(msg.Title, msg.Message, msg.Buttons...)
		m.showAlert = true
		return m, nil

	case alert.HideMsg:
		m.showAlert = false
		return m, nil

	case roominput.ShowMsg:
		m.roomInput = roominput.NewRoomInput()
		store := store.GetStore()
		m.roomInput, _ = m.roomInput.Update(tea.WindowSizeMsg{Width: store.Width, Height: store.Height})
		m.showRoomInput = true
		return m, m.roomInput.Init()

	case roominput.HideMsg:
		m.showRoomInput = false
		return m, nil

	case roominput.JoinRoomMsg:
		// TODO: Handle joining room with code
		m.showRoomInput = false
		return m, nil
	}

	if m.showAlert {
		var cmd tea.Cmd
		m.alert, cmd = m.alert.Update(msg)
		return m, cmd
	}

	if m.showRoomInput {
		var cmd tea.Cmd
		m.roomInput, cmd = m.roomInput.Update(msg)
		return m, cmd
	}

	switch m.screen {
	case tui.HomeScreen:
		return m.updateHome(msg)
	case tui.HostScreen:
		return m.updateHost(msg)
	case tui.PracticeScreen:
		return m.updatePractice(msg)
	default:
		return m, nil
	}
}

func (m Model) updateHome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tui.ScreenChangeMsg:
		m.screen = msg.Screen
		return m, nil
	}

	var cmd tea.Cmd
	m.home, cmd = m.home.Update(msg)
	return m, cmd
}

func (m Model) updateHost(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.screen = tui.HomeScreen
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.host, cmd = m.host.Update(msg)
	return m, cmd
}

func (m Model) updatePractice(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.screen = tui.HomeScreen
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.practice, cmd = m.practice.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var content string
	switch m.screen {
	case tui.HomeScreen:
		content = m.home.View()
	case tui.HostScreen:
		content = m.host.View()
	case tui.PracticeScreen:
		content = m.practice.View()
	default:
		content = m.home.View()
	}

	store := store.GetStore()
	if store.ConnectionStatus == tui.Connecting {
		spinnerView := lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(store.Width).
			Height(store.Height).
			Render("Connecting to server...\n" + m.spinner.View())
		content = lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(store.Width).
			Height(store.Height).
			Render(content + "\n\n" + spinnerView)
	}

	if m.showAlert {
		return zone.Scan(lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(store.Width).
			Height(store.Height).
			Render(m.alert.View()))
	}

	if m.showRoomInput {
		return zone.Scan(lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(store.Width).
			Height(store.Height).
			Render(m.roomInput.View()))
	}

	return zone.Scan(lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Width(store.Width).
		Height(store.Height).
		Render(content))
}
