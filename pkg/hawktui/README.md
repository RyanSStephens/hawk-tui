# HawkTUI - Terminal UI Toolkit

A comprehensive, batteries-included UI toolkit for Go terminal applications. HawkTUI combines the power of styling (like lipgloss), TUI framework components (like bubbletea), structured logging, and ready-to-use templates in one cohesive package.

## Features

- **Theme System**: Multiple built-in themes (Dark, Light, Nord, Dracula, Catppuccin) with easy customization
- **Rich Components**: Buttons, inputs, tables, lists, progress bars, spinners, panels, tabs, and more
- **Flexible Layouts**: Horizontal, vertical, grid, and flexbox-style layouts
- **Structured Logging**: Integrated, styled logging with multiple log levels
- **Pre-built Templates**: Dashboard, form, list view, and split view templates ready to use
- **Type-safe**: Full Go type safety with interfaces and strong typing
- **Extensible**: Easy to extend with custom components and themes

## Installation

\`\`\`bash
go get github.com/hawk-tui/hawk-tui/pkg/hawktui
\`\`\`

## Quick Start

\`\`\`go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/hawk-tui/hawk-tui/pkg/hawktui"
)

func main() {
    // Create HawkTUI app with a theme
    app := hawktui.New(hawktui.WithTheme(hawktui.ThemeDark()))

    // Create a button
    button := hawktui.NewButton("Click Me!", func() {
        app.Logger().Success("Button clicked!")
    })
    button.SetTheme(app.Theme())

    // Run your bubbletea program
    p := tea.NewProgram(yourModel, tea.WithAltScreen())
    p.Run()
}
\`\`\`

## Themes

HawkTUI includes 5 built-in themes:

\`\`\`go
hawktui.ThemeDark()        // Default dark theme
hawktui.ThemeLight()       // Light theme
hawktui.ThemeNord()        // Nord color scheme
hawktui.ThemeDracula()     // Dracula color scheme
hawktui.ThemeCatppuccin()  // Catppuccin Mocha
\`\`\`

See the full README for complete documentation on components, layouts, templates, and examples.

## Examples

Run the included examples:

\`\`\`bash
go run examples/hawktui/simple_demo.go
go run examples/hawktui/dashboard_demo.go
go run examples/hawktui/form_demo.go
go run examples/hawktui/list_demo.go
\`\`\`

## License

This project is dual-licensed under AGPL-3.0 and a commercial license.
