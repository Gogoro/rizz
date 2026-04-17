package main

import "strings"

func renderDiff(file FileDiff, width int) string {
	if file.IsBinary {
		return styleFileCounts.Render("binary file")
	}

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
			content := string(line.Kind) + line.Text
			switch line.Kind {
			case '+':
				b.WriteString(styleAddLine.Width(width).Render(content))
			case '-':
				b.WriteString(styleDelLine.Width(width).Render(content))
			default:
				b.WriteString(styleCtxLine.Render(content))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	return b.String()
}
