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
Set TypingSpeed 40ms
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
Sleep 400ms
Enter

# splash flashes, then file list
Sleep 2s

# scroll through a couple files in the sidebar
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 600ms

# dive into the diff
Enter
Sleep 1s

# scroll through the diff to show syntax highlighting
Type "d"
Sleep 600ms
Type "d"
Sleep 700ms
Type "u"
Sleep 500ms

# back to the list, mark as viewed
Escape
Sleep 400ms
Type "v"
Sleep 600ms

# jump to next unviewed
Type "U"
Sleep 700ms
Enter
Sleep 900ms

# word-level diff moment — scroll a little
Type "d"
Sleep 800ms

# filter
Escape
Sleep 300ms
Type "/"
Sleep 200ms
Type "go"
Sleep 1s
Escape
Sleep 400ms

# commit message suggestions — the vibey finale
Type "m"
Sleep 2s
Escape
Sleep 300ms

# clean exit
Type "q"
Sleep 500ms
