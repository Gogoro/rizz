package rizz

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderCommitMessages(msgs []string, width, height int) string {
	title := lipgloss.NewStyle().Foreground(colorGold).Bold(true).Render("💎  commit message suggestions  💎")
	subtitle := lipgloss.NewStyle().Foreground(colorGoldSoft).Italic(true).Render("press m or esc to close")

	itemStyle := lipgloss.NewStyle().Foreground(colorText)
	accent := lipgloss.NewStyle().Foreground(colorDiamond).Bold(true)

	var body strings.Builder
	body.WriteString(title)
	body.WriteString("\n")
	body.WriteString(subtitle)
	body.WriteString("\n\n")
	for i, msg := range msgs {
		body.WriteString("  ")
		body.WriteString(accent.Render("›"))
		body.WriteString(" ")
		body.WriteString(itemStyle.Render(msg))
		if i < len(msgs)-1 {
			body.WriteString("\n")
		}
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(colorGoldDeep).
		Padding(1, 3).
		Render(body.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box, lipgloss.WithWhitespaceChars(" "))
}
