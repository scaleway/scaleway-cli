package tui

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// LocalityMenu lists every zone or region so the user can pick the
// session's zone or region. It is a type-to-filter list (like the goto
// Menu): printable chars filter, arrows/pgup/pgdown move, enter selects
// (emitting a setLocalityMsg), esc closes.
type LocalityMenu struct {
	*listPicker
	isZone bool
}

// NewLocalityMenu builds the zone (isZone) or region picker, starting
// the cursor on the locality currently in use.
func NewLocalityMenu(isZone bool, loc *Locality, keys *Keymap) *LocalityMenu {
	var items []string
	var current, title string
	if isZone {
		title = "Select zone"
		items = make([]string, 0, len(scw.AllZones))
		for _, z := range scw.AllZones {
			items = append(items, string(z))
		}
		current = string(loc.GetZone())
	} else {
		title = "Select region"
		items = make([]string, 0, len(scw.AllRegions))
		for _, r := range scw.AllRegions {
			items = append(items, string(r))
		}
		current = string(loc.GetRegion())
	}
	sort.Strings(items)

	return &LocalityMenu{isZone: isZone, listPicker: newListPicker(title, items, current, keys)}
}

// Name implements Screen.
func (m *LocalityMenu) Name() string {
	if m.isZone {
		return "select-zone"
	}

	return "select-region"
}

// Hints implements Screen.
func (m *LocalityMenu) Hints() []string {
	label := "zone"
	if !m.isZone {
		label = "region"
	}
	k := m.keys

	return []string{
		"[ type ] filter",
		"[ ↑↓ ] move",
		hint(k.Key(ActionOpen), "select "+label),
		hint(k.Key(ActionBack), "cancel"),
	}
}

// Update implements Screen.
func (m *LocalityMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m, m.listPicker.handleKey(km, m.selectVisible)
}

// View implements Screen.
func (m *LocalityMenu) View() string {
	if m.width <= 0 {
		return ""
	}

	return strings.Join(m.listPicker.renderLines(), "\n")
}

func (m *LocalityMenu) selectVisible() tea.Cmd {
	vis := m.visible()
	if m.cursor < len(vis) {
		val := m.items[vis[m.cursor]]
		if m.isZone {
			return func() tea.Msg { return setLocalityMsg{isZone: true, zone: scw.Zone(val)} }
		}

		return func() tea.Msg { return setLocalityMsg{region: scw.Region(val)} }
	}

	return nil
}
