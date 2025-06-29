package components

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/types"
)

// DashboardView component for displaying dashboard widgets
type DashboardView struct {
	// Layout
	width  int
	height int
	styles interface{}

	// Data
	widgets         map[string]types.DashboardParams
	filteredWidgets []string

	// State
	filter         string
	selectedWidget int
	gridCols       int
	gridRows       int

	// Layout mode
	layoutMode DashboardLayout
}

// DashboardLayout represents different dashboard layout modes
type DashboardLayout int

const (
	DashboardLayoutAuto DashboardLayout = iota
	DashboardLayoutGrid
	DashboardLayoutCustom
)

// NewDashboardView creates a new dashboard view component
func NewDashboardView(styles interface{}) *DashboardView {
	return &DashboardView{
		styles:          styles,
		widgets:         make(map[string]types.DashboardParams),
		filteredWidgets: make([]string, 0),
		gridCols:        2,
		gridRows:        2,
		layoutMode:      DashboardLayoutAuto,
	}
}

// Init implements tea.Model
func (dv *DashboardView) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (dv *DashboardView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if dv.selectedWidget > 0 {
				dv.selectedWidget--
			}
		case "down", "j":
			if dv.selectedWidget < len(dv.filteredWidgets)-1 {
				dv.selectedWidget++
			}
		case "left", "h":
			if dv.gridCols > 1 {
				dv.gridCols--
				dv.updateLayout()
			}
		case "right", "l":
			if dv.gridCols < 4 {
				dv.gridCols++
				dv.updateLayout()
			}
		case "a":
			dv.layoutMode = DashboardLayoutAuto
		case "g":
			dv.layoutMode = DashboardLayoutGrid
		case "r":
			// Refresh all widgets
			return dv, dv.Refresh()
		}

	case tea.WindowSizeMsg:
		dv.SetSize(msg.Width, msg.Height)
	}

	return dv, nil
}

// View implements tea.Model
func (dv *DashboardView) View() string {
	if dv.width == 0 || dv.height == 0 {
		return "Loading dashboard..."
	}

	if len(dv.filteredWidgets) == 0 {
		return dv.renderEmptyState()
	}

	switch dv.layoutMode {
	case DashboardLayoutGrid:
		return dv.renderGridLayout()
	case DashboardLayoutCustom:
		return dv.renderCustomLayout()
	default: // DashboardLayoutAuto
		return dv.renderAutoLayout()
	}
}

// SetSize sets the dimensions of the dashboard view
func (dv *DashboardView) SetSize(width, height int) {
	dv.width = width
	dv.height = height
	dv.updateLayout()
}

// UpdateWidget updates or adds a dashboard widget
func (dv *DashboardView) UpdateWidget(widget types.DashboardParams) {
	dv.widgets[widget.WidgetID] = widget
	dv.applyFilter()
}

// SetFilter sets the filter for widget names
func (dv *DashboardView) SetFilter(filter string) {
	dv.filter = filter
	dv.applyFilter()
}

// Refresh refreshes the dashboard view
func (dv *DashboardView) Refresh() tea.Cmd {
	// In a real implementation, this might send refresh requests
	return nil
}

// applyFilter applies the current filter to widgets
func (dv *DashboardView) applyFilter() {
	dv.filteredWidgets = make([]string, 0)

	for id, widget := range dv.widgets {
		if dv.filter == "" || 
		   strings.Contains(strings.ToLower(widget.Title), strings.ToLower(dv.filter)) ||
		   strings.Contains(strings.ToLower(id), strings.ToLower(dv.filter)) {
			dv.filteredWidgets = append(dv.filteredWidgets, id)
		}
	}

	// Sort widgets by ID for consistent ordering
	sort.Strings(dv.filteredWidgets)

	// Ensure selected widget is in bounds
	if dv.selectedWidget >= len(dv.filteredWidgets) {
		dv.selectedWidget = len(dv.filteredWidgets) - 1
	}
	if dv.selectedWidget < 0 {
		dv.selectedWidget = 0
	}
}

// updateLayout recalculates layout based on current settings
func (dv *DashboardView) updateLayout() {
	if dv.layoutMode == DashboardLayoutAuto {
		// Auto-calculate optimal grid size
		numWidgets := len(dv.filteredWidgets)
		if numWidgets == 0 {
			return
		}

		// Try to maintain roughly square aspect ratio
		cols := 1
		for cols*cols < numWidgets && cols < 4 {
			cols++
		}
		dv.gridCols = cols
		dv.gridRows = (numWidgets + cols - 1) / cols
	}
}

