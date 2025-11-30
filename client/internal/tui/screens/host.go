package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HostModel struct{}

func NewHost() HostModel {
	return HostModel{}
}

func (m HostModel) Init() tea.Cmd {
	return nil
}

func (m HostModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m HostModel) View() string {
	return lipgloss.NewStyle().Padding(1).Render("Host Screen\n\nPress ESC to go back to Home.")
}
