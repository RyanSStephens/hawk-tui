# VHS tape file for HawkTUI List Demo
# This demonstrates list filtering and selection

Output list_demo.gif

Set FontSize 18
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set PlaybackSpeed 0.5

Type "go run examples/hawktui/list_demo.go"
Enter
Sleep 3s

# Navigate through the list
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Up
Sleep 500ms

# Filter the list (if supported)
Type "/"
Sleep 400ms
Type "task"
Sleep 1s

# Clear filter
Escape
Sleep 600ms

# Select an item
Enter
Sleep 1s

# Quit
Type "q"
Sleep 500ms
