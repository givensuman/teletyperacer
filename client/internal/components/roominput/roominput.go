// Package roominput defines the room code input component
package roominput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textInput textinput.Model
	width     int
	height    int
}

var _ tea.Model = Model{}

func NewRoomInput() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter room code"
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 30

	return Model{
		textInput: ti,
		width:     50,
		height:    10,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, func() tea.Msg {
				return JoinRoomMsg{Code: m.textInput.Value()}
			}
		case tea.KeyEsc:
			return m, func() tea.Msg { return HideMsg{} }
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	// Create the input box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.ANSIColor(4)).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	inputView := m.textInput.View()

	content := labelStyle.Render("Enter Room Code") + "\n\n" +
		inputView + "\n\n" +
		helpStyle.Render("enter submit • esc cancel")

	return boxStyle.Render(content)
}
