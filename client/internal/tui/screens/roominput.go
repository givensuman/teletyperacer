// Package screens defines the various screens
// renderable by the TUI
package screens

import (
	"strings"

	"github.com/charmbracelet/bubbletea"

	"github.com/givensuman/teletyperacer/client/internal/components/roominput"
	"github.com/givensuman/teletyperacer/client/internal/tui"
)

type RoomInputModel struct {
	roomInput tea.Model
}

func NewRoomInput() RoomInputModel {
	return RoomInputModel{
		roomInput: roominput.NewRoomInput(),
	}
}

func (m RoomInputModel) Init() tea.Cmd {
	return m.roomInput.Init()
}

func (m RoomInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case roominput.JoinRoomMsg:
		// Send join room message to server
		return m, func() tea.Msg {
			return tui.JoinRoomMsg{Code: strings.ToUpper(msg.Code)}
		}
	case tui.RoomJoinedMsg:
		// Successfully joined room, go back to home
		return m, func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.HomeScreen} }
	case roominput.HideMsg:
		return m, func() tea.Msg { return tui.ScreenChangeMsg{Screen: tui.HomeScreen} }
	}

	var cmd tea.Cmd
	m.roomInput, cmd = m.roomInput.Update(msg)
	return m, cmd
}

func (m RoomInputModel) View() string {
	return m.roomInput.View()
}