// renderAutoLayout renders widgets in automatic layout
func (dv *DashboardView) renderAutoLayout() string {
	header := dv.renderHeader("Dashboard (Auto)")
	content := dv.renderWidgetGrid()
	footer := dv.renderControls()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderGridLayout renders widgets in fixed grid layout
func (dv *DashboardView) renderGridLayout() string {
	header := dv.renderHeader(fmt.Sprintf("Dashboard (Grid %dx%d)", dv.gridCols, dv.gridRows))
	content := dv.renderWidgetGrid()
	footer := dv.renderControls()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderCustomLayout renders widgets using their custom layout settings
func (dv *DashboardView) renderCustomLayout() string {
	header := dv.renderHeader("Dashboard (Custom)")
	content := dv.renderCustomWidgetLayout()
	footer := dv.renderControls()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderWidgetGrid renders widgets in a grid
func (dv *DashboardView) renderWidgetGrid() string {
	if len(dv.filteredWidgets) == 0 {
		return dv.renderEmptyState()
	}

	availableHeight := dv.height - 6 // Account for header and footer
	widgetHeight := availableHeight / dv.gridRows
	if widgetHeight < 6 {
		widgetHeight = 6
	}

	widgetWidth := (dv.width - dv.gridCols - 1) / dv.gridCols
	if widgetWidth < 20 {
		widgetWidth = 20
	}

	var rows []string
	for row := 0; row < dv.gridRows; row++ {
		var rowWidgets []string
		for col := 0; col < dv.gridCols; col++ {
			widgetIndex := row*dv.gridCols + col
			if widgetIndex < len(dv.filteredWidgets) {
				widgetID := dv.filteredWidgets[widgetIndex]
				widget := dv.widgets[widgetID]
				selected := widgetIndex == dv.selectedWidget
				
				renderedWidget := dv.renderWidget(widget, widgetWidth, widgetHeight, selected)
				rowWidgets = append(rowWidgets, renderedWidget)
			} else {
				// Empty space
				emptyWidget := lipgloss.NewStyle().
					Width(widgetWidth).
					Height(widgetHeight).
					Render("")
				rowWidgets = append(rowWidgets, emptyWidget)
			}
		}
		if len(rowWidgets) > 0 {
			row := lipgloss.JoinHorizontal(lipgloss.Top, rowWidgets...)
			rows = append(rows, row)
		}
	}

	return strings.Join(rows, "\n")
}

// renderCustomWidgetLayout renders widgets using their custom layout positions
func (dv *DashboardView) renderCustomWidgetLayout() string {
	// This is a simplified implementation
	// In a full implementation, you'd create a 2D grid and place widgets at their specified positions
	return dv.renderWidgetGrid() // Fall back to grid for now
}

// renderWidget renders a single dashboard widget
func (dv *DashboardView) renderWidget(widget types.DashboardParams, width, height int, selected bool) string {
	borderStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1).
		Border(lipgloss.RoundedBorder())

	if selected {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("#00D7FF"))
	} else {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("#414868"))
	}

	// Render title
	title := widget.Title
	if title == "" {
		title = widget.WidgetID
	}
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5")).
		Width(width - 4).
		Align(lipgloss.Center)

	// Render content based on widget type
	var content string
	switch widget.Type {
	case types.WidgetTypeStatusGrid:
		content = dv.renderStatusGrid(widget, width-4, height-4)
	case types.WidgetTypeGauge:
		content = dv.renderGauge(widget, width-4, height-4)
	case types.WidgetTypeTable:
		content = dv.renderTable(widget, width-4, height-4)
	case types.WidgetTypeText:
		content = dv.renderText(widget, width-4, height-4)
	case types.WidgetTypeMetricChart:
		content = dv.renderMetricChart(widget, width-4, height-4)
	default:
		content = dv.renderUnsupportedWidget(widget, width-4, height-4)
	}

	widgetContent := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(title),
		content,
	)

	return borderStyle.Render(widgetContent)
}

// renderStatusGrid renders a status grid widget
func (dv *DashboardView) renderStatusGrid(widget types.DashboardParams, width, height int) string {
	if widget.Data == nil {
		return dv.renderNoData(width, height)
	}

	// Try to parse as StatusGridData
	var lines []string
	
	// In a real implementation, you'd properly unmarshal the data
	dataStr := fmt.Sprintf("%v", widget.Data)
	if len(dataStr) > width-4 {
		dataStr = dataStr[:width-7] + "..."
	}
	
	lines = append(lines, "Service Status:")
	lines = append(lines, dataStr)

	content := strings.Join(lines, "\n")
	
	return lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Render(content)
}

