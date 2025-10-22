package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/types"
)

// ConfigView component for displaying and managing configuration parameters
type ConfigView struct {
	// Layout
	width  int
	height int
	styles interface{}

	// Data
	configs         map[string]types.ConfigParams
	filteredConfigs []string
	categories      map[string][]string

	// State
	filter           string
	selectedConfig   int
	editMode         bool
	editInput        textinput.Model
	viewMode         ConfigViewMode
	selectedCategory string

	// Display options
	showCategories   bool
	showDescriptions bool
}

// ConfigViewMode represents different ways to view configuration
type ConfigViewMode int

const (
	ConfigViewModeList ConfigViewMode = iota
	ConfigViewModeCategories
	ConfigViewModeEdit
)

// NewConfigView creates a new configuration view component
func NewConfigView(styles interface{}) *ConfigView {
	input := textinput.New()
	input.Placeholder = "Enter new value..."

	return &ConfigView{
		styles:           styles,
		configs:          make(map[string]types.ConfigParams),
		filteredConfigs:  make([]string, 0),
		categories:       make(map[string][]string),
		editInput:        input,
		viewMode:         ConfigViewModeList,
		showCategories:   true,
		showDescriptions: true,
	}
}

// Init implements tea.Model
func (cv *ConfigView) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (cv *ConfigView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle edit mode input
	if cv.editMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				cv.applyEdit()
				cv.editMode = false
				return cv, nil
			case "esc":
				cv.editMode = false
				return cv, nil
			}
		}

		cv.editInput, cmd = cv.editInput.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if len(cmds) > 0 {
			return cv, tea.Batch(cmds...)
		}
		return cv, nil
	}

	// Normal mode handling
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if cv.selectedConfig > 0 {
				cv.selectedConfig--
			}
		case "down", "j":
			if cv.selectedConfig < len(cv.filteredConfigs)-1 {
				cv.selectedConfig++
			}
		case "left", "h":
			if cv.viewMode == ConfigViewModeCategories && cv.selectedCategory != "" {
				cv.selectedCategory = ""
				cv.applyFilter()
			}
		case "right", "l", "enter":
			if cv.viewMode == ConfigViewModeCategories && cv.selectedCategory == "" {
				// Select category
				categories := cv.getCategories()
				if cv.selectedConfig < len(categories) {
					cv.selectedCategory = categories[cv.selectedConfig]
					cv.applyFilter()
				}
			} else {
				// Edit selected config
				cv.startEdit()
			}
		case "L": // Shift+L for list view
			cv.viewMode = ConfigViewModeList
			cv.selectedCategory = ""
			cv.applyFilter()
		case "c":
			cv.viewMode = ConfigViewModeCategories
			cv.selectedCategory = ""
			cv.applyFilter()
		case "d":
			cv.showDescriptions = !cv.showDescriptions
		case "r":
			// Reset to default value
			cv.resetToDefault()
		case "esc":
			if cv.viewMode == ConfigViewModeCategories && cv.selectedCategory != "" {
				cv.selectedCategory = ""
				cv.applyFilter()
			}
		}

	case tea.WindowSizeMsg:
		cv.SetSize(msg.Width, msg.Height)
	}

	if len(cmds) > 0 {
		return cv, tea.Batch(cmds...)
	}

	return cv, nil
}

// View implements tea.Model
func (cv *ConfigView) View() string {
	if cv.width == 0 || cv.height == 0 {
		return "Loading configuration..."
	}

	if cv.editMode {
		return cv.renderEditMode()
	}

	switch cv.viewMode {
	case ConfigViewModeCategories:
		return cv.renderCategoriesView()
	default: // ConfigViewModeList
		return cv.renderListView()
	}
}

// SetSize sets the dimensions of the config view
func (cv *ConfigView) SetSize(width, height int) {
	cv.width = width
	cv.height = height
	cv.editInput.Width = width - 20
}

// UpdateConfig updates or adds a configuration parameter
func (cv *ConfigView) UpdateConfig(config types.ConfigParams) {
	cv.configs[config.Key] = config
	cv.updateCategories()
	cv.applyFilter()
}

