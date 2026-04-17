package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// fileIcon returns a 2-cell-wide emoji representing the file type.
// Kept consistent width so rows line up in the file list.
func fileIcon(path string) string {
	base := filepath.Base(path)
	lower := strings.ToLower(base)
	ext := strings.ToLower(filepath.Ext(base))

	// special filenames first
	switch {
	case lower == "dockerfile" || strings.HasPrefix(lower, "dockerfile."):
		return "🐳"
	case lower == "makefile":
		return "🔨"
	case lower == ".gitignore" || lower == ".gitattributes":
		return "🐙"
	case lower == "go.mod" || lower == "go.sum":
		return "🐹"
	case strings.HasSuffix(lower, "_test.go"),
		strings.Contains(lower, ".test."),
		strings.Contains(lower, ".spec."):
		return "🧪"
	}

	switch ext {
	case ".go":
		return "🐹"
	case ".ts", ".tsx":
		return "🟦"
	case ".js", ".jsx", ".mjs", ".cjs":
		return "🟨"
	case ".py":
		return "🐍"
	case ".rs":
		return "🦀"
	case ".rb":
		return "♦️ "
	case ".java", ".kt":
		return "☕"
	case ".swift":
		return "🦅"
	case ".css", ".scss", ".sass", ".less":
		return "🎨"
	case ".html", ".htm":
		return "🌐"
	case ".md", ".mdx", ".rst":
		return "📝"
	case ".json":
		return "📦"
	case ".yaml", ".yml":
		return "📋"
	case ".toml", ".ini", ".conf", ".cfg":
		return "⚙️ "
	case ".sh", ".bash", ".zsh", ".fish":
		return "🐚"
	case ".sql":
		return "🗄️ "
	case ".lock":
		return "🔒"
	case ".env":
		return "🔐"
	case ".svg", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico":
		return "🖼️ "
	case ".mp3", ".wav", ".ogg", ".flac":
		return "🎵"
	case ".mp4", ".mov", ".webm":
		return "🎬"
	case ".pdf":
		return "📕"
	case ".zip", ".tar", ".gz", ".tgz", ".7z":
		return "📦"
	case ".xml":
		return "📰"
	case ".graphql", ".gql":
		return "🔷"
	case ".proto":
		return "🧬"
	}
	return "📄"
}

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

		icon := fileIcon(f.Path)

		// 2 cells for mark + 1 space + 2 cells for icon + 1 space = 6 cells fixed
		// plus trailing space + counts
		countLen := len(fmt.Sprintf("+%d -%d ", f.Added, f.Removed))
		path := truncate(f.Path, width-countLen-7)
		pathStyle := styleFilePath
		if viewed {
			pathStyle = pathStyle.Foreground(colorMuted)
		}
		if i == cursor {
			pathStyle = styleFilePathSelected
		}

		lines = append(lines, fmt.Sprintf("%s %s %s %s", mark, icon, pathStyle.Render(path), counts))
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
