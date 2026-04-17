package main

import "strings"

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
		for _, line := range h.Lines {
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
		}
		b.WriteString("\n")
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
