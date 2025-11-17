package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	cursor  int
	choices []string
}

func NewHome() HomeModel {
	return HomeModel{
		cursor:  0,
		choices: []string{"Join", "Host", "Practice", "Quit"},
	}
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			switch m.cursor {
			case 0: // Join
				return m, nil
			case 1: // Host
				return m, func() tea.Msg { return ScreenChangeMsg{Screen: HostScreen} }
			case 2: // Practice
				return m, func() tea.Msg { return ScreenChangeMsg{Screen: PracticeScreen} }
			case 3: // Quit
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m HomeModel) View() string {
	s := "Welcome to TeleTypeRacer!\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += cursor + " " + choice + "\n"
	}

	s += "\nUse j/k or up/down to navigate, enter to select."

	return lipgloss.NewStyle().Padding(1).Render(s)
}
