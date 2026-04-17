# rizz

Terminal TUI for self-reviewing `git diff` before committing. Go + Bubble Tea + Lip Gloss + Chroma.

## project layout

- `main.go` — tiny shim that calls `rizz.Main()`
- `internal/rizz/` — all the actual code (one package, flat layout on purpose)
- `demo/setup.sh` — spins up a throwaway repo at `/tmp/rizz-demo` with varied changes
- `demo/tape/` — VHS tapes that drive rizz and capture screenshots
- `assets/logo.png` — project logo (gold/diamond/graffiti vibe)
- `assets/screenshots/` — PNGs embedded in the README

## screenshots

**Regenerate README screenshots with VHS, not a screen-recorder.** The tapes live in `demo/tape/` and output to `assets/screenshots/`.

```bash
# regenerate the main set (list, diff, word-diff, filter, help, commit-msgs, confetti, splash)
vhs demo/tape/main.tape

# regenerate the no-diff splash in an empty repo
vhs demo/tape/nodiff.tape
```

Both tapes use absolute paths to `/Users/ole/work/gogoro/rizz/rizz` and the demo setup script, so they only work on this machine as-is. Update the paths if working elsewhere.

When adding a new screenshot for the README:
1. Add a `Screenshot "assets/screenshots/<name>.png"` line at the right point in `demo/tape/main.tape` (quote the path — VHS requires string literals)
2. Re-run `vhs demo/tape/main.tape`
3. Verify the PNG looks right before embedding in the README
4. Keep README screenshot count low — lean on text for most features, screenshots only for the visually striking ones

If VHS isn't installed: `brew install vhs`. If it's installed but not on `$PATH`, it's at `/opt/homebrew/bin/vhs`.

## code style (per project owner)

- Prefer expressive, clear code over clever or abstract code
- Flat file layout inside `internal/rizz/` — resist splitting into more sub-packages
- Write out each thing explicitly; avoid code-generation patterns
- Full words over abbreviations
- No comments unless the WHY is non-obvious (rarely; let good names carry it)
- No trailing summary docstrings or feature explanation in comments — that belongs in the README or commit message

## commit messages

Short, lowercase, no conventional-commit prefix. Match existing style: `add --staged flag`, `fix sidebar width scaling`, `bump chroma`.
