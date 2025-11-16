// Package home describes the home view
package home

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/givensuman/teletyperacer/client/internal/components/button"
)

type Model struct {
	Buttons       [4]button.Model
	SelectedIndex int
}

// Model implements tea.Model
var _ tea.Model = Model{}

var buttons [4]button.Model = [4]button.Model{
	button.New(0, "Join", func() {}),
	button.New(1, "Host", func() {}),
	button.New(2, "Practice", func() {}),
	button.New(3, "Quit", func() {}),
}

func New() Model {
	buttons[0].Update(button.Focus)

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
			var prevIndex = m.SelectedIndex
			m.SelectedIndex = m.indexOfNearestButtonUpwards()

			m.Buttons[prevIndex].Update(button.Unfocus)
			m.Buttons[m.SelectedIndex].Update(button.Focus)

		// Select the nearest button downwards
		case tea.KeyDown.String(), "j":
			var prevIndex = m.SelectedIndex
			m.SelectedIndex = m.indexOfNearestButtonDownwards()

			m.Buttons[prevIndex].Update(button.Unfocus)
			m.Buttons[m.SelectedIndex].Update(button.Focus)
		}
	}

	return m, nil
}

func (m Model) View() string {
	style := lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center)

	var buttonViews []string
	for _, button := range m.Buttons {
		buttonViews = append(buttonViews, button.View())
	}

	buttonViews = append(buttonViews, fmt.Sprintf("selectedindex=%d", m.SelectedIndex))
	buttonViews = append(buttonViews, fmt.Sprintf("nearestbuttondown=%d", m.indexOfNearestButtonDownwards()))
	buttonViews = append(buttonViews, fmt.Sprintf("nearestbuttonup=%d", m.indexOfNearestButtonUpwards()))

	return style.Render(buttonViews...)
}
