package rizz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func renderDiffSideBySide(file *FileDiff, width int) string {
	gw := computeGutterWidth(file)
	gutterSize := gw + 4 // " N │ "
	overhead := 2 * gutterSize
	if width-overhead < 20 {
		return renderDiffInline(file, width)
	}
	halfWidth := (width - overhead) / 2

	var b strings.Builder

	header := file.Path
	if file.IsNew {
		header += "  (new file)"
	} else if file.IsDelete {
		header += "  (deleted)"
	} else if file.OldPath != "" && file.OldPath != file.Path {
		header = file.OldPath + " → " + file.Path
	}
	b.WriteString(styleHunkHeader.Render(header))
	b.WriteString("\n\n")

	for _, h := range file.Hunks {
		b.WriteString(styleHunkHeader.Render(h.Header))
		b.WriteString("\n")

		rows := pairHunkLines(h.Lines)
		for _, row := range rows {
			b.WriteString(renderSideRow(row, file, gw, halfWidth))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	return b.String()
}

type sideRow struct {
	left  *Line
	right *Line
	kind  byte // ' ' context · 'p' paired del+add · 'd' del only · 'a' add only
}

func pairHunkLines(lines []Line) []sideRow {
	var rows []sideRow
	i := 0
	for i < len(lines) {
		if lines[i].Kind == ' ' {
			ctx := lines[i]
			rows = append(rows, sideRow{left: &ctx, right: &ctx, kind: ' '})
			i++
			continue
		}

		var dels, adds []Line
		for i < len(lines) && lines[i].Kind == '-' {
			dels = append(dels, lines[i])
			i++
		}
		for i < len(lines) && lines[i].Kind == '+' {
			adds = append(adds, lines[i])
			i++
		}

		paired := len(dels) > 0 && len(dels) == len(adds)
		maxN := len(dels)
		if len(adds) > maxN {
			maxN = len(adds)
		}
		for r := 0; r < maxN; r++ {
			row := sideRow{}
			if r < len(dels) {
				d := dels[r]
				row.left = &d
			}
			if r < len(adds) {
				a := adds[r]
				row.right = &a
			}
			switch {
			case row.left != nil && row.right != nil:
				if paired {
					row.kind = 'p'
				} else {
					row.kind = 'P'
				}
			case row.left != nil:
				row.kind = 'd'
			default:
				row.kind = 'a'
			}
			rows = append(rows, row)
		}
	}
	return rows
}

func renderSideRow(row sideRow, file *FileDiff, gw, halfWidth int) string {
	switch row.kind {
	case ' ':
		content := sourceLineFor(*row.left, file.highlightedNew, file.highlightedOld)
		leftGut := renderSideGutter(row.left.OldLineNum, gw)
		rightGut := renderSideGutter(row.right.NewLineNum, gw)
		cell := renderPaddedCell(styleCtxLine, content, halfWidth)
		return leftGut + cell + rightGut + cell
	case 'p':
		oldSegs, newSegs := computeWordDiff(row.left.Text, row.right.Text)
		leftGut := renderSideGutter(row.left.OldLineNum, gw)
		rightGut := renderSideGutter(row.right.NewLineNum, gw)
		return leftGut + renderDelSide(oldSegs, halfWidth) + rightGut + renderAddSide(newSegs, halfWidth)
	case 'P':
		leftGut := renderSideGutter(row.left.OldLineNum, gw)
		rightGut := renderSideGutter(row.right.NewLineNum, gw)
		leftContent := sourceLineFor(*row.left, file.highlightedNew, file.highlightedOld)
		rightContent := sourceLineFor(*row.right, file.highlightedNew, file.highlightedOld)
		return leftGut + renderPaddedCell(styleDelBg, leftContent, halfWidth) +
			rightGut + renderPaddedCell(styleAddBg, rightContent, halfWidth)
	case 'd':
		leftGut := renderSideGutter(row.left.OldLineNum, gw)
		rightGut := renderSideGutter(0, gw)
		content := sourceLineFor(*row.left, file.highlightedNew, file.highlightedOld)
		return leftGut + renderPaddedCell(styleDelBg, content, halfWidth) +
			rightGut + strings.Repeat(" ", halfWidth)
	case 'a':
		leftGut := renderSideGutter(0, gw)
		rightGut := renderSideGutter(row.right.NewLineNum, gw)
		content := sourceLineFor(*row.right, file.highlightedNew, file.highlightedOld)
		return leftGut + strings.Repeat(" ", halfWidth) +
			rightGut + renderPaddedCell(styleAddBg, content, halfWidth)
	}
	return ""
}

// renderPaddedCell clips content to the given visible width and pads any
// remainder with style-applied spaces. We avoid lipgloss's Width(...) here
// because it reflows ANSI-heavy content and can wrap into the next row.
func renderPaddedCell(style lipgloss.Style, content string, width int) string {
	if width <= 0 {
		return ""
	}
	visible := lipgloss.Width(content)
	if visible > width {
		content = ansi.Truncate(content, width, "")
		visible = lipgloss.Width(content)
	}
	out := style.Render(content)
	if visible < width {
		out += style.Render(strings.Repeat(" ", width-visible))
	}
	return out
}

func renderSideGutter(n, gw int) string {
	num := strings.Repeat(" ", gw)
	if n > 0 {
		num = fmt.Sprintf("%*d", gw, n)
	}
	return styleGutter.Render(" " + num + " │ ")
}

func renderAddSide(segs []WordDiffSegment, width int) string {
	var b strings.Builder
	visible := 0
	for _, seg := range segs {
		if visible >= width {
			break
		}
		text := seg.Text
		segWidth := lipgloss.Width(text)
		if visible+segWidth > width {
			text = truncateToWidth(text, width-visible)
			segWidth = lipgloss.Width(text)
		}
		if seg.Kind == '+' {
			b.WriteString(styleAddIntra.Render(text))
		} else {
			b.WriteString(styleAddEq.Render(text))
		}
		visible += segWidth
	}
	if visible < width {
		b.WriteString(styleAddBg.Render(strings.Repeat(" ", width-visible)))
	}
	return b.String()
}

func renderDelSide(segs []WordDiffSegment, width int) string {
	var b strings.Builder
	visible := 0
	for _, seg := range segs {
		if visible >= width {
			break
		}
		text := seg.Text
		segWidth := lipgloss.Width(text)
		if visible+segWidth > width {
			text = truncateToWidth(text, width-visible)
			segWidth = lipgloss.Width(text)
		}
		if seg.Kind == '-' {
			b.WriteString(styleDelIntra.Render(text))
		} else {
			b.WriteString(styleDelEq.Render(text))
		}
		visible += segWidth
	}
	if visible < width {
		b.WriteString(styleDelBg.Render(strings.Repeat(" ", width-visible)))
	}
	return b.String()
}

func truncateToWidth(s string, w int) string {
	if w <= 0 {
		return ""
	}
	var b strings.Builder
	used := 0
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if used+rw > w {
			break
		}
		b.WriteRune(r)
		used += rw
	}
	return b.String()
}
