package main

import (
	"client/views"
	tea "github.com/charmbracelet/bubbletea"
)

func start() views.Model {
	return views.Model{
		CurrentView: views.Home,
		Width:       80,
		Height:      24,
	}
}

// https://github.com/givensuman/teletyperacer
func main() {
	p := tea.NewProgram(start(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
