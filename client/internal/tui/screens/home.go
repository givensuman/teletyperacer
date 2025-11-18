// Package screens defines the various screens
// renderable by the TUI
package screens

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone"

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

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			for i := range m.choices {
				if zone.Get(fmt.Sprintf("button-%d", i)).InBounds(msg) {
					return m, m.choices[i].GetAction()
				}
			}
		}

		switch msg.Button {
			case tea.MouseButtonWheelUp:
				// Simulate up
				prevBtn, _ := m.choices[m.cursor].Update(button.Unfocus)
				m.choices[m.cursor] = prevBtn.(button.Model)

				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.choices) - 1
				}
				btn, cmd := m.choices[m.cursor].Update(button.Focus)
				m.choices[m.cursor] = btn.(button.Model)

				return m, cmd

		case tea.MouseButtonWheelDown:
			// Simulate down
			prevBtn, _ := m.choices[m.cursor].Update(button.Unfocus)
			m.choices[m.cursor] = prevBtn.(button.Model)

			m.cursor = (m.cursor + 1) % len(m.choices)
			btn, cmd := m.choices[m.cursor].Update(button.Focus)
			m.choices[m.cursor] = btn.(button.Model)

			return m, cmd
		}
	}

	return m, nil
}

func (m HomeModel) View() string {
	var views []string
	for i, btn := range m.choices {
		views = append(views, zone.Mark(fmt.Sprintf("button-%d", i), btn.View()))
	}

	buttons := lipgloss.JoinVertical(lipgloss.Left, views...)

	renderedTitle := lipgloss.NewStyle().
		AlignVertical(lipgloss.Left).
		AlignHorizontal(lipgloss.Left).
		Render(`
  _______ ______ _      ______         
 |__   __|  ____| |    |  ____|        
    | |  | |__  | |    | |__           
    | |  |  __| | |    |  __|          
    | |  | |____| |____| |____         
  __|_|__|______|______|______|        
 |__   __\ \   / /  __ \|  ____|       
    | |   \ \_/ /| |__) | |__          
    | |    \   / |  ___/|  __|         
    | |     | |  | |    | |____        
  __|_|     |_|  |_|____|______|_____  
 |  __ \     /\   / ____|  ____|  __ \ 
 | |__) |   /  \ | |    | |__  | |__) |
 |  _  /   / /\ \| |    |  __| |  _  / 
 | | \ \  / ____ \ |____| |____| | \ \ 
 |_|  \_\/_/    \_\_____|______|_|  \_\
`)

	renderedButtons := lipgloss.NewStyle().
		Padding(1).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Right).
		Render(buttons)

	content := lipgloss.JoinHorizontal(lipgloss.Center, renderedTitle, renderedButtons)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("↑/k up • ↓/j down • enter select • q quit")

	fullContent := lipgloss.JoinVertical(lipgloss.Center, content, help)

	return lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Render(fullContent)
}
