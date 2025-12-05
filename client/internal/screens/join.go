// Package screens defines the various screens
// renderable by the TUI
package screens

import (
	"strings"

	"github.com/charmbracelet/bubbletea"

	"github.com/givensuman/teletyperacer/client/internal/components/input"
	"github.com/givensuman/teletyperacer/client/internal/tui"
)

type JoinModel struct {
	input tea.Model
}

func NewJoin() JoinModel {
	config := input.Config{
		Placeholder:    "Enter room code",
		Label:          "Enter Room Code",
		SubmittedLabel: "Joining Room...",
		SubmittedText:  "Attempting to join room " + strings.ToUpper(""), // This will be updated
		CharLimit:      6,
	}
	return JoinModel{
		input: input.NewInput(config),
	}
}

func (m JoinModel) Init() tea.Cmd {
	return m.input.Init()
}

func (m JoinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case input.SubmitMsg:
		// Send join room message to server
		return m, func() tea.Msg {
			return tui.JoinRoomMsg{Code: msg.Value}
		}
	case tui.AppRoomJoinedMsg:
		// Successfully joined room, go to lobby
		return m, func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.LobbyScreen} }
	case input.HideMsg:
		return m, func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.HomeScreen} }
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m JoinModel) View() string {
	return m.input.View()
}
