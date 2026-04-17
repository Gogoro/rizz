package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	statusHeight = 1
	headerHeight = 1

	listWidthMin     = 28 // enough room for a few filename chars + counts
	listWidthMax     = 56 // don't let the sidebar dwarf the diff
	minDiffWidth     = 40 // the diff pane always keeps at least this much
	listWidthPercent = 28 // target fraction of total width
)

type focusMode int

const (
	focusList focusMode = iota
	focusDiff
)

type model struct {
	files       []FileDiff
	visible     []int // indices into files after applying filter
	cursor      int   // index into visible
	focus       focusMode
	state       *State
	viewport    viewport.Model
	width       int
	height      int
	ready       bool
	showHelp    bool
	filter      string // current file path filter substring (case-insensitive)
	filterInput bool   // true while the user is actively typing the filter
	cmdInput    bool   // true while a vim-style : command is being typed
	cmdBuffer   string // current command text
	cmdError    string // feedback message shown after a command runs
	showCommit  bool   // commit message suggestions overlay
}

func Run(files []FileDiff, state *State) error {
	m := &model{files: files, state: state}
	m.recomputeVisible()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

func (m *model) recomputeVisible() {
	m.visible = m.visible[:0]
	needle := strings.ToLower(m.filter)
	for i, f := range m.files {
		if needle == "" || strings.Contains(strings.ToLower(f.Path), needle) {
			m.visible = append(m.visible, i)
		}
	}
	if m.cursor >= len(m.visible) {
		m.cursor = 0
	}
}

func (m *model) currentFile() *FileDiff {
	if len(m.visible) == 0 || m.cursor < 0 || m.cursor >= len(m.visible) {
		return nil
	}
	return &m.files[m.visible[m.cursor]]
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) listWidth() int {
	w := m.width * listWidthPercent / 100
	if w < listWidthMin {
		w = listWidthMin
	}
	if w > listWidthMax {
		w = listWidthMax
	}
	// make sure the diff pane always has breathing room
	if m.width-w-1 < minDiffWidth {
		w = m.width - minDiffWidth - 1
	}
	if w < 1 {
		w = 1
	}
	return w
}

func (m *model) diffPaneSize() (int, int) {
	return m.width - m.listWidth() - 1, m.height - statusHeight - headerHeight
}

func (m *model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.showHelp || m.filterInput {
		return m, nil
	}

	listW := m.listWidth()
	inSidebar := msg.X < listW
	inMainArea := msg.Y >= headerHeight && msg.Y < m.height-statusHeight

	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if inSidebar {
			if m.cursor > 0 {
				m.cursor--
				m.refreshDiff()
			}
		} else {
			m.viewport.LineUp(3)
		}
		return m, nil
	case tea.MouseButtonWheelDown:
		if inSidebar {
			if m.cursor < len(m.visible)-1 {
				m.cursor++
				m.refreshDiff()
			}
		} else {
			m.viewport.LineDown(3)
		}
		return m, nil
	case tea.MouseButtonLeft:
		if msg.Action != tea.MouseActionPress {
			return m, nil
		}
		if !inMainArea {
			return m, nil
		}
		if inSidebar {
			listHeight := m.height - statusHeight - headerHeight
			if m.filterInput || m.filter != "" {
				listHeight--
			}
			start := 0
			if m.cursor >= listHeight {
				start = m.cursor - listHeight + 1
			}
			row := msg.Y - headerHeight
			if row < 0 || row >= listHeight {
				return m, nil
			}
			idx := start + row
			if idx >= len(m.visible) {
				return m, nil
			}
			m.cursor = idx
			m.focus = focusDiff
			m.refreshDiff()
		} else {
			m.focus = focusDiff
		}
		return m, nil
	}
	return m, nil
}

// nextUnviewed walks the visible file list forward from the cursor, wrapping once,
// and returns the index of the next unviewed file.
func (m *model) nextUnviewed() (int, bool) {
	n := len(m.visible)
	for step := 1; step <= n; step++ {
		idx := (m.cursor + step) % n
		f := m.files[m.visible[idx]]
		if !m.state.IsViewed(f.Path, f.Hash) {
			return idx, true
		}
	}
	return 0, false
}

