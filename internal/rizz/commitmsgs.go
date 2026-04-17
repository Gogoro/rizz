package rizz

import (
	"path/filepath"
	"strings"
)

// suggestCommitMessages returns a shuffled, capped list of vibe-y commit
// message templates chosen from heuristics on the changed files.
func suggestCommitMessages(files []FileDiff) []string {
	var newCount, deleteCount int
	typeCounts := map[string]int{}
	for _, f := range files {
		if f.IsNew {
			newCount++
		}
		if f.IsDelete {
			deleteCount++
		}
		typeCounts[classifyForCommit(f.Path)]++
	}

	dominant := ""
	max := 0
	for t, c := range typeCounts {
		if c > max {
			max = c
			dominant = t
		}
	}

	var msgs []string

	if newCount > 0 && deleteCount == 0 && newCount == len(files) {
		msgs = append(msgs,
			"feat: add rizz",
			"feat: commit to rizz",
		)
	}
	if deleteCount > 0 {
		msgs = append(msgs,
			"fix: remove cringe code",
			"chore: clean up the mess",
		)
	}

	switch dominant {
	case "go":
		msgs = append(msgs, "refactor: gopher mode activated")
	case "ts":
		msgs = append(msgs, "feat: type-safe rizz")
	case "css":
		msgs = append(msgs,
			"style: drip check passed",
			"feat: aesthetic overhaul",
		)
	case "test":
		msgs = append(msgs, "test: vibes confirmed")
	case "docs":
		msgs = append(msgs, "docs: words about the rizz")
	case "config":
		msgs = append(msgs, "chore: config drip")
	}

	universals := []string{
		"wip: vibes not bugs",
		"chore: lgtm fr fr",
		"feat: certified clean",
		"refactor: pure vibes",
		"fix: it's giving fixed",
		"chore: prod approved",
	}
	for _, u := range universals {
		msgs = append(msgs, u)
	}

	// dedupe while preserving order
	seen := map[string]bool{}
	unique := msgs[:0]
	for _, m := range msgs {
		if !seen[m] {
			seen[m] = true
			unique = append(unique, m)
		}
	}

	vibeRng.Shuffle(len(unique), func(i, j int) { unique[i], unique[j] = unique[j], unique[i] })
	if len(unique) > 6 {
		unique = unique[:6]
	}
	return unique
}

func classifyForCommit(path string) string {
	base := strings.ToLower(filepath.Base(path))
	ext := strings.ToLower(filepath.Ext(base))
	if strings.Contains(base, "_test.") || strings.Contains(base, ".test.") || strings.Contains(base, ".spec.") {
		return "test"
	}
	switch ext {
	case ".go":
		return "go"
	case ".ts", ".tsx", ".js", ".jsx", ".mjs":
		return "ts"
	case ".py":
		return "py"
	case ".css", ".scss", ".sass", ".less":
		return "css"
	case ".md", ".rst", ".mdx":
		return "docs"
	case ".yaml", ".yml", ".toml", ".json", ".ini":
		return "config"
	}
	return "other"
}
