// Package home provides the view
// for the default, home state of the application
package home

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/types"
)

// Model represents the home view's state
type Model struct {
	helpModel helpModel
	SelectedIndex int
}

// The number of menu options currently supported.
const numberOfOptions = 4

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true).
			Align(lipgloss.Center)

	spacerStyle = lipgloss.NewStyle().
			Margin(1, 0)

	optionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)
)

func New() tea.Model {
	return Model{
		helpModel: newHelpModel(),
		SelectedIndex: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("teletyperacer")
}

// Update handles messages for the home view
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.SelectedIndex == 0 {
				m.SelectedIndex = numberOfOptions
			} else {
				m.SelectedIndex--
			}

		case "down", "j":
			if m.SelectedIndex+1 == numberOfOptions {
				m.SelectedIndex = 0
			} else {
				m.SelectedIndex++
			}

		case "p":
			cmd = func() tea.Msg { return types.Practice }
		}
	}

	return m, cmd
}

// View renders the home view
func (m Model) View() string {
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
		"[G] Join Game",
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

	return content
}
