# VHS tape file for HawkTUI Dashboard Demo
# This demonstrates the dashboard template with metrics and status widgets

Output demos/dashboard_demo.gif

Require go

Set FontSize 18
Set Width 1600
Set Height 1000
Set Theme "Dracula"
Set TypingSpeed 100ms
Set PlaybackSpeed 1.0

Type "go run examples/hawktui/dashboard_demo.go"
Enter
Sleep 3s

# Let the user view the dashboard
Sleep 2s

# Quit the application
Type "q"
Sleep 500ms
