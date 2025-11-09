package views

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/go-sockets/client"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true).
			Align(lipgloss.Center)

	spacerStyle = lipgloss.NewStyle().
			Margin(3, 0)

	optionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(2)
)

func HomeView(m Model) string {
	title := []string{
		`
⠀⢀⣀⣀⣀⠀⠀⠀⠀⢀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⢸⣿⣿⡿⢀⣠⣴⣾⣿⣿⣿⣿⣇⡀⠀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⢸⣿⣿⠟⢋⡙⠻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣶⣿⡿⠓⡐⠒⢶⣤⣄⡀⠀⠀
⠀⠸⠿⠇⢰⣿⣿⡆⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⣿⣿⡷⠈⣿⣿⣉⠁⠀
⠀⠀⠀⠀⠀⠈⠉⠀⠈⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠀⠈⠉⠁⠀⠈⠉⠉⠀⠀
`,
		"teletyperacer",
	}

	var styledTitle []string
	for _, s := range title {
		styledTitle = append(styledTitle, titleStyle.Render(s))
	}

	options := []string{
		"[P] Practice",
		"[H] Host Game",
		"[J] Join Game",
		"[Q] Quit",
	}

	var styledOptions []string
	for _, opt := range options {
		styledOptions = append(styledOptions, optionStyle.Render(opt))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, styledTitle...),
		spacerStyle.Render(""),
		lipgloss.JoinVertical(lipgloss.Left, styledOptions...),
	)

	return containerStyle.Width(m.Width).Height(m.Height).Render(content)
}

// connectToRoom establishes a WebSocket connection to a room.
func connectToRoom(roomID string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		socket, err := client.Connect("ws://localhost:3000/ws/", "/", func(s *client.Socket) {
			s.Emit("join", roomID)
		})
		if err != nil {
			// For simplicity, ignore error or handle later
			return nil
		}
		return ConnectedMsg{Socket: socket, RoomID: roomID}
	})
}

// ConnectedMsg is sent when a socket connection is established.
type ConnectedMsg struct {
	Socket *client.Socket
	RoomID string
}
