# VHS tape file for HawkTUI Simple Demo
# This demonstrates basic components: button, input, and spinner

Output demos/simple_demo.gif

Require go

Set FontSize 20
Set Width 1400
Set Height 900
Set Theme "Catppuccin Mocha"
Set TypingSpeed 100ms
Set PlaybackSpeed 1.0

Type "go run examples/hawktui/simple_demo.go"
Enter
Sleep 2s

# Type in the input field
Type "HawkTUI Demo"
Sleep 1s

# Press Tab to switch focus to button
Tab
Sleep 500ms

# Press Space/Enter to click the button
Space
Sleep 1s

# Tab back to input
Tab
Sleep 500ms

# Clear and type something else
Ctrl+U
Type "Terminal UI Toolkit"
Sleep 1s

# Quit the application
Type "q"
Sleep 500ms
