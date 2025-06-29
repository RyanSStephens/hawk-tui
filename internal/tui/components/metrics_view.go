package components

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/types"
)

// MetricsView component for displaying metrics and gauges
type MetricsView struct {
	// Layout
	width  int
	height int
	styles interface{}

	// Data
	metrics         map[string]types.MetricParams
	metricHistory   map[string][]MetricPoint
	filteredMetrics []string

	// State
	filter         string
	selectedMetric int
	viewMode       MetricViewMode
	maxHistory     int
	sortBy         MetricSortBy

	// Display options
	columns     int
	showCharts  bool
	chartHeight int
}

// MetricViewMode represents different ways to view metrics
type MetricViewMode int

const (
	MetricViewModeGrid MetricViewMode = iota
	MetricViewModeList
	MetricViewModeChart
)

// MetricSortBy represents sorting options for metrics
type MetricSortBy int

const (
	MetricSortByName MetricSortBy = iota
	MetricSortByValue
	MetricSortByLastUpdate
)

// MetricPoint represents a historical metric data point
type MetricPoint struct {
	Value     float64
	Timestamp time.Time
}

// NewMetricsView creates a new metrics view component
func NewMetricsView(styles interface{}) *MetricsView {
	return &MetricsView{
		styles:          styles,
		metrics:         make(map[string]types.MetricParams),
		metricHistory:   make(map[string][]MetricPoint),
		filteredMetrics: make([]string, 0),
		maxHistory:      100,
		columns:         3,
		showCharts:      true,
		chartHeight:     8,
		viewMode:        MetricViewModeGrid,
		sortBy:          MetricSortByName,
	}
}

// Init implements tea.Model
func (mv *MetricsView) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (mv *MetricsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if mv.selectedMetric > 0 {
				mv.selectedMetric--
			}
		case "down", "j":
			if mv.selectedMetric < len(mv.filteredMetrics)-1 {
				mv.selectedMetric++
			}
		case "left", "h":
			if mv.columns > 1 {
				mv.columns--
			}
		case "right", "l":
			if mv.columns < 6 {
				mv.columns++
			}
		case "g":
			mv.viewMode = MetricViewModeGrid
		case "L": // Shift+L for list
			mv.viewMode = MetricViewModeList
		case "c":
			mv.viewMode = MetricViewModeChart
		case "n":
			mv.sortBy = MetricSortByName
			mv.applyFilter()
		case "v":
			mv.sortBy = MetricSortByValue
			mv.applyFilter()
		case "t":
			mv.sortBy = MetricSortByLastUpdate
			mv.applyFilter()
		case "C": // Shift+C for toggle charts
			mv.showCharts = !mv.showCharts
		case "+", "=":
			if mv.chartHeight < 20 {
				mv.chartHeight++
			}
		case "-":
			if mv.chartHeight > 4 {
				mv.chartHeight--
			}
		}

	case tea.WindowSizeMsg:
		mv.SetSize(msg.Width, msg.Height)
	}

	return mv, nil
}

// View implements tea.Model
func (mv *MetricsView) View() string {
	if mv.width == 0 || mv.height == 0 {
		return "Loading metrics..."
	}

	switch mv.viewMode {
	case MetricViewModeGrid:
		return mv.renderGridView()
	case MetricViewModeList:
		return mv.renderListView()
	case MetricViewModeChart:
		return mv.renderChartView()
	default:
		return mv.renderGridView()
	}
}

// SetSize sets the dimensions of the metrics view
func (mv *MetricsView) SetSize(width, height int) {
	mv.width = width
	mv.height = height
	
	// Adjust columns based on width
	maxColumns := mv.width / 25 // Minimum 25 chars per column
	if mv.columns > maxColumns {
		mv.columns = maxColumns
	}
	if mv.columns < 1 {
		mv.columns = 1
	}
}

