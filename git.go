package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func runGit(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.Bytes(), nil
}

func IsGitRepo() bool {
	_, err := runGit("rev-parse", "--git-dir")
	return err == nil
}

func RepoRoot() (string, error) {
	out, err := runGit("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// RunDiff returns raw git diff output.
// If staged, diffs the index against HEAD.
// Else if base is empty, diffs uncommitted changes against HEAD.
// Else diffs current branch against base using merge-base (triple-dot).
func RunDiff(base string, staged bool) ([]byte, error) {
	if staged {
		return runGit("diff", "--cached")
	}
	if base == "" {
		return runGit("diff", "HEAD")
	}
	return runGit("diff", base+"...HEAD")
}
