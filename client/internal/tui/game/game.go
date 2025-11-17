package game

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/internal/components/button"
	"github.com/givensuman/teletyperacer/client/internal/constants/view"
)

type Model struct {
	Buttons       [1]button.Model
	SelectedIndex int
	Text          string
}

// Model implements tea.Model
var _ tea.Model = Model{}

var gameButtons [1]button.Model = [1]button.Model{
	button.New(0, "Back"),
}

func New() Model {
	focusedFirstButton, _ := gameButtons[0].Update(button.FocusMsg{0, -1})
	gameButtons[0] = focusedFirstButton.(button.Model)

	longestButtonLabel := len(gameButtons[0].Label)
	widthMsg := button.WidthMsg(longestButtonLabel)

	for i, oldButton := range gameButtons {
		newButton, _ := oldButton.Update(widthMsg)
		gameButtons[i] = newButton.(button.Model)
	}

	return Model{
		Buttons:       gameButtons,
		SelectedIndex: 0,
		Text:          "Game mode - Type the text here!",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Select the nearest button upwards
		case tea.KeyUp.String(), "k":
			// For single button, do nothing
		// Select the nearest button downwards
		case tea.KeyDown.String(), "j":
			// For single button, do nothing
		}

	case button.FocusMsg:
		// Pass focus message to all buttons
		for i := range m.Buttons {
			var cmd tea.Cmd
			model, cmd := m.Buttons[i].Update(msg)
			m.Buttons[i] = model.(button.Model)
			if cmd != nil {
				return m, cmd
			}
		}

	case button.ButtonClickedMsg:
		// Handle button clicks to navigate to different views
		switch msg.Index {
		case 0: // Back
			return m, view.NavigateTo(view.Home)
		}
	}

	// Pass all other messages to buttons
	for i := range m.Buttons {
		var cmd tea.Cmd
		model, cmd := m.Buttons[i].Update(msg)
		m.Buttons[i] = model.(button.Model)
		if cmd != nil {
			return m, cmd
		}
	}

	return m, nil
}

var gameStyle lipgloss.Style = lipgloss.NewStyle().
	AlignVertical(lipgloss.Center).
	AlignHorizontal(lipgloss.Center)

var gameTextStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FF0000")).
	AlignHorizontal(lipgloss.Center).
	MarginBottom(2)

func (m Model) View() string {
	var content []string
	content = append(content, gameTextStyle.Render(m.Text))
	for i, button := range m.Buttons {
		prefix := "  "
		if i == m.SelectedIndex {
			prefix = "> "
		}
		content = append(content, prefix+button.View())
	}

	return gameStyle.Render(content...)
}
