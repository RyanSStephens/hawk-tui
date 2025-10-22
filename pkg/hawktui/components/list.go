package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// ListItem represents an item in a list
type ListItem struct {
	Title       string
	Description string
	Value       interface{}
}

// List represents a selectable list
type List struct {
	*BaseComponent
	items      []ListItem
	cursor     int
	offset     int
	onSelect   func(ListItem)
	filterText string
	filtered   []int // Indices of filtered items
}

// NewList creates a new list
func NewList(items []ListItem) *List {
	l := &List{
		BaseComponent: NewBaseComponent(nil),
		items:         items,
		filtered:      make([]int, len(items)),
	}
	// Initialize filtered indices
	for i := range items {
		l.filtered[i] = i
	}
	return l
}

// NewListWithTheme creates a new list with a specific theme
func NewListWithTheme(items []ListItem, theme *styles.Theme) *List {
	l := &List{
		BaseComponent: NewBaseComponent(theme),
		items:         items,
		filtered:      make([]int, len(items)),
	}
	for i := range items {
		l.filtered[i] = i
	}
	return l
}

// SetItems sets the list items
func (l *List) SetItems(items []ListItem) {
	l.items = items
	l.cursor = 0
	l.offset = 0
	l.updateFiltered()
}

// AddItem adds a single item to the list
func (l *List) AddItem(item ListItem) {
	l.items = append(l.items, item)
	l.updateFiltered()
}

// SetOnSelect sets the selection callback
func (l *List) SetOnSelect(fn func(ListItem)) {
	l.onSelect = fn
}

// SetFilter sets the filter text
func (l *List) SetFilter(filter string) {
	l.filterText = filter
	l.cursor = 0
	l.offset = 0
	l.updateFiltered()
}

// SelectedItem returns the currently selected item
func (l *List) SelectedItem() *ListItem {
	if len(l.filtered) == 0 || l.cursor >= len(l.filtered) {
		return nil
	}
	idx := l.filtered[l.cursor]
	if idx >= len(l.items) {
		return nil
	}
	return &l.items[idx]
}

// updateFiltered updates the filtered items based on filter text
func (l *List) updateFiltered() {
	l.filtered = l.filtered[:0]

	if l.filterText == "" {
		for i := range l.items {
			l.filtered = append(l.filtered, i)
		}
		return
	}

	filterLower := strings.ToLower(l.filterText)
	for i, item := range l.items {
		titleLower := strings.ToLower(item.Title)
		descLower := strings.ToLower(item.Description)
		if strings.Contains(titleLower, filterLower) || strings.Contains(descLower, filterLower) {
			l.filtered = append(l.filtered, i)
		}
	}
}

// Init implements tea.Model
func (l *List) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (l *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !l.focused {
		return l, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if l.cursor > 0 {
				l.cursor--
				if l.cursor < l.offset {
					l.offset = l.cursor
				}
			}
		case "down", "j":
			if l.cursor < len(l.filtered)-1 {
				l.cursor++
				maxVisible := l.height - 2
				if l.cursor >= l.offset+maxVisible {
					l.offset = l.cursor - maxVisible + 1
				}
			}
		case "g":
			l.cursor = 0
			l.offset = 0
		case "G":
			l.cursor = len(l.filtered) - 1
			maxVisible := l.height - 2
			if len(l.filtered) > maxVisible {
				l.offset = len(l.filtered) - maxVisible
			}
		case "enter", " ":
			if item := l.SelectedItem(); item != nil && l.onSelect != nil {
				l.onSelect(*item)
			}
		}
	case tea.MouseMsg:
		// Handle mouse events
		if l.HandleMouse(msg) {
			// Calculate which item was clicked
			localY := msg.Y - l.y
			if localY >= 0 && msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionRelease {
				// Item index accounting for offset
				clickedItem := localY + l.offset
				if clickedItem >= 0 && clickedItem < len(l.filtered) {
					l.cursor = clickedItem
					if item := l.SelectedItem(); item != nil && l.onSelect != nil {
						l.onSelect(*item)
					}
				}
			}
		}
	}

	return l, nil
}

// View implements tea.Model
func (l *List) View() string {
	if !l.visible {
		return ""
	}

	var b strings.Builder

	maxVisible := l.height
	if maxVisible <= 0 {
		maxVisible = len(l.filtered)
	}

	endIdx := l.offset + maxVisible
	if endIdx > len(l.filtered) {
		endIdx = len(l.filtered)
	}

	for i := l.offset; i < endIdx; i++ {
		itemIdx := l.filtered[i]
		if itemIdx >= len(l.items) {
			continue
		}
		item := l.items[itemIdx]

		// Determine style
		var itemStyle lipgloss.Style
		if i == l.cursor && l.focused {
			itemStyle = lipgloss.NewStyle().
				Background(l.theme.BgSelected).
				Foreground(l.theme.TextPrimary).
				Bold(true).
				Padding(0, 1)
		} else {
			itemStyle = lipgloss.NewStyle().
				Background(l.theme.BgPrimary).
				Foreground(l.theme.TextPrimary).
				Padding(0, 1)
		}

		// Render cursor indicator
		cursor := " "
		if i == l.cursor && l.focused {
			cursor = ">"
		}

		// Render item
		title := lipgloss.NewStyle().
			Foreground(l.theme.Primary).
			Bold(true).
			Render(item.Title)

		desc := ""
		if item.Description != "" {
			desc = "\n  " + lipgloss.NewStyle().
				Foreground(l.theme.TextSecondary).
				Render(item.Description)
		}

		itemText := cursor + " " + title + desc

		if l.width > 0 {
			itemStyle = itemStyle.Width(l.width - 2)
		}

		b.WriteString(itemStyle.Render(itemText))
		if i < endIdx-1 {
			b.WriteString("\n")
		}
	}

	result := b.String()

	// Add border if focused
	if l.focused {
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(l.theme.BorderActive)

		if l.width > 0 {
			borderStyle = borderStyle.Width(l.width - 2)
		}
		if l.height > 0 {
			borderStyle = borderStyle.Height(l.height - 2)
		}

		result = borderStyle.Render(result)
	}

	return result
}
