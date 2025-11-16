package tui

import (
	"github.com/charmbracelet/bubbletea"

	"github.com/givensuman/teletyperacer/client/internal/tui/home"
)

type View int8

const (
	Home View = iota
	PracticeOptions
	Practice
	HostOptions
	Lobby
	Game
)

type Model struct {
	currentView View
	homeModel   home.Model
}

func New() Model {
	return Model{
		currentView: Home,
		homeModel:   home.New(),
	}
}

var _ tea.Model = Model{}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			return m, tea.Quit
		}
	}
	switch m.currentView {
	case Home:
		model, cmd := m.homeModel.Update(msg)
		m.homeModel = model.(home.Model)

		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	return m.homeModel.View()
}
