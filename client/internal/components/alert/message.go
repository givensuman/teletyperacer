package alert

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/givensuman/teletyperacer/client/internal/components/button"
)

// ShowMsg shows the alert
type ShowMsg struct {
	Title   string
	Message string
	Buttons []button.Model
}

// HideMsg hides the alert
type HideMsg struct{}

// ShowAlert creates a command to show the alert
func ShowAlert(title, message string, buttons ...button.Model) tea.Cmd {
	return func() tea.Msg {
		return ShowMsg{
			Title:   title,
			Message: message,
			Buttons: buttons,
		}
	}
}

// HideAlert creates a command to hide the alert
func HideAlert() tea.Msg {
	return HideMsg{}
}
