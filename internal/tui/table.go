package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// cursorGutter is the 2-char `>` marker column (k9s cursor marker).
const cursorGutter = 2

// maxColWidth caps the natural width of a column.
const maxColWidth = 40

// Row is a single table row. Obj holds the underlying object used by
// the details view. CellStyles maps a column index to its style.
type Row struct {
	ID         string
	Cells      []string
	CellStyles map[int]lipgloss.Style
	Obj        any
}

// Table is a generic scrollable table with a cursor and a substring
// filter. It is the k9s ui.Table/view.Table widget pair, idiomatic
// bubbletea flavored.
type Table struct {
	Columns []string
	Rows    []Row

	Cursor int
	Filter string

	// Offset is the first visible row.
	Offset int

	// Width and Height are set by the app. Height includes the header
	// line.
	Width  int
	Height int

	// Keys is the keymap for the navigation keys (nil: defaults).
	Keys *Keymap
}

// SetRows replaces the row data, resetting cursor and scroll.
func (t *Table) SetRows(rows []Row) {
	t.Rows = rows
	t.Cursor = 0
	t.Offset = 0
}

// SetFilter sets the (case-insensitive) substring filter.
func (t *Table) SetFilter(f string) {
	t.Filter = strings.ToLower(f)
}

// Visible returns the rows matching the filter.
func (t *Table) Visible() []Row {
	if t.Filter == "" {
		return t.Rows
	}

	visible := make([]Row, 0, len(t.Rows))
	for _, r := range t.Rows {
		var b strings.Builder
		b.WriteString(r.ID)
		for _, c := range r.Cells {
			b.WriteString(" ")
			b.WriteString(c)
		}
		if strings.Contains(strings.ToLower(b.String()), t.Filter) {
			visible = append(visible, r)
		}
	}

	return visible
}

// Selected returns the row under the cursor.
func (t *Table) Selected() Row {
	v := t.Visible()
	if t.Cursor < 0 || t.Cursor >= len(v) {
		return Row{}
	}

	return v[t.Cursor]
}

// Init implements tea.Model.
func (t *Table) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model. Only navigation keys are handled here;
// `enter` and the global keys are intercepted by the parent screen first
// (k9s input-capture layering).
func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return t, nil
	}

	n := len(t.Visible())
	switch t.Keys.Action(keyOf(keyMsg)) {
	case ActionUp:
		if t.Cursor > 0 {
			t.Cursor--
		}
	case ActionDown:
		if t.Cursor < n-1 {
			t.Cursor++
		}
	case ActionPageUp:
		t.Cursor = max(0, t.Cursor-t.rowsPerPage())
	case ActionPageDown:
		t.Cursor = min(n-1, t.Cursor+t.rowsPerPage())
	case ActionHome:
		t.Cursor = 0
	case ActionEnd:
		t.Cursor = max(0, n-1)
	}
	t.clampScroll()

	return t, nil
}

