package rizz

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHelp(width, height int) string {
	title := lipgloss.NewStyle().Foreground(colorGold).Bold(true).Render("⛓  R I Z Z  ⛓  keybindings")
	subtitle := lipgloss.NewStyle().Foreground(colorGoldSoft).Italic(true).Render("press ? or esc to close")

	sections := []helpSection{
		{
			name: "list mode",
			keys: []helpKey{
				{"j  k  ↑  ↓", "move between files"},
				{"g  G", "first · last file"},
				{"enter  l  →", "open diff view"},
			},
		},
		{
			name: "diff mode",
			keys: []helpKey{
				{"j  k  ↑  ↓", "scroll the diff"},
				{"^d  ^u", "half-page down · up"},
				{"^f  ^b", "full page down · up"},
				{"g  G", "top · bottom"},
				{"esc  h  ←", "back to list"},
			},
		},
		{
			name: "anywhere",
			keys: []helpKey{
				{"n  tab", "next file"},
				{"p  shift+tab", "previous file"},
				{"v  space", "toggle viewed 💎"},
				{"U", "jump to next unviewed"},
				{"/", "filter files (esc clears)"},
				{":", "vim-style command (:q, :help)"},
				{"m", "commit message suggestions"},
				{"a", "mark all viewed"},
				{"r", "reset all"},
				{"?", "toggle this help"},
				{"q  ^c", "quit"},
			},
		},
	}

	var body strings.Builder
	body.WriteString(title)
	body.WriteString("\n")
	body.WriteString(subtitle)
	body.WriteString("\n\n")
	for i, sec := range sections {
		body.WriteString(sec.render())
		if i < len(sections)-1 {
			body.WriteString("\n\n")
		}
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(colorGoldDeep).
		Padding(1, 3).
		Render(body.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box, lipgloss.WithWhitespaceChars(" "))
}

type helpKey struct {
	keys  string
	label string
}

type helpSection struct {
	name string
	keys []helpKey
}

func (s helpSection) render() string {
	header := lipgloss.NewStyle().Foreground(colorGold).Bold(true).Underline(true).Render(s.name)

	keyStyle := lipgloss.NewStyle().Foreground(colorDiamond).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(colorText)

	widest := 0
	for _, k := range s.keys {
		if lipgloss.Width(k.keys) > widest {
			widest = lipgloss.Width(k.keys)
		}
	}

	var b strings.Builder
	b.WriteString(header)
	b.WriteString("\n")
	for _, k := range s.keys {
		pad := widest - lipgloss.Width(k.keys)
		b.WriteString("  ")
		b.WriteString(keyStyle.Render(k.keys))
		b.WriteString(strings.Repeat(" ", pad))
		b.WriteString("   ")
		b.WriteString(labelStyle.Render(k.label))
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}