// SetFilter sets the filter for configuration keys
func (cv *ConfigView) SetFilter(filter string) {
	cv.filter = filter
	cv.applyFilter()
}

// Refresh refreshes the config view
func (cv *ConfigView) Refresh() tea.Cmd {
	return nil
}

// updateCategories organizes configs into categories
func (cv *ConfigView) updateCategories() {
	cv.categories = make(map[string][]string)

	for key, config := range cv.configs {
		category := config.Category
		if category == "" {
			// Extract category from key (everything before first dot)
			parts := strings.Split(key, ".")
			if len(parts) > 1 {
				category = parts[0]
			} else {
				category = "General"
			}
		}

		if _, exists := cv.categories[category]; !exists {
			cv.categories[category] = make([]string, 0)
		}
		cv.categories[category] = append(cv.categories[category], key)
	}

	// Sort configs within each category
	for category := range cv.categories {
		sort.Strings(cv.categories[category])
	}
}

// applyFilter applies the current filter and view mode
func (cv *ConfigView) applyFilter() {
	cv.filteredConfigs = make([]string, 0)

	if cv.viewMode == ConfigViewModeCategories && cv.selectedCategory == "" {
		// Show categories
		for category := range cv.categories {
			if cv.filter == "" || strings.Contains(strings.ToLower(category), strings.ToLower(cv.filter)) {
				cv.filteredConfigs = append(cv.filteredConfigs, category)
			}
		}
		sort.Strings(cv.filteredConfigs)
	} else {
		// Show configs
		var configsToCheck []string

		if cv.viewMode == ConfigViewModeCategories && cv.selectedCategory != "" {
			// Show configs from selected category
			configsToCheck = cv.categories[cv.selectedCategory]
		} else {
			// Show all configs
			for key := range cv.configs {
				configsToCheck = append(configsToCheck, key)
			}
		}

		for _, key := range configsToCheck {
			config := cv.configs[key]
			if cv.filter == "" ||
				strings.Contains(strings.ToLower(key), strings.ToLower(cv.filter)) ||
				strings.Contains(strings.ToLower(config.Description), strings.ToLower(cv.filter)) {
				cv.filteredConfigs = append(cv.filteredConfigs, key)
			}
		}
		sort.Strings(cv.filteredConfigs)
	}

	// Ensure selected config is in bounds
	if cv.selectedConfig >= len(cv.filteredConfigs) {
		cv.selectedConfig = len(cv.filteredConfigs) - 1
	}
	if cv.selectedConfig < 0 {
		cv.selectedConfig = 0
	}
}

