package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var vibeRng = rand.New(rand.NewSource(time.Now().UnixNano()))

var taglines = []string{
	"pure vibes",
	"diff so clean",
	"vibes not bugs",
	"prod approved",
	"certified drip",
	"no cringe code",
	"commit to rizz",
	"lgtm fr fr",
	"rizz detected",
	"clean diff gang",
	"drippin'",
	"ready to ship",
}

var celebrations = []string{
	"💎 CERTIFIED CLEAN 👑",
	"✨ ALL RIZZ NO CRINGE ✨",
	"🔥 LGTM FR FR 🔥",
	"💎 100% PURE RIZZ 💎",
	"👑 PROD APPROVED 👑",
	"⛓  DRIP CHECK PASSED ⛓ ",
}

var noDiffQuips = []string{
	"no diff, no problem",
	"nothing to review, pure rizz",
	"clean as a whistle",
	"no cringe detected",
	"diff empty, vibes full",
	"already shipped?",
}

var (
	vibeTagline     = pickVibe(taglines)
	vibeCelebration = pickVibe(celebrations)
)

func pickVibe(opts []string) string {
	return opts[vibeRng.Intn(len(opts))]
}

const rizzAscii = ` ██████╗ ██╗███████╗███████╗
 ██╔══██╗██║╚══███╔╝╚══███╔╝
 ██████╔╝██║  ███╔╝   ███╔╝
 ██╔══██╗██║ ███╔╝   ███╔╝
 ██║  ██║██║███████╗███████╗
 ╚═╝  ╚═╝╚═╝╚══════╝╚══════╝`

func renderNoDiff() string {
	gold := lipgloss.NewStyle().Foreground(colorGold).Bold(true)
	soft := lipgloss.NewStyle().Foreground(colorGoldSoft).Italic(true)

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(gold.Render(rizzAscii))
	b.WriteString("\n\n")
	b.WriteString("     ")
	b.WriteString(soft.Render(pickVibe(noDiffQuips)))
	b.WriteString(" 👑\n")
	return b.String()
}
