//go:build wasm

package tui

import "github.com/scaleway/scaleway-cli/v2/core"

// GetCommands returns the tui commands. The TUI is disabled on wasm:
// bubbletea needs a real terminal, which the web build does not have
// (and its dependencies do not compile for js/wasm).
func GetCommands() *core.Commands {
	return core.NewCommands()
}
