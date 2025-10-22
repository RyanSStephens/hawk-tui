package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/types"
)

// StatusBar component for displaying system status and information
type StatusBar struct {
	// Layout
	width  int
	height int
	styles interface{}

	// State
	status        StatusInfo
	progressItems map[string]types.ProgressParams
}

// StatusInfo contains information to display in the status bar
type StatusInfo struct {
	ViewMode     string
	MessageCount int64
	ErrorCount   int64
	LastUpdate   time.Time
	FPS          float64
	InputMode    bool
	SearchQuery  string
}

// NewStatusBar creates a new status bar component
func NewStatusBar(styles interface{}) *StatusBar {
	return &StatusBar{
		styles:        styles,
		progressItems: make(map[string]types.ProgressParams),
	}
}

// Init implements tea.Model
func (sb *StatusBar) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (sb *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sb.SetSize(msg.Width, msg.Height)
	}

	return sb, nil
}

// View implements tea.Model
func (sb *StatusBar) View() string {
	if sb.width == 0 {
		return ""
	}

	// Render progress bars if any
	var progressSection string
	if len(sb.progressItems) > 0 {
		progressSection = sb.renderProgressSection()
	}

	// Render main status line
	statusLine := sb.renderStatusLine()

	if progressSection != "" {
		return lipgloss.JoinVertical(lipgloss.Left, progressSection, statusLine)
	}

	return statusLine
}

// SetSize sets the dimensions of the status bar
func (sb *StatusBar) SetSize(width, height int) {
	sb.width = width
	sb.height = height
}

// UpdateStatus updates the status information
func (sb *StatusBar) UpdateStatus(status StatusInfo) {
	sb.status = status
}

// UpdateProgress updates or adds a progress item
func (sb *StatusBar) UpdateProgress(progress types.ProgressParams) {
	if progress.Status == types.ProgressStatusCompleted || progress.Status == types.ProgressStatusError {
		// Remove completed or errored progress items after a short delay
		// In a real implementation, you might want to show them briefly before removing
		delete(sb.progressItems, progress.ID)
	} else {
		sb.progressItems[progress.ID] = progress
	}
}

// RemoveProgress removes a progress item
func (sb *StatusBar) RemoveProgress(id string) {
	delete(sb.progressItems, id)
}

// renderStatusLine renders the main status line
func (sb *StatusBar) renderStatusLine() string {
	// Left section: view mode and input status
	leftSection := sb.renderLeftSection()
	
	// Center section: message counts and performance
	centerSection := sb.renderCenterSection()
	
	// Right section: time and system info
	rightSection := sb.renderRightSection()

	// Calculate spacing
	leftWidth := lipgloss.Width(leftSection)
	centerWidth := lipgloss.Width(centerSection)
	rightWidth := lipgloss.Width(rightSection)
	
	totalContentWidth := leftWidth + centerWidth + rightWidth
	availableSpace := sb.width - totalContentWidth
	
	var spacing string
	if availableSpace > 0 {
		leftSpacing := availableSpace / 2
		rightSpacing := availableSpace - leftSpacing
		spacing = strings.Repeat(" ", leftSpacing)
		centerSection += strings.Repeat(" ", rightSpacing)
	}

	statusLine := leftSection + spacing + centerSection + rightSection

	// Apply status bar styling
	statusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#24283B")).
		Foreground(lipgloss.Color("#C0CAF5")).
		Width(sb.width).
		Padding(0, 1)

	return statusStyle.Render(statusLine)
}

// renderLeftSection renders the left section of the status bar
func (sb *StatusBar) renderLeftSection() string {
	var parts []string

	// Current view mode
	viewStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D7FF"))
	parts = append(parts, viewStyle.Render(sb.status.ViewMode))

	// Input mode indicator
	if sb.status.InputMode {
		inputStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("#FFD93D")).
			Foreground(lipgloss.Color("#1A1B26")).
			Padding(0, 1).
			Bold(true)
		
		searchText := "SEARCH"
		if sb.status.SearchQuery != "" {
			searchText = fmt.Sprintf("SEARCH: %s", sb.status.SearchQuery)
		}
		parts = append(parts, inputStyle.Render(searchText))
	}

	return strings.Join(parts, " ")
}

// renderCenterSection renders the center section of the status bar
func (sb *StatusBar) renderCenterSection() string {
	var parts []string

	// Message count
	msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#51CF66"))
	parts = append(parts, msgStyle.Render(fmt.Sprintf("Msgs: %d", sb.status.MessageCount)))

	// Error count (if any)
	if sb.status.ErrorCount > 0 {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
		parts = append(parts, errStyle.Render(fmt.Sprintf("Errors: %d", sb.status.ErrorCount)))
	}

	// FPS indicator
	fpsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#74C0FC"))
	parts = append(parts, fpsStyle.Render(fmt.Sprintf("FPS: %.1f", sb.status.FPS)))

	return strings.Join(parts, " │ ")
}

