package rizz

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type State struct {
	path   string
	Viewed map[string]string `json:"viewed"`   // file path -> diff hash when marked viewed
	DiffMode string          `json:"diffMode"` // "side" or "inline"; empty means default
}

// stateFilePath resolves the per-worktree git dir via `git rev-parse --git-dir`
// so state is stored correctly in worktrees (where .git is a file, not a dir).
func stateFilePath(repoRoot string) string {
	out, err := runGit("rev-parse", "--git-dir")
	if err == nil {
		gitDir := string(out)
		if n := len(gitDir); n > 0 && gitDir[n-1] == '\n' {
			gitDir = gitDir[:n-1]
		}
		if !filepath.IsAbs(gitDir) {
			gitDir = filepath.Join(repoRoot, gitDir)
		}
		return filepath.Join(gitDir, "rizz-state.json")
	}
	return filepath.Join(repoRoot, ".git", "rizz-state.json")
}

func LoadState(repoRoot string) (*State, error) {
	p := stateFilePath(repoRoot)
	s := &State{path: p, Viewed: map[string]string{}}

	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	if s.Viewed == nil {
		s.Viewed = map[string]string{}
	}
	return s, nil
}

func (s *State) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *State) IsViewed(path, hash string) bool {
	return s.Viewed[path] == hash
}

func (s *State) ToggleViewed(path, hash string) {
	if s.IsViewed(path, hash) {
		delete(s.Viewed, path)
	} else {
		s.Viewed[path] = hash
	}
}

func (s *State) MarkAllViewed(files []FileDiff) {
	for _, f := range files {
		s.Viewed[f.Path] = f.Hash
	}
}

func (s *State) UnmarkAll() {
	s.Viewed = map[string]string{}
}
