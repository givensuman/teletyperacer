// Package view defines a shared API for
// the rendering of different views
package view

import "github.com/charmbracelet/bubbletea"

type View int8

const (
	Home View = iota
	PracticeOptions
	Practice
	HostOptions
	Lobby
	Game
)

type NavigateMsg View

func NavigateTo(view View) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg(view)
	}
}
