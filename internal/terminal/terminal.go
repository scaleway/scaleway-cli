package terminal

import (
	"os"

	"github.com/fatih/color"
	"golang.org/x/term"
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
