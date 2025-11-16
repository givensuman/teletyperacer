package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/givensuman/teletyperacer/client/views"
)

// https://github.com/givensuman/teletyperacer
func main() {
	p := tea.NewProgram(
		views.New(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
