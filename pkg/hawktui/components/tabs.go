package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Tab represents a single tab
type Tab struct {
	Title   string
	Content string
}

// Tabs represents a tabbed interface
type Tabs struct {
	*BaseComponent
	tabs        []Tab
	activeTab   int
	onTabChange func(int)
}

// NewTabs creates a new tabs component
func NewTabs(tabs []Tab) *Tabs {
	return &Tabs{
		BaseComponent: NewBaseComponent(nil),
		tabs:          tabs,
		activeTab:     0,
	}
}

// NewTabsWithTheme creates a new tabs component with a specific theme
func NewTabsWithTheme(tabs []Tab, theme *styles.Theme) *Tabs {
	return &Tabs{
		BaseComponent: NewBaseComponent(theme),
		tabs:          tabs,
		activeTab:     0,
	}
}

// SetTabs sets the tabs
func (t *Tabs) SetTabs(tabs []Tab) {
	t.tabs = tabs
	if t.activeTab >= len(tabs) {
		t.activeTab = len(tabs) - 1
	}
	if t.activeTab < 0 {
		t.activeTab = 0
	}
}

// AddTab adds a new tab
func (t *Tabs) AddTab(tab Tab) {
	t.tabs = append(t.tabs, tab)
}

// SetActiveTab sets the active tab by index
func (t *Tabs) SetActiveTab(index int) {
	if index >= 0 && index < len(t.tabs) {
		t.activeTab = index
		if t.onTabChange != nil {
			t.onTabChange(index)
		}
	}
}

// ActiveTab returns the index of the active tab
func (t *Tabs) ActiveTab() int {
	return t.activeTab
}

// SetOnTabChange sets the tab change callback
func (t *Tabs) SetOnTabChange(fn func(int)) {
	t.onTabChange = fn
}

// Init implements tea.Model
func (t *Tabs) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (t *Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !t.focused {
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "right", "l":
			t.SetActiveTab((t.activeTab + 1) % len(t.tabs))
		case "shift+tab", "left", "h":
			t.SetActiveTab((t.activeTab - 1 + len(t.tabs)) % len(t.tabs))
		}
	case tea.MouseMsg:
		// Handle mouse events
		if t.HandleMouse(msg) {
			// Calculate which tab was clicked based on X coordinate
			localX := msg.X - t.x
			if msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionRelease {
				// Calculate tab widths and positions
				tabWidth := t.width / len(t.tabs)
				if tabWidth > 0 {
					clickedTab := localX / tabWidth
					if clickedTab >= 0 && clickedTab < len(t.tabs) {
						t.SetActiveTab(clickedTab)
					}
				}
			}
		}
	}

	return t, nil
}

// View implements tea.Model
func (t *Tabs) View() string {
	if !t.visible || len(t.tabs) == 0 {
		return ""
	}

	// Render tab headers
	var tabHeaders []string
	for i, tab := range t.tabs {
		var tabStyle lipgloss.Style
		if i == t.activeTab {
			tabStyle = lipgloss.NewStyle().
				Background(t.theme.BgSelected).
				Foreground(t.theme.TextPrimary).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(t.theme.Primary)
		} else {
			tabStyle = lipgloss.NewStyle().
				Background(t.theme.BgTertiary).
				Foreground(t.theme.TextMuted).
				Padding(0, 2).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(t.theme.BorderPrimary)
		}
		tabHeaders = append(tabHeaders, tabStyle.Render(tab.Title))
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top, tabHeaders...)

	// Render content
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.theme.BorderPrimary).
		Padding(1)

	if t.width > 0 {
		contentStyle = contentStyle.Width(t.width - 4)
	}
	if t.height > 0 {
		contentStyle = contentStyle.Height(t.height - 5)
	}

	content := ""
	if t.activeTab < len(t.tabs) {
		content = t.tabs[t.activeTab].Content
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, contentStyle.Render(content))
}
