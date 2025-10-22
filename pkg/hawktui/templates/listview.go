package templates

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/components"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// ListView represents a list view template with search and actions
type ListView struct {
	theme      *styles.Theme
	title      string
	list       *components.List
	searchBar  *components.Input
	actions    []*components.Button
	showSearch bool
	width      int
	height     int
}

// NewListView creates a new list view template
func NewListView(title string, items []components.ListItem) *ListView {
	theme := styles.DefaultTheme()

	lv := &ListView{
		theme:      theme,
		title:      title,
		list:       components.NewListWithTheme(items, theme),
		searchBar:  components.NewInputWithTheme("Search...", theme),
		showSearch: true,
		actions:    make([]*components.Button, 0),
	}

	lv.searchBar.SetOnChange(func(value string) {
		lv.list.SetFilter(value)
	})

	return lv
}

// NewListViewWithTheme creates a new list view with a specific theme
func NewListViewWithTheme(title string, items []components.ListItem, theme *styles.Theme) *ListView {
	lv := &ListView{
		theme:      theme,
		title:      title,
		list:       components.NewListWithTheme(items, theme),
		searchBar:  components.NewInputWithTheme("Search...", theme),
		showSearch: true,
		actions:    make([]*components.Button, 0),
	}

	lv.searchBar.SetOnChange(func(value string) {
		lv.list.SetFilter(value)
	})

	return lv
}

// AddAction adds an action button
func (lv *ListView) AddAction(label string, onPress func()) {
	btn := components.NewButtonWithTheme(label, onPress, lv.theme)
	lv.actions = append(lv.actions, btn)
}

// SetShowSearch sets whether to show the search bar
func (lv *ListView) SetShowSearch(show bool) {
	lv.showSearch = show
}

// SetOnSelect sets the item selection callback
func (lv *ListView) SetOnSelect(fn func(components.ListItem)) {
	lv.list.SetOnSelect(fn)
}

// SetSize sets the list view dimensions
func (lv *ListView) SetSize(width, height int) {
	lv.width = width
	lv.height = height
	lv.list.SetSize(width-4, height-10)
}

// Init implements tea.Model
func (lv *ListView) Init() tea.Cmd {
	lv.list.SetFocused(true)
	return nil
}

// Update implements tea.Model
func (lv *ListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			lv.searchBar.SetFocused(true)
			lv.list.SetFocused(false)
			return lv, nil
		case "esc":
			if lv.searchBar.IsFocused() {
				lv.searchBar.SetFocused(false)
				lv.list.SetFocused(true)
				return lv, nil
			}
		}
	}

	// Update search bar if focused
	if lv.searchBar.IsFocused() {
		var model tea.Model
		var cmd tea.Cmd
		model, cmd = lv.searchBar.Update(msg)
		lv.searchBar = model.(*components.Input)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	} else {
		// Update list
		var model tea.Model
		var cmd tea.Cmd
		model, cmd = lv.list.Update(msg)
		lv.list = model.(*components.List)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) > 0 {
		return lv, tea.Batch(cmds...)
	}

	return lv, nil
}

// View implements tea.Model
func (lv *ListView) View() string {
	var sections []string

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lv.theme.Primary).
		Bold(true).
		Padding(1, 2).
		Background(lv.theme.BgSecondary).
		Width(lv.width)

	sections = append(sections, titleStyle.Render(lv.title))

	// Search bar
	if lv.showSearch {
		lv.searchBar.SetSize(lv.width-4, 1)
		searchSection := lipgloss.NewStyle().
			Padding(1, 2).
			Render(lv.searchBar.View())
		sections = append(sections, searchSection)
	}

	// List
	sections = append(sections, lv.list.View())

	// Actions
	if len(lv.actions) > 0 {
		var actionViews []string
		for _, action := range lv.actions {
			actionViews = append(actionViews, action.View())
		}
		actionsBar := lipgloss.NewStyle().
			Padding(1, 2).
			Render(lipgloss.JoinHorizontal(lipgloss.Left, actionViews...))
		sections = append(sections, actionsBar)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