// View implements tea.Model.
func (t *Table) View() string {
	if t.Width <= 0 {
		return ""
	}

	v := t.Visible()
	if len(v) == 0 {
		msg := "no data"
		if t.Filter != "" {
			msg = "no match for filter: " + t.Filter
		}

		lines := []string{dimStyle.Render(msg)}
		for len(lines) < t.Height {
			lines = append(lines, "")
		}

		return strings.Join(lines, "\n")
	}

	widths := t.columnWidths(v)
	lines := make([]string, 0, t.Height)
	lines = append(lines, headerStyle.Render(t.renderLine(widths, t.Columns, nil, false)))
	for i := t.Offset; i < len(v) && len(lines) < t.Height; i++ {
		lines = append(lines, t.renderLine(widths, v[i].Cells, v[i].CellStyles, i == t.Cursor))
	}
	for len(lines) < t.Height {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

func (t *Table) rowsPerPage() int {
	return max(1, t.Height-1)
}

func (t *Table) clampScroll() {
	n := len(t.Visible())
	if n == 0 {
		t.Offset = 0

		return
	}
	if t.Cursor < t.Offset {
		t.Offset = t.Cursor
	}
	if t.Cursor >= t.Offset+t.rowsPerPage() {
		t.Offset = t.Cursor - t.rowsPerPage() + 1
	}
	t.Offset = max(0, t.Offset)
}

func (t *Table) renderLine(
	widths []int,
	cells []string,
	styles map[int]lipgloss.Style,
	cursor bool,
) string {
	gutter := "  "
	if cursor {
		gutter = "> "
	}
	parts := make([]string, len(t.Columns))
	for c := range len(t.Columns) {
		cell := ""
		if c < len(cells) {
			cell = cells[c]
		}
		cell = pad(cell, widths[c])
		if cursor {
			// k9s inverts the whole cursor row; per-cell colors are dropped.
			parts[c] = cell

			continue
		}
		if st, ok := styles[c]; ok {
			parts[c] = st.Render(cell)
		} else {
			parts[c] = rowStyle.Render(cell)
		}
	}

	line := gutter + joinCells(parts)
	if cursor {
		line = cursorStyle.Render(line)
	}

	return line
}

// columnGrid fixes the width of shared columns so the grid lines up
// across resource views (server/volume/ip show ID, NAME, STATE, ... in
// the same positions). Columns not in the grid size to their content.
var columnGrid = map[string]int{
	"ID":        8,
	"NAME":      28,
	"STATE":     11,
	"TYPE":      14,
	"ZONE":      8,
	"REGION":    8,
	"SIZE":      9,
	"SERVER":    22,
	"ADDRESS":   15,
	"PUBLIC IP": 15,
	"CREATED":   16,
}

// columnWidths computes per-column widths: grid width for shared
// columns, natural size (capped at maxColWidth) otherwise, shrunk to
// fit t.Width.
// ponytail: greedy shave of the widest column, O(overflow) iterations;
// fine for <=100 rows, switch to a proportional split for wide tables.
func (t *Table) columnWidths(v []Row) []int {
	n := len(t.Columns)
	avail := t.Width - cursorGutter - (n - 1) // gutter + one space between columns
	avail = max(avail, n)
	widths := make([]int, n)
	for c, col := range t.Columns {
		if w, ok := columnGrid[col]; ok {
			widths[c] = min(w, maxColWidth)
		} else {
			widths[c] = min(len(col), maxColWidth)
		}
	}
	for _, r := range v {
		for c := 0; c < n && c < len(r.Cells); c++ {
			if _, ok := columnGrid[t.Columns[c]]; !ok {
				if l := min(len(r.Cells[c]), maxColWidth); l > widths[c] {
					widths[c] = l
				}
			}
		}
	}
	for sum(widths) > avail {
		widest := 0
		for c := 1; c < n; c++ {
			if widths[c] > widths[widest] {
				widest = c
			}
		}
		if widths[widest] <= 1 {
			break
		}
		widths[widest]--
	}

	return widths
}

func sum(widths []int) int {
	total := 0
	for _, w := range widths {
		total += w
	}

	return total
}

func joinCells(parts []string) string {
	var b strings.Builder
	for c, p := range parts {
		b.WriteString(p)
		if c < len(parts)-1 {
			b.WriteString(" ")
		}
	}

	return b.String()
}

// truncate shortens s to w runes, appending an ellipsis.
func truncate(s string, w int) string {
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	if w <= 1 {
		if w == 0 {
			return ""
		}

		return string(r[:1])
	}

	return string(r[:w-1]) + "…"
}

// pad right-pads s to w runes (truncating with an ellipsis if longer).
func pad(s string, w int) string {
	r := []rune(s)
	if len(r) >= w {
		return truncate(s, w)
	}

	return s + strings.Repeat(" ", w-len(r))
}
