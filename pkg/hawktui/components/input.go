package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Input represents a text input field
type Input struct {
	*BaseComponent
	value       string
	placeholder string
	cursorPos   int
	maxLength   int
	masked      bool
	onChange    func(string)
	validator   func(string) bool
	errorMsg    string
}

// NewInput creates a new input field
func NewInput(placeholder string) *Input {
	return &Input{
		BaseComponent: NewBaseComponent(nil),
		placeholder:   placeholder,
		maxLength:     -1, // No limit
	}
}

// NewInputWithTheme creates a new input field with a specific theme
func NewInputWithTheme(placeholder string, theme *styles.Theme) *Input {
	return &Input{
		BaseComponent: NewBaseComponent(theme),
		placeholder:   placeholder,
		maxLength:     -1,
	}
}

// SetValue sets the input value
func (i *Input) SetValue(value string) {
	if i.maxLength > 0 && len(value) > i.maxLength {
		value = value[:i.maxLength]
	}

	if i.validator != nil && !i.validator(value) {
		return
	}

	i.value = value
	i.cursorPos = len(value)

	if i.onChange != nil {
		i.onChange(value)
	}
}

// Value returns the current input value
func (i *Input) Value() string {
	return i.value
}

// SetPlaceholder sets the placeholder text
func (i *Input) SetPlaceholder(placeholder string) {
	i.placeholder = placeholder
}

// SetMasked sets whether the input should mask its value (for passwords)
func (i *Input) SetMasked(masked bool) {
	i.masked = masked
}

// SetMaxLength sets the maximum length
func (i *Input) SetMaxLength(length int) {
	i.maxLength = length
}

// SetOnChange sets the onChange callback
func (i *Input) SetOnChange(fn func(string)) {
	i.onChange = fn
}

// SetValidator sets the validation function
func (i *Input) SetValidator(fn func(string) bool) {
	i.validator = fn
}

// SetError sets an error message
func (i *Input) SetError(msg string) {
	i.errorMsg = msg
}

// ClearError clears the error message
func (i *Input) ClearError() {
	i.errorMsg = ""
}

// Init implements tea.Model
func (i *Input) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		// Handle mouse events - clicking focuses the input
		i.HandleMouse(msg)
		// Note: Focus management typically handled by parent
	case tea.KeyMsg:
		if !i.focused {
			return i, nil
		}
		switch msg.Type {
		case tea.KeyBackspace:
			if i.cursorPos > 0 {
				i.value = i.value[:i.cursorPos-1] + i.value[i.cursorPos:]
				i.cursorPos--
				if i.onChange != nil {
					i.onChange(i.value)
				}
			}
		case tea.KeyDelete:
			if i.cursorPos < len(i.value) {
				i.value = i.value[:i.cursorPos] + i.value[i.cursorPos+1:]
				if i.onChange != nil {
					i.onChange(i.value)
				}
			}
		case tea.KeyLeft:
			if i.cursorPos > 0 {
				i.cursorPos--
			}
		case tea.KeyRight:
			if i.cursorPos < len(i.value) {
				i.cursorPos++
			}
		case tea.KeyHome:
			i.cursorPos = 0
		case tea.KeyEnd:
			i.cursorPos = len(i.value)
		case tea.KeyRunes:
			if i.maxLength < 0 || len(i.value) < i.maxLength {
				newValue := i.value[:i.cursorPos] + string(msg.Runes) + i.value[i.cursorPos:]
				if i.validator == nil || i.validator(newValue) {
					i.value = newValue
					i.cursorPos += len(msg.Runes)
					if i.onChange != nil {
						i.onChange(i.value)
					}
					i.ClearError()
				}
			}
		}
	}

	return i, nil
}

// View implements tea.Model
func (i *Input) View() string {
	if !i.visible {
		return ""
	}

	style := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder())

	if i.focused {
		style = style.BorderForeground(i.theme.BorderActive)
	} else {
		style = style.BorderForeground(i.theme.BorderPrimary)
	}

	if i.errorMsg != "" {
		style = style.BorderForeground(i.theme.Error)
	}

	// Determine display value
	displayValue := i.value
	if i.masked && len(i.value) > 0 {
		displayValue = string(make([]rune, len(i.value)))
		for idx := range displayValue {
			displayValue = displayValue[:idx] + "*" + displayValue[idx+1:]
		}
	}

	// Add placeholder if empty
	if displayValue == "" && !i.focused {
		displayValue = i.placeholder
		style = style.Foreground(i.theme.TextMuted)
	} else {
		style = style.Foreground(i.theme.TextPrimary)
	}

	// Add cursor if focused
	if i.focused {
		before := displayValue[:i.cursorPos]
		after := ""
		if i.cursorPos < len(displayValue) {
			after = displayValue[i.cursorPos:]
		}
		cursor := lipgloss.NewStyle().
			Background(i.theme.Primary).
			Foreground(i.theme.TextInverse).
			Render(" ")
		displayValue = before + cursor + after
	}

	// Apply width if set
	if i.width > 0 {
		style = style.Width(i.width - 4) // Account for padding and border
	}

	result := style.Render(displayValue)

	// Add error message if present
	if i.errorMsg != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(i.theme.Error).
			Padding(0, 1)
		result = lipgloss.JoinVertical(lipgloss.Left, result, errorStyle.Render(i.errorMsg))
	}

	return result
}