// UpdateMetric updates or adds a metric
func (mv *MetricsView) UpdateMetric(metric types.MetricParams) {
	// Set timestamp if not provided
	if metric.Timestamp == nil {
		now := time.Now()
		metric.Timestamp = &now
	}

	// Store metric
	mv.metrics[metric.Name] = metric

	// Add to history
	point := MetricPoint{
		Value:     metric.Value,
		Timestamp: *metric.Timestamp,
	}

	if history, exists := mv.metricHistory[metric.Name]; exists {
		mv.metricHistory[metric.Name] = append(history, point)
	} else {
		mv.metricHistory[metric.Name] = []MetricPoint{point}
	}

	// Limit history size
	if len(mv.metricHistory[metric.Name]) > mv.maxHistory {
		mv.metricHistory[metric.Name] = mv.metricHistory[metric.Name][1:]
	}

	mv.applyFilter()
}

// SetFilter sets the filter for metric names
func (mv *MetricsView) SetFilter(filter string) {
	mv.filter = filter
	mv.applyFilter()
}

// Refresh refreshes the metrics view
func (mv *MetricsView) Refresh() tea.Cmd {
	return nil
}

// applyFilter applies the current filter and sorting
func (mv *MetricsView) applyFilter() {
	mv.filteredMetrics = make([]string, 0)

	for name := range mv.metrics {
		if mv.filter == "" || strings.Contains(strings.ToLower(name), strings.ToLower(mv.filter)) {
			mv.filteredMetrics = append(mv.filteredMetrics, name)
		}
	}

	// Sort metrics
	sort.Slice(mv.filteredMetrics, func(i, j int) bool {
		nameI := mv.filteredMetrics[i]
		nameJ := mv.filteredMetrics[j]
		metricI := mv.metrics[nameI]
		metricJ := mv.metrics[nameJ]

		switch mv.sortBy {
		case MetricSortByValue:
			return metricI.Value > metricJ.Value
		case MetricSortByLastUpdate:
			if metricI.Timestamp == nil && metricJ.Timestamp == nil {
				return nameI < nameJ
			}
			if metricI.Timestamp == nil {
				return false
			}
			if metricJ.Timestamp == nil {
				return true
			}
			return metricI.Timestamp.After(*metricJ.Timestamp)
		default: // MetricSortByName
			return nameI < nameJ
		}
	})

	// Ensure selected metric is in bounds
	if mv.selectedMetric >= len(mv.filteredMetrics) {
		mv.selectedMetric = len(mv.filteredMetrics) - 1
	}
	if mv.selectedMetric < 0 {
		mv.selectedMetric = 0
	}
}

