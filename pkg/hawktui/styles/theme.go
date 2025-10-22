package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color represents a color value
type Color = lipgloss.Color

// Theme defines the color scheme and styling for the UI toolkit
type Theme struct {
	Name string

	// Base colors
	Primary    Color
	Secondary  Color
	Accent     Color
	Success    Color
	Warning    Color
	Error      Color
	Info       Color
	Debug      Color

	// Background colors
	BgPrimary   Color
	BgSecondary Color
	BgTertiary  Color
	BgSelected  Color
	BgHover     Color

	// Text colors
	TextPrimary   Color
	TextSecondary Color
	TextMuted     Color
	TextInverse   Color

	// Border colors
	BorderPrimary   Color
	BorderSecondary Color
	BorderActive    Color
	BorderInactive  Color

	// Styles - pre-configured styles based on the theme
	Base          lipgloss.Style
	Focused       lipgloss.Style
	Blurred       lipgloss.Style
	Selected      lipgloss.Style
	ErrorStyle    lipgloss.Style
	WarningStyle  lipgloss.Style
	SuccessStyle  lipgloss.Style
	InfoStyle     lipgloss.Style
	DebugStyle    lipgloss.Style
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	Bold          lipgloss.Style
	Italic        lipgloss.Style
	Underline     lipgloss.Style
	Dimmed        lipgloss.Style
}

// NewTheme creates a new theme with the given colors
func NewTheme(name string) *Theme {
	t := &Theme{Name: name}
	t.initializeStyles()
	return t
}

// initializeStyles creates the pre-configured styles based on theme colors
func (t *Theme) initializeStyles() {
	t.Base = lipgloss.NewStyle().
		Foreground(t.TextPrimary).
		Background(t.BgPrimary)

	t.Focused = t.Base.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderActive)

	t.Blurred = t.Base.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderInactive)

	t.Selected = t.Base.Copy().
		Background(t.BgSelected).
		Foreground(t.TextPrimary).
		Bold(true)

	t.ErrorStyle = t.Base.Copy().
		Foreground(t.Error).
		Bold(true)

	t.WarningStyle = t.Base.Copy().
		Foreground(t.Warning).
		Bold(true)

	t.SuccessStyle = t.Base.Copy().
		Foreground(t.Success).
		Bold(true)

	t.InfoStyle = t.Base.Copy().
		Foreground(t.Info)

	t.DebugStyle = t.Base.Copy().
		Foreground(t.Debug)

	t.Title = t.Base.Copy().
		Foreground(t.Primary).
		Bold(true).
		Padding(0, 1)

	t.Subtitle = t.Base.Copy().
		Foreground(t.TextSecondary).
		Padding(0, 1)

	t.Bold = t.Base.Copy().Bold(true)
	t.Italic = t.Base.Copy().Italic(true)
	t.Underline = t.Base.Copy().Underline(true)
	t.Dimmed = t.Base.Copy().Foreground(t.TextMuted)
}

// DefaultTheme returns the default dark theme
func DefaultTheme() *Theme {
	return DarkTheme()
}

// DarkTheme returns a dark color scheme
func DarkTheme() *Theme {
	t := &Theme{
		Name:       "Dark",
		Primary:    Color("#00D7FF"),
		Secondary:  Color("#FF6B6B"),
		Accent:     Color("#BD93F9"),
		Success:    Color("#51CF66"),
		Warning:    Color("#FFD93D"),
		Error:      Color("#FF6B6B"),
		Info:       Color("#74C0FC"),
		Debug:      Color("#ADB5BD"),
		BgPrimary:  Color("#1A1B26"),
		BgSecondary: Color("#24283B"),
		BgTertiary: Color("#414868"),
		BgSelected: Color("#3D59A1"),
		BgHover:    Color("#292E42"),
		TextPrimary:   Color("#C0CAF5"),
		TextSecondary: Color("#9AA5CE"),
		TextMuted:     Color("#565F89"),
		TextInverse:   Color("#1A1B26"),
		BorderPrimary:   Color("#414868"),
		BorderSecondary: Color("#565F89"),
		BorderActive:    Color("#00D7FF"),
		BorderInactive:  Color("#292E42"),
	}
	t.initializeStyles()
	return t
}

