package input

import "github.com/charmbracelet/bubbletea"

// ShowMsg shows the input
type ShowMsg struct{}

// HideMsg hides the input
type HideMsg struct{}

// SubmitMsg is sent when input is submitted
type SubmitMsg struct {
	Value string
}

// ShowInput creates a command to show the input
func ShowInput() tea.Msg {
	return ShowMsg{}
}

// HideInput creates a command to hide the input
func HideInput() tea.Msg {
	return HideMsg{}
}
