// Package button defines the button component
package button

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/internal/constants/colors"
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

func New(index int, label string, action func()) Model {
	return Model{
		Index:      index,
		Label:      label,
		Action:     action,
		IsFocused:  false,
		IsDisabled: false,
		Style: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 1).
			Foreground(lipgloss.ANSIColor(colors.Black)).
			Background(lipgloss.ANSIColor(colors.Blue)),
		FocusStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 1).
			Foreground(lipgloss.ANSIColor(colors.White)).
			Background(lipgloss.ANSIColor(colors.Grey)),
	}
}

func (m *Model) WithStyle(style lipgloss.Style) {
	m.Style = style
}

func (m *Model) WithFocusStyle(style lipgloss.Style) {
	m.FocusStyle = style
}

type FocusMsg bool
const (
	Focus FocusMsg = true
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

	return m.Style.Render(m.Label, fmt.Sprintf("%b", m.IsFocused))
}
