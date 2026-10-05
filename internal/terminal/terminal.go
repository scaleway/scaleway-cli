//go:build !wasm

package terminal

import (
	"os"

	"github.com/fatih/color"
	"golang.org/x/term"
)

func Style(msg string, styles ...color.Attribute) string {
	if color.NoColor {
		return msg
	}
	return styleWithTheme(msg, styles...)
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
	return !color.NoColor
}
