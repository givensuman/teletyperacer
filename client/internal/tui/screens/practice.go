package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/teletyperacer/client/internal/components/typing"
)

type PracticeModel struct {
	typing typing.Model
}

func NewPractice() PracticeModel {
	// Sample text for practice
	sampleText := "The quick brown fox jumps over the lazy dog. This is a sample text for typing practice. Try to type as accurately and quickly as possible."
	return PracticeModel{
		typing: typing.NewTyping(sampleText),
	}
}

func (m PracticeModel) Init() tea.Cmd {
	return m.typing.Init()
}

func (m PracticeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	updatedTyping, cmd := m.typing.Update(msg)
	m.typing = updatedTyping.(typing.Model)
	return m, cmd
}

func (m PracticeModel) View() string {
	typingView := m.typing.View()
	if m.typing.IsCompleted() {
		return lipgloss.NewStyle().Padding(1).Render(typingView + "\n\nPress ESC to go back to Home.")
	}

	wpm := m.typing.GetWPM()
	accuracy := m.typing.GetAccuracy()
	progress := m.typing.GetProgress()

	stats := fmt.Sprintf("WPM: %.1f | Accuracy: %.1f%% | Progress: %.1f%%", wpm, accuracy, progress)
	return lipgloss.NewStyle().Padding(1).Render("Practice Screen\n\n" + stats + "\n\n" + typingView + "\n\nPress ESC to go back to Home.")
}
