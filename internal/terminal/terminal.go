//go:build !wasm

package terminal

import (
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/scaleway/scaleway-cli/v2/internal/color"
	"golang.org/x/term"
)

// colorEnabled plays the role of the former fatih/color NoColor flag: text
// attributes are only applied when stdout is a color-capable terminal
// (honoring NO_COLOR, TERM, CLICOLOR...).
var colorEnabled = colorprofile.Detect(os.Stdout, os.Environ()) >= colorprofile.ANSI

func Style(msg string, styles ...color.Attribute) string {
	if !colorEnabled {
		return msg
	}

	var style color.Attribute
	for _, s := range styles {
		style = append(style, s...)
	}

	return style.Styled(msg)
}

func GetWidth() int {
	w, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil {
		return -1
	}

	return w
}

func GetHeight() int {
	_, h, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil {
		return -1
	}

	return h
}

// IsTerm returns if stdout is considered a tty
func IsTerm() bool {
	return colorEnabled
}
