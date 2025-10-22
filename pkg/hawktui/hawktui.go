// Package hawktui provides a comprehensive UI toolkit for terminal applications.
// HawkTUI combines styling (lipgloss-like), TUI framework components (bubbletea-like),
// structured logging, and ready-to-use templates in one cohesive package.
//
// Features:
//   - Theme system with multiple built-in themes
//   - Reusable components (buttons, inputs, tables, lists, etc.)
//   - Flexible layout system
//   - Integrated styled logging
//   - Pre-built templates for common UI patterns
//
// Example usage:
//
//	import "github.com/hawk-tui/hawk-tui/pkg/hawktui"
//
//	// Create a HawkTUI application with a theme
//	app := hawktui.New(hawktui.WithTheme(hawktui.ThemeDark))
//
//	// Use components
//	button := hawktui.NewButton("Click Me", func() { /* handler */ })
//	table := hawktui.NewTable([]string{"Name", "Value"})
//
//	// Use templates
//	dashboard := hawktui.NewDashboardTemplate()
package hawktui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/components"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/layouts"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/logger"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/templates"
)

// Version is the current version of the UI toolkit
const Version = "0.1.0"

// App represents a UI toolkit application
type App struct {
	theme  *styles.Theme
	logger *logger.Logger
	width  int
	height int
}

// Option is a functional option for configuring the App
type Option func(*App)

// New creates a new UI toolkit application
func New(opts ...Option) *App {
	app := &App{
		theme:  styles.DefaultTheme(),
		logger: logger.New(),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// WithTheme sets the theme for the application
func WithTheme(theme *styles.Theme) Option {
	return func(app *App) {
		app.theme = theme
	}
}

// WithLogger sets a custom logger for the application
func WithLogger(log *logger.Logger) Option {
	return func(app *App) {
		app.logger = log
	}
}

// Theme returns the current theme
func (a *App) Theme() *styles.Theme {
	return a.theme
}

// Logger returns the application logger
func (a *App) Logger() *logger.Logger {
	return a.logger
}

// SetSize sets the application dimensions
func (a *App) SetSize(width, height int) {
	a.width = width
	a.height = height
}

// Re-export commonly used types and functions for convenience
type (
	// Model is the bubbletea model interface
	Model = tea.Model
	// Cmd is the bubbletea command type
	Cmd = tea.Cmd
	// Msg is the bubbletea message interface
	Msg = tea.Msg
	// Style is the lipgloss style type
	Style = lipgloss.Style
)

// Re-export component constructors
var (
	NewButton      = components.NewButton
	NewInput       = components.NewInput
	NewTable       = components.NewTable
	NewList        = components.NewList
	NewProgressBar = components.NewProgressBar
	NewSpinner     = components.NewSpinner
	NewPanel       = components.NewPanel
	NewTabs        = components.NewTabs
)

// Re-export layout functions
var (
	Horizontal   = layouts.Horizontal
	Vertical     = layouts.Vertical
	Grid         = layouts.Grid
	NewFlexbox   = layouts.NewFlexbox
	NewContainer = layouts.NewContainer
)

// Re-export template constructors
var (
	NewDashboard = templates.NewDashboard
	NewForm      = templates.NewForm
	NewListView  = templates.NewListView
	NewSplitView = templates.NewSplitView
)

// Re-export theme constructors
var (
	ThemeDark    = styles.DarkTheme
	ThemeLight   = styles.LightTheme
	ThemeNord    = styles.NordTheme
	ThemeDracula = styles.DraculaTheme
	ThemeCatppuccin = styles.CatppuccinTheme
)

// Helper functions

// Render is a shorthand for lipgloss rendering
func Render(style Style, text string) string {
	return style.Render(text)
}

// JoinVertical joins strings vertically
func JoinVertical(pos lipgloss.Position, strs ...string) string {
	return lipgloss.JoinVertical(pos, strs...)
}

// JoinHorizontal joins strings horizontally
func JoinHorizontal(pos lipgloss.Position, strs ...string) string {
	return lipgloss.JoinHorizontal(pos, strs...)
}
