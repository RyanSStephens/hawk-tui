package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Table represents a data table
type Table struct {
	*BaseComponent
	headers     []string
	rows        [][]string
	widths      []int
	cursor      int
	offset      int
	selectable  bool
	onSelect    func(int)
	highlighted map[int]bool
	sortable    bool
	sortColumn  int
	sortAsc     bool
}

// NewTable creates a new table
func NewTable(headers []string) *Table {
	return &Table{
		BaseComponent: NewBaseComponent(nil),
		headers:       headers,
		rows:          make([][]string, 0),
		widths:        make([]int, len(headers)),
		selectable:    true,
		highlighted:   make(map[int]bool),
		sortable:      false,
		sortAsc:       true,
	}
}

// NewTableWithTheme creates a new table with a specific theme
func NewTableWithTheme(headers []string, theme *styles.Theme) *Table {
	return &Table{
		BaseComponent: NewBaseComponent(theme),
		headers:       headers,
		rows:          make([][]string, 0),
		widths:        make([]int, len(headers)),
		selectable:    true,
		highlighted:   make(map[int]bool),
		sortable:      false,
		sortAsc:       true,
	}
}

// SetRows sets the table rows
func (t *Table) SetRows(rows [][]string) {
	t.rows = rows
	t.calculateWidths()
}

// AddRow adds a single row to the table
func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
	t.calculateWidths()
}

// SetColumnWidths sets custom column widths
func (t *Table) SetColumnWidths(widths []int) {
	t.widths = widths
}

// SetSelectable sets whether rows can be selected
func (t *Table) SetSelectable(selectable bool) {
	t.selectable = selectable
}

// SetOnSelect sets the selection callback
func (t *Table) SetOnSelect(fn func(int)) {
	t.onSelect = fn
}

// SetHighlighted sets a row as highlighted
func (t *Table) SetHighlighted(row int, highlighted bool) {
	if highlighted {
		t.highlighted[row] = true
	} else {
		delete(t.highlighted, row)
	}
}

// SelectedRow returns the currently selected row index
func (t *Table) SelectedRow() int {
	return t.cursor
}

// calculateWidths auto-calculates column widths based on content
func (t *Table) calculateWidths() {
	// Initialize with header widths
	for i, header := range t.headers {
		if i >= len(t.widths) {
			t.widths = append(t.widths, len(header))
		} else if t.widths[i] == 0 {
			t.widths[i] = len(header)
		}
	}

	// Update based on row content
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(t.widths) && len(cell) > t.widths[i] {
				t.widths[i] = len(cell)
			}
		}
	}
}

// Init implements tea.Model
func (t *Table) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !t.focused || !t.selectable {
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
				if t.cursor < t.offset {
					t.offset = t.cursor
				}
			}
		case "down", "j":
			if t.cursor < len(t.rows)-1 {
				t.cursor++
				maxVisible := t.height - 2 // Account for header
				if t.cursor >= t.offset+maxVisible {
					t.offset = t.cursor - maxVisible + 1
				}
			}
		case "g":
			t.cursor = 0
			t.offset = 0
		case "G":
			t.cursor = len(t.rows) - 1
			maxVisible := t.height - 2
			if len(t.rows) > maxVisible {
				t.offset = len(t.rows) - maxVisible
			}
		case "enter", " ":
			if t.onSelect != nil {
				t.onSelect(t.cursor)
			}
		}
	}

	return t, nil
}

// View implements tea.Model
func (t *Table) View() string {
	if !t.visible {
		return ""
	}

	var b strings.Builder

	// Header style
	headerStyle := lipgloss.NewStyle().
		Foreground(t.theme.TextPrimary).
		Background(t.theme.BgTertiary).
		Bold(true).
		Padding(0, 1)

	// Render headers
	var headerCells []string
	for i, header := range t.headers {
		width := t.widths[i] + 2 // Add padding
		cell := headerStyle.Width(width).Render(truncate(header, width))
		headerCells = append(headerCells, cell)
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))
	b.WriteString("\n")

	// Render rows
	maxVisible := t.height - 2
	if maxVisible <= 0 {
		maxVisible = len(t.rows)
	}

	endRow := t.offset + maxVisible
	if endRow > len(t.rows) {
		endRow = len(t.rows)
	}

	for idx := t.offset; idx < endRow; idx++ {
		row := t.rows[idx]

		// Determine row style
		var rowStyle lipgloss.Style
		if idx == t.cursor && t.focused {
			rowStyle = lipgloss.NewStyle().
				Background(t.theme.BgSelected).
				Foreground(t.theme.TextPrimary)
		} else if t.highlighted[idx] {
			rowStyle = lipgloss.NewStyle().
				Background(t.theme.BgHover).
				Foreground(t.theme.TextPrimary)
		} else if idx%2 == 0 {
			rowStyle = lipgloss.NewStyle().
				Background(t.theme.BgPrimary).
				Foreground(t.theme.TextPrimary)
		} else {
			rowStyle = lipgloss.NewStyle().
				Background(t.theme.BgSecondary).
				Foreground(t.theme.TextPrimary)
		}

		rowStyle = rowStyle.Padding(0, 1)

		// Render cells
		var cells []string
		for i, cell := range row {
			if i >= len(t.widths) {
				break
			}
			width := t.widths[i] + 2
			cellStr := rowStyle.Width(width).Render(truncate(cell, width))
			cells = append(cells, cellStr)
		}

		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cells...))
		if idx < endRow-1 {
			b.WriteString("\n")
		}
	}

	// Add border if focused
	result := b.String()
	if t.focused {
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.theme.BorderActive)
		result = borderStyle.Render(result)
	}

	return result
}

// truncate truncates text to fit within the given width
func truncate(text string, width int) string {
	if len(text) <= width {
		return text
	}
	if width <= 3 {
		return "..."
	}
	return text[:width-3] + "..."
}