// LightTheme returns a light color scheme
func LightTheme() *Theme {
	t := &Theme{
		Name:       "Light",
		Primary:    Color("#0066CC"),
		Secondary:  Color("#CC3333"),
		Accent:     Color("#6B4FBB"),
		Success:    Color("#2D9E4D"),
		Warning:    Color("#CC9900"),
		Error:      Color("#CC3333"),
		Info:       Color("#3399FF"),
		Debug:      Color("#666666"),
		BgPrimary:  Color("#FFFFFF"),
		BgSecondary: Color("#F5F5F5"),
		BgTertiary: Color("#E0E0E0"),
		BgSelected: Color("#CCE5FF"),
		BgHover:    Color("#EBEBEB"),
		TextPrimary:   Color("#1A1A1A"),
		TextSecondary: Color("#4D4D4D"),
		TextMuted:     Color("#999999"),
		TextInverse:   Color("#FFFFFF"),
		BorderPrimary:   Color("#CCCCCC"),
		BorderSecondary: Color("#E0E0E0"),
		BorderActive:    Color("#0066CC"),
		BorderInactive:  Color("#E8E8E8"),
	}
	t.initializeStyles()
	return t
}

// NordTheme returns the Nord color scheme
func NordTheme() *Theme {
	t := &Theme{
		Name:       "Nord",
		Primary:    Color("#88C0D0"),
		Secondary:  Color("#BF616A"),
		Accent:     Color("#B48EAD"),
		Success:    Color("#A3BE8C"),
		Warning:    Color("#EBCB8B"),
		Error:      Color("#BF616A"),
		Info:       Color("#81A1C1"),
		Debug:      Color("#4C566A"),
		BgPrimary:  Color("#2E3440"),
		BgSecondary: Color("#3B4252"),
		BgTertiary: Color("#434C5E"),
		BgSelected: Color("#5E81AC"),
		BgHover:    Color("#434C5E"),
		TextPrimary:   Color("#ECEFF4"),
		TextSecondary: Color("#D8DEE9"),
		TextMuted:     Color("#4C566A"),
		TextInverse:   Color("#2E3440"),
		BorderPrimary:   Color("#4C566A"),
		BorderSecondary: Color("#434C5E"),
		BorderActive:    Color("#88C0D0"),
		BorderInactive:  Color("#3B4252"),
	}
	t.initializeStyles()
	return t
}

// DraculaTheme returns the Dracula color scheme
func DraculaTheme() *Theme {
	t := &Theme{
		Name:       "Dracula",
		Primary:    Color("#BD93F9"),
		Secondary:  Color("#FF79C6"),
		Accent:     Color("#8BE9FD"),
		Success:    Color("#50FA7B"),
		Warning:    Color("#F1FA8C"),
		Error:      Color("#FF5555"),
		Info:       Color("#8BE9FD"),
		Debug:      Color("#6272A4"),
		BgPrimary:  Color("#282A36"),
		BgSecondary: Color("#343746"),
		BgTertiary: Color("#44475A"),
		BgSelected: Color("#6272A4"),
		BgHover:    Color("#383A59"),
		TextPrimary:   Color("#F8F8F2"),
		TextSecondary: Color("#E6E6E6"),
		TextMuted:     Color("#6272A4"),
		TextInverse:   Color("#282A36"),
		BorderPrimary:   Color("#6272A4"),
		BorderSecondary: Color("#44475A"),
		BorderActive:    Color("#BD93F9"),
		BorderInactive:  Color("#383A59"),
	}
	t.initializeStyles()
	return t
}

// CatppuccinTheme returns the Catppuccin Mocha color scheme
func CatppuccinTheme() *Theme {
	t := &Theme{
		Name:       "Catppuccin",
		Primary:    Color("#89B4FA"),
		Secondary:  Color("#F5C2E7"),
		Accent:     Color("#CBA6F7"),
		Success:    Color("#A6E3A1"),
		Warning:    Color("#F9E2AF"),
		Error:      Color("#F38BA8"),
		Info:       Color("#89DCEB"),
		Debug:      Color("#6C7086"),
		BgPrimary:  Color("#1E1E2E"),
		BgSecondary: Color("#313244"),
		BgTertiary: Color("#45475A"),
		BgSelected: Color("#585B70"),
		BgHover:    Color("#393B4D"),
		TextPrimary:   Color("#CDD6F4"),
		TextSecondary: Color("#BAC2DE"),
		TextMuted:     Color("#6C7086"),
		TextInverse:   Color("#1E1E2E"),
		BorderPrimary:   Color("#6C7086"),
		BorderSecondary: Color("#45475A"),
		BorderActive:    Color("#89B4FA"),
		BorderInactive:  Color("#393B4D"),
	}
	t.initializeStyles()
	return t
}
