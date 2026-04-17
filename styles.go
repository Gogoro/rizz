package main

import "github.com/charmbracelet/lipgloss"

var (
	// bling palette
	colorGold     = lipgloss.Color("#ffd700")
	colorGoldDeep = lipgloss.Color("#daa520")
	colorGoldSoft = lipgloss.Color("#f0c674")
	colorDiamond  = lipgloss.Color("#b9f2ff")

	colorAdd    = lipgloss.Color("#3fb950")
	colorDel    = lipgloss.Color("#f85149")
	colorMuted  = lipgloss.Color("#7d8590")
	colorAccent = colorGold
	colorText   = lipgloss.Color("#c9d1d9")
	colorBgAdd  = lipgloss.Color("#0d2819")
	colorBgDel  = lipgloss.Color("#2d0d10")
	colorSelBg  = lipgloss.Color("#3a2a05") // gold-tinted selection background
	colorViewed = colorDiamond
	colorBorder = colorGoldDeep
	colorStatus = lipgloss.Color("#1a1400")
	colorHeader = lipgloss.Color("#2a1f00")

	styleAddLine          = lipgloss.NewStyle().Foreground(colorAdd).Background(colorBgAdd)
	styleDelLine          = lipgloss.NewStyle().Foreground(colorDel).Background(colorBgDel)
	styleAddBg            = lipgloss.NewStyle().Background(colorBgAdd)
	styleDelBg            = lipgloss.NewStyle().Background(colorBgDel)
	styleAddPrefix        = lipgloss.NewStyle().Foreground(colorAdd).Background(colorBgAdd).Bold(true)
	styleDelPrefix        = lipgloss.NewStyle().Foreground(colorDel).Background(colorBgDel).Bold(true)
	styleCtxLine          = lipgloss.NewStyle().Foreground(colorText)
	styleHunkHeader       = lipgloss.NewStyle().Foreground(colorGold).Bold(true)
	styleFilePath         = lipgloss.NewStyle().Foreground(colorText)
	styleFilePathSelected    = lipgloss.NewStyle().Foreground(colorGold).Background(colorSelBg).Bold(true)
	styleFilePathSelectedDim = lipgloss.NewStyle().Foreground(colorGoldSoft).Underline(true)
	styleFileCounts       = lipgloss.NewStyle().Foreground(colorMuted)
	styleViewedMark       = lipgloss.NewStyle().Foreground(colorDiamond).Bold(true)
	styleStatusBar        = lipgloss.NewStyle().Foreground(colorGoldSoft).Background(colorStatus).Padding(0, 1)
	styleStatusAccent     = lipgloss.NewStyle().Foreground(colorGold).Background(colorStatus).Bold(true)
	styleHeaderBar        = lipgloss.NewStyle().Foreground(colorGold).Background(colorHeader).Bold(true).Padding(0, 1)
	styleHeaderTagline    = lipgloss.NewStyle().Foreground(colorGoldSoft).Background(colorHeader).Italic(true)
	styleHeaderChain      = lipgloss.NewStyle().Foreground(colorGoldDeep).Background(colorHeader)
)
