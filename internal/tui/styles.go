package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Color palette for the Hawk TUI
var (
	// Base colors
	ColorPrimary   = lipgloss.Color("#00D7FF") // Bright cyan
	ColorSecondary = lipgloss.Color("#FF6B6B") // Coral red
	ColorSuccess   = lipgloss.Color("#51CF66") // Green
	ColorWarning   = lipgloss.Color("#FFD93D") // Yellow
	ColorError     = lipgloss.Color("#FF6B6B") // Red
	ColorInfo      = lipgloss.Color("#74C0FC") // Light blue
	ColorDebug     = lipgloss.Color("#ADB5BD") // Gray

	// Background colors
	ColorBgPrimary   = lipgloss.Color("#1A1B26") // Dark blue-gray
	ColorBgSecondary = lipgloss.Color("#24283B") // Darker blue-gray
	ColorBgTertiary  = lipgloss.Color("#414868") // Medium blue-gray
	ColorBgSelected  = lipgloss.Color("#3D59A1") // Selected blue
	ColorBgHover     = lipgloss.Color("#292E42") // Hover state

	// Text colors
	ColorTextPrimary   = lipgloss.Color("#C0CAF5") // Light blue-white
	ColorTextSecondary = lipgloss.Color("#9AA5CE") // Muted blue-white
	ColorTextMuted     = lipgloss.Color("#565F89") // Dark blue-gray
	ColorTextInverse   = lipgloss.Color("#1A1B26") // Dark for light backgrounds

	// Border colors
	ColorBorderPrimary   = lipgloss.Color("#414868") // Medium blue-gray
	ColorBorderSecondary = lipgloss.Color("#565F89") // Darker blue-gray
	ColorBorderActive    = lipgloss.Color("#00D7FF") // Bright cyan
	ColorBorderInactive  = lipgloss.Color("#292E42") // Very dark
)

// Styles defines all the visual styles used in the TUI
type Styles struct {
	// Base styles
	Base        lipgloss.Style
	Focused     lipgloss.Style
	Blurred     lipgloss.Style
	Selected    lipgloss.Style
	Error       lipgloss.Style
	Warning     lipgloss.Style
	Success     lipgloss.Style
	Info        lipgloss.Style
	Debug       lipgloss.Style

	// Layout styles
	Container   lipgloss.Style
	Panel       lipgloss.Style
	PanelActive lipgloss.Style
	Header      lipgloss.Style
	Footer      lipgloss.Style
	Sidebar     lipgloss.Style
	Content     lipgloss.Style

	// Component styles
	Tab         lipgloss.Style
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
	Button      lipgloss.Style
	ButtonActive lipgloss.Style
	Input       lipgloss.Style
	InputActive lipgloss.Style

	// Log styles
	LogEntry     lipgloss.Style
	LogTimestamp lipgloss.Style
	LogLevel     lipgloss.Style
	LogMessage   lipgloss.Style
	LogContext   lipgloss.Style

	// Metric styles
	MetricCard   lipgloss.Style
	MetricValue  lipgloss.Style
	MetricLabel  lipgloss.Style
	MetricUnit   lipgloss.Style
	GaugeBar     lipgloss.Style
	GaugeFill    lipgloss.Style

	// Dashboard styles
	Widget       lipgloss.Style
	WidgetTitle  lipgloss.Style
	WidgetBorder lipgloss.Style

	// Table styles
	TableHeader lipgloss.Style
	TableRow    lipgloss.Style
	TableRowAlt lipgloss.Style
	TableCell   lipgloss.Style

	// Progress styles
	ProgressBar  lipgloss.Style
	ProgressFill lipgloss.Style
	ProgressText lipgloss.Style

	// Status styles
	StatusHealthy  lipgloss.Style
	StatusDegraded lipgloss.Style
	StatusDown     lipgloss.Style
	StatusUnknown  lipgloss.Style

	// Help styles
	HelpKey   lipgloss.Style
	HelpValue lipgloss.Style
}

