package rizz

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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
	inlineDiff := flag.Bool("inline", false, "start in inline diff mode instead of side-by-side")
	showVersion := flag.Bool("version", false, "print rizz version and exit")
	doUpdate := flag.Bool("update", false, "download and install the latest rizz release")
	flag.Parse()

	if *showVersion {
		fmt.Println("rizz", Version)
		return
	}

	if *doUpdate {
		if err := SelfUpdate(); err != nil {
			fmt.Fprintln(os.Stderr, "rizz: update failed:", err)
			os.Exit(1)
		}
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

	RefreshUpdateCacheBackground()
	if pending := PendingUpdateVersion(); pending != "" {
		promptForUpgrade(pending)
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

	// Only include untracked files for the default "uncommitted vs HEAD" view.
	// --staged and --base X are explicit requests that shouldn't implicitly pull in untracked files.
	if !*staged && *base == "" {
		untracked, err := ListUntracked(root)
		if err == nil {
			for _, path := range untracked {
				raw, err := DiffUntracked(root, path)
				if err != nil || len(raw) == 0 {
					continue
				}
				parsed, err := ParseDiff(raw)
				if err != nil {
					continue
				}
				for i := range parsed {
					parsed[i].IsUntracked = true
				}
				files = append(files, parsed...)
			}
		}
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
	if err := Run(files, state, showSplash, cfg.KeyRemap(), *inlineDiff); err != nil {
		fmt.Fprintln(os.Stderr, "rizz:", err)
		os.Exit(1)
	}
}

// promptForUpgrade asks the user if they want to install a newer version of
// rizz. On yes: runs self-update and exits so the next invocation picks up the
// new binary. On no: remembers the skipped version so we don't re-prompt for it.
func promptForUpgrade(newVersion string) {
	fmt.Printf("\n  \u26d3  rizz %s is available (you have v%s)\n", newVersion, Version)
	fmt.Print("  upgrade now? [Y/n] ")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer == "n" || answer == "no" {
		RememberSkippedVersion(newVersion)
		fmt.Println()
		return
	}

	fmt.Println()
	if err := SelfUpdate(); err != nil {
		fmt.Fprintln(os.Stderr, "rizz: update failed:", err)
		fmt.Println("continuing with current version...")
		fmt.Println()
		return
	}
	fmt.Println("\nrun rizz again to use the new version.")
	os.Exit(0)
}
