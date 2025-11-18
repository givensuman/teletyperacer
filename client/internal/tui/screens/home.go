// Package screens defines the various screens
// renderable by the TUI
package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/givensuman/teletyperacer/client/internal/components/button"
	"github.com/givensuman/teletyperacer/client/internal/tui"
)

type HomeModel struct {
	cursor       int
	choices      [4]button.Model
	status       tui.ConnectionStatus
	notification string
}

func NewHome() HomeModel {
	return HomeModel{
		cursor: 0,
		choices: [4]button.Model{
			button.NewFocusedButton("Join", nil),
			button.NewButton("Host", func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.HostScreen} }),
			button.NewButton("Practice", func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.PracticeScreen} }),
			button.NewButton("Quit", tea.Quit),
		},
		status:       tui.Connecting,
		notification: "",
	}
}

func (m HomeModel) Init() tea.Cmd {
	longestButtonLabel := len(m.choices[0].GetLabel())
	for _, btn := range m.choices[1:] {
		if longestButtonLabel < len(btn.GetLabel()) {
			longestButtonLabel = len(btn.GetLabel())
		}
	}

	return func() tea.Msg {
		return button.WidthMsg(longestButtonLabel)
	}
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tui.ConnectionStatusMsg:
		switch msg.Status {
		case tui.Connected:
			m.notification = "Connected to server successfully."

		case tui.Failed:
			m.notification = "Failed to connect to server. Join and Host are disabled."
		}

		return m, nil

	case button.WidthMsg:
		for i, btn := range m.choices {
			updatedBtn, _ := btn.Update(msg)
			m.choices[i] = updatedBtn.(button.Model)
		}

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			prevBtn, _ := m.choices[m.cursor].Update(button.Unfocus)
			m.choices[m.cursor] = prevBtn.(button.Model)

			m.cursor = (m.cursor + 1) % len(m.choices)
			btn, cmd := m.choices[m.cursor].Update(button.Focus)
			m.choices[m.cursor] = btn.(button.Model)

			return m, cmd

		case "k", "up":
			prevBtn, _ := m.choices[m.cursor].Update(button.Unfocus)
			m.choices[m.cursor] = prevBtn.(button.Model)

			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
			btn, cmd := m.choices[m.cursor].Update(button.Focus)
			m.choices[m.cursor] = btn.(button.Model)

			return m, cmd

		case "enter":
			return m, m.choices[m.cursor].GetAction()
		}
	}

	return m, nil
}

func (m HomeModel) View() string {
	var views []string
	for _, btn := range m.choices {
		views = append(views, btn.View())
	}

	buttons := lipgloss.JoinVertical(lipgloss.Center, views...)

	return lipgloss.NewStyle().
		Padding(1).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Render(buttons)
}
