package practice

import (
	"github.com/charmbracelet/bubbletea"
)

type OptionsModel struct {
}

// Model implements tea.Model
var _ tea.Model = OptionsModel{}

func (m OptionsModel) Init() tea.Cmd {
	return nil
}

func (m OptionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m OptionsModel) View() string {
	return ""
}
