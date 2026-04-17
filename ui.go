package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listPaneWidth = 36
	statusHeight  = 1
	headerHeight  = 1
)

type focusMode int

const (
	focusList focusMode = iota
	focusDiff
)

type model struct {
	files    []FileDiff
	cursor   int
	focus    focusMode
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
	return m.width - listPaneWidth - 1, m.height - statusHeight - headerHeight
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
		key := msg.String()

		// keys that work regardless of which pane has focus
		switch key {
		case "q", "ctrl+c":
			_ = m.state.Save()
			return m, tea.Quit
		case "v", " ":
			f := m.files[m.cursor]
			m.state.ToggleViewed(f.Path, f.Hash)
			_ = m.state.Save()
			return m, nil
		case "a":
			m.state.MarkAllViewed(m.files)
			_ = m.state.Save()
			return m, nil
		case "r":
			m.state.UnmarkAll()
			_ = m.state.Save()
			return m, nil
		case "n", "tab":
			if m.cursor < len(m.files)-1 {
				m.cursor++
				m.refreshDiff()
			}
			return m, nil
		case "p", "shift+tab":
			if m.cursor > 0 {
				m.cursor--
				m.refreshDiff()
			}
			return m, nil
		}

		if m.focus == focusList {
			switch key {
			case "j", "down":
				if m.cursor < len(m.files)-1 {
					m.cursor++
					m.refreshDiff()
				}
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
					m.refreshDiff()
				}
			case "g", "home":
				m.cursor = 0
				m.refreshDiff()
			case "G", "end":
				m.cursor = len(m.files) - 1
				m.refreshDiff()
			case "enter", "l", "right":
				m.focus = focusDiff
			}
			return m, nil
		}

		// focusDiff
		switch key {
		case "esc", "h", "left":
			m.focus = focusList
		case "j", "down":
			m.viewport.LineDown(1)
		case "k", "up":
			m.viewport.LineUp(1)
		case "d", "pgdown":
			m.viewport.HalfViewDown()
		case "u", "pgup":
			m.viewport.HalfViewUp()
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

	listHeight := m.height - statusHeight - headerHeight
	listFocused := m.focus == focusList
	listContent := renderFileList(m.files, m.cursor, m.state, listPaneWidth-2, listHeight, listFocused)

	borderColor := colorBorder
	if !listFocused {
		borderColor = colorMuted
	}

	listPane := lipgloss.NewStyle().
		Width(listPaneWidth).
		Height(listHeight).
		BorderStyle(lipgloss.ThickBorder()).
		BorderRight(true).
		BorderForeground(borderColor).
		Padding(0, 1).
		Render(listContent)

	main := lipgloss.JoinHorizontal(lipgloss.Top, listPane, m.viewport.View())
	return lipgloss.JoinVertical(lipgloss.Left, m.renderHeader(), main, m.renderStatus())
}

func (m *model) renderHeader() string {
	brand := styleHeaderBar.Render("⛓  R I Z Z  ⛓")
	tagline := styleHeaderTagline.Render(" pure vibes ")

	used := lipgloss.Width(brand) + lipgloss.Width(tagline)
	chainLen := m.width - used
	if chainLen < 0 {
		chainLen = 0
	}
	chain := styleHeaderChain.Render(strings.Repeat("━", chainLen))

	return lipgloss.NewStyle().
		Background(colorHeader).
		Width(m.width).
		Render(brand + chain + tagline)
}

func (m *model) renderStatus() string {
	viewed := 0
	for _, f := range m.files {
		if m.state.IsViewed(f.Path, f.Hash) {
			viewed++
		}
	}
	progress := styleStatusAccent.Render(fmt.Sprintf("💎 %d/%d", viewed, len(m.files)))

	var help string
	if m.focus == focusList {
		help = "j/k file · enter open · v view · a all · r reset · q quit"
	} else {
		help = "j/k scroll · d/u half · g/G top/bot · esc back · v view · q quit"
	}
	helpRendered := styleStatusBar.Render(help)

	inner := m.width - 2
	gap := inner - lipgloss.Width(progress) - lipgloss.Width(helpRendered)
	if gap < 1 {
		gap = 1
	}
	spacer := styleStatusBar.Render(strings.Repeat(" ", gap))
	return lipgloss.NewStyle().Background(colorStatus).Width(m.width).Render(progress + spacer + helpRendered)
}
