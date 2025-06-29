package tui

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/internal/protocol"
	"github.com/hawk-tui/hawk-tui/internal/tui/components"
	"github.com/hawk-tui/hawk-tui/pkg/types"
)

// ViewMode represents the current view mode of the TUI
type ViewMode int

const (
	ViewModeLogs ViewMode = iota
	ViewModeMetrics
	ViewModeDashboard
	ViewModeConfig
	ViewModeHelp
)

// Config holds the configuration for the TUI
type Config struct {
	AppName    string
	LogLevel   string
	ConfigFile string
	Debug      bool
}

// Model represents the main TUI model
type Model struct {
	// Configuration
	config Config
	styles *Styles

	// Layout
	width  int
	height int
	ready  bool

	// Current view state
	viewMode    ViewMode
	activePanel int
	showHelp    bool

	// Components
	logViewer    *components.LogViewer
	metricsView  *components.MetricsView
	dashboardView *components.DashboardView
	configView   *components.ConfigView
	helpView     *components.HelpView
	statusBar    *components.StatusBar

	// Protocol handling
	protocolHandler *protocol.ProtocolHandler
	responseWriter  io.Writer
	
	// Data storage
	logs     []types.LogParams
	metrics  map[string]types.MetricParams
	configs  map[string]types.ConfigParams
	progress map[string]types.ProgressParams
	widgets  map[string]types.DashboardParams
	events   []types.EventParams
	
	// State management
	mu       sync.RWMutex
	lastUpdate time.Time
	messageCount int64
	errorCount   int64

	// Input handling
	inputMode   bool
	searchQuery string
	
	// Performance tracking
	frameCount   int64
	lastFPSTime  time.Time
	currentFPS   float64
}

// NewModel creates a new TUI model
func NewModel(config Config) (*Model, error) {
	m := &Model{
		config:    config,
		styles:    NewStyles(),
		viewMode:  ViewModeLogs,
		logs:      make([]types.LogParams, 0),
		metrics:   make(map[string]types.MetricParams),
		configs:   make(map[string]types.ConfigParams),
		progress:  make(map[string]types.ProgressParams),
		widgets:   make(map[string]types.DashboardParams),
		events:    make([]types.EventParams, 0),
		lastUpdate: time.Now(),
		lastFPSTime: time.Now(),
	}

	// Initialize components
	m.logViewer = components.NewLogViewer(m.styles)
	m.metricsView = components.NewMetricsView(m.styles)
	m.dashboardView = components.NewDashboardView(m.styles)
	m.configView = components.NewConfigView(m.styles)
	m.helpView = components.NewHelpView(m.styles)
	m.statusBar = components.NewStatusBar(m.styles)

	// Setup protocol handler
	m.responseWriter = os.Stdout
	m.protocolHandler = protocol.NewProtocolHandler(
		os.Stdin,
		m.responseWriter,
		m, // MessageHandler interface
		m, // ResponseSender interface
	)

	// Start protocol handler
	if err := m.protocolHandler.Start(); err != nil {
		return nil, fmt.Errorf("failed to start protocol handler: %w", err)
	}

	return m, nil
}

// Init implements tea.Model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.logViewer.Init(),
		m.metricsView.Init(),
		m.dashboardView.Init(),
		m.configView.Init(),
		m.statusBar.Init(),
		tea.EnterAltScreen,
		tickCmd(),
	)
}

