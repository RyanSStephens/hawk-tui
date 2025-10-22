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
	// Position tracking for mouse events
	x        int
	y        int
	onClick  func(x, y int)
	onHover  func(x, y int)
	hovering bool
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

// SetPosition sets the component's position on screen
func (b *BaseComponent) SetPosition(x, y int) {
	b.x = x
	b.y = y
}

// Position returns the component's position
func (b *BaseComponent) Position() (int, int) {
	return b.x, b.y
}

// SetOnClick sets the click handler
func (b *BaseComponent) SetOnClick(handler func(x, y int)) {
	b.onClick = handler
}

// SetOnHover sets the hover handler
func (b *BaseComponent) SetOnHover(handler func(x, y int)) {
	b.onHover = handler
}

// IsHovering returns whether the mouse is hovering over the component
func (b *BaseComponent) IsHovering() bool {
	return b.hovering
}

// InBounds checks if the given coordinates are within the component bounds
func (b *BaseComponent) InBounds(x, y int) bool {
	return x >= b.x && x < b.x+b.width &&
		y >= b.y && y < b.y+b.height
}

// HandleMouse processes mouse events and returns true if handled
func (b *BaseComponent) HandleMouse(msg tea.MouseMsg) bool {
	if !b.visible {
		return false
	}

	inBounds := b.InBounds(msg.X, msg.Y)

	// Handle hover state
	if inBounds != b.hovering {
		b.hovering = inBounds
		if inBounds && b.onHover != nil {
			b.onHover(msg.X-b.x, msg.Y-b.y)
		}
	}

	// Handle clicks
	if inBounds && msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionRelease {
		if b.onClick != nil {
			b.onClick(msg.X-b.x, msg.Y-b.y)
			return true
		}
	}

	return inBounds
}
