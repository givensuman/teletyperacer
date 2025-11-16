// Package overlay provides a component
// which manages overlaying two components on top of each other.
//
// Adapted from: https://github.com/rmhubbert/bubbletea-overlay
package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Position represents a relative offset in the TUI.
type Position int

const (
	Top Position = iota
	Right
	Bottom
	Left
	Center
)

// Model implements tea.Model, and manages calculating and compositing
// the overlay UI from the backbround and foreground models.
type Model struct {
	Foreground tea.Model
	Background tea.Model
	XPosition  Position
	YPosition  Position
	XOffset    int
	YOffset    int
}

// New instantiates a new overlay model.
func New(foreground, background tea.Model, xPos, yPos Position, xOff, yOff int) *Model {
	return &Model{
		Foreground: foreground,
		Background: background,
		XPosition:  xPos,
		YPosition:  yPos,
		XOffset:    xOff,
		YOffset:    yOff,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View applies the compositing and handles rendering the view.
func (m *Model) View() string {
	if m.Foreground == nil && m.Background == nil {
		return ""
	}
	if m.Foreground == nil && m.Background != nil {
		return m.Background.View()
	}
	if m.Foreground != nil && m.Background == nil {
		return m.Foreground.View()
	}

	return composite(
		m.Foreground.View(),
		m.Background.View(),
		m.XPosition,
		m.YPosition,
		m.XOffset,
		m.YOffset,
	)
}
