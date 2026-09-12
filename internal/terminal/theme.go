package terminal

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	tint "github.com/lrstanley/bubbletint/v2"
)

// scalewayTint is the default Scaleway brand theme, using the official
// Scaleway purple/violet palette on a dark blue-gray background.
// Brand colors sourced from scaleway.com SVG assets and console CSS.
var scalewayTint = &tint.Tint{
	DisplayName: "Scaleway",
	ID:          "scaleway",
	Dark:        true,
	Fg:          tint.FromHex("#b8bac0"),
	Bg:          tint.FromHex("#1e2335"),
	SelectionBg: tint.FromHex("#252a3b"),
	Cursor:      tint.FromHex("#a365f6"),
	Black:       tint.FromHex("#303445"),
	Red:         tint.FromHex("#ff5555"),
	Green:       tint.FromHex("#50fa7b"),
	Yellow:      tint.FromHex("#f1fa8c"),
	Blue:        tint.FromHex("#a365f6"),
	Purple:      tint.FromHex("#4f0599"),
	Cyan:        tint.FromHex("#a060f6"),
	White:       tint.FromHex("#b8bac0"),
	BrightBlack:  tint.FromHex("#484b5a"),
	BrightRed:    tint.FromHex("#ff7777"),
	BrightGreen:  tint.FromHex("#73ffb0"),
	BrightYellow: tint.FromHex("#fff09e"),
	BrightBlue:   tint.FromHex("#c490f8"),
	BrightPurple: tint.FromHex("#521094"),
	BrightCyan:   tint.FromHex("#b884f8"),
	BrightWhite:  tint.FromHex("#f1eefc"),
}

// CuratedThemeIDs is the list of popular themes shown by default in 'scw config theme list'.
var CuratedThemeIDs = []string{
	"scaleway",
	"dracula_plus",
	"catppuccin_mocha",
	"catppuccin_frappe",
	"nord",
	"one_dark",
	"monokai_pro",
	"gruvbox_dark",
	"tokyo_night",
	"github",
	"solarized_dark___patched",
	"material_darker",
	"rose_pine",
	"cyberpunk",
	"snazzy",
	"dracula",
	"atom_one_light",
	"catppuccin_latte",
	"tomorrow_night",
}

// attributeToColor maps fatih/color.Attribute values to the corresponding
// bubbletint Tint color field accessor function.
// Style-only attributes (bold, italic) are handled separately.
var themeActive bool

var warnOnce sync.Once

// InitTheme initializes the bubbletint registry and activates the theme with the given ID.
// If themeID is empty, no theme is activated (legacy mode).
// If themeID is not found, falls back to legacy mode and emits a one-time warning to stderr.
func InitTheme(themeID string) {
	tint.NewDefaultRegistry()
	tint.Register(scalewayTint)

	if themeID == "" {
		return
	}

	if ok := tint.SetTintID(themeID); !ok {
		warnOnce.Do(func() {
			fmt.Fprintf(os.Stderr, "warning: theme %q not found, falling back to default. Run 'scw config theme list' to see available themes.\n", themeID)
		})
		return
	}

	if !color.NoColor {
		themeActive = true
	}
}

// IsThemeActive returns true when a non-empty tint has been set in the registry.
func IsThemeActive() bool {
	return themeActive && !color.NoColor
}

// SetThemeActive sets the themeActive flag, used by 'scw config theme set' to
// activate theming at runtime after InitTheme has been called.
func SetThemeActive(active bool) {
	themeActive = active
}

// styleWithTheme resolves color.Attribute values against the active bubbletint tint
// and emits true-color ANSI escape sequences. If no theme is active, falls back
// to standard fatih/color behavior.
func styleWithTheme(msg string, styles ...color.Attribute) string {
	if !IsThemeActive() {
		return color.New(styles...).Sprint(msg)
	}

	t := tint.Current()

	var prefix string
	for _, attr := range styles {
		var c *tint.Color
		switch attr {
		case color.FgGreen:
			c = t.Green
		case color.FgRed:
			c = t.Red
		case color.FgYellow:
			c = t.Yellow
		case color.FgBlue:
			c = t.Blue
		case color.FgCyan:
			c = t.Cyan
		case color.FgMagenta:
			c = t.Purple
		case color.FgWhite:
			c = t.White
		case color.FgHiGreen:
			c = t.BrightGreen
		case color.FgHiRed:
			c = t.BrightRed
		case color.FgHiYellow:
			c = t.BrightYellow
		case color.FgHiCyan:
			c = t.BrightCyan
		case color.FgHiBlue:
			c = t.BrightBlue
		case color.FgHiMagenta:
			c = t.BrightPurple
		case color.FgHiWhite:
			c = t.BrightWhite
		case color.Bold:
			prefix += "\x1b[1m"
			if t.BrightWhite != nil {
				prefix += fmt.Sprintf("\x1b[38;2;%d;%d;%dm", t.BrightWhite.R, t.BrightWhite.G, t.BrightWhite.B)
			} else if t.Fg != nil {
				prefix += fmt.Sprintf("\x1b[38;2;%d;%d;%dm", t.Fg.R, t.Fg.G, t.Fg.B)
			}
			continue
		case color.Italic:
			prefix += "\x1b[3m"
			continue
		default:
			prefix += fmt.Sprintf("\x1b[%dm", attr)
			continue
		}

		if c != nil {
			prefix += fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
		}
	}

	if prefix == "" {
		return msg
	}

	return prefix + msg + "\x1b[0m"
}
