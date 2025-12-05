// Package input defines a generic text input component
package input

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Placeholder    string
	Label          string
	SubmittedLabel string
	SubmittedText  string
	CharLimit      int
}

type Model struct {
	textInput textinput.Model
	width     int
	height    int
	submitted bool
	config    Config
}

var _ tea.Model = Model{}

func NewInput(config Config) Model {
	ti := textinput.New()
	ti.Placeholder = config.Placeholder
	ti.Focus()
	ti.CharLimit = config.CharLimit
	ti.Width = 30

	return Model{
		textInput: ti,
		width:     50,
		height:    10,
		submitted: false,
		config:    config,
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
			m.submitted = true
			return m, func() tea.Msg {
				return SubmitMsg{Value: strings.ToUpper(m.textInput.Value())}
			}
		case tea.KeyEsc:
			return m, func() tea.Msg { return HideMsg{} }
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			// Uppercase the input after updating
			currentValue := m.textInput.Value()
			if currentValue != strings.ToUpper(currentValue) {
				m.textInput.SetValue(strings.ToUpper(currentValue))
			}
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	// Uppercase the input after updating
	currentValue := m.textInput.Value()
	if currentValue != strings.ToUpper(currentValue) {
		m.textInput.SetValue(strings.ToUpper(currentValue))
	}
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

	var content string

	if m.submitted {
		content = labelStyle.Render(m.config.SubmittedLabel) + "\n\n" +
			m.config.SubmittedText + "\n\n" +
			helpStyle.Render("esc cancel")
	} else {
		inputView := m.textInput.View()
		content = labelStyle.Render(m.config.Label) + "\n\n" +
			inputView + "\n\n" +
			helpStyle.Render("enter submit • esc cancel")
	}

	return boxStyle.Render(content)
}
