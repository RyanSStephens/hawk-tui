package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Panel represents a bordered container for content
type Panel struct {
	*BaseComponent
	title       string
	content     string
	borderStyle lipgloss.Border
	padding     int
}

// NewPanel creates a new panel
func NewPanel(title string) *Panel {
	return &Panel{
		BaseComponent: NewBaseComponent(nil),
		title:         title,
		borderStyle:   lipgloss.RoundedBorder(),
		padding:       1,
	}
}

// NewPanelWithTheme creates a new panel with a specific theme
func NewPanelWithTheme(title string, theme *styles.Theme) *Panel {
	return &Panel{
		BaseComponent: NewBaseComponent(theme),
		title:         title,
		borderStyle:   lipgloss.RoundedBorder(),
		padding:       1,
	}
}

// SetTitle sets the panel title
func (p *Panel) SetTitle(title string) {
	p.title = title
}

// SetContent sets the panel content
func (p *Panel) SetContent(content string) {
	p.content = content
}

// SetBorderStyle sets the border style
func (p *Panel) SetBorderStyle(border lipgloss.Border) {
	p.borderStyle = border
}

// SetPadding sets the padding
func (p *Panel) SetPadding(padding int) {
	p.padding = padding
}

// Init implements tea.Model
func (p *Panel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (p *Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, nil
}

// View implements tea.Model
func (p *Panel) View() string {
	if !p.visible {
		return ""
	}

	style := lipgloss.NewStyle().
		Border(p.borderStyle).
		Padding(p.padding)

	if p.focused {
		style = style.BorderForeground(p.theme.BorderActive)
	} else {
		style = style.BorderForeground(p.theme.BorderPrimary)
	}

	if p.width > 0 {
		style = style.Width(p.width - (p.padding * 2) - 2)
	}

	if p.height > 0 {
		style = style.Height(p.height - (p.padding * 2) - 2)
	}

	// Render title if present
	content := p.content
	if p.title != "" {
		titleStyle := lipgloss.NewStyle().
			Foreground(p.theme.Primary).
			Bold(true).
			Background(p.theme.BgSecondary).
			Padding(0, 1)

		title := titleStyle.Render(p.title)
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", p.content)
	}

	return style.Render(content)
}
