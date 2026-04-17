package main

import "github.com/charmbracelet/lipgloss"

var (
	colorAdd    = lipgloss.Color("#3fb950")
	colorDel    = lipgloss.Color("#f85149")
	colorMuted  = lipgloss.Color("#7d8590")
	colorAccent = lipgloss.Color("#58a6ff")
	colorText   = lipgloss.Color("#c9d1d9")
	colorBgAdd  = lipgloss.Color("#0d2819")
	colorBgDel  = lipgloss.Color("#2d0d10")
	colorSelBg  = lipgloss.Color("#21262d")
	colorViewed = lipgloss.Color("#2ea043")
	colorBorder = lipgloss.Color("#30363d")
	colorStatus = lipgloss.Color("#161b22")

	styleAddLine          = lipgloss.NewStyle().Foreground(colorAdd).Background(colorBgAdd)
	styleDelLine          = lipgloss.NewStyle().Foreground(colorDel).Background(colorBgDel)
	styleCtxLine          = lipgloss.NewStyle().Foreground(colorText)
	styleHunkHeader       = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	styleFilePath         = lipgloss.NewStyle().Foreground(colorText)
	styleFilePathSelected = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(colorSelBg).Bold(true)
	styleFileCounts       = lipgloss.NewStyle().Foreground(colorMuted)
	styleViewedMark       = lipgloss.NewStyle().Foreground(colorViewed).Bold(true)
	styleStatusBar        = lipgloss.NewStyle().Foreground(colorMuted).Background(colorStatus).Padding(0, 1)
)
