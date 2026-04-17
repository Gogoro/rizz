package main

import (
	"fmt"
	"strings"
)

func renderFileList(files []FileDiff, cursor int, state *State, width, height int) string {
	// Scroll the visible window to keep cursor in view
	start := 0
	if cursor >= height {
		start = cursor - height + 1
	}
	end := start + height
	if end > len(files) {
		end = len(files)
	}

	var lines []string
	for i := start; i < end; i++ {
		f := files[i]
		viewed := state.IsViewed(f.Path, f.Hash)

		// "💎" renders 2 cells wide — pad placeholder so the path column doesn't shift
		var mark string
		if viewed {
			mark = styleViewedMark.Render("💎")
		} else {
			mark = styleFileCounts.Render("· ")
		}

		counts := styleFileCounts.Render(fmt.Sprintf("+%d -%d", f.Added, f.Removed))

		// "💎 " renders 2 cells wide, counts vary — leave generous headroom
		path := truncate(f.Path, width-len(fmt.Sprintf("+%d -%d ", f.Added, f.Removed))-3)
		pathStyle := styleFilePath
		if viewed {
			pathStyle = pathStyle.Foreground(colorMuted)
		}
		if i == cursor {
			pathStyle = styleFilePathSelected
		}

		lines = append(lines, fmt.Sprintf("%s %s %s", mark, pathStyle.Render(path), counts))
	}

	return strings.Join(lines, "\n")
}

func truncate(s string, max int) string {
	if max <= 1 {
		return s
	}
	if len(s) <= max {
		return s
	}
	return "…" + s[len(s)-(max-1):]
}
