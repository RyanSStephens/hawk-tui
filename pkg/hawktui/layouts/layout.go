package layouts

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Direction defines the layout direction
type Direction int

const (
	DirectionHorizontal Direction = iota
	DirectionVertical
)

// Alignment defines content alignment
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// Container represents a layout container
type Container struct {
	direction Direction
	alignment Alignment
	spacing   int
	children  []string
	width     int
	height    int
}

// NewContainer creates a new layout container
func NewContainer(direction Direction) *Container {
	return &Container{
		direction: direction,
		alignment: AlignStart,
		spacing:   0,
		children:  make([]string, 0),
	}
}

// SetAlignment sets the alignment
func (c *Container) SetAlignment(alignment Alignment) {
	c.alignment = alignment
}

// SetSpacing sets the spacing between children
func (c *Container) SetSpacing(spacing int) {
	c.spacing = spacing
}

// SetSize sets the container size
func (c *Container) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// AddChild adds a child element
func (c *Container) AddChild(child string) {
	c.children = append(c.children, child)
}

// Render renders the layout
func (c *Container) Render() string {
	if len(c.children) == 0 {
		return ""
	}

	// Add spacing between children
	var withSpacing []string
	if c.spacing > 0 && len(c.children) > 1 {
		spacer := strings.Repeat(" ", c.spacing)
		for i, child := range c.children {
			withSpacing = append(withSpacing, child)
			if i < len(c.children)-1 {
				if c.direction == DirectionHorizontal {
					withSpacing = append(withSpacing, spacer)
				} else {
					withSpacing = append(withSpacing, strings.Repeat("\n", c.spacing))
				}
			}
		}
	} else {
		withSpacing = c.children
	}

	var result string
	if c.direction == DirectionHorizontal {
		pos := lipgloss.Left
		switch c.alignment {
		case AlignCenter:
			pos = lipgloss.Center
		case AlignEnd:
			pos = lipgloss.Right
		}
		result = lipgloss.JoinHorizontal(pos, withSpacing...)
	} else {
		pos := lipgloss.Left
		switch c.alignment {
		case AlignCenter:
			pos = lipgloss.Center
		case AlignEnd:
			pos = lipgloss.Right
		}
		result = lipgloss.JoinVertical(pos, withSpacing...)
	}

	// Apply size constraints if set
	if c.width > 0 || c.height > 0 {
		style := lipgloss.NewStyle()
		if c.width > 0 {
			style = style.Width(c.width)
		}
		if c.height > 0 {
			style = style.Height(c.height)
		}
		result = style.Render(result)
	}

	return result
}

// Horizontal creates a horizontal layout
func Horizontal(children ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, children...)
}

// Vertical creates a vertical layout
func Vertical(children ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, children...)
}

// HorizontalCenter creates a horizontally centered layout
func HorizontalCenter(children ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Center, children...)
}

// VerticalCenter creates a vertically centered layout
func VerticalCenter(children ...string) string {
	return lipgloss.JoinVertical(lipgloss.Center, children...)
}

// Grid creates a grid layout
func Grid(columns int, children ...string) string {
	if columns <= 0 || len(children) == 0 {
		return ""
	}

	var rows []string
	for i := 0; i < len(children); i += columns {
		end := i + columns
		if end > len(children) {
			end = len(children)
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, children[i:end]...)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// Flexbox creates a flexible layout
type Flexbox struct {
	direction Direction
	wrap      bool
	justify   Alignment
	align     Alignment
	gap       int
	children  []FlexItem
	width     int
	height    int
}

// FlexItem represents an item in a flexbox
type FlexItem struct {
	Content string
	Grow    int
	Shrink  int
	Basis   int
}

// NewFlexbox creates a new flexbox layout
func NewFlexbox(direction Direction) *Flexbox {
	return &Flexbox{
		direction: direction,
		justify:   AlignStart,
		align:     AlignStart,
		children:  make([]FlexItem, 0),
	}
}

// SetWrap sets whether items should wrap
func (f *Flexbox) SetWrap(wrap bool) {
	f.wrap = wrap
}

// SetJustify sets the justify alignment
func (f *Flexbox) SetJustify(justify Alignment) {
	f.justify = justify
}

// SetAlign sets the align alignment
func (f *Flexbox) SetAlign(align Alignment) {
	f.align = align
}

// SetGap sets the gap between items
func (f *Flexbox) SetGap(gap int) {
	f.gap = gap
}

// SetSize sets the flexbox size
func (f *Flexbox) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// AddItem adds an item to the flexbox
func (f *Flexbox) AddItem(item FlexItem) {
	f.children = append(f.children, item)
}

// Render renders the flexbox
func (f *Flexbox) Render() string {
	if len(f.children) == 0 {
		return ""
	}

	// Simple flexbox implementation
	// For a full implementation, you'd need to calculate flex-grow, flex-shrink, etc.
	var items []string
	for _, child := range f.children {
		items = append(items, child.Content)
	}

	// Add gaps
	if f.gap > 0 && len(items) > 1 {
		var withGaps []string
		spacer := strings.Repeat(" ", f.gap)
		for i, item := range items {
			withGaps = append(withGaps, item)
			if i < len(items)-1 {
				withGaps = append(withGaps, spacer)
			}
		}
		items = withGaps
	}

	var result string
	if f.direction == DirectionHorizontal {
		var pos lipgloss.Position
		switch f.align {
		case AlignCenter:
			pos = lipgloss.Center
		case AlignEnd:
			pos = lipgloss.Bottom
		default:
			pos = lipgloss.Top
		}
		result = lipgloss.JoinHorizontal(pos, items...)
	} else {
		var pos lipgloss.Position
		switch f.justify {
		case AlignCenter:
			pos = lipgloss.Center
		case AlignEnd:
			pos = lipgloss.Right
		default:
			pos = lipgloss.Left
		}
		result = lipgloss.JoinVertical(pos, items...)
	}

	return result
}
