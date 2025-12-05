// Package typing defines the typing component
package typing

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/muesli/reflow/wordwrap"
)

type Styles struct {
	correct  lipgloss.Style
	toEnter  lipgloss.Style
	mistakes lipgloss.Style
	cursor   lipgloss.Style
}

type Model struct {
	text        string       // The full text to type
	runes       []rune       // Text as runes for easier manipulation
	inputBuffer []rune       // What the user has typed
	mistakes    map[int]bool // Positions where mistakes were made
	cursor      int          // Current position in the text
	styles      Styles       // Styles for different text segments
	width       int
	height      int
	completed   bool
	startTime   time.Time // When typing started
	lastKeyTime time.Time // Last key press time
	wpm         float64   // Current WPM
	wpmHistory  []float64 // WPM over time for graphing
}

var _ tea.Model = Model{}

func NewTyping(text string) Model {
	runes := []rune(text)
	now := time.Now()
	return Model{
		text:        text,
		runes:       runes,
		inputBuffer: []rune{},
		mistakes:    make(map[int]bool),
		cursor:      0,
		styles: Styles{
			correct:  lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)),                  // White
			toEnter:  lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(240)).Faint(true),     // Faint gray
			mistakes: lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1)).Underline(true),   // Red underlined
			cursor:   lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(240)).Underline(true), // Faint underlined
		},
		width:       80,
		height:      10,
		completed:   false,
		startTime:   now,
		lastKeyTime: now,
		wpm:         0,
		wpmHistory:  []float64{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyBackspace:
			if len(m.inputBuffer) > 0 {
				// Remove the last character
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				m.cursor--
				// Remove mistake if it was at this position
				delete(m.mistakes, m.cursor)
				m.lastKeyTime = time.Now()
				m.updateWPM()
				m.wpmHistory = append(m.wpmHistory, m.wpm)
			}
		case tea.KeySpace:
			if m.cursor < len(m.runes) {
				typed := ' '
				expected := m.runes[m.cursor]

				m.inputBuffer = append(m.inputBuffer, typed)
				if typed != expected {
					m.mistakes[m.cursor] = true
				}
				m.cursor++
				m.lastKeyTime = time.Now()
				m.updateWPM()
				m.wpmHistory = append(m.wpmHistory, m.wpm)

				if m.cursor >= len(m.runes) {
					m.completed = true
				}
			}
		case tea.KeyRunes:
			if m.cursor < len(m.runes) {
				typed := msg.Runes[0]
				expected := m.runes[m.cursor]

				m.inputBuffer = append(m.inputBuffer, typed)
				if typed != expected {
					m.mistakes[m.cursor] = true
				}
				m.cursor++
				m.lastKeyTime = time.Now()
				m.updateWPM()
				m.wpmHistory = append(m.wpmHistory, m.wpm)

				if m.cursor >= len(m.runes) {
					m.completed = true
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	if m.completed {
		// Display results with WPM graph
		graph := asciigraph.Plot(m.wpmHistory, asciigraph.Height(10), asciigraph.Width(m.width))
		wpm := m.GetWPM()
		accuracy := m.GetAccuracy()
		return fmt.Sprintf("Typing completed!\n\nWPM: %.1f\nAccuracy: %.1f%%\n\nWPM Graph:\n%s", wpm, accuracy, graph)
	}

	// Render the text with colors
	coloredText := m.renderText()

	// Use reflow for wrapping
	wrapped := wordwrap.String(coloredText, m.width)

	return wrapped
}

func (m Model) renderText() string {
	var result strings.Builder

	// Color typed text
	mistakePositions := make([]int, 0, len(m.mistakes))
	for pos := range m.mistakes {
		mistakePositions = append(mistakePositions, pos)
	}
	sort.Ints(mistakePositions)

	if len(mistakePositions) == 0 {
		result.WriteString(styleAllRunes(m.inputBuffer, m.styles.correct))
	} else {
		prevMistake := -1
		for _, pos := range mistakePositions {
			if pos >= len(m.inputBuffer) {
				break
			}
			// Correct text before mistake
			correctSlice := m.inputBuffer[prevMistake+1 : pos]
			result.WriteString(styleAllRunes(correctSlice, m.styles.correct))
			// Mistake
			result.WriteString(m.styles.mistakes.Render(string(m.inputBuffer[pos])))
			prevMistake = pos
		}
		// Remaining correct text
		remaining := m.inputBuffer[prevMistake+1:]
		result.WriteString(styleAllRunes(remaining, m.styles.correct))
	}

	// Color cursor position
	if m.cursor < len(m.runes) {
		cursorChar := string(m.runes[m.cursor])
		result.WriteString(m.styles.cursor.Render(cursorChar))
	}

	// Color remaining text
	if m.cursor+1 < len(m.runes) {
		remaining := string(m.runes[m.cursor+1:])
		result.WriteString(m.styles.toEnter.Render(remaining))
	}

	return result.String()
}

func styleAllRunes(runes []rune, style lipgloss.Style) string {
	var acc strings.Builder
	for _, char := range runes {
		acc.WriteString(style.Render(string(char)))
	}
	return acc.String()
}

// Getters for external use
func (m Model) IsCompleted() bool {
	return m.completed
}

func (m Model) GetAccuracy() float64 {
	if len(m.inputBuffer) == 0 {
		return 0
	}
	correct := len(m.inputBuffer) - len(m.mistakes)
	return float64(correct) / float64(len(m.inputBuffer)) * 100
}

func (m Model) GetProgress() float64 {
	if len(m.runes) == 0 {
		return 0
	}
	return float64(m.cursor) / float64(len(m.runes)) * 100
}

func (m Model) GetWPM() float64 {
	return m.wpm
}

func (m *Model) updateWPM() {
	if m.cursor == 0 {
		m.wpm = 0
		return
	}

	elapsed := m.lastKeyTime.Sub(m.startTime)
	if elapsed.Seconds() == 0 {
		m.wpm = 0
		return
	}

	// WPM = (characters typed / 5) / (time in minutes)
	// Using net characters typed (correct - incorrect)
	netChars := len(m.inputBuffer) - len(m.mistakes)
	m.wpm = (float64(netChars) / 5.0) / (elapsed.Minutes())
}
