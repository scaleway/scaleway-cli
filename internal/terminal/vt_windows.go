//go:build windows

package terminal

import "golang.org/x/sys/windows"

func init() {
	// Enable ANSI (VT) processing on the Windows console so that the SGR
	// sequences printed by the CLI are rendered. This replaces the former
	// mattn/go-colorable stdout/stderr wrappers.
	enableVT(windows.Stdout)
	enableVT(windows.Stderr)
}

func enableVT(h windows.Handle) {
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err == nil {
		_ = windows.SetConsoleMode(h, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}
