# dcctl terminal demo — run from repository root: vhs docs/demo.tape
# Requires: dcctl on PATH (e.g. make build && export PATH="$PWD:$PATH")

Output docs/demo.gif

Require dcctl

Set Shell "bash"
Set FontSize 14
Set Width 900
Set Height 500
Set Padding 20

# Clean prompt for recording
Enter
Sleep 300ms

# Help
Type "dcctl -h"
Enter
Sleep 2s

# Version
Type "dcctl version -s"
Enter
Sleep 1s

# Config template (first few lines)
Type "dcctl config | head -20"
Enter
Sleep 2s

# Clean exit
Type "exit"
Enter
Sleep 500ms
