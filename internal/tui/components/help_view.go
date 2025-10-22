package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpView component for displaying help and keyboard shortcuts
type HelpView struct {
	// Layout
	width    int
	height   int
	styles   interface{}
	viewport viewport.Model

	// Content
	helpContent string
}

// NewHelpView creates a new help view component
func NewHelpView(styles interface{}) *HelpView {
	v := viewport.New(0, 0)
	v.Style = lipgloss.NewStyle()

	hv := &HelpView{
		styles:   styles,
		viewport: v,
	}

	hv.generateHelpContent()
	return hv
}

// Init implements tea.Model
func (hv *HelpView) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (hv *HelpView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			hv.viewport.LineUp(1)
		case "down", "j":
			hv.viewport.LineDown(1)
		case "pgup":
			hv.viewport.HalfViewUp()
		case "pgdown":
			hv.viewport.HalfViewDown()
		case "home", "g":
			hv.viewport.GotoTop()
		case "end", "G":
			hv.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
		hv.SetSize(msg.Width, msg.Height)
	}

	hv.viewport, cmd = hv.viewport.Update(msg)
	return hv, cmd
}

// View implements tea.Model
func (hv *HelpView) View() string {
	if hv.width == 0 || hv.height == 0 {
		return "Loading help..."
	}

	header := hv.renderHeader()
	footer := hv.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		hv.viewport.View(),
		footer,
	)
}

// SetSize sets the dimensions of the help view
func (hv *HelpView) SetSize(width, height int) {
	hv.width = width
	hv.height = height

	// Update viewport size (account for header and footer)
	hv.viewport.Width = width - 2   // Account for borders
	hv.viewport.Height = height - 4 // Account for header, footer, and borders

	// Update content
	hv.viewport.SetContent(hv.helpContent)
}

