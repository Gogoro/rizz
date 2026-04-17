package rizz

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

// runGitAllowDiff runs git but treats exit code 1 as success.
// Used with `git diff --no-index`, which returns 1 when files differ.
func runGitAllowDiff(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return stdout.Bytes(), nil
		}
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

// ListUntracked returns paths of untracked, non-ignored files, relative to the repo root.
func ListUntracked(repoRoot string) ([]string, error) {
	out, err := runGit("-C", repoRoot, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths, nil
}

// DiffUntracked produces a git-style diff for a single untracked file,
// treating it as a new file added against /dev/null. Paths are always
// interpreted relative to the repo root.
func DiffUntracked(repoRoot, path string) ([]byte, error) {
	return runGitAllowDiff("-C", repoRoot, "diff", "--no-index", "--", "/dev/null", path)
}
