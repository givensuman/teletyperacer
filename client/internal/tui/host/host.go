package host

import (
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
}

// Model implements tea.Model
var _ tea.Model = Model{}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m Model) View() string {
	return ""
}
