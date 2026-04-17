package main

import (
	"os"
	"path/filepath"
	"strings"
)

type refPair struct {
	oldRef string // ref for old-side content
	newRef string // ref for new-side content; empty = working tree
}

func resolveRefs(base string) (refPair, error) {
	if base == "" {
		return refPair{oldRef: "HEAD", newRef: ""}, nil
	}
	out, err := runGit("merge-base", base, "HEAD")
	if err != nil {
		return refPair{}, err
	}
	return refPair{oldRef: strings.TrimSpace(string(out)), newRef: "HEAD"}, nil
}

func readAtRef(ref, path string) []byte {
	if path == "" || path == "/dev/null" {
		return nil
	}
	data, err := runGit("show", ref+":"+path)
	if err != nil {
		return nil
	}
	return data
}

func readWorkingTree(repoRoot, path string) []byte {
	data, err := os.ReadFile(filepath.Join(repoRoot, path))
	if err != nil {
		return nil
	}
	return data
}

// LoadFileSources populates NewContent and OldContent on each file.
// Failures are silent — rendering falls back to the raw diff text.
func LoadFileSources(files []FileDiff, base, repoRoot string) []FileDiff {
	refs, err := resolveRefs(base)
	if err != nil {
		return files
	}
	for i := range files {
		f := &files[i]
		if !f.IsNew {
			f.OldContent = readAtRef(refs.oldRef, f.OldPath)
		}
		if !f.IsDelete {
			if refs.newRef == "" {
				f.NewContent = readWorkingTree(repoRoot, f.Path)
			} else {
				f.NewContent = readAtRef(refs.newRef, f.Path)
			}
		}
	}
	return files
}