// generateHelpContent generates the help content
func (hv *HelpView) generateHelpContent() {
	var sections []string

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D7FF")).
		Align(lipgloss.Center).
		Render("🦅 Hawk TUI - Help & Documentation")

	sections = append(sections, title)
	sections = append(sections, "")

	// Overview
	overview := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5")).
		Render("OVERVIEW")

	sections = append(sections, overview)
	sections = append(sections, "")
	sections = append(sections, "Hawk TUI is a universal Terminal User Interface framework for monitoring")
	sections = append(sections, "applications, services, and systems in real-time. It accepts JSON-RPC")
	sections = append(sections, "messages via stdin and displays them in organized, interactive views.")
	sections = append(sections, "")

	// Global Navigation
	sections = append(sections, hv.renderSection("GLOBAL NAVIGATION", []HelpItem{
		{"q, Ctrl+C", "Quit application"},
		{"1", "Switch to Logs view"},
		{"2", "Switch to Metrics view"},
		{"3", "Switch to Dashboard view"},
		{"4", "Switch to Configuration view"},
		{"h", "Show this help"},
		{"Tab", "Navigate between panels"},
		{"Shift+Tab", "Navigate backwards between panels"},
		{"/", "Start search/filter mode"},
		{"Esc", "Exit search mode or go back"},
		{"r", "Refresh current view"},
	}))

	// Logs View
	sections = append(sections, hv.renderSection("LOGS VIEW", []HelpItem{
		{"↑↓, k j", "Navigate log entries"},
		{"PgUp/PgDn", "Page up/down through logs"},
		{"Home/End, g G", "Go to first/last log entry"},
		{"c", "Toggle context display"},
		{"s", "Toggle sidebar"},
		{"a", "Toggle auto-scroll"},
		{"d", "Filter by DEBUG level"},
		{"i", "Filter by INFO level"},
		{"w", "Filter by WARN level"},
		{"e", "Filter by ERROR level"},
		{"x", "Clear level filter"},
	}))

	// Metrics View
	sections = append(sections, hv.renderSection("METRICS VIEW", []HelpItem{
		{"↑↓, k j", "Navigate metrics"},
		{"←→, h l", "Adjust columns (grid view)"},
		{"g", "Switch to grid view"},
		{"L", "Switch to list view"},
		{"c", "Switch to chart view"},
		{"n", "Sort by name"},
		{"v", "Sort by value"},
		{"t", "Sort by last update time"},
		{"C", "Toggle chart display"},
		{"+/-", "Adjust chart height"},
	}))

	// Dashboard View
	sections = append(sections, hv.renderSection("DASHBOARD VIEW", []HelpItem{
		{"↑↓, k j", "Navigate widgets"},
		{"←→, h l", "Adjust grid columns"},
		{"a", "Auto layout mode"},
		{"g", "Grid layout mode"},
		{"r", "Refresh dashboard"},
		{"Enter", "Select/activate widget"},
	}))

	// Configuration View
	sections = append(sections, hv.renderSection("CONFIGURATION VIEW", []HelpItem{
		{"↑↓, k j", "Navigate configuration items"},
		{"←→, h l", "Navigate categories/back"},
		{"Enter", "Edit selected configuration"},
		{"r", "Reset to default value"},
		{"c", "Categories view"},
		{"L", "List view"},
		{"d", "Toggle descriptions"},
	}))

	// Editing Mode
	sections = append(sections, hv.renderSection("EDIT MODE (Configuration)", []HelpItem{
		{"Enter", "Save changes"},
		{"Esc", "Cancel editing"},
		{"Backspace", "Delete character"},
		{"←→", "Move cursor"},
		{"Home/End", "Move to start/end"},
	}))

	// Search Mode
	sections = append(sections, hv.renderSection("SEARCH/FILTER MODE", []HelpItem{
		{"/", "Start searching"},
		{"Type text", "Enter search query"},
		{"Enter", "Apply search"},
		{"Esc", "Cancel search"},
		{"Backspace", "Delete character"},
	}))

	// Status Indicators
	sections = append(sections, hv.renderSection("STATUS INDICATORS", []HelpItem{
		{"●", "Active/healthy status"},
		{"◐", "Warning/degraded status"},
		{"○", "Inactive/error status"},
		{"[R]", "Restart required (config)"},
		{"[M]", "Modified from default (config)"},
		{"SEARCH", "Search mode active"},
		{"FPS", "Real-time performance indicator"},
	}))

	// Log Levels
	sections = append(sections, hv.renderSection("LOG LEVELS", []HelpItem{
		{"DEBUG", "Detailed debugging information"},
		{"INFO", "General information messages"},
		{"WARN", "Warning conditions"},
		{"ERROR", "Error conditions"},
		{"SUCCESS", "Success/completion messages"},
	}))

	// Metric Types
	sections = append(sections, hv.renderSection("METRIC TYPES", []HelpItem{
		{"Counter", "Monotonically increasing values"},
		{"Gauge", "Point-in-time measurements"},
		{"Histogram", "Distribution of values"},
	}))

	// Widget Types
	sections = append(sections, hv.renderSection("DASHBOARD WIDGET TYPES", []HelpItem{
		{"Status Grid", "Service health overview"},
		{"Gauge", "Single metric with visual indicator"},
		{"Chart", "Time-series data visualization"},
		{"Table", "Tabular data display"},
		{"Text", "Rich text content"},
		{"Histogram", "Value distribution charts"},
	}))

	// Protocol Information
	sections = append(sections, hv.renderSection("PROTOCOL INFORMATION", []HelpItem{
		{"Transport", "JSON-RPC 2.0 over stdin/stdout"},
		{"Message Format", "Line-delimited JSON"},
		{"Supported Methods", "hawk.log, hawk.metric, hawk.config, hawk.progress, hawk.dashboard, hawk.event"},
		{"Bidirectional", "TUI can send configuration updates and commands back to client"},
		{"Error Handling", "Malformed messages are gracefully ignored"},
		{"Performance", "Optimized for high-frequency updates (60 FPS)"},
	}))

	// Examples
	sections = append(sections, hv.renderSection("EXAMPLE USAGE", []HelpItem{
		{"Log Message", `{"jsonrpc":"2.0","method":"hawk.log","params":{"message":"Server started","level":"INFO"}}`},
		{"Metric Update", `{"jsonrpc":"2.0","method":"hawk.metric","params":{"name":"cpu.usage","value":45.2,"type":"gauge","unit":"%"}}`},
		{"Configuration", `{"jsonrpc":"2.0","method":"hawk.config","params":{"key":"app.port","value":8080,"type":"integer"}}`},
		{"Progress Bar", `{"jsonrpc":"2.0","method":"hawk.progress","params":{"id":"upload","label":"Uploading","current":75,"total":100}}`},
	}))

	// Performance Tips
	sections = append(sections, hv.renderSection("PERFORMANCE TIPS", []HelpItem{
		{"Batching", "Send multiple messages in JSON arrays for better performance"},
		{"Rate Limiting", "Default limit is 1000 messages/second"},
		{"Memory Usage", "Log history is limited to 1000 entries by default"},
		{"Auto-scroll", "Disable auto-scroll in logs when reviewing historical data"},
		{"Filtering", "Use filters to reduce visual noise and improve performance"},
	}))

	// Troubleshooting
	sections = append(sections, hv.renderSection("TROUBLESHOOTING", []HelpItem{
		{"No Data", "Ensure your application is sending messages to stdout in JSON-RPC format"},
		{"Slow Performance", "Check message rate and consider batching updates"},
		{"Layout Issues", "Try resizing terminal or adjusting view modes"},
		{"Search Not Working", "Press '/' to enter search mode, type query, then Enter"},
		{"Config Changes", "Some changes require application restart (marked with [R])"},
	}))

	// About
	sections = append(sections, "")
	about := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5")).
		Render("ABOUT")

	sections = append(sections, about)
	sections = append(sections, "")
	sections = append(sections, "Hawk TUI v1.0.0")
	sections = append(sections, "Universal TUI framework for any programming language")
	sections = append(sections, "Built with Bubble Tea and Lipgloss")
	sections = append(sections, "")
	sections = append(sections, "For more information:")
	sections = append(sections, "• Documentation: https://github.com/hawk-tui/hawk-tui")
	sections = append(sections, "• Examples: See examples/ directory")
	sections = append(sections, "• Issues: https://github.com/hawk-tui/hawk-tui/issues")

	hv.helpContent = strings.Join(sections, "\n")
}

