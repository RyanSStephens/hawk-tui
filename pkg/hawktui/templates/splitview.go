package templates

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// SplitOrientation defines the split orientation
type SplitOrientation int

const (
	SplitHorizontal SplitOrientation = iota
	SplitVertical
)

// SplitView represents a split view template
type SplitView struct {
	theme       *styles.Theme
	orientation SplitOrientation
	left        string
	right       string
	splitRatio  float64
	width       int
	height      int
	focusLeft   bool
}

// NewSplitView creates a new split view template
func NewSplitView(orientation SplitOrientation) *SplitView {
	return &SplitView{
		theme:       styles.DefaultTheme(),
		orientation: orientation,
		splitRatio:  0.5,
		focusLeft:   true,
	}
}

// NewSplitViewWithTheme creates a new split view with a specific theme
func NewSplitViewWithTheme(orientation SplitOrientation, theme *styles.Theme) *SplitView {
	return &SplitView{
		theme:       theme,
		orientation: orientation,
		splitRatio:  0.5,
		focusLeft:   true,
	}
}

// SetLeft sets the left/top content
func (sv *SplitView) SetLeft(content string) {
	sv.left = content
}

// SetRight sets the right/bottom content
func (sv *SplitView) SetRight(content string) {
	sv.right = content
}

// SetSplitRatio sets the split ratio (0.0 to 1.0)
func (sv *SplitView) SetSplitRatio(ratio float64) {
	if ratio < 0.1 {
		ratio = 0.1
	}
	if ratio > 0.9 {
		ratio = 0.9
	}
	sv.splitRatio = ratio
}

// SetSize sets the split view dimensions
func (sv *SplitView) SetSize(width, height int) {
	sv.width = width
	sv.height = height
}

// SetFocus sets which pane has focus
func (sv *SplitView) SetFocus(left bool) {
	sv.focusLeft = left
}

// Init implements tea.Model
func (sv *SplitView) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (sv *SplitView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			sv.focusLeft = !sv.focusLeft
		}
	}
	return sv, nil
}

// View implements tea.Model
func (sv *SplitView) View() string {
	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	rightStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	if sv.focusLeft {
		leftStyle = leftStyle.BorderForeground(sv.theme.BorderActive)
		rightStyle = rightStyle.BorderForeground(sv.theme.BorderPrimary)
	} else {
		leftStyle = leftStyle.BorderForeground(sv.theme.BorderPrimary)
		rightStyle = rightStyle.BorderForeground(sv.theme.BorderActive)
	}

	if sv.orientation == SplitHorizontal {
		// Horizontal split (side by side)
		leftWidth := int(float64(sv.width) * sv.splitRatio)
		rightWidth := sv.width - leftWidth

		leftStyle = leftStyle.Width(leftWidth - 4).Height(sv.height - 2)
		rightStyle = rightStyle.Width(rightWidth - 4).Height(sv.height - 2)

		leftPane := leftStyle.Render(sv.left)
		rightPane := rightStyle.Render(sv.right)

		return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	} else {
		// Vertical split (top and bottom)
		topHeight := int(float64(sv.height) * sv.splitRatio)
		bottomHeight := sv.height - topHeight

		leftStyle = leftStyle.Width(sv.width - 4).Height(topHeight - 2)
		rightStyle = rightStyle.Width(sv.width - 4).Height(bottomHeight - 2)

		topPane := leftStyle.Render(sv.left)
		bottomPane := rightStyle.Render(sv.right)

		return lipgloss.JoinVertical(lipgloss.Left, topPane, bottomPane)
	}
}
