package tui

import "github.com/charmbracelet/lipgloss"

// Palette taken from the k9s "stock" skin (k9scli.io default look):
// black background, blue/aqua text, orange accent, inverted aqua cursor.
const (
	colorBg         = "0"
	colorDodgerblue = "33"
	colorAqua       = "51"
	colorCadetblue  = "72"
	colorSteelblue  = "67"
	colorLightsky   = "75"
	colorWhite      = "231"
	colorOrange     = "208"
	colorGreen      = "42"
	colorRed        = "203"
	colorYellow     = "220"
	colorDim        = "245"
)

var (
	// titleStyle is the k9s title bar (aqua).
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorAqua))
	// hintStyle renders the key hints in the title bar.
	hintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorDodgerblue))
	// headerStyle is the table header row (white, bold).
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorWhite))
	// rowStyle is a normal table row (dodgerblue).
	rowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorDodgerblue))
	// cursorStyle is the k9s inverted cursor row: aqua bg, black fg.
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorBg)).
			Background(lipgloss.Color(colorAqua))
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorDim))
	// promptStyle is the `:`/`/` command bar (cadetblue).
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorCadetblue))
	// statusStyle is the bottom status bar text.
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorLightsky))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(colorRed))
	// flashStyle is the status-bar flash message (orange).
	flashStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorOrange))
	// menu styles (k9s menu: white label).
	menuLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorWhite))
	// detail styles (k9s yaml view: steelblue keys, papayawhip values).
	detailKeyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSteelblue))
	detailValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
)

// stateStyle colorizes a resource state cell (k9s-style health coloring).
func stateStyle(state string) lipgloss.Style {
	switch state {
	case "running", "active":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen))
	case "stopped", "dead", "error":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorRed))
	case "standby", "locked", "in_transfer", "migrating", "resize":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorYellow))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorDodgerblue))
	}
}
