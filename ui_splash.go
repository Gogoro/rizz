package main

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	splashTotalFrames = 12
	splashTickMs      = 80 * time.Millisecond
)

type splashTickMsg struct{}

func splashTick() tea.Cmd {
	return tea.Tick(splashTickMs, func(time.Time) tea.Msg {
		return splashTickMsg{}
	})
}

func renderSplash(width, height, frame int) string {
	r := rand.New(rand.NewSource(int64(frame) * 31))
	sparkles := []string{"✨", "💫", "⭐"}

	rows := make([]string, height)
	for y := 0; y < height; y++ {
		var b strings.Builder
		x := 0
		for x < width-1 {
			if r.Intn(80) == 0 {
				b.WriteString(sparkles[r.Intn(len(sparkles))])
				x += 2
			} else {
				b.WriteString(" ")
				x++
			}
		}
		if x < width {
			b.WriteString(" ")
		}
		rows[y] = b.String()
	}

	banner := lipgloss.NewStyle().Foreground(colorGold).Bold(true).Render(rizzAscii)
	bannerLines := strings.Split(banner, "\n")
	bannerH := len(bannerLines)
	startY := (height - bannerH - 2) / 2

	for i, bl := range bannerLines {
		y := startY + i
		if y < 0 || y >= len(rows) {
			continue
		}
		padL := (width - lipgloss.Width(bl)) / 2
		if padL < 0 {
			padL = 0
		}
		rows[y] = strings.Repeat(" ", padL) + bl
	}

	tagline := lipgloss.NewStyle().Foreground(colorGoldSoft).Italic(true).Render("⛓  " + vibeTagline + "  ⛓")
	tagY := startY + bannerH + 1
	if tagY >= 0 && tagY < len(rows) {
		padL := (width - lipgloss.Width(tagline)) / 2
		if padL < 0 {
			padL = 0
		}
		rows[tagY] = strings.Repeat(" ", padL) + tagline
	}

	return strings.Join(rows, "\n")
}
