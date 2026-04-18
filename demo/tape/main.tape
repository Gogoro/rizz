# VHS tape that captures PNG screenshots of rizz in action.
# Run from the repo root with: vhs demo/tape/main.tape

Output "/tmp/rizz-vhs-main.gif"

Set Shell "bash"
Set FontSize 14
Set Width 1400
Set Height 800
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 30ms

# setup a throwaway demo repo
Hide
Type "bash /Users/ole/work/gogoro/rizz/demo/setup.sh > /dev/null 2>&1"
Enter
Sleep 2s
Type "cd /tmp/rizz-demo"
Enter
Sleep 200ms
Type "clear"
Enter
Sleep 200ms
Show

# launch rizz and catch the splash
Type "/Users/ole/work/gogoro/rizz/rizz"
Enter
Sleep 500ms
Screenshot "assets/screenshots/splash.png"
Sleep 1500ms
Screenshot "assets/screenshots/list-mode.png"

# into the first file's diff (side-by-side is the default)
Enter
Sleep 800ms
Screenshot "assets/screenshots/diff-mode.png"

# jump a few files to hit the CSS file (intra-line word diff territory)
Escape
Sleep 200ms
Type "jjjjj"
Sleep 300ms
Enter
Sleep 800ms
Screenshot "assets/screenshots/side-by-side.png"
Sleep 200ms
Screenshot "assets/screenshots/word-diff.png"

# toggle to inline mode for comparison
Type "s"
Sleep 600ms
Screenshot "assets/screenshots/inline-mode.png"

# back to side-by-side so the rest of the tape runs in the default view
Type "s"
Sleep 300ms

# filter
Escape
Sleep 200ms
Type "/"
Sleep 200ms
Type "go"
Sleep 800ms
Screenshot "assets/screenshots/filter.png"
Escape
Sleep 300ms

# help overlay
Type "?"
Sleep 800ms
Screenshot "assets/screenshots/help.png"
Escape
Sleep 300ms

# commit message suggestions
Type "m"
Sleep 800ms
Screenshot "assets/screenshots/commit-msgs.png"
Escape
Sleep 300ms

# easter egg
Type "rizz"
Sleep 600ms
Screenshot "assets/screenshots/confetti.png"
Sleep 1500ms

# clean exit
Type "q"
Sleep 400ms