func (m *model) refreshDiff() {
	f := m.currentFile()
	if f == nil {
		m.viewport.SetContent("")
		return
	}
	w, _ := m.diffPaneSize()
	m.viewport.SetContent(renderDiff(f, w))
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

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.KeyMsg:
		key := msg.String()

		// help overlay swallows most input; only allow close keys
		if m.showHelp {
			switch key {
			case "?", "esc", "q":
				m.showHelp = false
			case "ctrl+c":
				_ = m.state.Save()
				return m, tea.Quit
			}
			return m, nil
		}

		// commit message overlay — similar behaviour
		if m.showCommit {
			switch key {
			case "m", "esc", "q":
				m.showCommit = false
			case "ctrl+c":
				_ = m.state.Save()
				return m, tea.Quit
			}
			return m, nil
		}

		// filter input mode — typing builds up the filter substring
		if m.filterInput {
			return m.updateFilterInput(key)
		}

		// command mode — typing builds up a :command
		if m.cmdInput {
			return m.updateCmdInput(key)
		}

		// clear any lingering command feedback on the next real key press
		m.cmdError = ""

		// keys that work regardless of which pane has focus
		switch key {
		case "?":
			m.showHelp = true
			return m, nil
		case "/":
			m.filterInput = true
			return m, nil
		case ":":
			m.cmdInput = true
			m.cmdBuffer = ""
			return m, nil
		case "m":
			m.showCommit = true
			return m, nil
		case "q", "ctrl+c":
			_ = m.state.Save()
			return m, tea.Quit
		case "v", " ":
			if f := m.currentFile(); f != nil {
				m.state.ToggleViewed(f.Path, f.Hash)
				_ = m.state.Save()
			}
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
			if m.cursor < len(m.visible)-1 {
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
		case "U":
			if idx, ok := m.nextUnviewed(); ok {
				m.cursor = idx
				m.refreshDiff()
			}
			return m, nil
		}

		if m.focus == focusList {
			switch key {
			case "j", "down":
				if m.cursor < len(m.visible)-1 {
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
				if len(m.visible) > 0 {
					m.cursor = len(m.visible) - 1
					m.refreshDiff()
				}
			case "enter", "l", "right":
				if m.currentFile() != nil {
					m.focus = focusDiff
				}
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
		case "d", "ctrl+d":
			m.viewport.HalfViewDown()
		case "u", "ctrl+u":
			m.viewport.HalfViewUp()
		case "ctrl+f", "pgdown":
			m.viewport.ViewDown()
		case "ctrl+b", "pgup":
			m.viewport.ViewUp()
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

	if m.showHelp {
		return renderHelp(m.width, m.height)
	}
	if m.showCommit {
		return renderCommitMessages(suggestCommitMessages(m.files), m.width, m.height)
	}

	listHeight := m.height - statusHeight - headerHeight
	listFocused := m.focus == focusList
	listW := m.listWidth()

	// Carve a row off the bottom for the filter prompt when it's relevant
	promptHeight := 0
	if m.filterInput || m.filter != "" {
		promptHeight = 1
	}

	visibleFiles := make([]FileDiff, 0, len(m.visible))
	for _, idx := range m.visible {
		visibleFiles = append(visibleFiles, m.files[idx])
	}

	listContent := renderFileList(visibleFiles, m.cursor, m.state, listW-2, listHeight-promptHeight, listFocused)
	if len(visibleFiles) == 0 {
		listContent = lipgloss.NewStyle().Foreground(colorMuted).Italic(true).Render("no matches")
	}

	if promptHeight > 0 {
		listContent += "\n" + renderFilterPrompt(m.filter, m.filterInput, listW-2)
	}

	borderColor := colorBorder
	if !listFocused {
		borderColor = colorMuted
	}

	listPane := lipgloss.NewStyle().
		Width(listW).
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
	tagline := styleHeaderTagline.Render(" " + vibeTagline + " ")

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
	var progress string
	if viewed == len(m.files) && viewed > 0 {
		progress = styleStatusAccent.Render(vibeCelebration)
	} else {
		progress = styleStatusAccent.Render(fmt.Sprintf("💎 %d/%d", viewed, len(m.files)))
	}

	var helpRendered string
	switch {
	case m.cmdInput:
		cmdLine := ":" + m.cmdBuffer + "▎"
		helpRendered = lipgloss.NewStyle().Foreground(colorGold).Bold(true).Background(colorStatus).Render(cmdLine)
	case m.cmdError != "":
		helpRendered = lipgloss.NewStyle().Foreground(colorDel).Bold(true).Background(colorStatus).Render(m.cmdError)
	case m.focus == focusList:
		helpRendered = styleStatusBar.Render("j/k file · enter open · v view · ? help · q quit")
	default:
		helpRendered = styleStatusBar.Render("j/k scroll · ^d/^u half · esc back · ? help · q quit")
	}

	inner := m.width - 2
	gap := inner - lipgloss.Width(progress) - lipgloss.Width(helpRendered)
	if gap < 1 {
		gap = 1
	}
	spacer := styleStatusBar.Render(strings.Repeat(" ", gap))
	return lipgloss.NewStyle().Background(colorStatus).Width(m.width).Render(progress + spacer + helpRendered)
}
