// Package button defines the button component
package button

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/internal/constants"
)

type Model struct {
	Index      int
	Label      string
	Action     func()
	IsFocused  bool
	IsDisabled bool
	Style      lipgloss.Style
	FocusStyle lipgloss.Style
}

var _ tea.Model = Model{}

func NewButton(index int, label string, action func()) Model {
	return Model{
		Index:      index,
		Label:      label,
		Action:     action,
		IsFocused:  false,
		IsDisabled: false,
		Style: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 1).
			Foreground(lipgloss.ANSIColor(constants.Black)).
			Background(lipgloss.ANSIColor(constants.Blue)).
			Faint(true),

		FocusStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 1).
			Foreground(lipgloss.ANSIColor(constants.White)).
			Background(lipgloss.ANSIColor(constants.Grey)).
			Bold(true),
	}
}

type FocusMsg bool

const (
	Focus   FocusMsg = true
	Unfocus FocusMsg = false
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && m.IsFocused && m.Action != nil {
			go m.Action()

			return m, nil
		}

	case FocusMsg:
		switch msg {
		case Focus:
			m.IsFocused = true
		case Unfocus:
			m.IsFocused = false
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