// HelpItem represents a single help item
type HelpItem struct {
	Key         string
	Description string
}

// renderSection renders a help section with items
func (hv *HelpView) renderSection(title string, items []HelpItem) string {
	var lines []string

	// Section title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5"))
	lines = append(lines, titleStyle.Render(title))
	lines = append(lines, "")

	// Items
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D7FF")).
		Width(20)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9AA5CE"))

	for _, item := range items {
		line := lipgloss.JoinHorizontal(
			lipgloss.Top,
			keyStyle.Render(item.Key),
			descStyle.Render(item.Description),
		)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

// renderHeader renders the help view header
func (hv *HelpView) renderHeader() string {
	title := "Help & Documentation"

	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#24283B")).
		Foreground(lipgloss.Color("#C0CAF5")).
		Padding(0, 1).
		Width(hv.width).
		Bold(true)

	return headerStyle.Render(title)
}

// renderFooter renders the help view footer
func (hv *HelpView) renderFooter() string {
	controls := "↑↓ Scroll │ PgUp/PgDn Page │ Home/End Top/Bottom │ Esc Close Help"

	footerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#24283B")).
		Foreground(lipgloss.Color("#565F89")).
		Padding(0, 1).
		Width(hv.width).
		Italic(true)

	return footerStyle.Render(controls)
}
