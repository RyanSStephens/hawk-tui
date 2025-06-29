package components

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/types"
)

// LogViewer component for displaying and filtering log messages
type LogViewer struct {
	// Layout
	width    int
	height   int
	styles   interface{} // Will be *Styles from parent package

	// Components
	viewport viewport.Model
	
	// Data
	logs         []types.LogParams
	filteredLogs []types.LogParams
	events       []types.EventParams
	
	// State
	filter       string
	levelFilter  string
	showContext  bool
	autoScroll   bool
	selectedLog  int
	
	// Display options
	maxLogs      int
	showSidebar  bool
	sidebarWidth int
}

// NewLogViewer creates a new log viewer component
func NewLogViewer(styles interface{}) *LogViewer {
	v := viewport.New(0, 0)
	v.Style = lipgloss.NewStyle()
	
	return &LogViewer{
		styles:       styles,
		viewport:     v,
		logs:         make([]types.LogParams, 0),
		filteredLogs: make([]types.LogParams, 0),
		events:       make([]types.EventParams, 0),
		maxLogs:      1000,
		autoScroll:   true,
		sidebarWidth: 25,
		showSidebar:  true,
		showContext:  false,
	}
}

// Init implements tea.Model
func (lv *LogViewer) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (lv *LogViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			lv.viewport.LineUp(1)
			lv.autoScroll = false
		case "down", "j":
			lv.viewport.LineDown(1)
			lv.autoScroll = false
		case "pgup":
			lv.viewport.HalfViewUp()
			lv.autoScroll = false
		case "pgdown":
			lv.viewport.HalfViewDown()
			lv.autoScroll = false
		case "home", "g":
			lv.viewport.GotoTop()
			lv.autoScroll = false
		case "end", "G":
			lv.viewport.GotoBottom()
			lv.autoScroll = true
		case "c":
			lv.showContext = !lv.showContext
			lv.refreshView()
		case "s":
			lv.showSidebar = !lv.showSidebar
			lv.updateLayout()
		case "a":
			lv.autoScroll = !lv.autoScroll
			if lv.autoScroll {
				lv.viewport.GotoBottom()
			}
		case "d":
			lv.levelFilter = "DEBUG"
			lv.applyFilter()
		case "i":
			lv.levelFilter = "INFO"
			lv.applyFilter()
		case "w":
			lv.levelFilter = "WARN"
			lv.applyFilter()
		case "e":
			lv.levelFilter = "ERROR"
			lv.applyFilter()
		case "x":
			lv.levelFilter = "" // Clear filter
			lv.applyFilter()
		}

	case tea.WindowSizeMsg:
		lv.SetSize(msg.Width, msg.Height)
	}

	// Update viewport
	lv.viewport, cmd = lv.viewport.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if len(cmds) > 0 {
		return lv, tea.Batch(cmds...)
	}

	return lv, nil
}

// View implements tea.Model
func (lv *LogViewer) View() string {
	if lv.width == 0 || lv.height == 0 {
		return "Loading logs..."
	}

	// Update content if auto-scroll is enabled
	if lv.autoScroll && len(lv.filteredLogs) > 0 {
		lv.refreshView()
		lv.viewport.GotoBottom()
	}

	if lv.showSidebar {
		return lv.renderWithSidebar()
	}
	
	return lv.renderFullWidth()
}

// SetSize sets the dimensions of the log viewer
func (lv *LogViewer) SetSize(width, height int) {
	lv.width = width
	lv.height = height
	lv.updateLayout()
}

// AddLog adds a new log entry
func (lv *LogViewer) AddLog(log types.LogParams) {
	// Set timestamp if not provided
	if log.Timestamp == nil {
		now := time.Now()
		log.Timestamp = &now
	}
	
	// Add to logs
	lv.logs = append(lv.logs, log)
	
	// Limit log count
	if len(lv.logs) > lv.maxLogs {
		lv.logs = lv.logs[len(lv.logs)-lv.maxLogs:]
	}
	
	lv.applyFilter()
}

// AddEvent adds an event that will be displayed as a log entry
func (lv *LogViewer) AddEvent(event types.EventParams) {
	lv.events = append(lv.events, event)
	
	// Convert event to log format for display
	logLevel := types.LogLevelInfo
	switch event.Severity {
	case types.EventSeverityError, types.EventSeverityCritical:
		logLevel = types.LogLevelError
	case types.EventSeverityWarning:
		logLevel = types.LogLevelWarn
	case types.EventSeveritySuccess:
		logLevel = types.LogLevelSuccess
	}
	
	log := types.LogParams{
		Message:   fmt.Sprintf("[EVENT] %s: %s", event.Title, event.Message),
		Level:     logLevel,
		Timestamp: event.Timestamp,
		Tags:      []string{"event", event.Type},
		Context:   event.Data,
	}
	
	lv.AddLog(log)
}

