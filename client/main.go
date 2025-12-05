package main

import (
	"github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/givensuman/teletyperacer/client/internal/tui"
)

// https://github.com/givensuman/teletyperacer
func main() {
	zone.NewGlobal()

	p := tea.NewProgram(
		root.New(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
