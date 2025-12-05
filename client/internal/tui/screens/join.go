package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/givensuman/teletyperacer/client/internal/tui/components/input"
	"github.com/givensuman/teletyperacer/client/internal/types"
)

type JoinModel struct {
	input input.Model
}

func NewJoin() JoinModel {
	config := input.Config{
		Placeholder:    "Enter room code",
		Label:          "Enter Room Code",
		SubmittedLabel: "Joining Room...",
		SubmittedText:  "Attempting to join room",
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
		return m, func() tea.Msg {
			return types.JoinRoomMsg{Code: strings.ToUpper(msg.Value)}
		}
	case input.HideMsg:
		return m, func() tea.Msg { return types.ScreenChangeMsg{Screen: types.HomeScreen} }
	case types.RoomJoinedMsg:
		// Successfully joined room, go to lobby as player
		return m, func() tea.Msg { return types.ScreenChangeMsg{Screen: types.LobbyScreen} }
	default:
		var cmd tea.Cmd
		updatedInput, cmd := m.input.Update(msg)
		m.input = updatedInput.(input.Model)
		return m, cmd
	}
}

func (m JoinModel) View() string {
	return m.input.View()
}
