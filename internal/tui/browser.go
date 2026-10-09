package tui

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// fetchTimeout bounds a single list call.
const fetchTimeout = 30 * time.Second

// fetchMsg carries the result of a list fetch (k9s: TableDataChanged).
type fetchMsg struct {
	columns []string
	rows    []Row
	err     error
}

// refreshMsg asks a screen to re-fetch its data (k9s: the poll tick of
// model.Table.updater).
type refreshMsg struct{}

// Browser lists a Resource and lets the user open row details. It is the
// k9s view.Browser (the extender chain is not needed here: Scaleway
// resources share one generic list behavior).
type Browser struct {
	res   *Resource
	table *Table
	keys  *Keymap

	loading    bool
	err        string
	lastUpdate time.Time
}

// NewBrowser creates a Browser for the given resource. Columns are set
// from the first fetch (they are derived from the response shape).
func NewBrowser(res *Resource, keys *Keymap) *Browser {
	return &Browser{res: res, keys: keys, table: &Table{Keys: keys}}
}

// Name implements Screen.
func (b *Browser) Name() string {
	return b.res.Name
}

// Hints implements Screen. The keys come from the keymap so the hints
// track cli.yaml overrides.
func (b *Browser) Hints() []string {
	k := b.keys

	return []string{
		hint(k.Key(ActionOpen), "open"),
		hint(k.Key(ActionFilter), "filter"),
		hint(k.Key(ActionCopy), "copy id"),
		hint(k.Key(ActionCommand), "command"),
		hint(k.Key(ActionMenu), "menu"),
		hint(k.Key(ActionRefresh), "refresh"),
		hint(k.Key(ActionZone), "zone"),
		hint(k.Key(ActionRegion), "region"),
		hint(k.Key(ActionProfile), "profile"),
		hint(k.Key(ActionProject), "project"),
		hint(k.Key(ActionBack), "back"),
		hint(k.Key(ActionQuit), "quit"),
	}
}

// Size implements Screen. The browser reserves one line each for its
// title (the resource name) and its status line.
func (b *Browser) Size(width, height int) {
	b.table.Width = width
	b.table.Height = max(1, height-2)
}

// Table returns the underlying table (used by the app to apply filters).
func (b *Browser) Table() *Table {
	return b.table
}

// Init implements Screen.
func (b *Browser) Init() tea.Cmd {
	return b.fetch()
}

// Update implements Screen.
func (b *Browser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case fetchMsg:
		b.loading = false
		if m.err != nil {
			b.err = m.err.Error()

			return b, nil
		}
		b.err = ""
		b.lastUpdate = time.Now()
		b.table.Columns = m.columns
		b.table.SetRows(m.rows)

		return b, nil
	case refreshMsg:
		return b, b.fetch()
	case tea.KeyMsg:
		switch b.keys.Action(keyOf(m)) {
		case ActionOpen:
			if r := b.table.Selected(); r.Obj != nil {
				details := NewDetails(b.res.Name+":"+shortID(r.ID), r.Obj, b.keys)

				return b, func() tea.Msg { return pushScreenMsg{screen: details} }
			}
		case ActionCopy:
			if r := b.table.Selected(); r.ID != "" {
				return b, func() tea.Msg { return copyMsg{value: r.ID, label: "id"} }
			}
		}
		mdl, cmd := b.table.Update(msg)
		b.table = mdl.(*Table)

		return b, cmd
	}

	return b, nil
}

// View implements Screen. The first line is the browsed resource name
// (k9s title bar), above the table columns.
func (b *Browser) View() string {
	return titleStyle.Render(
		truncate(b.res.Title, b.table.Width),
	) + "\n" + b.table.View() + "\n" + b.statusLine()
}

// ponytail: no in-flight guard (k9s has an atomic one); a double fetch
// is harmless since the last fetchMsg wins.
func (b *Browser) fetch() tea.Cmd {
	b.loading = true

	return func() tea.Msg {
		return b.doFetch()
	}
}

// doFetch runs the fetch with a recover: 200+ heterogeneous commands
// mean one bad Run closure should degrade to an error line, not kill
// the TUI process.
func (b *Browser) doFetch() (m fetchMsg) {
	defer func() {
		if r := recover(); r != nil {
			m = fetchMsg{err: fmt.Errorf("fetch panicked: %v", r)}
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
	defer cancel()
	res, err := b.res.Fetch(ctx)
	if err != nil {
		return fetchMsg{err: err}
	}

	return fetchMsg{columns: res.columns, rows: res.rows}
}

func (b *Browser) statusLine() string {
	switch {
	case b.loading:
		return dimStyle.Render("loading…")
	case b.err != "":
		return errorStyle.Render("error: " + b.err)
	case !b.lastUpdate.IsZero():
		return dimStyle.Render(time.Since(b.lastUpdate).Round(time.Second).String() + " ago")
	default:
		return ""
	}
}
