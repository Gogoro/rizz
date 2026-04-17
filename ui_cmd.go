package main

import (
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

// updateCmdInput handles keys while the user is typing a : command.
func (m *model) updateCmdInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c", "esc":
		m.cmdInput = false
		m.cmdBuffer = ""
		return m, nil
	case "enter":
		return m.runCmd(m.cmdBuffer)
	case "backspace":
		if len(m.cmdBuffer) > 0 {
			runes := []rune(m.cmdBuffer)
			m.cmdBuffer = string(runes[:len(runes)-1])
		}
		return m, nil
	}
	runes := []rune(key)
	if len(runes) == 1 && unicode.IsPrint(runes[0]) {
		m.cmdBuffer += key
	}
	return m, nil
}

func (m *model) runCmd(cmd string) (tea.Model, tea.Cmd) {
	cmd = strings.TrimSpace(cmd)
	m.cmdInput = false
	m.cmdBuffer = ""

	switch cmd {
	case "":
		return m, nil
	case "q", "quit", "exit":
		_ = m.state.Save()
		return m, tea.Quit
	case "h", "help":
		m.showHelp = true
	case "w", "write":
		m.cmdError = "rizz doesn't write. we only review."
	case "wq", "x":
		m.cmdError = "rizz doesn't write — try :q instead."
	case "a", "all":
		m.state.MarkAllViewed(m.files)
		_ = m.state.Save()
	case "r", "reset":
		m.state.UnmarkAll()
		_ = m.state.Save()
	default:
		m.cmdError = "unknown command: " + cmd
	}
	return m, nil
}