// renderGridView renders metrics in a grid layout
func (mv *MetricsView) renderGridView() string {
	if len(mv.filteredMetrics) == 0 {
		return mv.renderEmptyState()
	}

	header := mv.renderHeader("Metrics Grid")
	
	var rows []string
	var currentRow []string

	for i, name := range mv.filteredMetrics {
		metric := mv.metrics[name]
		card := mv.renderMetricCard(metric, i == mv.selectedMetric)
		currentRow = append(currentRow, card)

		if len(currentRow) >= mv.columns || i == len(mv.filteredMetrics)-1 {
			row := lipgloss.JoinHorizontal(lipgloss.Top, currentRow...)
			rows = append(rows, row)
			currentRow = []string{}
		}
	}

	content := strings.Join(rows, "\n")
	
	// Add controls footer
	footer := mv.renderGridControls()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderListView renders metrics in a list layout
func (mv *MetricsView) renderListView() string {
	if len(mv.filteredMetrics) == 0 {
		return mv.renderEmptyState()
	}

	header := mv.renderHeader("Metrics List")
	
	// Calculate column widths
	nameWidth := 30
	valueWidth := 15
	unitWidth := 10
	timeWidth := 19

	// Headers
	headerRow := mv.renderListHeader(nameWidth, valueWidth, unitWidth, timeWidth)
	
	var rows []string
	rows = append(rows, headerRow)

	for i, name := range mv.filteredMetrics {
		metric := mv.metrics[name]
		row := mv.renderMetricRow(metric, i == mv.selectedMetric, nameWidth, valueWidth, unitWidth, timeWidth)
		rows = append(rows, row)
	}

	content := strings.Join(rows, "\n")
	footer := mv.renderListControls()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderChartView renders the selected metric as a chart
func (mv *MetricsView) renderChartView() string {
	if len(mv.filteredMetrics) == 0 {
		return mv.renderEmptyState()
	}

	selectedName := mv.filteredMetrics[mv.selectedMetric]
	metric := mv.metrics[selectedName]
	history := mv.metricHistory[selectedName]

	header := mv.renderHeader(fmt.Sprintf("Chart: %s", selectedName))
	
	var content string
	if len(history) > 1 {
		content = mv.renderChart(selectedName, history)
	} else {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565F89")).
			Italic(true).
			Render("Not enough data points for chart (need at least 2)")
	}

	// Add current value
	valueDisplay := mv.renderCurrentValue(metric)
	
	footer := mv.renderChartControls()

	return lipgloss.JoinVertical(lipgloss.Left, header, valueDisplay, content, footer)
}

// renderMetricCard renders a single metric as a card
func (mv *MetricsView) renderMetricCard(metric types.MetricParams, selected bool) string {
	cardWidth := (mv.width / mv.columns) - 2
	if cardWidth < 20 {
		cardWidth = 20
	}

	// Card style
	cardStyle := lipgloss.NewStyle().
		Width(cardWidth).
		Height(6).
		Padding(1).
		Margin(0, 1, 1, 0).
		Border(lipgloss.RoundedBorder())

	if selected {
		cardStyle = cardStyle.BorderForeground(lipgloss.Color("#00D7FF"))
	} else {
		cardStyle = cardStyle.BorderForeground(lipgloss.Color("#414868"))
	}

	// Name (truncated if too long)
	name := metric.Name
	if len(name) > cardWidth-4 {
		name = name[:cardWidth-7] + "..."
	}
	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5")).
		Width(cardWidth - 2).
		Align(lipgloss.Center)

	// Value with unit
	valueText := mv.formatValue(metric.Value)
	if metric.Unit != "" {
		valueText += " " + metric.Unit
	}
	
	valueStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(mv.getValueColor(metric)).
		Width(cardWidth - 2).
		Align(lipgloss.Center)

	// Gauge (if it's a gauge type)
	var gauge string
	if metric.Type == types.MetricTypeGauge {
		gauge = mv.renderMiniGauge(metric.Value, 0, 100, cardWidth-2)
	}

	// Last update
	lastUpdate := "Unknown"
	if metric.Timestamp != nil {
		lastUpdate = metric.Timestamp.Format("15:04:05")
	}
	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Width(cardWidth - 2).
		Align(lipgloss.Center)

	var content []string
	content = append(content, nameStyle.Render(name))
	content = append(content, valueStyle.Render(valueText))
	if gauge != "" {
		content = append(content, gauge)
	}
	content = append(content, timeStyle.Render(lastUpdate))

	return cardStyle.Render(strings.Join(content, "\n"))
}

// renderMetricRow renders a single metric as a table row
func (mv *MetricsView) renderMetricRow(metric types.MetricParams, selected bool, nameWidth, valueWidth, unitWidth, timeWidth int) string {
	rowStyle := lipgloss.NewStyle().Padding(0, 1)
	
	if selected {
		rowStyle = rowStyle.Background(lipgloss.Color("#3D59A1"))
	}

	// Name
	name := metric.Name
	if len(name) > nameWidth {
		name = name[:nameWidth-3] + "..."
	}
	nameCell := lipgloss.NewStyle().Width(nameWidth).Render(name)

	// Value
	valueText := mv.formatValue(metric.Value)
	valueCell := lipgloss.NewStyle().
		Width(valueWidth).
		Foreground(mv.getValueColor(metric)).
		Bold(true).
		Align(lipgloss.Right).
		Render(valueText)

	// Unit
	unit := metric.Unit
	if len(unit) > unitWidth {
		unit = unit[:unitWidth-3] + "..."
	}
	unitCell := lipgloss.NewStyle().Width(unitWidth).Render(unit)

	// Timestamp
	timestamp := "Unknown"
	if metric.Timestamp != nil {
		timestamp = metric.Timestamp.Format("2006-01-02 15:04:05")
	}
	timeCell := lipgloss.NewStyle().
		Width(timeWidth).
		Foreground(lipgloss.Color("#565F89")).
		Render(timestamp)

	row := lipgloss.JoinHorizontal(lipgloss.Top, nameCell, valueCell, unitCell, timeCell)
	return rowStyle.Render(row)
}

// renderChart renders a simple ASCII chart
func (mv *MetricsView) renderChart(name string, history []MetricPoint) string {
	if len(history) < 2 {
		return "Insufficient data"
	}

	chartWidth := mv.width - 10
	if chartWidth < 40 {
		chartWidth = 40
	}

	// Find min/max values
	minVal, maxVal := history[0].Value, history[0].Value
	for _, point := range history {
		if point.Value < minVal {
			minVal = point.Value
		}
		if point.Value > maxVal {
			maxVal = point.Value
		}
	}

	// Add some padding to the range
	valueRange := maxVal - minVal
	if valueRange == 0 {
		valueRange = 1
	}
	padding := valueRange * 0.1
	minVal -= padding
	maxVal += padding

	var lines []string

	// Y-axis labels and chart lines
	for row := mv.chartHeight - 1; row >= 0; row-- {
		// Calculate value for this row
		rowValue := minVal + (maxVal-minVal)*float64(row)/float64(mv.chartHeight-1)
		
		// Y-axis label
		label := fmt.Sprintf("%8.1f │", rowValue)
		
		// Chart line
		var chartLine strings.Builder
		for i := 0; i < len(history) && i < chartWidth; i++ {
			point := history[len(history)-chartWidth+i:]
			if i >= len(point) {
				break
			}
			
			// Normalize value to chart height
			normalizedValue := (point[i].Value - minVal) / (maxVal - minVal)
			pointRow := int(normalizedValue * float64(mv.chartHeight-1))
			
			if pointRow == row {
				chartLine.WriteString("●")
			} else if pointRow < row {
				chartLine.WriteString(" ")
			} else {
				chartLine.WriteString(" ")
			}
		}
		
		line := label + chartLine.String()
		lines = append(lines, line)
	}

	// X-axis
	xAxis := strings.Repeat(" ", 9) + "└" + strings.Repeat("─", chartWidth)
	lines = append(lines, xAxis)

	return strings.Join(lines, "\n")
}

// renderMiniGauge renders a small gauge bar
func (mv *MetricsView) renderMiniGauge(value, min, max float64, width int) string {
	if width <= 0 {
		return ""
	}

	percentage := (value - min) / (max - min)
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 1 {
		percentage = 1
	}

	filled := int(percentage * float64(width))
	empty := width - filled

	gauge := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	
	gaugeStyle := lipgloss.NewStyle().
		Foreground(mv.getGaugeColor(percentage))
	
	return gaugeStyle.Render(gauge)
}

// formatValue formats a metric value for display
func (mv *MetricsView) formatValue(value float64) string {
	if math.Abs(value) >= 1e9 {
		return fmt.Sprintf("%.1fG", value/1e9)
	} else if math.Abs(value) >= 1e6 {
		return fmt.Sprintf("%.1fM", value/1e6)
	} else if math.Abs(value) >= 1e3 {
		return fmt.Sprintf("%.1fK", value/1e3)
	} else if math.Abs(value) >= 100 {
		return fmt.Sprintf("%.0f", value)
	} else if math.Abs(value) >= 10 {
		return fmt.Sprintf("%.1f", value)
	} else {
		return fmt.Sprintf("%.2f", value)
	}
}

// getValueColor returns appropriate color for metric value
func (mv *MetricsView) getValueColor(metric types.MetricParams) lipgloss.Color {
	switch metric.Type {
	case types.MetricTypeCounter:
		return lipgloss.Color("#51CF66") // Green for counters
	case types.MetricTypeGauge:
		// Color based on value (assuming 0-100 range for gauges)
		if metric.Value >= 80 {
			return lipgloss.Color("#FF6B6B") // Red
		} else if metric.Value >= 60 {
			return lipgloss.Color("#FFD93D") // Yellow
		} else {
			return lipgloss.Color("#51CF66") // Green
		}
	default:
		return lipgloss.Color("#00D7FF") // Default cyan
	}
}

// getGaugeColor returns color for gauge fill based on percentage
func (mv *MetricsView) getGaugeColor(percentage float64) lipgloss.Color {
	if percentage >= 0.8 {
		return lipgloss.Color("#FF6B6B") // Red
	} else if percentage >= 0.6 {
		return lipgloss.Color("#FFD93D") // Yellow
	} else {
		return lipgloss.Color("#51CF66") // Green
	}
}

// renderHeader renders the view header
func (mv *MetricsView) renderHeader(title string) string {
	var parts []string
	parts = append(parts, title)
	
	// Add filter info
	if mv.filter != "" {
		parts = append(parts, fmt.Sprintf("(Filter: %s)", mv.filter))
	}
	
	// Add count
	parts = append(parts, fmt.Sprintf("(%d metrics)", len(mv.filteredMetrics)))
	
	// Add sort info
	sortInfo := ""
	switch mv.sortBy {
	case MetricSortByName:
		sortInfo = "Sort: Name"
	case MetricSortByValue:
		sortInfo = "Sort: Value"
	case MetricSortByLastUpdate:
		sortInfo = "Sort: Time"
	}
	parts = append(parts, sortInfo)
	
	headerText := strings.Join(parts, " ")
	
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#24283B")).
		Foreground(lipgloss.Color("#C0CAF5")).
		Padding(0, 1).
		Width(mv.width).
		Bold(true).
		Render(headerText)
}

// renderListHeader renders the table header for list view
func (mv *MetricsView) renderListHeader(nameWidth, valueWidth, unitWidth, timeWidth int) string {
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#414868")).
		Foreground(lipgloss.Color("#C0CAF5")).
		Bold(true).
		Padding(0, 1)

	nameHeader := lipgloss.NewStyle().Width(nameWidth).Render("Metric Name")
	valueHeader := lipgloss.NewStyle().Width(valueWidth).Align(lipgloss.Right).Render("Value")
	unitHeader := lipgloss.NewStyle().Width(unitWidth).Render("Unit")
	timeHeader := lipgloss.NewStyle().Width(timeWidth).Render("Last Updated")

	header := lipgloss.JoinHorizontal(lipgloss.Top, nameHeader, valueHeader, unitHeader, timeHeader)
	return headerStyle.Render(header)
}

// renderCurrentValue renders current value display for chart view
func (mv *MetricsView) renderCurrentValue(metric types.MetricParams) string {
	valueText := mv.formatValue(metric.Value)
	if metric.Unit != "" {
		valueText += " " + metric.Unit
	}

	return lipgloss.NewStyle().
		Bold(true).
		Foreground(mv.getValueColor(metric)).
		Align(lipgloss.Center).
		Width(mv.width).
		Render(fmt.Sprintf("Current Value: %s", valueText))
}

// renderEmptyState renders the empty state
func (mv *MetricsView) renderEmptyState() string {
	message := "No metrics to display"
	if mv.filter != "" {
		message = fmt.Sprintf("No metrics match filter: %s", mv.filter)
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Italic(true).
		Align(lipgloss.Center).
		Width(mv.width).
		Height(mv.height).
		Render(message)
}

// renderGridControls renders controls for grid view
func (mv *MetricsView) renderGridControls() string {
	controls := []string{
		"↑↓ Navigate",
		"← → Columns",
		"g Grid",
		"L List", 
		"c Chart",
		"n/v/t Sort",
	}
	
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Render(strings.Join(controls, " │ "))
}

// renderListControls renders controls for list view
func (mv *MetricsView) renderListControls() string {
	controls := []string{
		"↑↓ Navigate",
		"g Grid",
		"L List",
		"c Chart",
		"n/v/t Sort",
	}
	
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Render(strings.Join(controls, " │ "))
}

// renderChartControls renders controls for chart view
func (mv *MetricsView) renderChartControls() string {
	controls := []string{
		"↑↓ Navigate",
		"g Grid",
		"L List",
		"c Chart",
		"+/- Chart Height",
	}
	
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Render(strings.Join(controls, " │ "))
}