// Package root describes the root of
// the TUI application
package root

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	sockets "github.com/givensuman/go-sockets/client"
	"github.com/givensuman/teletyperacer/client/internal/tui"
	"github.com/givensuman/teletyperacer/client/internal/tui/screens"
)

type Model struct {
	width  int
	height int
	// Currently rendered screen
	screen tui.Screen
	// Child models
	home,
	host,
	practice tea.Model
	// WebSocket connection
	socket  *sockets.Socket
	status  tui.ConnectionStatus
	spinner spinner.Model
}

func New() Model {
	socket, err := sockets.Connect("ws://localhost:3000/ws/", "/", nil)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	status := tui.Connected
	if err != nil {
		status = tui.Failed
	}

	return Model{
		screen:   tui.HomeScreen,
		home:     screens.NewHome(),
		host:     screens.NewHost(),
		practice: screens.NewPractice(),
		socket:   socket,
		status:   status,
		spinner:  s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
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
		switch msg.String() {
		case "q", "ctrl+c":
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
		case "q", "ctrl+c":
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
		return m.home.View()
	}

	if m.status == tui.Connecting {
		spinnerView := lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(m.width).
			Height(m.height).
			Render("Connecting to server...\n" + m.spinner.View())
		return lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Width(m.width).
			Height(m.height).
			Render(content + "\n\n" + spinnerView)
	}

	return lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Width(m.width).
		Height(m.height).
		Render(content)
}
