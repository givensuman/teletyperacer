// Package button defines the button component
package button

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	label         string
	action        tea.Cmd
	isFocused     bool
	isDisabled    bool
	style         lipgloss.Style
	focusStyle    lipgloss.Style
	disabledStyle lipgloss.Style
}

var _ tea.Model = Model{}

func NewButton(label string, action tea.Cmd) Model {
	return Model{
		label:      label,
		action:     action,
		isFocused:  false,
		isDisabled: false,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(15)).
			Faint(true).
			PaddingLeft(1).
			Align(lipgloss.Center),

		focusStyle: lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(4)).
			Bold(true).
			Border(lipgloss.ASCIIBorder(), false, false, false, true).
			BorderForeground(lipgloss.ANSIColor(4)).
			Align(lipgloss.Center),

		disabledStyle: lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(8)).
			Faint(true).
			PaddingLeft(1).
			Align(lipgloss.Center),
	}
}

func NewFocusedButton(label string, action tea.Cmd) Model {
	m := NewButton(label, action)
	m.isFocused = true

	return m
}

func (m Model) GetAction() tea.Cmd {
	return m.action
}

func (m Model) GetLabel() string {
	return m.label
}

func (m Model) IsFocused() bool {
	return m.isFocused
}

func (m Model) IsDisabled() bool {
	return m.isDisabled
}

func (m Model) Focus() Model {
	m.isFocused = true
	return m
}

func (m Model) Unfocus() Model {
	m.isFocused = false
	return m
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
		m.disabledStyle = m.disabledStyle.Width(int(msg))
	}

	return m, nil
}

func (m Model) View() string {
	if m.isDisabled {
		return m.disabledStyle.Render(m.label)
	}

	if m.isFocused {
		return m.focusStyle.Render(m.label)
	}

	return m.style.Render(m.label)
}