// renderRightSection renders the right section of the status bar
func (sb *StatusBar) renderRightSection() string {
	var parts []string

	// Last update time
	if !sb.status.LastUpdate.IsZero() {
		timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9AA5CE"))
		lastUpdate := sb.status.LastUpdate.Format("15:04:05")
		parts = append(parts, timeStyle.Render(fmt.Sprintf("Last: %s", lastUpdate)))
	}

	// Current time
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#C0CAF5"))
	currentTime := time.Now().Format("15:04:05")
	parts = append(parts, timeStyle.Render(currentTime))

	return strings.Join(parts, " │ ")
}

// renderProgressSection renders active progress bars
func (sb *StatusBar) renderProgressSection() string {
	if len(sb.progressItems) == 0 {
		return ""
	}

	var progressLines []string

	for _, progress := range sb.progressItems {
		progressLine := sb.renderProgressBar(progress)
		progressLines = append(progressLines, progressLine)
	}

	// Style the progress section
	progressStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#414868")).
		Foreground(lipgloss.Color("#C0CAF5")).
		Width(sb.width).
		Padding(0, 1)

	content := strings.Join(progressLines, "\n")
	return progressStyle.Render(content)
}

// renderProgressBar renders a single progress bar
func (sb *StatusBar) renderProgressBar(progress types.ProgressParams) string {
	// Calculate available width for the progress bar
	labelWidth := 30
	percentWidth := 8
	detailsWidth := 20
	progressBarWidth := sb.width - labelWidth - percentWidth - detailsWidth - 10 // Account for spacing and padding

	if progressBarWidth < 10 {
		progressBarWidth = 10
	}

	// Label
	label := progress.Label
	if len(label) > labelWidth-2 {
		label = label[:labelWidth-5] + "..."
	}
	labelStyle := lipgloss.NewStyle().
		Width(labelWidth).
		Foreground(lipgloss.Color("#C0CAF5"))

	// Progress percentage
	percentage := float64(progress.Current) / float64(progress.Total) * 100
	if progress.Total == 0 {
		percentage = 0
	}
	
	percentText := fmt.Sprintf("%.1f%%", percentage)
	if progress.Unit != "" && progress.Unit != "%" {
		percentText = fmt.Sprintf("%.0f/%0.f %s", progress.Current, progress.Total, progress.Unit)
	}
	
	percentStyle := lipgloss.NewStyle().
		Width(percentWidth).
		Align(lipgloss.Right).
		Foreground(sb.getProgressColor(progress.Status))

	// Progress bar
	filled := int(percentage / 100 * float64(progressBarWidth))
	if filled > progressBarWidth {
		filled = progressBarWidth
	}
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", progressBarWidth-filled)
	barStyle := lipgloss.NewStyle().
		Foreground(sb.getProgressColor(progress.Status))

	// Details/status
	details := string(progress.Status)
	if progress.Details != "" {
		details = progress.Details
	}
	if len(details) > detailsWidth-2 {
		details = details[:detailsWidth-5] + "..."
	}
	detailsStyle := lipgloss.NewStyle().
		Width(detailsWidth).
		Foreground(lipgloss.Color("#9AA5CE"))

	// Combine all parts
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		labelStyle.Render(label),
		" ",
		barStyle.Render(bar),
		" ",
		percentStyle.Render(percentText),
		" ",
		detailsStyle.Render(details),
	)
}

// getProgressColor returns appropriate color for progress status
func (sb *StatusBar) getProgressColor(status types.ProgressStatus) lipgloss.Color {
	switch status {
	case types.ProgressStatusCompleted:
		return lipgloss.Color("#51CF66") // Green
	case types.ProgressStatusError:
		return lipgloss.Color("#FF6B6B") // Red
	case types.ProgressStatusInProgress:
		return lipgloss.Color("#00D7FF") // Cyan
	case types.ProgressStatusPending:
		return lipgloss.Color("#FFD93D") // Yellow
	default:
		return lipgloss.Color("#9AA5CE") // Gray
	}
}

// GetHeight returns the current height needed for the status bar
func (sb *StatusBar) GetHeight() int {
	height := 1 // Base status line
	
	// Add height for progress bars
	if len(sb.progressItems) > 0 {
		height += len(sb.progressItems)
	}
	
	return height
}