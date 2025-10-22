package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/templates"
)

// Dashboard demo showing template usage
type model struct {
	app       *hawktui.App
	dashboard *templates.Dashboard
	width     int
	height    int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dashboard.SetSize(m.width, m.height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	dashModel, cmd := m.dashboard.Update(msg)
	m.dashboard = dashModel.(*templates.Dashboard)
	return m, cmd
}

func (m model) View() string {
	return m.dashboard.View()
}

func main() {
	// Create app with Dracula theme
	app := hawktui.New(hawktui.WithTheme(hawktui.ThemeDracula()))

	// Create dashboard
	dashboard := templates.NewDashboardWithTheme(app.Theme())
	dashboard.SetTitle("System Dashboard")
	dashboard.SetColumns(3)

	// Add widgets
	dashboard.AddWidget(templates.CreateMetricWidget(
		"CPU Usage",
		"45.2%",
		"8 cores",
		app.Theme(),
	))

	dashboard.AddWidget(templates.CreateMetricWidget(
		"Memory",
		"8.5 GB",
		"16 GB total",
		app.Theme(),
	))

	dashboard.AddWidget(templates.CreateMetricWidget(
		"Disk I/O",
		"234 MB/s",
		"read/write",
		app.Theme(),
	))

	dashboard.AddWidget(templates.CreateStatusWidget(
		"API Server",
		"healthy",
		"All endpoints responding normally",
		app.Theme(),
	))

	dashboard.AddWidget(templates.CreateStatusWidget(
		"Database",
		"healthy",
		"Connection pool: 45/100",
		app.Theme(),
	))

	dashboard.AddWidget(templates.CreateStatusWidget(
		"Cache",
		"warning",
		"Hit rate: 72% (low)",
		app.Theme(),
	))

	// Add chart widget
	chartData := map[string]float64{
		"API":     234,
		"Web":     189,
		"Mobile":  156,
		"Desktop": 98,
	}
	dashboard.AddWidget(templates.CreateChartWidget(
		"Requests by Platform",
		chartData,
		250,
		app.Theme(),
	))

	m := model{
		app:       app,
		dashboard: dashboard,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
