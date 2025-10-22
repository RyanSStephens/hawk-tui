# VHS tape file for HawkTUI Form Demo
# This demonstrates form validation and input handling

Output form_demo.gif

Set FontSize 18
Set Width 1200
Set Height 800
Set Theme "Nord"
Set TypingSpeed 50ms
Set PlaybackSpeed 0.5

Type "go run examples/hawktui/form_demo.go"
Enter
Sleep 3s

# Fill in the form fields
Type "John Doe"
Sleep 600ms
Tab
Sleep 400ms

Type "john@example.com"
Sleep 600ms
Tab
Sleep 400ms

Type "password123"
Sleep 600ms
Tab
Sleep 400ms

# Navigate to submit button and press it
Tab
Sleep 400ms
Space
Sleep 1.5s

# Quit
Type "q"
Sleep 500ms
