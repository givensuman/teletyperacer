// Package tui controls rendering
// of the TUI
package tui

type Screen int

const (
	HomeScreen Screen = iota
	HostScreen
	PracticeScreen
)

type ScreenChangeMsg struct {
	Screen Screen
}

type ConnectionStatus int

const (
	Connecting ConnectionStatus = iota
	Connected
	Failed
)

type ConnectionStatusMsg struct {
	Status ConnectionStatus
}
