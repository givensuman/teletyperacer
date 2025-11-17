package button

import "github.com/charmbracelet/bubbletea"

type FocusMsg struct {
	Focus   int
	Unfocus int
}

// WidthMsg dictates the button's width
type WidthMsg int

type DisableMsg bool

const (
	// Disable the button
	Disable DisableMsg = true
	// Enable the button
	Enable DisableMsg = false
)

// ButtonClickedMsg is sent when a button is clicked
type ButtonClickedMsg struct {
	Index int
}

// NewFocusCmd creates a command to focus/unfocus buttons
func NewFocusCmd(focus, unfocus int) tea.Cmd {
	return func() tea.Msg {
		return FocusMsg{Focus: focus, Unfocus: unfocus}
	}
}
