// Package home describes the home view
package home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/internal/components/button"
	"github.com/givensuman/teletyperacer/client/internal/constants/view"
)

type Model struct {
	Buttons       [4]button.Model
	SelectedIndex int
}

// Model implements tea.Model
var _ tea.Model = Model{}

func New() Model {
	buttons := [4]button.Model{
		button.New(0, "Join"),
		button.New(1, "Host"),
		button.New(2, "Practice"),
		button.New(3, "Quit"),
	}

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

	return Model{
		Buttons:       buttons,
		SelectedIndex: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) thereAreUsableButtons() bool {
	for _, button := range m.Buttons {
		if !button.IsDisabled {
			return true
		}
	}

	return false
}

func (m Model) indexOfNearestButtonUpwards() int {
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

func (m Model) indexOfNearestButtonDownwards() int {
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case 0: // Join
			return m, view.NavigateTo(view.Lobby)
		case 1: // Host
			return m, view.NavigateTo(view.HostOptions)
		case 2: // Practice
			return m, view.NavigateTo(view.Practice)
		case 3: // Quit
			return m, tea.Quit
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

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FF69B4")).
	AlignHorizontal(lipgloss.Center).
	MarginBottom(2)

var instructionsStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#888888")).
	AlignHorizontal(lipgloss.Center).
	MarginTop(2)

func (m Model) View() string {
	var content []string
	content = append(content, titleStyle.Render("TeleType Racer"))
	for i, button := range m.Buttons {
		prefix := "  "
		if i == m.SelectedIndex {
			prefix = "> "
		}
		content = append(content, prefix+button.View())
	}
	content = append(content, instructionsStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Ctrl+C to quit"))

	return style.Render(content...)
}
