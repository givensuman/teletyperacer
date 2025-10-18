package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/givensuman/teletyperacer/client/views"
)

func start() views.Model {
	return views.Model{
		CurrentView: views.Home,
		Width:       80,
		Height:      24,
		Socket:      nil,
		RoomID:      "",
	}
}

// https://github.com/givensuman/teletyperacer
func main() {
	p := tea.NewProgram(start(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
