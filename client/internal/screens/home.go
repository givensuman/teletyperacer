// Package screens defines the various screens
// renderable by the TUI
package screens

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/givensuman/teletyperacer/client/internal/components/button"
	"github.com/givensuman/teletyperacer/client/internal/tui"
)

type HomeModel struct {
	cursor           int
	choices          [4]button.Model
	notification     string
	spinner          spinner.Model
	connectionStatus tui.ConnectionStatus
	connectionError  error
}

func NewHome() HomeModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	joinBtn := button.NewFocusedButton("Join", func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.JoinScreen} })
	hostBtn := button.NewButton("Host", func() tea.Msg { return tui.HostRoomMsg{} })

	// Initially disable Join and Host since connection starts as Connecting
	disabledJoinBtn, _ := joinBtn.Update(button.Disable)
	disabledHostBtn, _ := hostBtn.Update(button.Disable)
	joinBtn = disabledJoinBtn.(button.Model)
	hostBtn = disabledHostBtn.(button.Model)

	return HomeModel{
		cursor: 0,
		choices: [4]button.Model{
			joinBtn,
			hostBtn,
			button.NewButton("Practice", func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.PracticeScreen} }),
			button.NewButton("Quit", tea.Quit),
		},
		notification:     "",
		spinner:          s,
		connectionStatus: tui.Connecting,
		connectionError:  nil,
	}
}

func (m HomeModel) Init() tea.Cmd {
	longestButtonLabel := len(m.choices[0].GetLabel())
	for _, btn := range m.choices[1:] {
		if longestButtonLabel < len(btn.GetLabel()) {
			longestButtonLabel = len(btn.GetLabel())
		}
	}

	return tea.Batch(
		func() tea.Msg {
			return button.WidthMsg(longestButtonLabel)
		},
		m.spinner.Tick,
	)
}

