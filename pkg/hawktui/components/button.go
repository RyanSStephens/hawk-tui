package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Button represents a clickable button
type Button struct {
	*BaseComponent
	label    string
	onPress  func()
	disabled bool
	style    ButtonStyle
}

// ButtonStyle defines the visual style of a button
type ButtonStyle int

const (
	ButtonStylePrimary ButtonStyle = iota
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleWarning
	ButtonStyleError
	ButtonStyleGhost
)

// NewButton creates a new button
func NewButton(label string, onPress func()) *Button {
	return &Button{
		BaseComponent: NewBaseComponent(nil),
		label:         label,
		onPress:       onPress,
		style:         ButtonStylePrimary,
	}
}

// NewButtonWithTheme creates a new button with a specific theme
func NewButtonWithTheme(label string, onPress func(), theme *styles.Theme) *Button {
	return &Button{
		BaseComponent: NewBaseComponent(theme),
		label:         label,
		onPress:       onPress,
		style:         ButtonStylePrimary,
	}
}

// SetLabel sets the button label
func (b *Button) SetLabel(label string) {
	b.label = label
}

// SetStyle sets the button style
func (b *Button) SetStyle(style ButtonStyle) {
	b.style = style
}

// SetDisabled sets the disabled state
func (b *Button) SetDisabled(disabled bool) {
	b.disabled = disabled
}

// IsDisabled returns the disabled state
func (b *Button) IsDisabled() bool {
	return b.disabled
}

// Press triggers the button's onPress handler
func (b *Button) Press() {
	if !b.disabled && b.onPress != nil {
		b.onPress()
	}
}

// Init implements tea.Model
func (b *Button) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if b.focused && !b.disabled {
			switch msg.String() {
			case "enter", " ":
				b.Press()
				return b, nil
			}
		}
	case tea.MouseMsg:
		// Handle mouse events
		if b.HandleMouse(msg) {
			// If mouse was handled (clicked within bounds)
			if !b.disabled && msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionRelease {
				b.Press()
				return b, nil
			}
		}
	}
	return b, nil
}

// View implements tea.Model
func (b *Button) View() string {
	if !b.visible {
		return ""
	}

	style := b.getButtonStyle()

	if b.disabled {
		style = style.Copy().
			Foreground(b.theme.TextMuted).
			Background(b.theme.BgTertiary)
	}

	text := fmt.Sprintf(" %s ", b.label)
	return style.Render(text)
}

// getButtonStyle returns the appropriate style based on button state
func (b *Button) getButtonStyle() lipgloss.Style {
	baseStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder())

	var bgColor, fgColor, borderColor styles.Color

	switch b.style {
	case ButtonStylePrimary:
		bgColor = b.theme.Primary
		fgColor = b.theme.TextInverse
		borderColor = b.theme.Primary
	case ButtonStyleSecondary:
		bgColor = b.theme.Secondary
		fgColor = b.theme.TextInverse
		borderColor = b.theme.Secondary
	case ButtonStyleSuccess:
		bgColor = b.theme.Success
		fgColor = b.theme.TextInverse
		borderColor = b.theme.Success
	case ButtonStyleWarning:
		bgColor = b.theme.Warning
		fgColor = b.theme.TextInverse
		borderColor = b.theme.Warning
	case ButtonStyleError:
		bgColor = b.theme.Error
		fgColor = b.theme.TextInverse
		borderColor = b.theme.Error
	case ButtonStyleGhost:
		bgColor = b.theme.BgPrimary
		fgColor = b.theme.Primary
		borderColor = b.theme.Primary
	}

	if b.focused {
		borderColor = b.theme.BorderActive
		baseStyle = baseStyle.Bold(true)
	}

	return baseStyle.
		Background(bgColor).
		Foreground(fgColor).
		BorderForeground(borderColor)
}
