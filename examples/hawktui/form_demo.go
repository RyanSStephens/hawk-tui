package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/templates"
)

// Form demo
type model struct {
	app    *hawktui.App
	form   *templates.Form
	done   bool
	result map[string]string
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form.SetSize(m.width, m.height)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || (msg.String() == "q" && m.done) {
			return m, tea.Quit
		}
	}

	if !m.done {
		formModel, cmd := m.form.Update(msg)
		m.form = formModel.(*templates.Form)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		theme := m.app.Theme()

		successMsg := theme.SuccessStyle.Render("Form submitted successfully!")

		var fields []string
		fields = append(fields, successMsg, "")
		for key, value := range m.result {
			fields = append(fields, fmt.Sprintf("%s: %s", key, value))
		}
		fields = append(fields, "", "Press q to quit")

		return hawktui.JoinVertical(0, fields...)
	}

	return m.form.View()
}

func main() {
	// Create app with Nord theme
	app := hawktui.New(hawktui.WithTheme(hawktui.ThemeNord()))

	// Create form
	form := templates.NewFormWithTheme("User Registration", app.Theme())

	// Add fields
	form.AddField("Username", "Enter your username", true)
	form.AddField("Email", "your@email.com", true)
	form.AddPasswordField("Password", "Enter password", true)
	form.AddPasswordField("Confirm Password", "Confirm password", true)
	form.AddField("Full Name", "John Doe", false)

	m := model{
		app:  app,
		form: form,
	}

	// Set callbacks
	form.SetOnSubmit(func(values map[string]string) {
		m.result = values
		m.done = true
		app.Logger().Success("Form submitted!")
	})

	form.SetOnCancel(func() {
		app.Logger().Info("Form cancelled")
	})

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
