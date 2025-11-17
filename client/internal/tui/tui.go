// Package tui represents the root of the TUI application
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/internal/constants/view"
	"github.com/givensuman/teletyperacer/client/internal/tui/game"
	"github.com/givensuman/teletyperacer/client/internal/tui/home"
	"github.com/givensuman/teletyperacer/client/internal/tui/host"
	"github.com/givensuman/teletyperacer/client/internal/tui/lobby"
	"github.com/givensuman/teletyperacer/client/internal/tui/practice"
)

type Model struct {
	width       int
	height      int
	currentView view.View
	models      map[view.View]tea.Model
}

var models map[view.View]tea.Model = map[view.View]tea.Model{
	view.Home:            home.New(),
	view.HostOptions:     host.NewOptions(),
	view.Lobby:           lobby.New(),
	view.PracticeOptions: practice.NewOptions(),
	view.Practice:        practice.New(),
	view.Game:            game.New(),
}

func New() Model {
	return Model{
		currentView: view.Home,
		models:      models,
	}
}

var _ tea.Model = Model{}

func (m Model) getCurrentModel() tea.Model {
	return m.models[m.currentView]
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			return m, tea.Quit
		} else {
			model, _ := m.getCurrentModel().Update(msg)
			m.models[m.currentView] = model
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case view.NavigateMsg:
		m.currentView = view.View(msg)
	}

	return m, nil
}

var style lipgloss.Style = lipgloss.NewStyle().
	AlignHorizontal(lipgloss.Center)

func (m Model) View() string {
	return style.Width(m.width).Height(m.height).Render(m.getCurrentModel().View())
}
