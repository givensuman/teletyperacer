package practice

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

var practiceButtons [1]button.Model = [1]button.Model{
	button.New(0, "Back"),
}

func New() Model {
	focusedFirstButton, _ := practiceButtons[0].Update(button.FocusMsg{0, -1})
	practiceButtons[0] = focusedFirstButton.(button.Model)

	longestButtonLabel := len(practiceButtons[0].Label)
	widthMsg := button.WidthMsg(longestButtonLabel)

	for i, oldButton := range practiceButtons {
		newButton, _ := oldButton.Update(widthMsg)
		practiceButtons[i] = newButton.(button.Model)
	}

	return Model{
		Buttons:       practiceButtons,
		SelectedIndex: 0,
		Text:          "Practice mode - Type the text here!",
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

var practiceStyle lipgloss.Style = lipgloss.NewStyle().
	AlignVertical(lipgloss.Center).
	AlignHorizontal(lipgloss.Center)

var textStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00FF00")).
	AlignHorizontal(lipgloss.Center).
	MarginBottom(2)

func (m Model) View() string {
	var content []string
	content = append(content, textStyle.Render(m.Text))
	for i, button := range m.Buttons {
		prefix := "  "
		if i == m.SelectedIndex {
			prefix = "> "
		}
		content = append(content, prefix+button.View())
	}

	return practiceStyle.Render(content...)
}
