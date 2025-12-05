package roominput

import "github.com/charmbracelet/bubbletea"

// ShowMsg shows the room input
type ShowMsg struct{}

// HideMsg hides the room input
type HideMsg struct{}

// JoinRoomMsg is sent when joining a room
type JoinRoomMsg struct {
	Code string
}

// ShowRoomInput creates a command to show the room input
func ShowRoomInput() tea.Msg {
	return ShowMsg{}
}

// HideRoomInput creates a command to hide the room input
func HideRoomInput() tea.Msg {
	return HideMsg{}
}
