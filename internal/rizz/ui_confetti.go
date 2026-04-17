package rizz

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	confettiTotalFrames = 24
	confettiTickMs      = 80 * time.Millisecond
)

type confettiTickMsg struct{}

func confettiTick() tea.Cmd {
	return tea.Tick(confettiTickMs, func(time.Time) tea.Msg {
		return confettiTickMsg{}
	})
}

var confettiGlyphs = []string{"💎", "✨", "👑", "🔥", "⭐", "💫", "💰", "⛓ "}

func renderConfetti(width, height, frame int) string {
	r := rand.New(rand.NewSource(int64(frame) * 7919))

	rows := make([]string, height)
	for y := 0; y < height; y++ {
		var b strings.Builder
		x := 0
		for x < width-1 {
			if r.Intn(22) == 0 {
				b.WriteString(confettiGlyphs[r.Intn(len(confettiGlyphs))])
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

	if height > 0 {
		banner := styleHeaderBar.Render(" 💎  P U R E   R I Z Z  💎 ")
		bannerW := lipgloss.Width(banner)
		padL := (width - bannerW) / 2
		if padL < 0 {
			padL = 0
		}
		midY := height / 2
		if midY >= len(rows) {
			midY = len(rows) - 1
		}
		rows[midY] = strings.Repeat(" ", padL) + banner
	}

	return strings.Join(rows, "\n")
}
