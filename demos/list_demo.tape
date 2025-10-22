# VHS tape file for HawkTUI List Demo
# This demonstrates list filtering and selection

Output demos/list_demo.gif

Require go

Set FontSize 20
Set Width 1400
Set Height 900
Set Theme "Catppuccin Mocha"
Set TypingSpeed 100ms
Set PlaybackSpeed 1.0

Type "go run examples/hawktui/list_demo.go"
Enter
Sleep 2s

# Navigate through the list
Down
Sleep 400ms
Down
Sleep 400ms
Down
Sleep 400ms
Up
Sleep 400ms

# Filter the list
Type "/"
Sleep 300ms
Type "task"
Sleep 1s

# Clear filter
Escape
Sleep 500ms

# Select an item
Enter
Sleep 1s

# Quit
Type "q"
Sleep 500ms
