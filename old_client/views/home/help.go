package home

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Select, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys(tea.KeyUp.String(), "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys(tea.KeyDown.String(), "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys(         tea.KeyEnter.String()),
		key.WithHelp("↳", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", tea.KeyEsc.String(), tea.KeyCtrlC.String()),
		key.WithHelp("q", "quit"),
	),
}

type helpModel struct {
	keys     keyMap
	help     help.Model
	quitting bool
}

func newHelpModel() helpModel {
	return helpModel{
		keys:     keys,
		help:     help.New(),
		quitting: false,
	}
}

func (hm helpModel) Init() tea.Cmd {
	return nil
}

func (hm helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		hm.help.Width = msg.Width
	}

	return hm, nil
}

func (hm helpModel) View() string {
	return hm.help.View(hm.keys)
}
