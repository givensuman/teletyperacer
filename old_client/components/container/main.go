// Package container provides a component
// which manages wrapping an entire view.
package container

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/teletyperacer/client/types"
)

var style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Padding(2)

// Model represents the container's state.
type Model struct {
	types.SingleParent[tea.Model]
	Width  int
	Height int
}

// New instantiates a new container model.
func New(child tea.Model) Model {
	m := Model{
		Width:  0,
		Height: 0,
	}

	m.Child = child

	return m
}

// Init initializes the container
func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

// Update handles messages for the container.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Child, cmd = m.Child.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Automatically handle window resizing.
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, cmd
}

// View renders the container.
func (m Model) View() string {
	return style.
		Width(m.Width).
		Height(m.Height).
		Render(m.Child.View())
}
