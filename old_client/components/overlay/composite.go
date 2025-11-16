package overlay

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// composite merges and flattens the background and foreground views into a single view.
//
// This implementation is based off of the one used by Superfile:
// https://github.com/yorukot/superfile/blob/main/src/pkg/string_function/overplace.go
func composite(foreground, background string, xPos, yPos Position, xOff, yOff int) string {
	fgWidth, fgHeight := lipgloss.Size(foreground)
	bgWidth, bgHeight := lipgloss.Size(background)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return foreground
	}

	x, y := offsets(foreground, background, xPos, yPos, xOff, yOff)
	x = clamp(x, 0, bgWidth-fgWidth)
	y = clamp(y, 0, bgHeight-fgHeight)

	fgLines := lines(foreground)
	bgLines := lines(background)
	var sb strings.Builder

	for i, bgLine := range bgLines {
		if i > 0 {
			sb.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			sb.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := ansi.Truncate(bgLine, x, "")
			pos = ansi.StringWidth(left)
			sb.WriteString(left)
			if pos < x {
				sb.WriteString(whitespace(x - pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		sb.WriteString(fgLine)

		pos += ansi.StringWidth(fgLine)
		right := ansi.TruncateLeft(bgLine, pos, "")
		bgWidth := ansi.StringWidth(bgLine)
		rightWidth := ansi.StringWidth(right)
		if rightWidth <= bgWidth-pos {
			sb.WriteString(whitespace(bgWidth - rightWidth - pos))
		}

		sb.WriteString(right)
	}

	return sb.String()
}

// offsets calculates the actual vertical and horizontal offsets used to position the foreground
// tea.Model relative to the background tea.Model.
func offsets(foreground, background string, xPos, yPos Position, xOff, yOff int) (int, int) {
	var x, y int
	switch xPos {
	case Center:
		halfBackgroundWidth := lipgloss.Width(background) / 2
		halfForegroundWidth := lipgloss.Width(foreground) / 2
		x = halfBackgroundWidth - halfForegroundWidth
	case Right:
		x = lipgloss.Width(background) - lipgloss.Width(foreground)
	}

	switch yPos {
	case Center:
		halfBackgroundHeight := lipgloss.Height(background) / 2
		halfForegroundHeight := lipgloss.Height(foreground) / 2
		y = halfBackgroundHeight - halfForegroundHeight
	case Bottom:
		y = lipgloss.Height(background) - lipgloss.Height(foreground)
	}

	return x + xOff, y + yOff
}

// clamp calculates the lowest possible number between the given boundaries.
func clamp(v, lower, upper int) int {
	if upper < lower {
		return min(max(v, upper), lower)
	}

	return min(max(v, lower), upper)
}

// lines normalises any non standard new lines within a string, and then splits and returns a slice
// of strings split on the new lines.
func lines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.Split(s, "\n")
}

// whitespace returns a string of whitespace characters of the requested length.
func whitespace(length int) string {
	return strings.Repeat(" ", length)
}
