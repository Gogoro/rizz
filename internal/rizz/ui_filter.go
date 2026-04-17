package rizz

import (
	"unicode"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

// updateFilterInput handles keystrokes while the filter prompt is active.
func (m *model) updateFilterInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c":
		_ = m.state.Save()
		return m, tea.Quit
	case "esc":
		m.filter = ""
		m.filterInput = false
		m.recomputeVisible()
		m.cursor = 0
		m.refreshDiff()
		return m, nil
	case "enter":
		m.filterInput = false
		return m, nil
	case "backspace":
		if len(m.filter) > 0 {
			runes := []rune(m.filter)
			m.filter = string(runes[:len(runes)-1])
			m.recomputeVisible()
			m.cursor = 0
			m.refreshDiff()
		}
		return m, nil
	}

	runes := []rune(key)
	if len(runes) == 1 && unicode.IsPrint(runes[0]) {
		m.filter += key
		m.recomputeVisible()
		m.cursor = 0
		m.refreshDiff()
	}
	return m, nil
}

func renderFilterPrompt(filter string, editing bool, width int) string {
	icon := "/"
	text := filter
	if editing {
		text += "▎" // thin cursor bar
	}

	var style lipgloss.Style
	if editing {
		style = lipgloss.NewStyle().Foreground(colorGold).Bold(true)
	} else {
		style = lipgloss.NewStyle().Foreground(colorGoldSoft)
	}

	// Truncate if it somehow overflows
	content := icon + " " + text
	if lipgloss.Width(content) > width {
		content = content[:width]
	}
	return style.Render(content)
}
