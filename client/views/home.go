package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// View defines the application views which may
// be rendered at any given time.
type View int64

const (
	None View = iota // Identity case
	Home
	Practice
	Host
	Join
	Lobby
	Play
)

// Model represents application state for the TUI.
type Model struct {
	CurrentView View
	Width       int
	Height      int
}

func HomeView(m Model) string {
	title := `
'########:'########:'##:::::::'########:'########:'##:::'##:'########::'########:'########:::::'###:::::'######::'########:'########::
... ##..:: ##.....:: ##::::::: ##.....::... ##..::. ##:'##:: ##.... ##: ##.....:: ##.... ##:::'## ##:::'##... ##: ##.....:: ##.... ##:
::: ##:::: ##::::::: ##::::::: ##:::::::::: ##:::::. ####::: ##:::: ##: ##::::::: ##:::: ##::'##:. ##:: ##:::..:: ##::::::: ##:::: ##:
::: ##:::: ######::: ##::::::: ######:::::: ##::::::. ##:::: ########:: ######::: ########::'##:::. ##: ##::::::: ######::: ########::
::: ##:::: ##...:::: ##::::::: ##...::::::: ##::::::: ##:::: ##.....::: ##...:::: ##.. ##::: #########: ##::::::: ##...:::: ##.. ##:::
::: ##:::: ##::::::: ##::::::: ##:::::::::: ##::::::: ##:::: ##:::::::: ##::::::: ##::. ##:: ##.... ##: ##::: ##: ##::::::: ##::. ##::
::: ##:::: ########: ########: ########:::: ##::::::: ##:::: ##:::::::: ########: ##:::. ##: ##:::: ##:. ######:: ########: ##:::. ##:
:::..:::::........::........::........:::::..::::::::..:::::..:::::::::........::..:::::..::..:::::..:::......:::........::..:::::..::
`

	options := []string{
		"[P] Practice",
		"[H] Host Game",
		"[J] Join Game",
	}
	centeredOptions := centerText(strings.Join(options, "  "), m.Width)

	// Calculate vertical centering
	totalHeight := 5 // title + subtitle + 3 options + footer
	topPadding := max(0, (m.Height-totalHeight)/2)

	padding := strings.Repeat("\n", topPadding)

	return padding +
		title + "\n\n" +
		centeredOptions + "\n\n"
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			m.CurrentView = Home
		case "p":
			m.CurrentView = Practice
		case "o":
			m.CurrentView = Host
		case "j":
			m.CurrentView = Join
		case "l":
			m.CurrentView = Lobby
		case "y":
			m.CurrentView = Play
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.CurrentView {
	case Home:
		return HomeView(m)
	case Practice:
		return "Practice View\nPress 'h' for Home, 'q' to Quit"
	case Host:
		return "Host View\nPress 'h' for Home, 'j' to Join, 'q' to Quit"
	case Join:
		return "Join View\nPress 'h' for Home, 'o' to Host, 'q' to Quit"
	case Lobby:
		return "Lobby View\nPress 'h' for Home, 'y' to Play, 'q' to Quit"
	case Play:
		return "Play View\nPress 'h' for Home, 'l' for Lobby, 'q' to Quit"
	default:
		return "Unknown View\nPress 'h' for Home, 'q' to Quit"
	}
}
