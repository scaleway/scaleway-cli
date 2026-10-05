//go:build wasm

package terminal

import (
	"github.com/fatih/color"
)

var (
	Width  int = 0
	Height int = 0
)

func Style(msg string, styles ...color.Attribute) string {
	if color.NoColor {
		return msg
	}
	return styleWithTheme(msg, styles...)
}

func GetWidth() int {
	return Width
}

func GetHeight() int {
	return Height
}

// IsTerm returns if stdout is considered a tty
func IsTerm() bool {
	return true
}
