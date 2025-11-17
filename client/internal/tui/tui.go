// Package tui defines the different
// screens which may be rendered
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Screen int

const (
	HomeScreen Screen = iota
	HostScreen
	PracticeScreen
)

var viewableScreens map[Screen]tea.Model = map[Screen]tea.Model{
	HomeScreen:     NewHome(),
	HostScreen:     NewHost(),
	PracticeScreen: NewPractice(),
}

type Model struct {
	width  int
	height int
	// Currently rendered screen
	screen Screen
	// Child models
	home,
	host,
	practice tea.Model
}

func New() Model {
	return Model{
		screen:   HomeScreen,
		home:     NewHome(),
		host:     NewHost(),
		practice: NewPractice(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.home.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	}

	switch m.screen {
	case HomeScreen:
		return m.updateHome(msg)
	case HostScreen:
		return m.updateHost(msg)
	case PracticeScreen:
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
	case ScreenChangeMsg:
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
			m.screen = HomeScreen
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
			m.screen = HomeScreen
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
	case HomeScreen:
		content = m.home.View()
	case HostScreen:
		content = m.host.View()
	case PracticeScreen:
		content = m.practice.View()
	default:
		return m.home.View()
	}

	return lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Width(m.width).
		Height(m.height).
		Render(content)
}

type ScreenChangeMsg struct {
	Screen Screen
}
