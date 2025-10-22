package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// SpinnerStyle defines different spinner animations
type SpinnerStyle int

const (
	SpinnerDots SpinnerStyle = iota
	SpinnerLine
	SpinnerCircle
	SpinnerBounce
	SpinnerPulse
)

// Spinner represents a loading spinner
type Spinner struct {
	*BaseComponent
	style   SpinnerStyle
	frame   int
	frames  []string
	label   string
	running bool
}

// Spinner frame definitions
var spinnerFrames = map[SpinnerStyle][]string{
	SpinnerDots:   {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	SpinnerLine:   {"-", "\\", "|", "/"},
	SpinnerCircle: {"◐", "◓", "◑", "◒"},
	SpinnerBounce: {"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"},
	SpinnerPulse:  {"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	return &Spinner{
		BaseComponent: NewBaseComponent(nil),
		style:         SpinnerDots,
		frames:        spinnerFrames[SpinnerDots],
		running:       false,
	}
}

// NewSpinnerWithTheme creates a new spinner with a specific theme
func NewSpinnerWithTheme(theme *styles.Theme) *Spinner {
	return &Spinner{
		BaseComponent: NewBaseComponent(theme),
		style:         SpinnerDots,
		frames:        spinnerFrames[SpinnerDots],
		running:       false,
	}
}

// SetStyle sets the spinner style
func (s *Spinner) SetStyle(style SpinnerStyle) {
	s.style = style
	s.frames = spinnerFrames[style]
	s.frame = 0
}

// SetLabel sets the spinner label
func (s *Spinner) SetLabel(label string) {
	s.label = label
}

// Start starts the spinner animation
func (s *Spinner) Start() tea.Cmd {
	s.running = true
	return s.tick()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.running = false
}

// IsRunning returns whether the spinner is running
func (s *Spinner) IsRunning() bool {
	return s.running
}

// tick creates a tick command for animation
func (s *Spinner) tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg(t)
	})
}

type spinnerTickMsg time.Time

// Init implements tea.Model
func (s *Spinner) Init() tea.Cmd {
	if s.running {
		return s.tick()
	}
	return nil
}

// Update implements tea.Model
func (s *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerTickMsg:
		if s.running {
			s.frame = (s.frame + 1) % len(s.frames)
			return s, s.tick()
		}
	case tea.MouseMsg:
		// Handle mouse events for click detection
		s.HandleMouse(msg)
	}
	return s, nil
}

// View implements tea.Model
func (s *Spinner) View() string {
	if !s.visible || !s.running {
		return ""
	}

	spinnerStyle := lipgloss.NewStyle().
		Foreground(s.theme.Primary).
		Bold(true)

	spinner := spinnerStyle.Render(s.frames[s.frame])

	if s.label != "" {
		labelStyle := lipgloss.NewStyle().
			Foreground(s.theme.TextPrimary).
			Padding(0, 1)
		return lipgloss.JoinHorizontal(lipgloss.Center, spinner, labelStyle.Render(s.label))
	}

	return spinner
}
