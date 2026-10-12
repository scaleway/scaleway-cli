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
// `from` identifies the browser that issued the fetch: the app forwards
// fetchMsg to the current screen, so a result still in flight when the
// user navigates away must be dropped by the new screen, not applied
// (a stale row would open the old resource's details on enter).
type fetchMsg struct {
	from    *Browser
	columns []string
	rows    []Row
	err     error
}

// refreshMsg asks a screen to re-fetch its data (k9s: the poll tick of
// model.Table.updater).
type refreshMsg struct{}

// subListMsg carries the result of a drill-down fetch for one row
// (k9s: drill-down fetch). `from` identifies the browser that issued
// it, for the same stale-result reason as fetchMsg.
type subListMsg struct {
	from *Browser
	id   string
	res  listResult
	err  error
}

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
// track cli.yaml overrides. The connect hint is only shown when the
// resource has a registered connector.
func (b *Browser) Hints() []string {
	k := b.keys
	hints := []string{hint(k.Key(ActionOpen), "open")}
	if b.res.Editable() {
		hints = append(hints, hint(k.Key(ActionEdit), "edit"))
	}
	if b.res.Deleteable() {
		hints = append(hints, hint(k.Key(ActionDelete), "delete"))
	}
	if b.res.Sub != nil {
		hints = append(hints, hint(k.Key(ActionDrill), "drill"))
	}
	if conn, ok := ConnectorFor(b.res.Name); ok {
		hints = append(hints, hint(k.Key(ActionConnect), conn.Hint))
	}

	return append(hints,
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
	)
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
		if m.from != b { // stale: issued by another (previous) browser
			return b, nil
		}
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
	case subListMsg:
		if m.from != b { // stale: issued by another (previous) browser
			return b, nil
		}
		b.loading = false
		if m.err != nil {
			return b, func() tea.Msg { return flashMsg{msg: "drill: " + m.err.Error()} }
		}
		if len(m.res.rows) == 0 {
			return b, func() tea.Msg {
				return flashMsg{msg: "no " + b.res.Sub.Res.Name + " for " + shortID(m.id)}
			}
		}
		screen := NewBrowser(b.drillResource(m.id), b.keys)

		return b, func() tea.Msg { return pushScreenMsg{screen: screen} }
	case refreshMsg:
		return b, b.fetch()
	case tea.KeyMsg:
		switch b.keys.Action(keyOf(m)) {
		case ActionOpen:
			if r := b.table.Selected(); r.Obj != nil {
				details := NewDetails(b.res.Name+":"+shortID(r.ID), r.Obj, b.keys)

				return b, func() tea.Msg { return pushScreenMsg{screen: details} }
			}
		case ActionConnect:
			if r := b.table.Selected(); r.ID != "" {
				return b, func() tea.Msg {
					return connectRequestMsg{resource: b.res.Name, row: r}
				}
			}
		case ActionCopy:
			if r := b.table.Selected(); r.ID != "" {
				return b, func() tea.Msg { return copyMsg{value: r.ID, label: "id"} }
			}
		case ActionEdit:
			if r := b.table.Selected(); r.ID != "" {
				return b, func() tea.Msg { return editRequestMsg{res: b.res, id: r.ID} }
			}
		case ActionDelete:
			if r := b.table.Selected(); r.ID != "" {
				return b, func() tea.Msg { return deleteRequestMsg{res: b.res, id: r.ID} }
			}
		case ActionDrill:
			if b.res.Sub != nil {
				if r := b.table.Selected(); r.ID != "" {
					return b, b.drill(r.ID)
				}
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
// is harmless since the last fetchMsg wins. Cross-browser stale results
// are dropped by the `from` check in Update.
func (b *Browser) fetch() tea.Cmd {
	b.loading = true

	return func() tea.Msg {
		return b.doFetch()
	}
}

// drill fetches the subresource list of one row (k9s drill-down); the
// result arrives as subListMsg, which pushes the sub browser when the
// list is non-empty.
//
// ponytail: no in-flight guard, like fetch (the last subListMsg wins);
// the sub browser re-fetches on its own, so the double list call is
// accepted.
func (b *Browser) drill(id string) tea.Cmd {
	b.loading = true

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
		defer cancel()
		res, err := b.res.Sub.Scoped(id)(ctx)

		return subListMsg{from: b, id: id, res: res, err: err}
	}
}

// drillResource builds the scoped child Resource for the drill
// browser: the child resource's commands (edit/delete/connect), a
// fetch pinned to the parent row, and a title naming the parent.
func (b *Browser) drillResource(parentID string) *Resource {
	sub := b.res.Sub

	return &Resource{
		Name:      sub.Res.Name,
		Title:     sub.Res.Name + " · " + shortID(parentID),
		Fetch:     sub.Scoped(parentID),
		GetCmd:    sub.Res.GetCmd,
		UpdateCmd: sub.Res.UpdateCmd,
		DeleteCmd: sub.Res.DeleteCmd,
	}
}

// doFetch runs the fetch with a recover: 200+ heterogeneous commands
// mean one bad Run closure should degrade to an error line, not kill
// the TUI process.
func (b *Browser) doFetch() (m fetchMsg) {
	defer func() {
		if r := recover(); r != nil {
			m = fetchMsg{from: b, err: fmt.Errorf("fetch panicked: %v", r)}
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
	defer cancel()
	res, err := b.res.Fetch(ctx)
	if err != nil {
		return fetchMsg{from: b, err: err}
	}

	return fetchMsg{from: b, columns: res.columns, rows: res.rows}
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