// getCategories returns sorted list of categories
func (cv *ConfigView) getCategories() []string {
	categories := make([]string, 0, len(cv.categories))
	for category := range cv.categories {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// startEdit starts editing the selected configuration
func (cv *ConfigView) startEdit() {
	if len(cv.filteredConfigs) == 0 {
		return
	}

	key := cv.filteredConfigs[cv.selectedConfig]
	config := cv.configs[key]

	cv.editMode = true
	cv.editInput.SetValue(fmt.Sprintf("%v", config.Value))
	cv.editInput.Focus()
}

// applyEdit applies the edited value
func (cv *ConfigView) applyEdit() {
	if len(cv.filteredConfigs) == 0 {
		return
	}

	key := cv.filteredConfigs[cv.selectedConfig]
	config := cv.configs[key]
	newValue := cv.editInput.Value()

	// In a real implementation, you'd:
	// 1. Validate the new value based on config.Type
	// 2. Send the update via the protocol handler
	// 3. Handle validation errors

	// For now, just update locally
	config.Value = newValue
	cv.configs[key] = config

	cv.editInput.Blur()
}

// resetToDefault resets the selected config to its default value
func (cv *ConfigView) resetToDefault() {
	if len(cv.filteredConfigs) == 0 {
		return
	}

	key := cv.filteredConfigs[cv.selectedConfig]
	config := cv.configs[key]

	if config.Default != nil {
		config.Value = config.Default
		cv.configs[key] = config
		// In a real implementation, send this change via protocol
	}
}

// renderListView renders configurations in list view
func (cv *ConfigView) renderListView() string {
	if len(cv.filteredConfigs) == 0 {
		return cv.renderEmptyState()
	}

	header := cv.renderHeader("Configuration (List)")
	content := cv.renderConfigList()
	footer := cv.renderListControls()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderCategoriesView renders configurations organized by categories
func (cv *ConfigView) renderCategoriesView() string {
	header := cv.renderHeader("Configuration (Categories)")

	var content string
	if cv.selectedCategory == "" {
		content = cv.renderCategoryList()
	} else {
		content = cv.renderConfigList()
	}

	footer := cv.renderCategoryControls()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderEditMode renders the edit mode interface
func (cv *ConfigView) renderEditMode() string {
	if len(cv.filteredConfigs) == 0 {
		return "No configuration selected"
	}

	key := cv.filteredConfigs[cv.selectedConfig]
	config := cv.configs[key]

	var content []string

	// Title
	content = append(content, lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5")).
		Render(fmt.Sprintf("Editing: %s", key)))

	content = append(content, "")

	// Description
	if config.Description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9AA5CE")).
			Italic(true)
		content = append(content, descStyle.Render(config.Description))
		content = append(content, "")
	}

	// Current value
	content = append(content, fmt.Sprintf("Current value: %v", config.Value))

	// Default value
	if config.Default != nil {
		content = append(content, fmt.Sprintf("Default value: %v", config.Default))
	}

	// Type and constraints
	constraints := cv.getConstraintsText(config)
	if constraints != "" {
		content = append(content, constraints)
	}

	content = append(content, "")

	// Input field
	content = append(content, "New value:")
	content = append(content, cv.editInput.View())

	content = append(content, "")
	content = append(content, "Press Enter to save, Esc to cancel")

	return strings.Join(content, "\n")
}

// renderConfigList renders the list of configurations
func (cv *ConfigView) renderConfigList() string {
	var lines []string

	for i, key := range cv.filteredConfigs {
		config := cv.configs[key]
		line := cv.renderConfigLine(config, i == cv.selectedConfig)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderCategoryList renders the list of categories
func (cv *ConfigView) renderCategoryList() string {
	var lines []string

	categories := cv.getCategories()
	for i, category := range categories {
		selected := i == cv.selectedConfig
		count := len(cv.categories[category])

		line := cv.renderCategoryLine(category, count, selected)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderConfigLine renders a single configuration line
func (cv *ConfigView) renderConfigLine(config types.ConfigParams, selected bool) string {
	lineStyle := lipgloss.NewStyle().Padding(0, 1)

	if selected {
		lineStyle = lineStyle.Background(lipgloss.Color("#3D59A1"))
	}

	// Key
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5")).
		Width(30)

	key := config.Key
	if len(key) > 27 {
		key = key[:24] + "..."
	}

	// Value
	valueStyle := lipgloss.NewStyle().
		Foreground(cv.getValueColor(config)).
		Width(20)

	value := fmt.Sprintf("%v", config.Value)
	if len(value) > 17 {
		value = value[:14] + "..."
	}

	// Type
	typeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#74C0FC")).
		Width(10)

	// Status indicators
	var indicators []string
	if config.RestartRequired {
		indicators = append(indicators, "R") // Restart required
	}
	if config.Value != config.Default {
		indicators = append(indicators, "M") // Modified
	}

	indicatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD93D")).
		Width(5)

	indicatorText := strings.Join(indicators, "")

	line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		keyStyle.Render(key),
		valueStyle.Render(value),
		typeStyle.Render(string(config.Type)),
		indicatorStyle.Render(indicatorText),
	)

	// Add description if enabled
	if cv.showDescriptions && config.Description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565F89")).
			Italic(true).
			MarginLeft(2)

		description := config.Description
		if len(description) > cv.width-10 {
			description = description[:cv.width-13] + "..."
		}

		line = lipgloss.JoinVertical(
			lipgloss.Left,
			line,
			descStyle.Render(description),
		)
	}

	return lineStyle.Render(line)
}

// renderCategoryLine renders a single category line
func (cv *ConfigView) renderCategoryLine(category string, count int, selected bool) string {
	lineStyle := lipgloss.NewStyle().Padding(0, 1)

	if selected {
		lineStyle = lineStyle.Background(lipgloss.Color("#3D59A1"))
	}

	categoryStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5")).
		Width(40)

	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Width(10)

	line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		categoryStyle.Render(category),
		countStyle.Render(fmt.Sprintf("(%d)", count)),
	)

	return lineStyle.Render(line)
}

// renderHeader renders the configuration header
func (cv *ConfigView) renderHeader(title string) string {
	var parts []string

	parts = append(parts, title)

	if cv.filter != "" {
		parts = append(parts, fmt.Sprintf("(Filter: %s)", cv.filter))
	}

	parts = append(parts, fmt.Sprintf("(%d items)", len(cv.filteredConfigs)))

	headerText := strings.Join(parts, " ")

	return lipgloss.NewStyle().
		Background(lipgloss.Color("#24283B")).
		Foreground(lipgloss.Color("#C0CAF5")).
		Padding(0, 1).
		Width(cv.width).
		Bold(true).
		Render(headerText)
}

// renderListControls renders controls for list view
func (cv *ConfigView) renderListControls() string {
	controls := []string{
		"↑↓ Navigate",
		"Enter Edit",
		"r Reset",
		"c Categories",
		"d Toggle Desc",
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Render(strings.Join(controls, " │ "))
}

// renderCategoryControls renders controls for category view
func (cv *ConfigView) renderCategoryControls() string {
	var controls []string

	if cv.selectedCategory == "" {
		controls = []string{
			"↑↓ Navigate",
			"→ Enter",
			"L List",
		}
	} else {
		controls = []string{
			"↑↓ Navigate",
			"Enter Edit",
			"← Back",
			"r Reset",
		}
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Render(strings.Join(controls, " │ "))
}

// renderEmptyState renders the empty state
func (cv *ConfigView) renderEmptyState() string {
	message := "No configuration parameters to display"
	if cv.filter != "" {
		message = fmt.Sprintf("No configuration matches filter: %s", cv.filter)
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Italic(true).
		Align(lipgloss.Center).
		Width(cv.width).
		Height(cv.height).
		Render(message)
}

// Helper functions

// getValueColor returns appropriate color for config value
func (cv *ConfigView) getValueColor(config types.ConfigParams) lipgloss.Color {
	if config.Value != config.Default {
		return lipgloss.Color("#FFD93D") // Yellow for modified values
	}

	switch config.Type {
	case types.ConfigTypeBoolean:
		if value, ok := config.Value.(bool); ok && value {
			return lipgloss.Color("#51CF66") // Green for true
		}
		return lipgloss.Color("#FF6B6B") // Red for false
	case types.ConfigTypeString:
		return lipgloss.Color("#74C0FC") // Blue for strings
	case types.ConfigTypeInteger, types.ConfigTypeFloat:
		return lipgloss.Color("#00D7FF") // Cyan for numbers
	default:
		return lipgloss.Color("#C0CAF5") // Default
	}
}

// getConstraintsText returns text describing value constraints
func (cv *ConfigView) getConstraintsText(config types.ConfigParams) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Type: %s", config.Type))

	if config.Min != nil {
		parts = append(parts, fmt.Sprintf("Min: %.2f", *config.Min))
	}

	if config.Max != nil {
		parts = append(parts, fmt.Sprintf("Max: %.2f", *config.Max))
	}

	if len(config.Options) > 0 {
		optionStrs := make([]string, len(config.Options))
		for i, opt := range config.Options {
			optionStrs[i] = fmt.Sprintf("%v", opt)
		}
		parts = append(parts, fmt.Sprintf("Options: %s", strings.Join(optionStrs, ", ")))
	}

	if config.RestartRequired {
		parts = append(parts, "Restart required")
	}

	if len(parts) > 1 {
		return strings.Join(parts, " | ")
	}

	return ""
}
