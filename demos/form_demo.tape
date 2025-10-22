# VHS tape file for HawkTUI Form Demo
# This demonstrates form validation and input handling

Output demos/form_demo.gif

Require go

Set FontSize 20
Set Width 1400
Set Height 900
Set Theme "Nord"
Set TypingSpeed 100ms
Set PlaybackSpeed 1.0

Type "go run examples/hawktui/form_demo.go"
Enter
Sleep 2s

# Fill in the form fields
Type "John Doe"
Sleep 500ms
Tab
Sleep 300ms

Type "john.doe@example.com"
Sleep 500ms
Tab
Sleep 300ms

Type "mypassword123"
Sleep 500ms
Tab
Sleep 300ms

# Navigate to submit button and press it
Tab
Sleep 300ms
Space
Sleep 1.5s

# Quit
Type "q"
Sleep 500ms
