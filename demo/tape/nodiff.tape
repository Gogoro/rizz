# Captures the "no diff" splash in a freshly initialised repo.
# Run from the repo root with: vhs demo/tape/nodiff.tape

Output "/tmp/rizz-vhs-nodiff.gif"

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 500
Set Padding 20
Set Theme "Catppuccin Mocha"

Hide
Type "rm -rf /tmp/rizz-empty && mkdir /tmp/rizz-empty && cd /tmp/rizz-empty"
Enter
Type "git init -q && git -c commit.gpgsign=false commit -q --allow-empty -m init"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 200ms
Show

Type "/Users/ole/work/gogoro/rizz/rizz"
Enter
Sleep 800ms
Screenshot "assets/screenshots/no-diff.png"
Sleep 300ms
