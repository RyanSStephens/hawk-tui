package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Component is the base interface for all UI components
type Component interface {
	tea.Model
	// SetTheme applies a theme to the component
	SetTheme(*styles.Theme)
	// SetSize sets the component dimensions
	SetSize(width, height int)
	// SetFocused sets the focus state
	SetFocused(bool)
	// IsFocused returns the focus state
	IsFocused() bool
}

// BaseComponent provides common functionality for all components
type BaseComponent struct {
	theme   *styles.Theme
	width   int
	height  int
	focused bool
	visible bool
}

// NewBaseComponent creates a new base component
func NewBaseComponent(theme *styles.Theme) *BaseComponent {
	if theme == nil {
		theme = styles.DefaultTheme()
	}
	return &BaseComponent{
		theme:   theme,
		visible: true,
	}
}

// SetTheme sets the theme
func (b *BaseComponent) SetTheme(theme *styles.Theme) {
	b.theme = theme
}

// Theme returns the current theme
func (b *BaseComponent) Theme() *styles.Theme {
	return b.theme
}

// SetSize sets the dimensions
func (b *BaseComponent) SetSize(width, height int) {
	b.width = width
	b.height = height
}

// Width returns the width
func (b *BaseComponent) Width() int {
	return b.width
}

// Height returns the height
func (b *BaseComponent) Height() int {
	return b.height
}

// SetFocused sets the focus state
func (b *BaseComponent) SetFocused(focused bool) {
	b.focused = focused
}

// IsFocused returns the focus state
func (b *BaseComponent) IsFocused() bool {
	return b.focused
}

// SetVisible sets the visibility
func (b *BaseComponent) SetVisible(visible bool) {
	b.visible = visible
}

// IsVisible returns the visibility
func (b *BaseComponent) IsVisible() bool {
	return b.visible
}

// Style returns the appropriate style based on focus state
func (b *BaseComponent) Style() lipgloss.Style {
	if b.focused {
		return b.theme.Focused
	}
	return b.theme.Blurred
}
