// Package color provides the named text attributes (colors and styles) used
// across the CLI. They are built on charmbracelet/x/ansi, the SGR layer of
// the lipgloss stack, replacing the former fatih/color dependency.
package color

import "github.com/charmbracelet/x/ansi"

// Attribute is an ANSI SGR attribute (color or style) applied to text,
// see terminal.Style.
type Attribute = ansi.Style

// Named foreground colors (ANSI 4-bit).
var (
	FgRed    = ansi.Style{}.ForegroundColor(ansi.Red)
	FgGreen  = ansi.Style{}.ForegroundColor(ansi.Green)
	FgYellow = ansi.Style{}.ForegroundColor(ansi.Yellow)
	FgBlue   = ansi.Style{}.ForegroundColor(ansi.Blue)
	FgCyan   = ansi.Style{}.ForegroundColor(ansi.Cyan)

	FgHiGreen  = ansi.Style{}.ForegroundColor(ansi.BrightGreen)
	FgHiBlue   = ansi.Style{}.ForegroundColor(ansi.BrightBlue)
	FgHiYellow = ansi.Style{}.ForegroundColor(ansi.BrightYellow)

	Bold   = ansi.Style{}.Bold()
	Faint  = ansi.Style{}.Faint()
	Italic = ansi.Style{}.Italic(true)
)