// Update implements tea.Model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle window resize
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		
		// Update component sizes
		m.updateComponentSizes()
		
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tickMsg:
		m.updateFPS()
		m.updateStatusBar()
		cmds = append(cmds, tickCmd())

	case dataUpdateMsg:
		m.handleDataUpdate(msg)
	}

	// Update active component
	switch m.viewMode {
	case ViewModeLogs:
		var newModel tea.Model
		newModel, cmd = m.logViewer.Update(msg)
		m.logViewer = newModel.(*components.LogViewer)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case ViewModeMetrics:
		var newModel tea.Model
		newModel, cmd = m.metricsView.Update(msg)
		m.metricsView = newModel.(*components.MetricsView)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case ViewModeDashboard:
		var newModel tea.Model
		newModel, cmd = m.dashboardView.Update(msg)
		m.dashboardView = newModel.(*components.DashboardView)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case ViewModeConfig:
		var newModel tea.Model
		newModel, cmd = m.configView.Update(msg)
		m.configView = newModel.(*components.ConfigView)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case ViewModeHelp:
		var newModel tea.Model
		newModel, cmd = m.helpView.Update(msg)
		m.helpView = newModel.(*components.HelpView)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update status bar
	var statusModel tea.Model
	statusModel, cmd = m.statusBar.Update(msg)
	m.statusBar = statusModel.(*components.StatusBar)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View implements tea.Model
func (m *Model) View() string {
	if !m.ready {
		return "Initializing Hawk TUI..."
	}

	// Calculate layout dimensions
	headerHeight := 3
	footerHeight := 2
	contentHeight := m.height - headerHeight - footerHeight

	// Render header
	header := m.renderHeader()

	// Render main content based on current view
	var content string
	switch m.viewMode {
	case ViewModeLogs:
		content = m.logViewer.View()
	case ViewModeMetrics:
		content = m.metricsView.View()
	case ViewModeDashboard:
		content = m.dashboardView.View()
	case ViewModeConfig:
		content = m.configView.View()
	case ViewModeHelp:
		content = m.helpView.View()
	}

	// Ensure content fits within available space
	content = lipgloss.NewStyle().
		Width(m.width).
		Height(contentHeight).
		Render(content)

	// Render footer
	footer := m.statusBar.View()

	// Combine all parts
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

// renderHeader renders the top navigation bar
func (m *Model) renderHeader() string {
	// Create tab styles
	activeTabStyle := m.styles.TabActive
	inactiveTabStyle := m.styles.TabInactive

	// Define tabs
	tabs := []struct {
		name     string
		mode     ViewMode
		shortcut string
	}{
		{"Logs", ViewModeLogs, "1"},
		{"Metrics", ViewModeMetrics, "2"},
		{"Dashboard", ViewModeDashboard, "3"},
		{"Config", ViewModeConfig, "4"},
		{"Help", ViewModeHelp, "h"},
	}

	// Render tabs
	var tabViews []string
	for _, tab := range tabs {
		style := inactiveTabStyle
		if tab.mode == m.viewMode {
			style = activeTabStyle
		}

		tabText := fmt.Sprintf("%s (%s)", tab.name, tab.shortcut)
		tabViews = append(tabViews, style.Render(tabText))
	}

	// Join tabs horizontally
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)

	// Add app title on the right
	title := fmt.Sprintf("Hawk TUI - %s", m.config.AppName)
	titleStyle := m.styles.Header.Copy().
		Align(lipgloss.Right).
		Width(m.width - lipgloss.Width(tabBar))

	// Combine tab bar and title
	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabBar,
		titleStyle.Render(title),
	)

	// Add border
	return m.styles.Header.Copy().
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorBorderPrimary).
		Render(headerLine)
}

// handleKeyPress handles keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Global shortcuts
	switch msg.String() {
	case "ctrl+c", "q":
		if !m.inputMode {
			m.protocolHandler.Stop()
			return m, tea.Quit
		}

	case "1":
		if !m.inputMode {
			m.viewMode = ViewModeLogs
			return m, nil
		}

	case "2":
		if !m.inputMode {
			m.viewMode = ViewModeMetrics
			return m, nil
		}

	case "3":
		if !m.inputMode {
			m.viewMode = ViewModeDashboard
			return m, nil
		}

	case "4":
		if !m.inputMode {
			m.viewMode = ViewModeConfig
			return m, nil
		}

	case "h":
		if !m.inputMode {
			m.viewMode = ViewModeHelp
			return m, nil
		}

	case "r":
		if !m.inputMode {
			// Refresh current view
			cmd := m.refreshCurrentView()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

	case "/":
		if !m.inputMode {
			m.inputMode = true
			m.searchQuery = ""
			return m, nil
		}

	case "esc":
		if m.inputMode {
			m.inputMode = false
			m.searchQuery = ""
			return m, nil
		}
		if m.viewMode == ViewModeHelp {
			m.viewMode = ViewModeLogs
			return m, nil
		}

	case "tab":
		if !m.inputMode {
			// Cycle through panels within current view
			m.activePanel = (m.activePanel + 1) % 3
			return m, nil
		}

	case "shift+tab":
		if !m.inputMode {
			// Cycle backwards through panels
			m.activePanel = (m.activePanel - 1 + 3) % 3
			return m, nil
		}
	}

	// Handle search input
	if m.inputMode {
		switch msg.Type {
		case tea.KeyRunes:
			m.searchQuery += string(msg.Runes)
			m.applySearch()
			return m, nil
		case tea.KeyBackspace:
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.applySearch()
			}
			return m, nil
		case tea.KeyEnter:
			m.inputMode = false
			return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}

// updateComponentSizes updates the sizes of all components based on terminal dimensions
func (m *Model) updateComponentSizes() {
	// Calculate available content area
	contentWidth := m.width
	contentHeight := m.height - 5 // Header (3) + Footer (2)

	// Update all components
	m.logViewer.SetSize(contentWidth, contentHeight)
	m.metricsView.SetSize(contentWidth, contentHeight)
	m.dashboardView.SetSize(contentWidth, contentHeight)
	m.configView.SetSize(contentWidth, contentHeight)
	m.helpView.SetSize(contentWidth, contentHeight)
	m.statusBar.SetSize(contentWidth, 2)
}

// refreshCurrentView refreshes the current view's data
func (m *Model) refreshCurrentView() tea.Cmd {
	switch m.viewMode {
	case ViewModeLogs:
		return m.logViewer.Refresh()
	case ViewModeMetrics:
		return m.metricsView.Refresh()
	case ViewModeDashboard:
		return m.dashboardView.Refresh()
	case ViewModeConfig:
		return m.configView.Refresh()
	}
	return nil
}

// applySearch applies the current search query to the active view
func (m *Model) applySearch() {
	switch m.viewMode {
	case ViewModeLogs:
		m.logViewer.SetFilter(m.searchQuery)
	case ViewModeMetrics:
		m.metricsView.SetFilter(m.searchQuery)
	case ViewModeConfig:
		m.configView.SetFilter(m.searchQuery)
	}
}

// updateFPS calculates and updates the current FPS
func (m *Model) updateFPS() {
	m.frameCount++
	now := time.Now()
	
	if now.Sub(m.lastFPSTime) >= time.Second {
		m.currentFPS = float64(m.frameCount) / now.Sub(m.lastFPSTime).Seconds()
		m.frameCount = 0
		m.lastFPSTime = now
	}
}

// updateStatusBar updates the status bar with current information
func (m *Model) updateStatusBar() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := components.StatusInfo{
		ViewMode:     m.getViewModeString(),
		MessageCount: m.messageCount,
		ErrorCount:   m.errorCount,
		LastUpdate:   m.lastUpdate,
		FPS:          m.currentFPS,
		InputMode:    m.inputMode,
		SearchQuery:  m.searchQuery,
	}

	m.statusBar.UpdateStatus(status)
}