// SetFilter sets the text filter for log messages
func (lv *LogViewer) SetFilter(filter string) {
	lv.filter = filter
	lv.applyFilter()
}

// Refresh refreshes the log view
func (lv *LogViewer) Refresh() tea.Cmd {
	lv.refreshView()
	return nil
}

// applyFilter applies current filters to the log list
func (lv *LogViewer) applyFilter() {
	lv.filteredLogs = make([]types.LogParams, 0)
	
	for _, log := range lv.logs {
		// Apply level filter
		if lv.levelFilter != "" && string(log.Level) != lv.levelFilter {
			continue
		}
		
		// Apply text filter
		if lv.filter != "" {
			filterLower := strings.ToLower(lv.filter)
			messageLower := strings.ToLower(log.Message)
			componentLower := strings.ToLower(log.Component)
			
			// Check message, component, and tags
			matches := strings.Contains(messageLower, filterLower) ||
				strings.Contains(componentLower, filterLower)
			
			if !matches && log.Tags != nil {
				for _, tag := range log.Tags {
					if strings.Contains(strings.ToLower(tag), filterLower) {
						matches = true
						break
					}
				}
			}
			
			if !matches {
				continue
			}
		}
		
		lv.filteredLogs = append(lv.filteredLogs, log)
	}
	
	// Sort by timestamp
	sort.Slice(lv.filteredLogs, func(i, j int) bool {
		if lv.filteredLogs[i].Timestamp == nil || lv.filteredLogs[j].Timestamp == nil {
			return i < j
		}
		return lv.filteredLogs[i].Timestamp.Before(*lv.filteredLogs[j].Timestamp)
	})
	
	lv.refreshView()
}

// refreshView updates the viewport content
func (lv *LogViewer) refreshView() {
	content := lv.renderLogs()
	lv.viewport.SetContent(content)
}

// updateLayout updates the layout dimensions
func (lv *LogViewer) updateLayout() {
	if lv.showSidebar {
		lv.viewport.Width = lv.width - lv.sidebarWidth - 2 // Account for borders
	} else {
		lv.viewport.Width = lv.width - 2
	}
	lv.viewport.Height = lv.height - 4 // Account for borders and header
	
	lv.refreshView()
}

// renderWithSidebar renders the log viewer with sidebar
func (lv *LogViewer) renderWithSidebar() string {
	sidebar := lv.renderSidebar()
	main := lv.renderMain()
	
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		main,
	)
}

// renderFullWidth renders the log viewer at full width
func (lv *LogViewer) renderFullWidth() string {
	return lv.renderMain()
}

// renderSidebar renders the sidebar with filters and stats
func (lv *LogViewer) renderSidebar() string {
	baseStyle := lipgloss.NewStyle().
		Width(lv.sidebarWidth).
		Height(lv.height).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#414868"))

	var content []string
	
	// Title
	content = append(content, lipgloss.NewStyle().Bold(true).Render("Log Filters"))
	content = append(content, "")
	
	// Level filter status
	levelStatus := "All"
	if lv.levelFilter != "" {
		levelStatus = lv.levelFilter
	}
	content = append(content, fmt.Sprintf("Level: %s", levelStatus))
	
	// Text filter status
	textStatus := "None"
	if lv.filter != "" {
		textStatus = lv.filter
		if len(textStatus) > 15 {
			textStatus = textStatus[:12] + "..."
		}
	}
	content = append(content, fmt.Sprintf("Text: %s", textStatus))
	content = append(content, "")
	
	// Statistics
	content = append(content, lipgloss.NewStyle().Bold(true).Render("Statistics"))
	content = append(content, "")
	content = append(content, fmt.Sprintf("Total: %d", len(lv.logs)))
	content = append(content, fmt.Sprintf("Filtered: %d", len(lv.filteredLogs)))
	
	// Level counts
	levelCounts := lv.getLevelCounts()
	for level, count := range levelCounts {
		if count > 0 {
			content = append(content, fmt.Sprintf("%s: %d", level, count))
		}
	}
	
	content = append(content, "")
	
	// Controls
	content = append(content, lipgloss.NewStyle().Bold(true).Render("Controls"))
	content = append(content, "")
	content = append(content, "d - Debug logs")
	content = append(content, "i - Info logs")
	content = append(content, "w - Warn logs")
	content = append(content, "e - Error logs")
	content = append(content, "x - Clear filter")
	content = append(content, "c - Toggle context")
	content = append(content, "a - Toggle auto-scroll")
	content = append(content, "s - Toggle sidebar")
	
	return baseStyle.Render(strings.Join(content, "\n"))
}

// renderMain renders the main log content area
func (lv *LogViewer) renderMain() string {
	var width int
	if lv.showSidebar {
		width = lv.width - lv.sidebarWidth - 2
	} else {
		width = lv.width
	}
	
	// Create header
	header := lv.renderHeader()
	
	// Create bordered viewport
	mainStyle := lipgloss.NewStyle().
		Width(width).
		Height(lv.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#414868"))
	
	// Combine header and viewport
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lv.viewport.View(),
	)
	
	return mainStyle.Render(content)
}

