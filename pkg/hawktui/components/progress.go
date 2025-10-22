package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// ProgressBar represents a progress indicator
type ProgressBar struct {
	*BaseComponent
	current     float64
	max         float64
	label       string
	showValue   bool
	showPercent bool
	barWidth    int
	fillChar    string
	emptyChar   string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(max float64) *ProgressBar {
	return &ProgressBar{
		BaseComponent: NewBaseComponent(nil),
		max:           max,
		barWidth:      30,
		fillChar:      "█",
		emptyChar:     "░",
		showPercent:   true,
	}
}

// NewProgressBarWithTheme creates a new progress bar with a specific theme
func NewProgressBarWithTheme(max float64, theme *styles.Theme) *ProgressBar {
	return &ProgressBar{
		BaseComponent: NewBaseComponent(theme),
		max:           max,
		barWidth:      30,
		fillChar:      "█",
		emptyChar:     "░",
		showPercent:   true,
	}
}

// SetValue sets the current progress value
func (p *ProgressBar) SetValue(value float64) {
	if value < 0 {
		value = 0
	}
	if value > p.max {
		value = p.max
	}
	p.current = value
}

// SetMax sets the maximum value
func (p *ProgressBar) SetMax(max float64) {
	p.max = max
}

// SetLabel sets the progress bar label
func (p *ProgressBar) SetLabel(label string) {
	p.label = label
}

// SetShowValue sets whether to show the numeric value
func (p *ProgressBar) SetShowValue(show bool) {
	p.showValue = show
}

// SetShowPercent sets whether to show the percentage
func (p *ProgressBar) SetShowPercent(show bool) {
	p.showPercent = show
}

// SetBarWidth sets the width of the progress bar
func (p *ProgressBar) SetBarWidth(width int) {
	p.barWidth = width
}

// Percent returns the current percentage
func (p *ProgressBar) Percent() float64 {
	if p.max == 0 {
		return 0
	}
	return (p.current / p.max) * 100
}

// IsComplete returns whether the progress is complete
func (p *ProgressBar) IsComplete() bool {
	return p.current >= p.max
}

// Init implements tea.Model
func (p *ProgressBar) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (p *ProgressBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, nil
}

// View implements tea.Model
func (p *ProgressBar) View() string {
	if !p.visible {
		return ""
	}

	percent := p.Percent()
	fillWidth := int(float64(p.barWidth) * (percent / 100.0))
	emptyWidth := p.barWidth - fillWidth

	// Create the bar
	fillStyle := lipgloss.NewStyle().
		Foreground(p.theme.Success)

	emptyStyle := lipgloss.NewStyle().
		Foreground(p.theme.TextMuted)

	filled := fillStyle.Render(strings.Repeat(p.fillChar, fillWidth))
	empty := emptyStyle.Render(strings.Repeat(p.emptyChar, emptyWidth))
	bar := filled + empty

	// Add brackets
	bar = "[" + bar + "]"

	// Add label if present
	var parts []string
	if p.label != "" {
		labelStyle := lipgloss.NewStyle().
			Foreground(p.theme.TextPrimary).
			Bold(true)
		parts = append(parts, labelStyle.Render(p.label))
	}

	parts = append(parts, bar)

	// Add percentage/value
	if p.showPercent || p.showValue {
		var info string
		if p.showPercent && p.showValue {
			info = fmt.Sprintf("%.1f%% (%.0f/%.0f)", percent, p.current, p.max)
		} else if p.showPercent {
			info = fmt.Sprintf("%.1f%%", percent)
		} else {
			info = fmt.Sprintf("%.0f/%.0f", p.current, p.max)
		}

		infoStyle := lipgloss.NewStyle().
			Foreground(p.theme.TextSecondary)
		parts = append(parts, infoStyle.Render(info))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}
