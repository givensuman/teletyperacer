package main

import (
	"github.com/charmbracelet/bubbletea"

	"github.com/givensuman/teletyperacer/client/internal/tui/root"
)

// https://github.com/givensuman/teletyperacer
func main() {
	p := tea.NewProgram(
		root.New(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
