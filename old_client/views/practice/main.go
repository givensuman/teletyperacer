// Package practice provides the view
// for the practice mode of the game.
package practice

import "github.com/charmbracelet/bubbletea"

type Model struct {

}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "Hello"
}