// findNextEnabledButton finds the next enabled button index in the given direction
func (m HomeModel) findNextEnabledButton(currentIndex int, direction int) int {
	len := len(m.choices)
	for i := 1; i < len; i++ {
		nextIndex := (currentIndex + direction*i) % len
		if nextIndex < 0 {
			nextIndex += len
		}
		if !m.choices[nextIndex].IsDisabled() {
			return nextIndex
		}
	}
	return currentIndex // fallback to current if no enabled buttons found
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tui.ConnectionStatusMsg:
		m.connectionStatus = msg.Status
		m.connectionError = msg.Error
		switch msg.Status {
		case tui.Connected:
			m.notification = "Connected to server successfully."
			// Enable Join and Host buttons
			enabledJoin, _ := m.choices[0].Update(button.Enable)
			enabledHost, _ := m.choices[1].Update(button.Enable)
			m.choices[0] = enabledJoin.(button.Model)
			m.choices[1] = enabledHost.(button.Model)
		case tui.Connecting:
			m.notification = "Connecting to server..."
		case tui.ServerUnreachable:
			m.notification = "Server unreachable. Join and Host are disabled."
			// Disable Join and Host buttons
			disabledJoin, _ := m.choices[0].Update(button.Disable)
			disabledHost, _ := m.choices[1].Update(button.Disable)
			m.choices[0] = disabledJoin.(button.Model)
			m.choices[1] = disabledHost.(button.Model)
		case tui.ClientError:
			if msg.Error != nil {
				m.notification = fmt.Sprintf("Client error: %v. Join and Host are disabled.", msg.Error)
			} else {
				m.notification = "Client configuration error. Join and Host are disabled."
			}
			// Disable Join and Host buttons
			disabledJoin, _ := m.choices[0].Update(button.Disable)
			disabledHost, _ := m.choices[1].Update(button.Disable)
			m.choices[0] = disabledJoin.(button.Model)
			m.choices[1] = disabledHost.(button.Model)
		case tui.Failed:
			// Keep backward compatibility - treat as server unreachable
			m.notification = "Connection failed. Join and Host are disabled."
			// Disable Join and Host buttons
			disabledJoin, _ := m.choices[0].Update(button.Disable)
			disabledHost, _ := m.choices[1].Update(button.Disable)
			m.choices[0] = disabledJoin.(button.Model)
			m.choices[1] = disabledHost.(button.Model)

			// If current cursor is on a disabled button, move to next enabled one
			if m.choices[m.cursor].IsDisabled() {
				m.cursor = m.findNextEnabledButton(m.cursor, 1)
				if m.choices[m.cursor].IsDisabled() {
					m.cursor = m.findNextEnabledButton(m.cursor, -1)
				}
				// Focus the new cursor position
				btn, cmd := m.choices[m.cursor].Update(button.Focus)
				m.choices[m.cursor] = btn.(button.Model)
				cmds = append(cmds, cmd)
			}
		}

	case button.WidthMsg:
		for i, btn := range m.choices {
			updatedBtn, _ := btn.Update(msg)
			m.choices[i] = updatedBtn.(button.Model)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			prevBtn, _ := m.choices[m.cursor].Update(button.Unfocus)
			m.choices[m.cursor] = prevBtn.(button.Model)

			m.cursor = m.findNextEnabledButton(m.cursor, 1)
			btn, cmd := m.choices[m.cursor].Update(button.Focus)
			m.choices[m.cursor] = btn.(button.Model)
			cmds = append(cmds, cmd)

		case "k", "up":
			prevBtn, _ := m.choices[m.cursor].Update(button.Unfocus)
			m.choices[m.cursor] = prevBtn.(button.Model)

			m.cursor = m.findNextEnabledButton(m.cursor, -1)
			btn, cmd := m.choices[m.cursor].Update(button.Focus)
			m.choices[m.cursor] = btn.(button.Model)
			cmds = append(cmds, cmd)

		case "enter":
			if !m.choices[m.cursor].IsDisabled() {
				return m, tea.Batch(m.choices[m.cursor].GetAction(), m.spinner.Tick)
			}

		case "q":
			return m, tea.Quit
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			for i := range m.choices {
				if zone.Get(fmt.Sprintf("button-%d", i)).InBounds(msg) && !m.choices[i].IsDisabled() {
					return m, tea.Batch(m.choices[i].GetAction(), m.spinner.Tick)
				}
			}
		}

		switch msg.Button {
		case tea.MouseButtonWheelUp:
			prevBtn, _ := m.choices[m.cursor].Update(button.Unfocus)
			m.choices[m.cursor] = prevBtn.(button.Model)

			m.cursor = m.findNextEnabledButton(m.cursor, -1)
			btn, cmd := m.choices[m.cursor].Update(button.Focus)
			m.choices[m.cursor] = btn.(button.Model)
			cmds = append(cmds, cmd)

		case tea.MouseButtonWheelDown:
			prevBtn, _ := m.choices[m.cursor].Update(button.Unfocus)
			m.choices[m.cursor] = prevBtn.(button.Model)

			m.cursor = m.findNextEnabledButton(m.cursor, 1)
			btn, cmd := m.choices[m.cursor].Update(button.Focus)
			m.choices[m.cursor] = btn.(button.Model)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
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

	var status string
	switch m.connectionStatus {
	case tui.Connecting:
		status = m.spinner.View() + " Connecting..."
	case tui.Connected:
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Render("✓ Connected")
	case tui.ServerUnreachable:
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Render("✗ Server unreachable")
	case tui.ClientError:
		if m.connectionError != nil {
			status = lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Render(fmt.Sprintf("✗ Client error: %v", m.connectionError))
		} else {
			status = lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Render("✗ Client configuration error")
		}
	case tui.Failed:
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Render("✗ Connection failed")
	}

	statusNotifier := lipgloss.NewStyle().
		Padding(1, 0).
		AlignHorizontal(lipgloss.Left).
		Render(status)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("↑/k up • ↓/j down • enter select • q quit")

	fullContent := lipgloss.JoinVertical(lipgloss.Center, content, statusNotifier, help)

	return lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Render(fullContent)
}
