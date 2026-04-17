package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type State struct {
	path   string
	Viewed map[string]string `json:"viewed"` // file path -> diff hash when marked viewed
}

func stateFilePath(repoRoot string) string {
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