// renderHeader renders the log viewer header
func (lv *LogViewer) renderHeader() string {
	var parts []string
	
	// Title
	title := "Logs"
	if lv.filter != "" || lv.levelFilter != "" {
		title += " (Filtered)"
	}
	
	parts = append(parts, lipgloss.NewStyle().Bold(true).Render(title))
	
	// Status indicators
	var indicators []string
	if lv.autoScroll {
		indicators = append(indicators, "Auto-scroll")
	}
	if lv.showContext {
		indicators = append(indicators, "Context")
	}
	
	if len(indicators) > 0 {
		parts = append(parts, fmt.Sprintf("[%s]", strings.Join(indicators, ", ")))
	}
	
	// Count
	count := fmt.Sprintf("(%d/%d)", len(lv.filteredLogs), len(lv.logs))
	parts = append(parts, count)
	
	headerText := strings.Join(parts, " ")
	
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#24283B")).
		Foreground(lipgloss.Color("#C0CAF5")).
		Padding(0, 1).
		Width(lv.viewport.Width).
		Bold(true)
	
	return headerStyle.Render(headerText)
}

// renderLogs renders the log entries
func (lv *LogViewer) renderLogs() string {
	if len(lv.filteredLogs) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565F89")).
			Italic(true).
			Align(lipgloss.Center).
			Width(lv.viewport.Width).
			Height(lv.viewport.Height).
			Render("No logs to display")
	}
	
	var lines []string
	
	for _, log := range lv.filteredLogs {
		lines = append(lines, lv.renderLogEntry(log))
	}
	
	return strings.Join(lines, "\n")
}

// renderLogEntry renders a single log entry
func (lv *LogViewer) renderLogEntry(log types.LogParams) string {
	var parts []string
	
	// Timestamp
	timestamp := "----"
	if log.Timestamp != nil {
		timestamp = log.Timestamp.Format("15:04:05")
	}
	
	timestampStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Width(8)
	parts = append(parts, timestampStyle.Render(timestamp))
	
	// Level
	levelStyle := lv.getLevelStyle(string(log.Level))
	level := string(log.Level)
	if level == "" {
		level = "INFO"
	}
	parts = append(parts, levelStyle.Width(5).Align(lipgloss.Center).Render(level))
	
	// Component
	if log.Component != "" {
		componentStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9AA5CE")).
			Width(12)
		component := log.Component
		if len(component) > 12 {
			component = component[:9] + "..."
		}
		parts = append(parts, componentStyle.Render(component))
	}
	
	// Message
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C0CAF5"))
	parts = append(parts, messageStyle.Render(log.Message))
	
	line := strings.Join(parts, " ")
	
	// Add context if enabled and available
	if lv.showContext && log.Context != nil && len(log.Context) > 0 {
		contextLines := []string{line}
		for key, value := range log.Context {
			contextStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#565F89")).
				Italic(true)
			contextLine := fmt.Sprintf("  %s: %v", key, value)
			contextLines = append(contextLines, contextStyle.Render(contextLine))
		}
		line = strings.Join(contextLines, "\n")
	}
	
	// Add tags if available
	if log.Tags != nil && len(log.Tags) > 0 {
		tagStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#74C0FC")).
			Background(lipgloss.Color("#1A1B26")).
			Padding(0, 1).
			MarginLeft(1)
		
		tagStrings := make([]string, len(log.Tags))
		for i, tag := range log.Tags {
			tagStrings[i] = tagStyle.Render(tag)
		}
		
		line += " " + strings.Join(tagStrings, " ")
	}
	
	return line
}

// getLevelStyle returns the appropriate style for a log level
func (lv *LogViewer) getLevelStyle(level string) lipgloss.Style {
	baseStyle := lipgloss.NewStyle().Bold(true)
	
	switch level {
	case "DEBUG":
		return baseStyle.Foreground(lipgloss.Color("#ADB5BD"))
	case "INFO":
		return baseStyle.Foreground(lipgloss.Color("#74C0FC"))
	case "WARN":
		return baseStyle.Foreground(lipgloss.Color("#FFD93D"))
	case "ERROR":
		return baseStyle.Foreground(lipgloss.Color("#FF6B6B"))
	case "SUCCESS":
		return baseStyle.Foreground(lipgloss.Color("#51CF66"))
	default:
		return baseStyle.Foreground(lipgloss.Color("#9AA5CE"))
	}
}

// getLevelCounts returns counts for each log level
func (lv *LogViewer) getLevelCounts() map[string]int {
	counts := make(map[string]int)
	
	for _, log := range lv.logs {
		level := string(log.Level)
		if level == "" {
			level = "INFO"
		}
		counts[level]++
	}
	
	return counts
}