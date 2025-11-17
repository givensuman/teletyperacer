package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PracticeModel struct{}

func NewPractice() PracticeModel {
	return PracticeModel{}
}

func (m PracticeModel) Init() tea.Cmd {
	return nil
}

func (m PracticeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m PracticeModel) View() string {
	return lipgloss.NewStyle().Padding(1).Render("Practice Screen\n\nPress ESC to go back to Home.")
}
