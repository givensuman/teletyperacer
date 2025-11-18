// Package button defines the button component
package button

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
)

type Model struct {
	label      string
	action     tea.Cmd
	isFocused  bool
	isDisabled bool
	style      lipgloss.Style
	focusStyle lipgloss.Style
}

var _ tea.Model = Model{}

func NewButton(label string, action tea.Cmd) Model {
	return Model{
		label:      label,
		action:     action,
		isFocused:  false,
		isDisabled: false,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(lipgloss.White)).
			Background(lipgloss.ANSIColor(lipgloss.BrightBlack)).
			Faint(true),

		focusStyle: lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(lipgloss.Black)).
			Background(lipgloss.ANSIColor(lipgloss.BrightBlue)).
			Bold(true),
	}
}

func (m Model) GetAction() tea.Cmd {
	return m.action
}

func (m Model) GetLabel() string {
	return m.label
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && m.isFocused && m.action != nil {
			return m, m.action
		}

	case FocusMsg:
		switch msg {
		case Focus:
			m.isFocused = true
		case Unfocus:
			m.isFocused = false
		}

	case DisableMsg:
		switch msg {
		case Disable:
			m.isDisabled = true
		case Enable:
			m.isDisabled = false
		}

	case WidthMsg:
		m.style = m.style.Width(int(msg))
		m.focusStyle = m.focusStyle.Width(int(msg))
	}

	return m, nil
}

func (m Model) View() string {
	if m.isFocused {
		return m.focusStyle.Render(m.label)
	}

	return m.style.Render(m.label)
}
