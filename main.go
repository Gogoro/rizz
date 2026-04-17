package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	base := flag.String("base", "", "compare current branch vs this ref (e.g. main). default: uncommitted changes vs HEAD")
	flag.Parse()

	if !IsGitRepo() {
		fmt.Fprintln(os.Stderr, "rizz: not inside a git repository")
		os.Exit(1)
	}

	root, err := RepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "rizz:", err)
		os.Exit(1)
	}

	raw, err := RunDiff(*base)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rizz: git diff failed:", err)
		os.Exit(1)
	}

	files, err := ParseDiff(raw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rizz: parse diff failed:", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("rizz: no changes to review")
		return
	}

	state, err := LoadState(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rizz: load state:", err)
		os.Exit(1)
	}

	if err := Run(files, state); err != nil {
		fmt.Fprintln(os.Stderr, "rizz:", err)
		os.Exit(1)
	}
}
