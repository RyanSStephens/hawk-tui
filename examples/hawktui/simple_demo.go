package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/components"
)

// Simple demo showing basic components
type model struct {
	app     *hawktui.App
	button  *components.Button
	input   *components.Input
	spinner *components.Spinner
	width   int
	height  int
}

func (m model) Init() tea.Cmd {
	return m.spinner.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.input.IsFocused() {
				m.input.SetFocused(false)
				m.button.SetFocused(true)
			} else {
				m.input.SetFocused(true)
				m.button.SetFocused(false)
			}
			return m, nil
		}
	}

	var cmds []tea.Cmd

	// Update input
	if m.input.IsFocused() {
		inputModel, cmd := m.input.Update(msg)
		m.input = inputModel.(*components.Input)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update button
	if m.button.IsFocused() {
		btnModel, cmd := m.button.Update(msg)
		m.button = btnModel.(*components.Button)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update spinner
	spinnerModel, cmd := m.spinner.Update(msg)
	m.spinner = spinnerModel.(*components.Spinner)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) View() string {
	theme := m.app.Theme()

	title := theme.Title.Render("UITK Simple Demo")
	subtitle := theme.Subtitle.Render("Press Tab to switch focus, q to quit")

	m.input.SetSize(40, 1)
	input := m.input.View()

	button := m.button.View()

	spinner := m.spinner.View()

	content := hawktui.JoinVertical(
		0, // lipgloss.Left
		title,
		subtitle,
		"",
		input,
		"",
		button,
		"",
		spinner,
	)

	return theme.Base.
		Width(m.width).
		Height(m.height).
		Render(content)
}

func main() {
	// Create app with dark theme
	app := hawktui.New(hawktui.WithTheme(hawktui.ThemeDark()))

	// Create button
	button := hawktui.NewButton("Click Me!", func() {
		app.Logger().Success("Button clicked!")
	})
	button.SetTheme(app.Theme())

	// Create input
	input := hawktui.NewInput("Enter your name...")
	input.SetTheme(app.Theme())
	input.SetFocused(true)
	input.SetOnChange(func(value string) {
		app.Logger().Infof("Input changed to: %s", value)
	})

	// Create spinner
	spinner := hawktui.NewSpinner()
	spinner.SetTheme(app.Theme())
	spinner.SetLabel("Loading...")

	m := model{
		app:     app,
		button:  button,
		input:   input,
		spinner: spinner,
	}

	// Start spinner
	spinner.Start()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
