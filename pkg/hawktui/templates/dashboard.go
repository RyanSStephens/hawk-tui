package templates

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/components"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/layouts"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Widget represents a dashboard widget
type Widget struct {
	Title   string
	Content string
	Width   int
	Height  int
}

// Dashboard represents a dashboard template
type Dashboard struct {
	theme   *styles.Theme
	title   string
	widgets []Widget
	width   int
	height  int
	columns int
}

// NewDashboard creates a new dashboard template
func NewDashboard() *Dashboard {
	return &Dashboard{
		theme:   styles.DefaultTheme(),
		columns: 2,
		widgets: make([]Widget, 0),
	}
}

// NewDashboardWithTheme creates a new dashboard with a specific theme
func NewDashboardWithTheme(theme *styles.Theme) *Dashboard {
	return &Dashboard{
		theme:   theme,
		columns: 2,
		widgets: make([]Widget, 0),
	}
}

// SetTitle sets the dashboard title
func (d *Dashboard) SetTitle(title string) {
	d.title = title
}

// SetColumns sets the number of columns
func (d *Dashboard) SetColumns(columns int) {
	d.columns = columns
}

// AddWidget adds a widget to the dashboard
func (d *Dashboard) AddWidget(widget Widget) {
	d.widgets = append(d.widgets, widget)
}

// SetSize sets the dashboard dimensions
func (d *Dashboard) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// Init implements tea.Model
func (d *Dashboard) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return d, nil
}

// View implements tea.Model
func (d *Dashboard) View() string {
	var sections []string

	// Render title
	if d.title != "" {
		titleStyle := lipgloss.NewStyle().
			Foreground(d.theme.Primary).
			Bold(true).
			Padding(1, 2).
			Background(d.theme.BgSecondary).
			Width(d.width)

		sections = append(sections, titleStyle.Render(d.title))
	}

	// Render widgets in grid layout
	if len(d.widgets) > 0 {
		var widgetViews []string
		for _, widget := range d.widgets {
			panel := components.NewPanelWithTheme(widget.Title, d.theme)
			panel.SetContent(widget.Content)
			if widget.Width > 0 {
				panel.SetSize(widget.Width, widget.Height)
			}
			widgetViews = append(widgetViews, panel.View())
		}

		grid := layouts.Grid(d.columns, widgetViews...)
		sections = append(sections, grid)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// Example helper functions for creating common widgets

// CreateMetricWidget creates a metric display widget
func CreateMetricWidget(title, value, unit string, theme *styles.Theme) Widget {
	if theme == nil {
		theme = styles.DefaultTheme()
	}

	valueStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 0)

	unitStyle := lipgloss.NewStyle().
		Foreground(theme.TextMuted).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		valueStyle.Render(value),
		unitStyle.Render(unit),
	)

	return Widget{
		Title:   title,
		Content: content,
		Width:   25,
		Height:  8,
	}
}

// CreateStatusWidget creates a status display widget
func CreateStatusWidget(title, status, description string, theme *styles.Theme) Widget {
	if theme == nil {
		theme = styles.DefaultTheme()
	}

	var statusColor styles.Color
	switch status {
	case "healthy", "online", "active":
		statusColor = theme.Success
	case "warning", "degraded":
		statusColor = theme.Warning
	case "error", "offline", "down":
		statusColor = theme.Error
	default:
		statusColor = theme.TextMuted
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		Bold(true).
		Padding(0, 0, 1, 0)

	descStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		statusStyle.Render(status),
		descStyle.Render(description),
	)

	return Widget{
		Title:   title,
		Content: content,
		Width:   25,
		Height:  8,
	}
}

// CreateChartWidget creates a simple bar chart widget
func CreateChartWidget(title string, data map[string]float64, max float64, theme *styles.Theme) Widget {
	if theme == nil {
		theme = styles.DefaultTheme()
	}

	var lines []string
	barWidth := 20

	for label, value := range data {
		percent := value / max
		fillWidth := int(float64(barWidth) * percent)
		emptyWidth := barWidth - fillWidth

		fillStyle := lipgloss.NewStyle().Foreground(theme.Primary)
		emptyStyle := lipgloss.NewStyle().Foreground(theme.TextMuted)

		bar := fillStyle.Render(lipgloss.NewStyle().Width(fillWidth).Render("█")) +
			emptyStyle.Render(lipgloss.NewStyle().Width(emptyWidth).Render("░"))

		labelStyle := lipgloss.NewStyle().
			Foreground(theme.TextPrimary).
			Width(15)

		valueStyle := lipgloss.NewStyle().
			Foreground(theme.TextSecondary).
			Width(6).
			Align(lipgloss.Right)

		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStyle.Render(label),
			bar,
			valueStyle.Render(fmt.Sprintf("%.1f", value)),
		)

		lines = append(lines, line)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	return Widget{
		Title:   title,
		Content: content,
		Width:   50,
		Height:  len(lines) + 4,
	}
}
