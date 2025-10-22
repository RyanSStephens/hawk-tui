package templates

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/components"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Field represents a form field
type Field struct {
	Label    string
	Input    *components.Input
	Required bool
}

// Form represents a form template
type Form struct {
	theme      *styles.Theme
	title      string
	fields     []Field
	focusIndex int
	submitBtn  *components.Button
	cancelBtn  *components.Button
	onSubmit   func(map[string]string)
	onCancel   func()
	width      int
	height     int
}

// NewForm creates a new form template
func NewForm(title string) *Form {
	f := &Form{
		theme:  styles.DefaultTheme(),
		title:  title,
		fields: make([]Field, 0),
	}

	f.submitBtn = components.NewButtonWithTheme("Submit", func() {
		if f.onSubmit != nil {
			values := make(map[string]string)
			for _, field := range f.fields {
				values[field.Label] = field.Input.Value()
			}
			f.onSubmit(values)
		}
	}, f.theme)

	f.cancelBtn = components.NewButtonWithTheme("Cancel", func() {
		if f.onCancel != nil {
			f.onCancel()
		}
	}, f.theme)
	f.cancelBtn.SetStyle(components.ButtonStyleSecondary)

	return f
}

// NewFormWithTheme creates a new form with a specific theme
func NewFormWithTheme(title string, theme *styles.Theme) *Form {
	f := &Form{
		theme:  theme,
		title:  title,
		fields: make([]Field, 0),
	}

	f.submitBtn = components.NewButtonWithTheme("Submit", func() {
		if f.onSubmit != nil {
			values := make(map[string]string)
			for _, field := range f.fields {
				values[field.Label] = field.Input.Value()
			}
			f.onSubmit(values)
		}
	}, theme)

	f.cancelBtn = components.NewButtonWithTheme("Cancel", func() {
		if f.onCancel != nil {
			f.onCancel()
		}
	}, theme)
	f.cancelBtn.SetStyle(components.ButtonStyleSecondary)

	return f
}

// AddField adds a field to the form
func (f *Form) AddField(label string, placeholder string, required bool) {
	input := components.NewInputWithTheme(placeholder, f.theme)
	field := Field{
		Label:    label,
		Input:    input,
		Required: required,
	}
	f.fields = append(f.fields, field)
}

// AddPasswordField adds a password field to the form
func (f *Form) AddPasswordField(label string, placeholder string, required bool) {
	input := components.NewInputWithTheme(placeholder, f.theme)
	input.SetMasked(true)
	field := Field{
		Label:    label,
		Input:    input,
		Required: required,
	}
	f.fields = append(f.fields, field)
}

// SetOnSubmit sets the submit callback
func (f *Form) SetOnSubmit(fn func(map[string]string)) {
	f.onSubmit = fn
}

// SetOnCancel sets the cancel callback
func (f *Form) SetOnCancel(fn func()) {
	f.onCancel = fn
}

// SetSize sets the form dimensions
func (f *Form) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// Validate validates the form
func (f *Form) Validate() bool {
	valid := true
	for _, field := range f.fields {
		if field.Required && field.Input.Value() == "" {
			field.Input.SetError("This field is required")
			valid = false
		}
	}
	return valid
}

// Init implements tea.Model
func (f *Form) Init() tea.Cmd {
	if len(f.fields) > 0 {
		f.fields[f.focusIndex].Input.SetFocused(true)
	}
	return nil
}

// Update implements tea.Model
func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			// Move to next field
			if f.focusIndex < len(f.fields) {
				f.fields[f.focusIndex].Input.SetFocused(false)
			}
			f.focusIndex++
			if f.focusIndex > len(f.fields)+1 {
				f.focusIndex = 0
			}
			if f.focusIndex < len(f.fields) {
				f.fields[f.focusIndex].Input.SetFocused(true)
			}
			f.updateButtonFocus()

		case "shift+tab", "up":
			// Move to previous field
			if f.focusIndex < len(f.fields) {
				f.fields[f.focusIndex].Input.SetFocused(false)
			}
			f.focusIndex--
			if f.focusIndex < 0 {
				f.focusIndex = len(f.fields) + 1
			}
			if f.focusIndex < len(f.fields) {
				f.fields[f.focusIndex].Input.SetFocused(true)
			}
			f.updateButtonFocus()

		case "enter":
			if f.focusIndex == len(f.fields) {
				// Submit button
				if f.Validate() {
					f.submitBtn.Press()
				}
			} else if f.focusIndex == len(f.fields)+1 {
				// Cancel button
				f.cancelBtn.Press()
			}
		}
	}

	// Update focused field
	if f.focusIndex < len(f.fields) {
		var model tea.Model
		var cmd tea.Cmd
		model, cmd = f.fields[f.focusIndex].Input.Update(msg)
		f.fields[f.focusIndex].Input = model.(*components.Input)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) > 0 {
		return f, tea.Batch(cmds...)
	}

	return f, nil
}

// updateButtonFocus updates the button focus states
func (f *Form) updateButtonFocus() {
	f.submitBtn.SetFocused(f.focusIndex == len(f.fields))
	f.cancelBtn.SetFocused(f.focusIndex == len(f.fields)+1)
}

// View implements tea.Model
func (f *Form) View() string {
	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(f.theme.Primary).
		Bold(true).
		Padding(1, 0).
		Align(lipgloss.Center)

	title := titleStyle.Render(f.title)

	// Fields
	var fieldViews []string
	for _, field := range f.fields {
		labelStyle := lipgloss.NewStyle().
			Foreground(f.theme.TextPrimary).
			Bold(true).
			Padding(0, 0, 0, 0)

		label := labelStyle.Render(field.Label)
		if field.Required {
			requiredStyle := lipgloss.NewStyle().
				Foreground(f.theme.Error)
			label += requiredStyle.Render(" *")
		}

		field.Input.SetSize(40, 1)
		input := field.Input.View()

		fieldView := lipgloss.JoinVertical(
			lipgloss.Left,
			label,
			input,
		)
		fieldViews = append(fieldViews, fieldView)
	}

	fields := lipgloss.JoinVertical(lipgloss.Left, fieldViews...)

	// Buttons
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Left,
		f.submitBtn.View(),
		"  ",
		f.cancelBtn.View(),
	)

	// Combine everything
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		fields,
		"",
		buttons,
	)

	// Add border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(f.theme.BorderPrimary).
		Padding(2)

	if f.width > 0 {
		borderStyle = borderStyle.Width(f.width - 6)
	}

	return borderStyle.Render(content)
}
