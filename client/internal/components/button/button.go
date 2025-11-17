// Package button defines the button component
package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/internal/constants/colors"
)

type Model struct {
	Index      int
	Label      string
	IsFocused  bool
	IsDisabled bool
	Style      lipgloss.Style
	FocusStyle lipgloss.Style
}

var _ tea.Model = Model{}

func New(index int, label string) Model {
	return Model{
		Index:      index,
		Label:      label,
		IsFocused:  false,
		IsDisabled: false,
		Style: lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			Margin(1, 1).
			Foreground(lipgloss.ANSIColor(colors.White)).
			Faint(true),
		FocusStyle: lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			Margin(1, 1).
			Foreground(lipgloss.ANSIColor(colors.White)).
			Bold(true),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg {

	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter && m.IsFocused && !m.IsDisabled {
				return m, func() tea.Msg {
					return ButtonClickedMsg{Index: m.Index}
				}
			}

		case FocusMsg:
			if msg.Focus == m.Index {
				m.IsFocused = true
			} else if msg.Unfocus == m.Index {
				m.IsFocused = false
			}

		case DisableMsg:
			m.IsDisabled = bool(msg)

		case WidthMsg:
			m.Style = m.Style.Width(int(msg))
			m.FocusStyle = m.FocusStyle.Width(int(msg))
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.IsFocused {
		return m.FocusStyle.Render(m.Label)
	}

	return m.Style.Render(m.Label)
}
