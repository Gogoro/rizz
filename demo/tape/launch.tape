# Launch GIF — short, punchy demo for Show HN / social posts.
# Run from the repo root with: vhs demo/tape/launch.tape
# Output lands at assets/launch.gif.

Output "assets/launch.gif"

Set Shell "bash"
Set FontSize 14
Set Width 1400
Set Height 800
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set PlaybackSpeed 1.0

# setup demo repo silently
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

# the user types the command
Type "rizz"
Sleep 300ms
Enter

# splash flashes by quickly, then file list
Sleep 900ms

# quick scroll through the sidebar
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 300ms

# dive into the diff — the money shot (syntax highlight)
Enter
Sleep 600ms
Type "d"
Sleep 400ms
Type "d"
Sleep 500ms

# back, mark viewed 💎
Escape
Sleep 250ms
Type "v"
Sleep 500ms

# easter egg finale 🎉
Type "rizz"
Sleep 1800ms

# clean exit
Type "q"
Sleep 400ms
