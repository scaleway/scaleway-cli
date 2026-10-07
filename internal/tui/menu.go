package tui

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// gotoResourceMsg asks the app to navigate to a resource, replacing the
// stack (k9s gotoResource).
type gotoResourceMsg struct{ res *Resource }

// Menu is the k9s "goto" menu: the list of resources to switch between.
// Opened with the m/g keys. Typing filters the list with a fuzzy
// (subsequence) match, matched characters are highlighted.
type Menu struct {
	items  []menuItem
	query  string
	cursor int // index into the filtered list
	off    int // first visible row
	width  int
	height int
	keys   *Keymap
}

type menuItem struct {
	label string
	res   *Resource
}

// NewMenu builds the goto menu from the available resources.
func NewMenu(resources map[string]*Resource, keys *Keymap) *Menu {
	names := make([]string, 0, len(resources))
	for name := range resources {
		names = append(names, name)
	}
	sort.Strings(names)

	m := &Menu{keys: keys}
	for _, name := range names {
		m.items = append(m.items, menuItem{label: name, res: resources[name]})
	}

	return m
}

// Name implements Screen.
func (m *Menu) Name() string {
	return "goto"
}

// Hints implements Screen.
func (m *Menu) Hints() []string {
	return []string{
		"[ type ] filter", "[ ↑↓ ] move", "[ enter ] open",
		"[ esc ] back/clear", "[ ctrl+c ] quit",
	}
}

// Size implements Screen.
func (m *Menu) Size(width, height int) {
	m.width, m.height = width, height
}

// Init implements Screen.
func (m *Menu) Init() tea.Cmd {
	return nil
}

// Update implements Screen.
func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	// The menu is a type-to-navigate widget (k9s command buffer):
	// readline keys always edit the query, every other printable rune
	// feeds the fuzzy query — including k, j, space, q — except the
	// runes explicitly bound to open or back. The remaining keys
	// (enter, esc, arrows, pgup/pgdown) are bound through the keymap.
	action := m.keys.Action(keyOf(km))
	if edited, ok := lineEdit(km.Type, m.query); ok {
		if edited != m.query {
			m.query = edited
			m.cursor, m.off = 0, 0
		}

		return m, nil
	}
	if km.Type == tea.KeyRunes && action != ActionOpen && action != ActionBack {
		if len(km.String()) == 1 && len(m.query) < 120 {
			m.query += km.String()
			m.cursor, m.off = 0, 0
		}

		return m, nil
	}

	switch action {
	case ActionOpen:
		return m, m.selectVisible()
	case ActionBack:
		if m.query != "" {
			m.query = ""
			m.cursor, m.off = 0, 0

			return m, nil
		}

		return m, func() tea.Msg { return popScreenMsg{} }
	case ActionUp:
		m.move(-1)
	case ActionDown:
		m.move(1)
	case ActionPageUp:
		m.move(-m.rows())
	case ActionPageDown:
		m.move(m.rows())
	}

	return m, nil
}

// View implements Screen.
func (m *Menu) View() string {
	if m.width <= 0 {
		return ""
	}
	head := titleStyle.Render("Goto")
	if m.query != "" {
		head += "  " + promptStyle.Render(m.query)
	}
	lines := []string{head}
	vis := m.visible()
	if len(vis) == 0 {
		lines = append(lines, dimStyle.Render("no match"))
	}
	off := m.off
	if off > len(vis)-m.rows() {
		off = max(0, len(vis)-m.rows())
	}
	for i := off; i < len(vis) && len(lines) < m.height; i++ {
		lines = append(lines, m.renderItem(m.items[vis[i]], i == m.cursor))
	}
	for len(lines) < m.height {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

func (m *Menu) rows() int {
	return max(1, m.height-1) // one line for the title
}

func (m *Menu) move(delta int) {
	n := len(m.visible())
	if n == 0 {
		return
	}
	m.cursor = min(max(0, m.cursor+delta), n-1)
	if m.cursor < m.off {
		m.off = m.cursor
	}
	if m.cursor >= m.off+m.rows() {
		m.off = m.cursor - m.rows() + 1
	}
}

// visible returns the indices into items matching the fuzzy query.
func (m *Menu) visible() []int {
	if m.query == "" {
		idx := make([]int, len(m.items))
		for i := range idx {
			idx[i] = i
		}

		return idx
	}
	var idx []int
	for i, item := range m.items {
		if fuzzyMatch(item.label, m.query) != nil {
			idx = append(idx, i)
		}
	}

	return idx
}

func (m *Menu) selectVisible() tea.Cmd {
	vis := m.visible()
	if m.cursor < len(vis) {
		res := m.items[vis[m.cursor]].res

		return func() tea.Msg { return gotoResourceMsg{res: res} }
	}

	return nil
}

func (m *Menu) renderItem(item menuItem, cursor bool) string {
	label := item.label
	if m.query != "" {
		if pos := fuzzyMatch(label, m.query); pos != nil {
			label = highlightMatch(label, pos)
		}
	}
	if cursor {
		return cursorStyle.Render("> " + label)
	}

	return "  " + menuLabelStyle.Render(label)
}

// fuzzyMatch returns the positions of query's runes in label (ordered
// subsequence, case-insensitive), or nil on no match.
func fuzzyMatch(label, query string) []int {
	ql := strings.ToLower(query)
	var pos []int
	qi := 0
	for i, r := range label {
		if qi < len(ql) && strings.ToLower(string(r)) == string(ql[qi]) {
			pos = append(pos, i)
			qi++
		}
	}
	if qi < len(ql) {
		return nil
	}

	return pos
}

// highlightMatch wraps the matched runes in the orange flash style
// (k9s highlights the matched menu substring).
func highlightMatch(s string, pos []int) string {
	matched := make(map[int]bool, len(pos))
	for _, p := range pos {
		matched[p] = true
	}
	var b strings.Builder
	for i, r := range s {
		if matched[i] {
			b.WriteString(flashStyle.Render(string(r)))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}
