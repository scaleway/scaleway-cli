//go:build !wasm

package terminal

import (
	"os"

	"golang.org/x/term"

	"github.com/fatih/color"
)

func Style(msg string, styles ...color.Attribute) string {
	return color.New(styles...).Sprint(msg)
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
