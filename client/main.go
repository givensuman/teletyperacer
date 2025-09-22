package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// View defines the application views which may
// be rendered at any given time.
type View int64

const (
	None View = iota // Identity case
	Home
	Host
	Join
	Lobby
	Play
)

// Model represents application state for the TUI.
type Model struct {
	currentView View
}

func start() Model {
	return Model{
		currentView: Home,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			m.currentView = Home
		case "o":
			m.currentView = Host
		case "j":
			m.currentView = Join
		case "l":
			m.currentView = Lobby
		case "p":
			m.currentView = Play
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.currentView {
	case Home:
		return "Home View\nPress 'o' to Host, 'j' to Join, 'q' to Quit"
	case Host:
		return "Host View\nPress 'h' for Home, 'j' to Join, 'q' to Quit"
	case Join:
		return "Join View\nPress 'h' for Home, 'o' to Host, 'q' to Quit"
	case Lobby:
		return "Lobby View\nPress 'h' for Home, 'p' to Play, 'q' to Quit"
	case Play:
		return "Play View\nPress 'h' for Home, 'l' for Lobby, 'q' to Quit"
	default:
		return "Unknown View\nPress 'h' for Home, 'q' to Quit"
	}
}

// https://github.com/givensuman/teletyperacer
func main() {
	p := tea.NewProgram(start())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
