package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// setProfileMsg asks the app to switch the session to the named scw
// profile (re-creating the API client for it).
type setProfileMsg struct{ profile string }

// ProfileMenu lists the profiles defined in the scw config file so the
// user can switch the session's profile (ctrl+p).
type ProfileMenu struct {
	*listPicker
}

// NewProfileMenu builds the profile picker, starting the cursor on the
// profile currently in use.
func NewProfileMenu(items []string, current string, keys *Keymap) *ProfileMenu {
	return &ProfileMenu{listPicker: newListPicker("Select profile", items, current, keys)}
}

// Name implements Screen.
func (m *ProfileMenu) Name() string {
	return "select-profile"
}

// Hints implements Screen.
func (m *ProfileMenu) Hints() []string {
	k := m.keys

	return []string{
		"[ type ] filter",
		"[ ↑↓ ] move",
		hint(k.Key(ActionOpen), "select profile"),
		hint(k.Key(ActionBack), "cancel"),
	}
}

// Update implements Screen.
func (m *ProfileMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m, m.listPicker.handleKey(km, m.selectVisible)
}

// View implements Screen.
func (m *ProfileMenu) View() string {
	if m.width <= 0 {
		return ""
	}

	return strings.Join(m.listPicker.renderLines(), "\n")
}

func (m *ProfileMenu) selectVisible() tea.Cmd {
	vis := m.visible()
	if m.cursor < len(vis) {
		val := m.items[vis[m.cursor]]

		return func() tea.Msg { return setProfileMsg{profile: val} }
	}

	return nil
}
