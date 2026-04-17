package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

const (
	// Skip highlighting for anything bigger than this — keeps large files snappy.
	maxHighlightBytes = 500_000

	// Reset foreground only; leaves background (set by outer style) intact.
	ansiFgReset = "\x1b[39m"
)

var chromaStyle = pickChromaStyle()

func pickChromaStyle() *chroma.Style {
	for _, name := range []string{"catppuccin-mocha", "monokai", "dracula", "github-dark"} {
		if s := styles.Get(name); s != nil && s.Name == name {
			return s
		}
	}
	return styles.Fallback
}

// highlightLines returns ANSI-colored per-line strings for content.
// Emits only foreground codes (never \x1b[0m), so the caller can wrap each
// line with its own background without the background getting reset mid-line.
// Returns nil for empty, binary, or too-large content.
func highlightLines(content []byte, filename string) []string {
	if len(content) == 0 || len(content) > maxHighlightBytes {
		return nil
	}
	if bytes.IndexByte(content, 0) >= 0 {
		return nil // looks binary
	}

	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		return strings.Split(string(content), "\n")
	}

	var tokens []chroma.Token
	for t := iterator(); t != chroma.EOF; t = iterator() {
		tokens = append(tokens, t)
	}

	lineGroups := chroma.SplitTokensIntoLines(tokens)
	out := make([]string, 0, len(lineGroups))
	for _, group := range lineGroups {
		var buf strings.Builder
		for _, t := range group {
			v := strings.TrimRight(t.Value, "\n")
			if v == "" {
				continue
			}
			entry := chromaStyle.Get(t.Type)
			if entry.Colour.IsSet() {
				c := entry.Colour
				buf.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.Red(), c.Green(), c.Blue()))
				buf.WriteString(v)
				buf.WriteString(ansiFgReset)
			} else {
				buf.WriteString(v)
			}
		}
		out = append(out, buf.String())
	}
	return out
}
