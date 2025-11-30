// Package alert defines the alert component
package alert

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/teletyperacer/client/internal/components/button"
)

type Model struct {
	title   string
	message string
	buttons []button.Model
	width   int
	height  int
}

var _ tea.Model = Model{}

func NewAlert(title, message string, buttons ...button.Model) Model {
	// Focus the first button
	if len(buttons) > 0 {
		buttons[0] = buttons[0].Focus()
	}
	return Model{
		title:   title,
		message: message,
		buttons: buttons,
		width:   50,
		height:  10,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Trigger the focused button's action
			for _, btn := range m.buttons {
				if btn.IsFocused() {
					return m, btn.GetAction()
				}
			}
		case tea.KeyTab:
			// Cycle focus through buttons
			focusedIndex := -1
			for i, btn := range m.buttons {
				if btn.IsFocused() {
					focusedIndex = i
					break
				}
			}
			nextIndex := (focusedIndex + 1) % len(m.buttons)
			for i := range m.buttons {
				if i == nextIndex {
					m.buttons[i] = m.buttons[i].Focus()
				} else {
					m.buttons[i] = m.buttons[i].Unfocus()
				}
			}
		}
	}

	// Update buttons
	var cmds []tea.Cmd
	for i, btn := range m.buttons {
		updated, cmd := btn.Update(msg)
		m.buttons[i] = updated.(button.Model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// Create the alert box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.ANSIColor(4)).
		Padding(1, 2).
		Width(m.width - 10).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.ANSIColor(4)).
		Align(lipgloss.Center)

	messageStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		MarginBottom(1)

	buttonsView := ""
	for _, btn := range m.buttons {
		buttonsView += btn.View() + " "
	}

	content := titleStyle.Render(m.title) + "\n\n" +
		messageStyle.Render(m.message) + "\n\n" +
		buttonsView

	return boxStyle.Render(content)
}
