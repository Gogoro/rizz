package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listPaneWidth = 36
	statusHeight  = 1
)

type model struct {
	files    []FileDiff
	cursor   int
	state    *State
	viewport viewport.Model
	width    int
	height   int
	ready    bool
}

func Run(files []FileDiff, state *State) error {
	m := &model{files: files, state: state}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) diffPaneSize() (int, int) {
	return m.width - listPaneWidth - 1, m.height - statusHeight
}

func (m *model) refreshDiff() {
	if len(m.files) == 0 {
		return
	}
	w, _ := m.diffPaneSize()
	m.viewport.SetContent(renderDiff(m.files[m.cursor], w))
	m.viewport.GotoTop()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		w, h := m.diffPaneSize()
		if !m.ready {
			m.viewport = viewport.New(w, h)
			m.ready = true
		} else {
			m.viewport.Width = w
			m.viewport.Height = h
		}
		m.refreshDiff()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			_ = m.state.Save()
			return m, tea.Quit
		case "j", "right":
			m.viewport.LineDown(1)
		case "k", "left":
			m.viewport.LineUp(1)
		case "d", "pgdown":
			m.viewport.HalfViewDown()
		case "u", "pgup":
			m.viewport.HalfViewUp()
		case "n", "down", "tab":
			if m.cursor < len(m.files)-1 {
				m.cursor++
				m.refreshDiff()
			}
		case "p", "up", "shift+tab":
			if m.cursor > 0 {
				m.cursor--
				m.refreshDiff()
			}
		case "v", " ":
			f := m.files[m.cursor]
			m.state.ToggleViewed(f.Path, f.Hash)
			_ = m.state.Save()
		case "a":
			m.state.MarkAllViewed(m.files)
			_ = m.state.Save()
		case "r":
			m.state.UnmarkAll()
			_ = m.state.Save()
		case "g", "home":
			m.viewport.GotoTop()
		case "G", "end":
			m.viewport.GotoBottom()
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if !m.ready {
		return "loading..."
	}

	listHeight := m.height - statusHeight
	listContent := renderFileList(m.files, m.cursor, m.state, listPaneWidth-2, listHeight)

	listPane := lipgloss.NewStyle().
		Width(listPaneWidth).
		Height(listHeight).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(colorBorder).
		Padding(0, 1).
		Render(listContent)

	main := lipgloss.JoinHorizontal(lipgloss.Top, listPane, m.viewport.View())
	return lipgloss.JoinVertical(lipgloss.Left, main, m.renderStatus())
}

func (m *model) renderStatus() string {
	viewed := 0
	for _, f := range m.files {
		if m.state.IsViewed(f.Path, f.Hash) {
			viewed++
		}
	}
	progress := fmt.Sprintf("%d/%d viewed", viewed, len(m.files))
	help := "↑↓ file · ←→/jk scroll · v view · a all · r reset · g/G top/bot · q quit"

	inner := m.width - 2
	gap := inner - lipgloss.Width(progress) - lipgloss.Width(help)
	if gap < 1 {
		gap = 1
	}
	spacer := fmt.Sprintf("%*s", gap, "")
	return styleStatusBar.Width(m.width).Render(progress + spacer + help)
}
