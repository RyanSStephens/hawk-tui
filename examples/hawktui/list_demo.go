package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/components"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/templates"
)

// List view demo
type model struct {
	app      *hawktui.App
	listView *templates.ListView
	width    int
	height   int
}

func (m model) Init() tea.Cmd {
	return m.listView.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.listView.SetSize(m.width, m.height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	lvModel, cmd := m.listView.Update(msg)
	m.listView = lvModel.(*templates.ListView)
	return m, cmd
}

func (m model) View() string {
	return m.listView.View()
}

func main() {
	// Create app with Catppuccin theme
	app := hawktui.New(hawktui.WithTheme(hawktui.ThemeCatppuccin()))

	// Create list items
	items := []components.ListItem{
		{
			Title:       "Project Alpha",
			Description: "High priority - Due tomorrow",
			Value:       "alpha",
		},
		{
			Title:       "Project Beta",
			Description: "Medium priority - Due next week",
			Value:       "beta",
		},
		{
			Title:       "Project Gamma",
			Description: "Low priority - Due next month",
			Value:       "gamma",
		},
		{
			Title:       "Project Delta",
			Description: "Critical - Overdue!",
			Value:       "delta",
		},
		{
			Title:       "Project Epsilon",
			Description: "On hold - Waiting for approval",
			Value:       "epsilon",
		},
	}

	// Create list view
	listView := templates.NewListViewWithTheme("Projects", items, app.Theme())

	// Set selection handler
	listView.SetOnSelect(func(item components.ListItem) {
		app.Logger().Infof("Selected: %s", item.Title)
	})

	// Add action buttons
	listView.AddAction("New", func() {
		app.Logger().Info("Creating new project...")
	})

	listView.AddAction("Edit", func() {
		app.Logger().Info("Editing project...")
	})

	listView.AddAction("Delete", func() {
		app.Logger().Warn("Deleting project...")
	})

	m := model{
		app:      app,
		listView: listView,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
