package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderDiff(file *FileDiff, width int) string {
	if file.IsBinary {
		return styleFileCounts.Render("binary file")
	}

	ensureHighlighted(file)

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
		i := 0
		for i < len(h.Lines) {
			line := h.Lines[i]
			// Detect a simple paired change: one '-' immediately followed by one '+',
			// with no other '-'/'+' around them in a larger group.
			if line.Kind == '-' && i+1 < len(h.Lines) && h.Lines[i+1].Kind == '+' &&
				isIsolatedPair(h.Lines, i) {
				oldText := line.Text
				newText := h.Lines[i+1].Text
				oldSegs, newSegs := computeWordDiff(oldText, newText)
				b.WriteString(renderDelWithIntra(oldSegs, width))
				b.WriteString("\n")
				b.WriteString(renderAddWithIntra(newSegs, width))
				b.WriteString("\n")
				i += 2
				continue
			}

			content := sourceLineFor(line, file.highlightedNew, file.highlightedOld)
			prefix := string(line.Kind)
			switch line.Kind {
			case '+':
				b.WriteString(styleAddPrefix.Render(prefix))
				b.WriteString(styleAddBg.Width(width - 1).Render(content))
			case '-':
				b.WriteString(styleDelPrefix.Render(prefix))
				b.WriteString(styleDelBg.Width(width - 1).Render(content))
			default:
				b.WriteString(styleCtxLine.Render(prefix + content))
			}
			b.WriteString("\n")
			i++
		}
		b.WriteString("\n")
	}
	return b.String()
}

// isIsolatedPair reports whether the - line at position i is part of a 1:1
// removal/addition pair — not part of a larger block of consecutive removals
// or additions.
func isIsolatedPair(lines []Line, i int) bool {
	// the line before must not be '-' or '+'
	if i > 0 {
		prev := lines[i-1].Kind
		if prev == '-' || prev == '+' {
			return false
		}
	}
	// the line after the '+' must not be '-' or '+'
	if i+2 < len(lines) {
		next := lines[i+2].Kind
		if next == '-' || next == '+' {
			return false
		}
	}
	return true
}

func renderAddWithIntra(segs []WordDiffSegment, width int) string {
	var b strings.Builder
	b.WriteString(styleAddPrefix.Render("+"))
	visible := 1
	for _, seg := range segs {
		if seg.Kind == '+' {
			b.WriteString(styleAddIntra.Render(seg.Text))
		} else {
			b.WriteString(styleAddEq.Render(seg.Text))
		}
		visible += lipgloss.Width(seg.Text)
	}
	if visible < width {
		b.WriteString(styleAddBg.Render(strings.Repeat(" ", width-visible)))
	}
	return b.String()
}

func renderDelWithIntra(segs []WordDiffSegment, width int) string {
	var b strings.Builder
	b.WriteString(styleDelPrefix.Render("-"))
	visible := 1
	for _, seg := range segs {
		if seg.Kind == '-' {
			b.WriteString(styleDelIntra.Render(seg.Text))
		} else {
			b.WriteString(styleDelEq.Render(seg.Text))
		}
		visible += lipgloss.Width(seg.Text)
	}
	if visible < width {
		b.WriteString(styleDelBg.Render(strings.Repeat(" ", width-visible)))
	}
	return b.String()
}

func ensureHighlighted(file *FileDiff) {
	if file.highlighted {
		return
	}
	file.highlightedNew = highlightLines(file.NewContent, file.Path)
	file.highlightedOld = highlightLines(file.OldContent, file.OldPath)
	file.highlighted = true
}

// sourceLineFor returns the syntax-highlighted content for a diff line,
// falling back to the raw diff text when highlighted content isn't available.
func sourceLineFor(line Line, newLines, oldLines []string) string {
	lookup := func(arr []string, idx int) (string, bool) {
		if idx < 1 || idx > len(arr) {
			return "", false
		}
		return arr[idx-1], true
	}

	switch line.Kind {
	case '+', ' ':
		if s, ok := lookup(newLines, line.NewLineNum); ok {
			return s
		}
	case '-':
		if s, ok := lookup(oldLines, line.OldLineNum); ok {
			return s
		}
	}
	return line.Text
}