// NewStyles creates a new Styles instance with default styling
func NewStyles() *Styles {
	s := &Styles{}

	// Detect color profile
	profile := termenv.EnvColorProfile()
	lipgloss.SetColorProfile(profile)

	// Base styles
	s.Base = lipgloss.NewStyle().
		Foreground(ColorTextPrimary).
		Background(ColorBgPrimary)

	s.Focused = s.Base.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorderActive)

	s.Blurred = s.Base.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorderInactive)

	s.Selected = s.Base.Copy().
		Background(ColorBgSelected).
		Foreground(ColorTextPrimary).
		Bold(true)

	s.Error = s.Base.Copy().
		Foreground(ColorError).
		Bold(true)

	s.Warning = s.Base.Copy().
		Foreground(ColorWarning).
		Bold(true)

	s.Success = s.Base.Copy().
		Foreground(ColorSuccess).
		Bold(true)

	s.Info = s.Base.Copy().
		Foreground(ColorInfo)

	s.Debug = s.Base.Copy().
		Foreground(ColorDebug)

	// Layout styles
	s.Container = s.Base.Copy().
		Padding(0).
		Margin(0)

	s.Panel = s.Base.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorderPrimary).
		Padding(1).
		Margin(0, 1)

	s.PanelActive = s.Panel.Copy().
		BorderForeground(ColorBorderActive).
		BorderStyle(lipgloss.ThickBorder())

	s.Header = s.Base.Copy().
		Background(ColorBgSecondary).
		Foreground(ColorTextPrimary).
		Padding(0, 1).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorBorderPrimary)

	s.Footer = s.Header.Copy().
		Border(lipgloss.NormalBorder(), true, false, false, false)

	s.Sidebar = s.Panel.Copy().
		Width(20).
		Background(ColorBgSecondary)

	s.Content = s.Panel.Copy().
		Background(ColorBgPrimary)

	// Component styles
	s.Tab = s.Base.Copy().
		Padding(0, 2).
		Background(ColorBgSecondary).
		Foreground(ColorTextSecondary)

	s.TabActive = s.Tab.Copy().
		Background(ColorBgSelected).
		Foreground(ColorTextPrimary).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorPrimary)

	s.TabInactive = s.Tab.Copy().
		Background(ColorBgTertiary).
		Foreground(ColorTextMuted)

	s.Button = s.Base.Copy().
		Padding(0, 2).
		Background(ColorBgTertiary).
		Foreground(ColorTextPrimary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorderPrimary)

	s.ButtonActive = s.Button.Copy().
		Background(ColorPrimary).
		Foreground(ColorTextInverse).
		BorderForeground(ColorPrimary).
		Bold(true)

	s.Input = s.Base.Copy().
		Padding(0, 1).
		Background(ColorBgSecondary).
		Foreground(ColorTextPrimary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorderSecondary)

	s.InputActive = s.Input.Copy().
		BorderForeground(ColorPrimary)

	// Log styles
	s.LogEntry = s.Base.Copy().
		Padding(0, 1).
		MarginBottom(0)

	s.LogTimestamp = s.Base.Copy().
		Foreground(ColorTextMuted).
		Width(19) // "2006-01-02 15:04:05"

	s.LogLevel = s.Base.Copy().
		Width(7).
		Align(lipgloss.Center).
		Bold(true)

	s.LogMessage = s.Base.Copy().
		Foreground(ColorTextPrimary)

	s.LogContext = s.Base.Copy().
		Foreground(ColorTextSecondary).
		Italic(true)

	// Metric styles
	s.MetricCard = s.Panel.Copy().
		Width(20).
		Height(6).
		Align(lipgloss.Center)

	s.MetricValue = s.Base.Copy().
		Foreground(ColorPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(18)

	s.MetricLabel = s.Base.Copy().
		Foreground(ColorTextSecondary).
		Align(lipgloss.Center).
		Width(18)

	s.MetricUnit = s.Base.Copy().
		Foreground(ColorTextMuted).
		Align(lipgloss.Center)

	s.GaugeBar = s.Base.Copy().
		Background(ColorBgTertiary).
		Height(1)

	s.GaugeFill = s.Base.Copy().
		Background(ColorPrimary).
		Height(1)

	// Dashboard styles
	s.Widget = s.Panel.Copy().
		Background(ColorBgSecondary)

	s.WidgetTitle = s.Base.Copy().
		Foreground(ColorTextPrimary).
		Bold(true).
		Padding(0, 1).
		Background(ColorBgTertiary).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorBorderPrimary)

	s.WidgetBorder = s.Base.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorderPrimary)

	// Table styles
	s.TableHeader = s.Base.Copy().
		Foreground(ColorTextPrimary).
		Background(ColorBgTertiary).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorBorderPrimary)

	s.TableRow = s.Base.Copy().
		Foreground(ColorTextPrimary).
		Padding(0, 1)

	s.TableRowAlt = s.TableRow.Copy().
		Background(ColorBgSecondary)

	s.TableCell = s.Base.Copy().
		Padding(0, 1)

	// Progress styles
	s.ProgressBar = s.Base.Copy().
		Background(ColorBgTertiary).
		Height(1)

	s.ProgressFill = s.Base.Copy().
		Background(ColorSuccess).
		Height(1)

	s.ProgressText = s.Base.Copy().
		Foreground(ColorTextPrimary).
		Align(lipgloss.Center)

	// Status styles
	s.StatusHealthy = s.Base.Copy().
		Foreground(ColorSuccess).
		Bold(true)

	s.StatusDegraded = s.Base.Copy().
		Foreground(ColorWarning).
		Bold(true)

	s.StatusDown = s.Base.Copy().
		Foreground(ColorError).
		Bold(true)

	s.StatusUnknown = s.Base.Copy().
		Foreground(ColorTextMuted).
		Bold(true)

	// Help styles
	s.HelpKey = s.Base.Copy().
		Foreground(ColorPrimary).
		Bold(true)

	s.HelpValue = s.Base.Copy().
		Foreground(ColorTextSecondary)

	return s
}

