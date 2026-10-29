//go:build wasm

package terminal

import (
	"github.com/scaleway/scaleway-cli/v2/internal/color"
)

var (
	Width  int = 0
	Height int = 0

	// colorEnabled is set by the wasm host, see wasm.ConfigureOutput.
	colorEnabled = false
)

// SetColorEnabled enables or disables text attributes. It is called by the
// wasm host once it knows whether the embedding terminal renders colors.
func SetColorEnabled(enabled bool) {
	colorEnabled = enabled
}

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
	return Width
}

func GetHeight() int {
	return Height
}

// IsTerm returns if stdout is considered a tty
func IsTerm() bool {
	return true
}
