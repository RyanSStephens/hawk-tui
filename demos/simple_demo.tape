# VHS tape file for HawkTUI Simple Demo
# This demonstrates basic components: button, input, and spinner

Output simple_demo.gif

Set FontSize 18
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set PlaybackSpeed 0.5

Type "go run examples/hawktui/simple_demo.go"
Enter
Sleep 3s

# Type in the input field
Type "HawkTUI"
Sleep 800ms

# Press Tab to switch focus to button
Tab
Sleep 600ms

# Press Space to click the button
Space
Sleep 1s

# Tab back to input
Tab
Sleep 600ms

# Type more
Backspace@3
Type "Toolkit"
Sleep 1s

# Quit the application
Type "q"
Sleep 500ms