// LogLevelStyle returns the appropriate style for a log level
func (s *Styles) LogLevelStyle(level string) lipgloss.Style {
	switch level {
	case "DEBUG":
		return s.LogLevel.Copy().Foreground(ColorDebug)
	case "INFO":
		return s.LogLevel.Copy().Foreground(ColorInfo)
	case "WARN":
		return s.LogLevel.Copy().Foreground(ColorWarning)
	case "ERROR":
		return s.LogLevel.Copy().Foreground(ColorError)
	case "SUCCESS":
		return s.LogLevel.Copy().Foreground(ColorSuccess)
	default:
		return s.LogLevel.Copy().Foreground(ColorTextMuted)
	}
}

// StatusStyle returns the appropriate style for a status
func (s *Styles) StatusStyle(status string) lipgloss.Style {
	switch status {
	case "healthy", "online", "active", "running":
		return s.StatusHealthy
	case "degraded", "warning", "slow":
		return s.StatusDegraded
	case "down", "offline", "error", "failed":
		return s.StatusDown
	default:
		return s.StatusUnknown
	}
}

// ProgressStyle returns the appropriate style for progress status
func (s *Styles) ProgressStyle(status string) lipgloss.Style {
	switch status {
	case "completed":
		return s.ProgressFill.Copy().Background(ColorSuccess)
	case "error":
		return s.ProgressFill.Copy().Background(ColorError)
	case "in_progress":
		return s.ProgressFill.Copy().Background(ColorPrimary)
	default:
		return s.ProgressFill.Copy().Background(ColorTextMuted)
	}
}

// GetMetricColor returns appropriate color for metric values
func GetMetricColor(value, min, max float64) lipgloss.Color {
	if max == min {
		return ColorPrimary
	}

	ratio := (value - min) / (max - min)
	
	switch {
	case ratio >= 0.8:
		return ColorError
	case ratio >= 0.6:
		return ColorWarning
	case ratio >= 0.4:
		return ColorInfo
	default:
		return ColorSuccess
	}
}

// RenderBorder renders a border with the given title
func RenderBorder(content string, title string, style lipgloss.Style) string {
	// Note: BorderTitle is not available in this version of lipgloss
	// We'll implement title rendering differently
	return style.Render(content)
}

// Truncate truncates text to the specified width with ellipsis
func Truncate(text string, width int) string {
	if len(text) <= width {
		return text
	}
	if width <= 3 {
		return "..."
	}
	return text[:width-3] + "..."
}

// PadRight pads text to the right with spaces
func PadRight(text string, width int) string {
	if len(text) >= width {
		return text
	}
	return text + lipgloss.NewStyle().Width(width-len(text)).Render("")
}

// PadLeft pads text to the left with spaces
func PadLeft(text string, width int) string {
	if len(text) >= width {
		return text
	}
	return lipgloss.NewStyle().Width(width-len(text)).Render("") + text
}

// Center centers text within the given width
func Center(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	leftPad := lipgloss.NewStyle().Width(padding).Render("")
	rightPad := lipgloss.NewStyle().Width(width-len(text)-padding).Render("")
	return leftPad + text + rightPad
}