// getViewModeString returns a string representation of the current view mode
func (m *Model) getViewModeString() string {
	switch m.viewMode {
	case ViewModeLogs:
		return "Logs"
	case ViewModeMetrics:
		return "Metrics"
	case ViewModeDashboard:
		return "Dashboard"
	case ViewModeConfig:
		return "Config"
	case ViewModeHelp:
		return "Help"
	default:
		return "Unknown"
	}
}

// handleDataUpdate processes data update messages
func (m *Model) handleDataUpdate(msg dataUpdateMsg) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastUpdate = time.Now()
	m.messageCount++

	// Forward update to appropriate component
	switch msg.Type {
	case "log":
		if logData, ok := msg.Data.(types.LogParams); ok {
			m.logViewer.AddLog(logData)
		}
	case "metric":
		if metricData, ok := msg.Data.(types.MetricParams); ok {
			m.metricsView.UpdateMetric(metricData)
		}
	case "config":
		if configData, ok := msg.Data.(types.ConfigParams); ok {
			m.configView.UpdateConfig(configData)
		}
	case "dashboard":
		if dashboardData, ok := msg.Data.(types.DashboardParams); ok {
			m.dashboardView.UpdateWidget(dashboardData)
		}
	case "progress":
		if progressData, ok := msg.Data.(types.ProgressParams); ok {
			// Update progress in relevant components
			m.statusBar.UpdateProgress(progressData)
		}
	case "event":
		if eventData, ok := msg.Data.(types.EventParams); ok {
			// Handle events (could show notifications, etc.)
			m.logViewer.AddEvent(eventData)
		}
	}
}

// Message types for internal communication
type tickMsg time.Time
type dataUpdateMsg struct {
	Type string
	Data interface{}
}

// tickCmd generates tick messages for regular updates
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg { // ~60 FPS
		return tickMsg(t)
	})
}

// MessageHandler interface implementation
func (m *Model) HandleLog(params types.LogParams, msgID interface{}) error {
	// Send data update message
	// Note: In a real implementation, you'd want to use a proper message queue
	// For now, we'll store and update directly
	m.mu.Lock()
	m.logs = append(m.logs, params)
	m.mu.Unlock()

	// Trigger UI update (this would be better done with proper channels)
	go func() {
		// In practice, you'd send this through a channel that Update() listens to
	}()

	return nil
}

func (m *Model) HandleMetric(params types.MetricParams, msgID interface{}) error {
	m.mu.Lock()
	m.metrics[params.Name] = params
	m.mu.Unlock()
	return nil
}

func (m *Model) HandleConfig(params types.ConfigParams, msgID interface{}) error {
	m.mu.Lock()
	m.configs[params.Key] = params
	m.mu.Unlock()
	return nil
}

func (m *Model) HandleProgress(params types.ProgressParams, msgID interface{}) error {
	m.mu.Lock()
	m.progress[params.ID] = params
	m.mu.Unlock()
	return nil
}

func (m *Model) HandleDashboard(params types.DashboardParams, msgID interface{}) error {
	m.mu.Lock()
	m.widgets[params.WidgetID] = params
	m.mu.Unlock()
	return nil
}

func (m *Model) HandleEvent(params types.EventParams, msgID interface{}) error {
	m.mu.Lock()
	m.events = append(m.events, params)
	m.mu.Unlock()
	return nil
}

// ResponseSender interface implementation
func (m *Model) SendResponse(id interface{}, result interface{}) error {
	response := types.JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = m.responseWriter.Write(append(data, '\n'))
	return err
}

func (m *Model) SendError(id interface{}, rpcErr *types.RPCError) error {
	response := types.JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error:   rpcErr,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = m.responseWriter.Write(append(data, '\n'))
	return err
}

func (m *Model) SendNotification(method string, params interface{}) error {
	notification := types.JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	_, err = m.responseWriter.Write(append(data, '\n'))
	return err
}