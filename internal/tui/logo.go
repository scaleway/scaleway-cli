package tui

import "github.com/charmbracelet/lipgloss"

// scalewayLogo is the Scaleway "S" mark as ASCII art (18x8), derived
// from the official logo. Rendered in the Scaleway brand purple.
var scalewayLogo = []string{
	"  #############   ",
	" ###.      :##### ",
	" ### .#####   :## ",
	" ### ###   ##  ###",
	" ### .##  ###  ###",
	" ###    #####  ###",
	"  ###:  ####.  ###",
	"   ############## ",
}

// logoStyle is the Scaleway brand purple.
var logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#55139a"))
