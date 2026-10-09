package tui

import tea "github.com/charmbracelet/bubbletea"

// lineEdit applies readline-style editing keys to a text buffer. It
// returns the new buffer and whether the key was handled.
//
// Every text input of the TUI (the goto menu query, the `:` command
// prompt, the `/` filter and find prompts) routes its keys through
// this function, so a readline shortcut added to the switch below
// becomes available in all of them at once.
func lineEdit(key tea.KeyType, buffer string) (string, bool) {
	switch key {
	case tea.KeyCtrlU: // kill the whole line
		return "", true
	case tea.KeyCtrlH, tea.KeyBackspace:
		if buffer == "" {
			return buffer, true
		}

		return buffer[:len(buffer)-1], true
	}

	return buffer, false
}
