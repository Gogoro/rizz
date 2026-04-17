package rizz

import (
	"flag"
	"fmt"
	"os"
)

// Version is set at build time via -ldflags "-X github.com/Gogoro/rizz/internal/rizz.Version=..."
var Version = "dev"

// Main is the CLI entry point. It parses flags, loads config, discovers the
// current git diff, and hands everything off to the TUI.
func Main() {
	base := flag.String("base", "", "compare current branch vs this ref (e.g. main). default: uncommitted changes vs HEAD")
	staged := flag.Bool("staged", false, "review only staged changes (git diff --cached)")
	theme := flag.String("theme", "", "chroma syntax theme (e.g. monokai, dracula, nord). use 'list' to see all")
	noSplash := flag.Bool("no-splash", false, "skip the boot splash animation")
	showVersion := flag.Bool("version", false, "print rizz version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("rizz", Version)
		return
	}

	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "rizz: config:", err)
		os.Exit(1)
	}
	if *theme == "" && cfg.Theme != "" {
		*theme = cfg.Theme
	}

	if *theme == "list" {
		for _, name := range AvailableThemes() {
			fmt.Println(name)
		}
		return
	}
	if *theme != "" {
		SetSyntaxTheme(*theme)
	}

	if !IsGitRepo() {
		fmt.Fprintln(os.Stderr, "rizz: not inside a git repository")
		os.Exit(1)
	}

	root, err := RepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "rizz:", err)
		os.Exit(1)
	}

	raw, err := RunDiff(*base, *staged)
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
		fmt.Println(renderNoDiff())
		return
	}

	files = LoadFileSources(files, *base, *staged, root)

	state, err := LoadState(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rizz: load state:", err)
		os.Exit(1)
	}

	showSplash := !*noSplash && !cfg.NoSplash
	if err := Run(files, state, showSplash, cfg.KeyRemap()); err != nil {
		fmt.Fprintln(os.Stderr, "rizz:", err)
		os.Exit(1)
	}
}