// renderGauge renders a gauge widget
func (dv *DashboardView) renderGauge(widget types.DashboardParams, width, height int) string {
	if widget.Data == nil {
		return dv.renderNoData(width, height)
	}

	// Simple gauge representation
	// In a real implementation, you'd parse types.GaugeData
	value := 75.0 // Default value
	
	// Create a simple gauge bar
	gaugeWidth := width - 2
	if gaugeWidth < 10 {
		gaugeWidth = 10
	}
	
	percentage := value / 100.0
	filled := int(percentage * float64(gaugeWidth))
	empty := gaugeWidth - filled
	
	gauge := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	
	gaugeStyle := lipgloss.NewStyle().
		Foreground(dv.getGaugeColor(percentage))
	
	valueText := fmt.Sprintf("%.1f%%", value)
	
	return lipgloss.JoinVertical(
		lipgloss.Center,
		gaugeStyle.Render(gauge),
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(valueText),
	)
}

// renderTable renders a table widget
func (dv *DashboardView) renderTable(widget types.DashboardParams, width, height int) string {
	if widget.Data == nil {
		return dv.renderNoData(width, height)
	}

	// Simple table representation
	// In a real implementation, you'd parse types.TableData
	return lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Align(lipgloss.Center).
		Render("Table Data\n(Implementation pending)")
}

// renderText renders a text widget
func (dv *DashboardView) renderText(widget types.DashboardParams, width, height int) string {
	if widget.Data == nil {
		return dv.renderNoData(width, height)
	}

	// Simple text display
	// In a real implementation, you'd parse types.TextData
	content := fmt.Sprintf("%v", widget.Data)
	
	// Word wrap the content
	wrapped := dv.wordWrap(content, width)
	
	return lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Render(wrapped)
}

// renderMetricChart renders a metric chart widget
func (dv *DashboardView) renderMetricChart(widget types.DashboardParams, width, height int) string {
	if widget.Data == nil {
		return dv.renderNoData(width, height)
	}

	// Simple chart representation
	// In a real implementation, you'd parse types.ChartData and render actual charts
	return lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Align(lipgloss.Center).
		Render("Chart Data\n(Implementation pending)")
}

// renderUnsupportedWidget renders an unsupported widget type
func (dv *DashboardView) renderUnsupportedWidget(widget types.DashboardParams, width, height int) string {
	message := fmt.Sprintf("Unsupported widget type:\n%s", widget.Type)
	
	return lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Foreground(lipgloss.Color("#FF6B6B")).
		Align(lipgloss.Center).
		Render(message)
}

// renderNoData renders a no data message
func (dv *DashboardView) renderNoData(width, height int) string {
	return lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Foreground(lipgloss.Color("#565F89")).
		Italic(true).
		Align(lipgloss.Center).
		Render("No data available")
}

// renderHeader renders the dashboard header
func (dv *DashboardView) renderHeader(title string) string {
	var parts []string
	parts = append(parts, title)
	
	if dv.filter != "" {
		parts = append(parts, fmt.Sprintf("(Filter: %s)", dv.filter))
	}
	
	parts = append(parts, fmt.Sprintf("(%d widgets)", len(dv.filteredWidgets)))
	
	headerText := strings.Join(parts, " ")
	
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#24283B")).
		Foreground(lipgloss.Color("#C0CAF5")).
		Padding(0, 1).
		Width(dv.width).
		Bold(true).
		Render(headerText)
}

// renderControls renders the control instructions
func (dv *DashboardView) renderControls() string {
	controls := []string{
		"↑↓ Navigate",
		"← → Columns",
		"a Auto",
		"g Grid",
		"r Refresh",
	}
	
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Render(strings.Join(controls, " │ "))
}

// renderEmptyState renders the empty state
func (dv *DashboardView) renderEmptyState() string {
	message := "No dashboard widgets to display"
	if dv.filter != "" {
		message = fmt.Sprintf("No widgets match filter: %s", dv.filter)
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Italic(true).
		Align(lipgloss.Center).
		Width(dv.width).
		Height(dv.height).
		Render(message)
}

// Helper functions

// getGaugeColor returns color for gauge based on percentage
func (dv *DashboardView) getGaugeColor(percentage float64) lipgloss.Color {
	if percentage >= 0.8 {
		return lipgloss.Color("#FF6B6B") // Red
	} else if percentage >= 0.6 {
		return lipgloss.Color("#FFD93D") // Yellow
	} else {
		return lipgloss.Color("#51CF66") // Green
	}
}

// wordWrap wraps text to fit within the specified width
func (dv *DashboardView) wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// If adding this word would exceed the width, start a new line
		if currentLine.Len()+len(word)+1 > width && currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}