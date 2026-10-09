package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// listPicker is the shared state and behavior of the "Select ..." list
// pickers (zone, region, profile): a scrollable list with a
// type-to-filter query. The screen embedding it only provides its
// title, its items and its selection message.
type listPicker struct {
	title   string
	items   []string
	current string // the active item, shown as "(current)"
	query   string
	cursor  int
	off     int
	width   int
	height  int
	keys    *Keymap
}

// newListPicker builds a picker, starting the cursor on the current item.
func newListPicker(title string, items []string, current string, keys *Keymap) *listPicker {
	cursor := 0
	for i, it := range items {
		if it != "" && it == current {
			cursor = i

			break
		}
	}

	return &listPicker{title: title, items: items, current: current, cursor: cursor, keys: keys}
}

// Size implements Screen.
func (p *listPicker) Size(width, height int) {
	p.width, p.height = width, height
}

// Init implements Screen.
func (p *listPicker) Init() tea.Cmd {
	return nil
}

// handleKey processes a key against the picker. onOpen is called for
// the open action (each picker emits its own selection message).
func (p *listPicker) handleKey(km tea.KeyMsg, onOpen func() tea.Cmd) tea.Cmd {
	action := p.keys.Action(keyOf(km))
	if p.feedQuery(km, action) {
		return nil
	}
	switch action {
	case ActionOpen:
		return onOpen()
	case ActionBack:
		if p.query != "" {
			p.query = ""
			p.cursor, p.off = 0, 0

			return nil
		}

		return func() tea.Msg { return popScreenMsg{} }
	case ActionUp:
		p.move(-1)
	case ActionDown:
		p.move(1)
	case ActionPageUp:
		p.move(-p.rows())
	case ActionPageDown:
		p.move(p.rows())
	case ActionHome:
		p.cursor, p.off = 0, 0
	case ActionEnd:
		if n := len(p.visible()); n > 0 {
			p.cursor = n - 1
			p.off = max(0, n-p.rows())
		}
	}

	return nil
}

// feedQuery treats the key as a type-to-filter keystroke. It returns
// true when the key was consumed. Keys bound to open/back are not
// typed into the query (cli.yaml keybinding overrides).
func (p *listPicker) feedQuery(km tea.KeyMsg, action string) bool {
	if edited, ok := lineEdit(km.Type, p.query); ok {
		if edited != p.query {
			p.query = edited
			p.cursor, p.off = 0, 0
		}

		return true
	}
	if km.Type == tea.KeyRunes && action != ActionOpen && action != ActionBack {
		if len(km.String()) == 1 && len(p.query) < 120 {
			p.query += km.String()
			p.cursor, p.off = 0, 0
		}

		return true
	}

	return false
}

// renderLines renders the "Select ..." screen body.
func (p *listPicker) renderLines() []string {
	lines := []string{titleStyle.Render(p.title)}
	if p.query != "" {
		lines[0] += "  " + promptStyle.Render(p.query)
	}
	vis := p.visible()
	if len(vis) == 0 {
		lines = append(lines, dimStyle.Render("no match"))
	}
	off := p.off
	if off > len(vis)-p.rows() {
		off = max(0, len(vis)-p.rows())
	}
	for i := off; i < len(vis) && len(lines) < p.height; i++ {
		lines = append(lines, p.renderItem(p.items[vis[i]], i == p.cursor))
	}
	for len(lines) < p.height {
		lines = append(lines, "")
	}

	return lines
}

func (p *listPicker) rows() int {
	return max(1, p.height-1)
}

func (p *listPicker) move(delta int) {
	n := len(p.visible())
	if n == 0 {
		return
	}
	p.cursor = min(max(0, p.cursor+delta), n-1)
	if p.cursor < p.off {
		p.off = p.cursor
	}
	if p.cursor >= p.off+p.rows() {
		p.off = p.cursor - p.rows() + 1
	}
}

func (p *listPicker) visible() []int {
	if p.query == "" {
		idx := make([]int, len(p.items))
		for i := range idx {
			idx[i] = i
		}

		return idx
	}
	var idx []int
	for i, item := range p.items {
		if fuzzyMatch(item, p.query) != nil {
			idx = append(idx, i)
		}
	}

	return idx
}

func (p *listPicker) renderItem(item string, cursor bool) string {
	label := item
	if p.query != "" {
		if pos := fuzzyMatch(label, p.query); pos != nil {
			label = highlightMatch(label, pos)
		}
	}
	if p.current != "" && item == p.current {
		label += dimStyle.Render("  (current)")
	}
	if cursor {
		return cursorStyle.Render("> " + label)
	}

	return "  " + menuLabelStyle.Render(label)
}
