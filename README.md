<h1 align="center">
  <img src="assets/logo.png" alt="rizz" width="420"/>
</h1>

<p align="center">
  <b>Review your own diffs with a little extra rizz.</b><br/>
  <i>A terminal code review tool for when you've got drip but no PR.</i>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/status-alpha-gold?style=flat-square" alt="status: alpha"/>
  <img src="https://img.shields.io/badge/built%20with-bubble%20tea-ff69b4?style=flat-square" alt="built with bubble tea"/>
  <img src="https://img.shields.io/badge/license-MIT-brightgreen?style=flat-square" alt="license: MIT"/>
  <img src="https://img.shields.io/badge/vibes-immaculate-yellow?style=flat-square" alt="vibes: immaculate"/>
</p>

---

## diff so clean, even prod approved it

`rizz` is a terminal TUI that gives you GitHub's "Files Changed" experience for your own uncommitted changes. Scroll the diff, check files off as you review them, ship with confidence — no PR required.

Built for the solo dev, the vibe-coder, the one who ships straight to main and still wants to look over 14 AI-generated files before `git push`.

## why tho

You're deep in the zone. Your AI bestie just wrote half a feature. You need to actually *read the diff* before committing. Your options:

- ❌ Open a PR, review yourself, merge, pull main, delete branch. That's **five steps** for solo work.
- ❌ Blast through `git diff` in a dumb pager and hope you don't miss a file.
- ❌ Boot lazygit, get overwhelmed by 47 panels, rage-quit.
- ✅ Run `rizz`. Scroll. Check. Ship.

## install

```bash
go install github.com/Gogoro/rizz@latest
```

Or from source:

```bash
git clone https://github.com/Gogoro/rizz.git
cd rizz
go build
mv rizz /usr/local/bin/   # or wherever your PATH points
```

Requires Go 1.22+.

## usage

```bash
# review your uncommitted changes (default)
rizz

# review your branch vs main (or any ref)
rizz --base main
rizz --base origin/develop
rizz --base v1.2.0
```

## keybindings

`rizz` has two focus modes: the **file list** on the left, and the **diff view** on the right. Navigate the sidebar, `enter` to jump into the diff, `esc` to jump back. Same muscle memory, context-aware.

**list mode** (default)

| key | action |
|---|---|
| `j` · `k` · `↑` · `↓` | move between files |
| `g` · `G` | first · last file |
| `enter` · `l` · `→` | open diff view |

**diff mode**

| key | action |
|---|---|
| `j` · `k` · `↑` · `↓` | scroll the diff |
| `d` · `u` · `pgdn` · `pgup` | half-page scroll |
| `g` · `G` | top · bottom of diff |
| `esc` · `h` · `←` | back to list |

**works in both modes**

| key | action |
|---|---|
| `n` · `tab` | next file |
| `p` · `shift+tab` | previous file |
| `v` · `space` | toggle viewed 💎 |
| `a` | mark all viewed |
| `r` | reset all |
| `q` · `ctrl+c` | quit |

## viewed tracking that doesn't lie

Mark files as viewed and `rizz` remembers. State lives in `.git/rizz-state.json` per repo — no global mess, no stale marks across projects.

Each viewed mark is keyed to a hash of that file's diff content. When the diff changes (you commit, you edit, anything), the viewed mark auto-invalidates. Same behavior as GitHub's "Viewed" checkbox on PRs. No lies, no gaslighting yourself into thinking you already looked at that file.

## file type flex

`rizz` drops an emoji next to each file so you can skim what's changing at a glance. Heavy on the stylesheets? Pure Go? You'll see it before you read it.

| ext | icon | ext | icon | ext | icon |
|---|---|---|---|---|---|
| `.go` | 🐹 | `.ts` | 🟦 | `.js` | 🟨 |
| `.py` | 🐍 | `.rs` | 🦀 | `.rb` | ♦️ |
| `.css` `.scss` | 🎨 | `.html` | 🌐 | `.md` | 📝 |
| `.json` | 📦 | `.yaml` | 📋 | `.toml` `.ini` | ⚙️ |
| `.sh` `.bash` | 🐚 | `.sql` | 🗄 | `Dockerfile` | 🐳 |
| `Makefile` | 🔨 | `*_test.go` | 🧪 | images | 🖼 |
| `.env` | 🔐 | `.lock` | 🔒 | anything else | 📄 |

## what's NOT here (yet)

This is a weekend build. On purpose:

- No side-by-side diff (unified only)
- No syntax highlighting inside the diff
- No inline comments or annotations
- No staging, committing, or any git mutation — `rizz` is strictly read-only
- No config file — sensible defaults, take it or fork it

If any of these would genuinely make your life better, open an issue.

## status

🚧 **alpha**. Sharp edges, missing features, opinions baked into the defaults. Works on macOS and Linux. Windows? Probably. Report back.

## credits

Built with 💛 using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) by [Charm](https://charm.sh/). Diff parsing by [sourcegraph/go-diff](https://github.com/sourcegraph/go-diff).

## license

MIT. Do whatever you want.

---

<p align="center">
  <i>commit to rizz.</i> 👑
</p>
