// Package views contains the logic for rendering the various
// screens rendered for the TUI.
package views

import (
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/givensuman/go-sockets/client"

	"github.com/givensuman/teletyperacer/client/components/container"
	"github.com/givensuman/teletyperacer/client/types"
	"github.com/givensuman/teletyperacer/client/views/home"
	"github.com/givensuman/teletyperacer/client/views/practice"
)

type Model struct {
	// The view model parents a container.
	types.SingleParent[container.Model]
	// The currently requested view.
	CurrentView types.View
	// TODO: Networking
	RoomID string
	Socket *client.Socket
}

// New instantiates a new view model.
func New() tea.Model {
	m := Model{
		CurrentView: types.Home,
	}

	m.Child = container.New(home.New())

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return m.Child.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	child, cmd := m.Child.Update(msg)
	m.Child = child.(container.Model)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
			case tea.KeyCtrlC:
				os.Exit(0)
		}

	case types.View:
		switch msg {
		case types.Home:
			m.CurrentView = types.Home
			m.Child.Replace(home.New())

		case types.Practice:
			m.CurrentView = types.Practice
			m.Child.Replace(practice.New())

		default:
			print("Not supported yet!")
		}

	case home.ConnectedMsg:
		m.Socket = msg.Socket
		m.RoomID = msg.RoomID
		m.Socket.On("game_started", func(args ...any) {
			// Handle game start
		})
		m.CurrentView = types.Lobby
	}

	return m, cmd
}
