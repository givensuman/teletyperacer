package host

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/internal/components/button"
	"github.com/givensuman/teletyperacer/client/internal/constants/view"
)

type OptionsModel struct {
	Buttons       [2]button.Model
	SelectedIndex int
}

// Model implements tea.Model
var _ tea.Model = OptionsModel{}

var buttons [2]button.Model = [2]button.Model{
	button.New(0, "Start Hosting"),
	button.New(1, "Back"),
}

func NewOptions() OptionsModel {
	focusedFirstButton, _ := buttons[0].Update(button.FocusMsg{0, -1})
	buttons[0] = focusedFirstButton.(button.Model)

	longestButtonLabel := len(buttons[0].Label)
	for _, button := range buttons[1:] {
		longestButtonLabel = max(longestButtonLabel, len(button.Label))
	}
	widthMsg := button.WidthMsg(longestButtonLabel)

	for i, oldButton := range buttons {
		newButton, _ := oldButton.Update(widthMsg)
		buttons[i] = newButton.(button.Model)
	}

	return OptionsModel{
		Buttons:       buttons,
		SelectedIndex: 0,
	}
}

func (m OptionsModel) Init() tea.Cmd {
	return nil
}

func (m OptionsModel) thereAreUsableButtons() bool {
	for _, button := range m.Buttons {
		if !button.IsDisabled {
			return true
		}
	}

	return false
}

func (m OptionsModel) indexOfNearestButtonUpwards() int {
	if !m.thereAreUsableButtons() {
		return -1
	}

	currentIndex := m.SelectedIndex
	for {
		if currentIndex == 0 {
			currentIndex = len(m.Buttons) - 1
		} else {
			currentIndex--
		}

		if !m.Buttons[currentIndex].IsDisabled {
			return currentIndex
		}

		if currentIndex == m.SelectedIndex {
			break
		}
	}

	return m.SelectedIndex
}

func (m OptionsModel) indexOfNearestButtonDownwards() int {
	if !m.thereAreUsableButtons() {
		return -1
	}

	currentIndex := m.SelectedIndex
	for {
		if currentIndex == len(m.Buttons)-1 {
			currentIndex = 0
		} else {
			currentIndex++
		}

		if !m.Buttons[currentIndex].IsDisabled {
			return currentIndex
		}

		if currentIndex == m.SelectedIndex {
			break
		}
	}

	return m.SelectedIndex
}

func (m OptionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Select the nearest button upwards
		case tea.KeyUp.String(), "k":
			prevIndex := m.SelectedIndex
			m.SelectedIndex = m.indexOfNearestButtonUpwards()
			return m, button.NewFocusCmd(m.SelectedIndex, prevIndex)

		// Select the nearest button downwards
		case tea.KeyDown.String(), "j":
			prevIndex := m.SelectedIndex
			m.SelectedIndex = m.indexOfNearestButtonDownwards()
			return m, button.NewFocusCmd(m.SelectedIndex, prevIndex)
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
		case 0: // Start Hosting
			return m, view.NavigateTo(view.Lobby) // For now, go to lobby after hosting
		case 1: // Back
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

var style lipgloss.Style = lipgloss.NewStyle().
	AlignVertical(lipgloss.Center).
	AlignHorizontal(lipgloss.Center)

func (m OptionsModel) View() string {
	var buttonViews []string
	for i, button := range m.Buttons {
		prefix := "  "
		if i == m.SelectedIndex {
			prefix = "> "
		}
		buttonViews = append(buttonViews, prefix+button.View())
	}

	return style.Render(buttonViews...)
}
