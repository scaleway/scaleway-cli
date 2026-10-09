package tui

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTableFilterAndCursor(t *testing.T) {
	tbl := &Table{
		Columns: []string{"ID", "NAME", "STATE"},
		Rows: []Row{
			{ID: "aaa", Cells: []string{"aaa", "web-1", "running"}},
			{ID: "bbb", Cells: []string{"bbb", "db-1", "stopped"}},
			{ID: "ccc", Cells: []string{"ccc", "web-2", "running"}},
		},
		Width:  80,
		Height: 5,
	}

	// Filter matches on cells.
	tbl.SetFilter("web")
	if got := len(tbl.Visible()); got != 2 {
		t.Fatalf("Visible() = %d, want 2", got)
	}

	// Cursor movement stays in bounds.
	tbl.SetFilter("")
	tbl.Cursor = 0
	if r := tbl.Selected(); r.ID != "aaa" {
		t.Fatalf("Selected() = %+v, want row aaa", r)
	}
	tbl.Cursor = 10
	tbl.clampScroll()
	tbl.Cursor = 2
	if r := tbl.Selected(); r.ID != "ccc" {
		t.Fatalf("Selected() = %+v, want row ccc", r)
	}
}

func TestTableScrollClamp(t *testing.T) {
	tbl := &Table{Columns: []string{"A"}, Width: 20, Height: 3}
	for range 10 {
		tbl.Rows = append(tbl.Rows, Row{Cells: []string{"x"}})
	}
	// rowsPerPage = Height-1 = 2.
	if got := tbl.rowsPerPage(); got != 2 {
		t.Fatalf("rowsPerPage() = %d, want 2", got)
	}
	tbl.Cursor = 5
	tbl.clampScroll()
	if tbl.Offset != 4 {
		t.Fatalf("Offset = %d, want 4", tbl.Offset)
	}
	tbl.Cursor = 1
	tbl.clampScroll()
	if tbl.Offset != 1 {
		t.Fatalf("Offset = %d, want 1", tbl.Offset)
	}
}

func TestColumnWidthsFit(t *testing.T) {
	tbl := &Table{
		Columns: []string{"ID", "NAME", "STATE"},
		Rows: []Row{
			{Cells: []string{"abcdefgh", "a-very-very-long-name-here", "running"}},
		},
		Width: 30,
	}
	widths := tbl.columnWidths(tbl.Visible())
	if got := sum(widths); got > 30-2 { // 3 cols -> 2 spaces of padding.
		t.Fatalf("widths sum = %d, want <= 28: %v", got, widths)
	}
	for _, w := range widths {
		if w < 1 {
			t.Fatalf("width too small: %v", widths)
		}
	}
}

func TestTruncate(t *testing.T) {
	if got := truncate("hello", 10); got != "hello" {
		t.Fatalf("truncate = %q", got)
	}
	if got := truncate("hello world", 6); got != "hello…" {
		t.Fatalf("truncate = %q", got)
	}
	if got := truncate("ab", 0); got != "" {
		t.Fatalf("truncate = %q", got)
	}
	if got := truncate("ab", 1); got != "a" {
		t.Fatalf("truncate = %q", got)
	}
}

func TestDetailsLines(t *testing.T) {
	d := NewDetails("test", struct {
		Name    string
		Count   int
		Tags    []string
		Zone    string
		Missing *string
	}{Name: "foo", Count: 2, Tags: []string{"a", "b"}}, nil)
	found := map[string]bool{}
	for _, l := range d.lines {
		if k, v, _ := cutTab(l); k != "" {
			found[k+"="+v] = true
		}
	}
	for _, want := range []string{"name=foo", "count=2", "tags=a, b"} {
		if !found[want] {
			t.Fatalf("missing line %q in %v", want, d.lines)
		}
	}
}

// TestDetailsViewTitle verifies the details view leads with the
// browsed object name (title line) and sizes the body below it.
func TestDetailsViewTitle(t *testing.T) {
	d := NewDetails("instance server:abcd1234", struct {
		Name string
	}{Name: "foo"}, nil)
	d.Size(80, 10)
	if d.height != 9 { // 10 - title line
		t.Fatalf("height = %d, want 9", d.height)
	}
	plain := ansiRe.ReplaceAllString(d.View(), "")
	if !strings.HasPrefix(plain, "instance server:abcd1234") {
		t.Fatalf("view should lead with the object title, got:\n%s", plain)
	}
}

// TestDetailsLongKeyFits verifies a long attribute name is shown in
// full when the view is wide enough, and truncated when it is not.
func TestDetailsLongKeyFits(t *testing.T) {
	longKey := "container_spec.container_image.reference.digest" // 46 chars
	d := NewDetails("t", struct {
		A int
	}{A: 1}, nil)
	d.lines = append(d.lines, longKey+"\tsome-value")

	d.Size(120, 10)
	plain := ansiRe.ReplaceAllString(d.View(), "")
	if !strings.Contains(plain, longKey) {
		t.Fatalf("long key should fit in full, got:\n%s", plain)
	}

	d.Size(40, 10)
	plain = ansiRe.ReplaceAllString(d.View(), "")
	if strings.Contains(plain, longKey) {
		t.Fatalf("long key must be truncated without space, got:\n%s", plain)
	}
}

// TestDetailsSearchAndCopy verifies the / attribute search: it selects
// the matching attribute and yields a copyMsg with its value.
func TestDetailsSearchAndCopy(t *testing.T) {
	d := NewDetails("test", struct {
		MacAddress string
		Name       string
	}{MacAddress: "de:00:00:3b:d7:02", Name: "foo"}, nil)
	d.Size(80, 10)

	cmd := d.SearchAndCopy("MAC")
	m, ok := cmd().(copyMsg)
	if !ok {
		t.Fatalf("expected copyMsg, got %T", cmd())
	}
	if m.value != "de:00:00:3b:d7:02" || m.label != "mac_address" {
		t.Fatalf("copyMsg = %+v", m)
	}
	if d.sel != 0 {
		t.Fatalf("sel = %d, want 0", d.sel)
	}

	if m2, ok := d.SearchAndCopy("does-not-exist")().(flashMsg); !ok ||
		m2.msg != "no attribute: does-not-exist" {
		t.Fatalf("expected no-match flashMsg, got %v", m2)
	}
}

// TestDetailsNextPrevMatch verifies less-style n/p navigation over the
// multiple matches of a `/` search: each stop selects the line and
// copies its value, and the cycle wraps around.
func TestDetailsNextPrevMatch(t *testing.T) {
	d := NewDetails("test", struct {
		IP        string
		PublicIP  string
		PrivateIP string
	}{IP: "1.2.3.4", PublicIP: "5.6.7.8", PrivateIP: "9.10.11.12"}, nil)
	d.Size(80, 10)

	// n/p before any search: a hint, not a copy.
	if m, ok := d.nextMatch(1)().(flashMsg); !ok || m.msg == "" {
		t.Fatalf("expected a hint flashMsg, got %v", d.nextMatch(1)())
	}

	// / finds all three "ip" attributes and copies the first.
	m, ok := d.SearchAndCopy("ip")().(copyMsg)
	if !ok || m.label != "ip" || m.value != "1.2.3.4" {
		t.Fatalf("first match = %+v, want ip=1.2.3.4", m)
	}
	if d.sel != 0 || len(d.matches) != 3 {
		t.Fatalf("sel = %d matches = %v", d.sel, d.matches)
	}

	// n walks forward: public_ip, private_ip, then wraps to ip.
	wantSeq := []struct {
		label, value string
		sel          int
	}{
		{"public_ip", "5.6.7.8", 1},
		{"private_ip", "9.10.11.12", 2},
		{"ip", "1.2.3.4", 0},
	}
	for i, want := range wantSeq {
		m, ok := d.nextMatch(1)().(copyMsg)
		if !ok || m.label != want.label || m.value != want.value {
			t.Fatalf("n #%d = %+v, want %s=%s", i+1, m, want.label, want.value)
		}
		if d.sel != want.sel {
			t.Fatalf("n #%d sel = %d, want %d", i+1, d.sel, want.sel)
		}
	}

	// p wraps the other way: ip -> private_ip.
	m, ok = d.nextMatch(-1)().(copyMsg)
	if !ok || m.label != "private_ip" || d.sel != 2 {
		t.Fatalf("p = %+v sel=%d, want private_ip at line 2", m, d.sel)
	}
}

// TestPageKeysScroll pipes the raw xterm Page Up/Down sequences that a
// real terminal sends (\x1b[5~ / \x1b[6~) through a real tea.Program,
// covering the full path: terminal sequence -> bubbletea key parser ->
// screen scroll state.
func TestPageKeysScroll(t *testing.T) {
	d := &Details{title: "t", width: 80, height: 10}
	for i := range 40 {
		d.lines = append(d.lines, fmt.Sprintf("attr_%02d\tvalue%02d", i, i))
	}
	app := &App{
		loc:       NewLocality("fr-par-1"),
		resources: map[string]*Resource{},
		stack:     []Screen{d},
		width:     100,
		height:    30,
	}
	// Page Down, Page Down, Page Up.
	input := "\x1b[6~\x1b[6~\x1b[5~"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p := tea.NewProgram(
		app,
		tea.WithContext(ctx),
		tea.WithInput(strings.NewReader(input)),
		tea.WithOutput(io.Discard),
		tea.WithoutRenderer(),
	)
	_, _ = p.Run()
	if d.sel != 10 || d.off != 10 {
		t.Fatalf("sel=%d off=%d, want 10/10", d.sel, d.off)
	}
}

func TestTablePageKeys(t *testing.T) {
	tbl := &Table{Columns: []string{"A"}, Width: 40, Height: 5}
	for range 20 {
		tbl.Rows = append(tbl.Rows, Row{Cells: []string{"x"}})
	}
	tbl.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	if tbl.Cursor != 4 { // rowsPerPage = Height-1 = 4
		t.Fatalf("cursor = %d, want 4", tbl.Cursor)
	}
	tbl.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	if tbl.Cursor != 0 {
		t.Fatalf("cursor = %d, want 0", tbl.Cursor)
	}
}

// TestAppGotoMenuNavigation verifies the full goto flow at the App
// level: m opens the menu, typing fuzzy-filters it, enter navigates to
// the matching resource browser.
func TestAppGotoMenuNavigation(t *testing.T) {
	mk := func(name string) *Resource {
		return &Resource{
			Name:  name,
			Title: name,
			Fetch: func(ctx context.Context) (listResult, error) {
				return listResult{columns: []string{"ID"}, rows: []Row{{ID: name}}}, nil
			},
		}
	}
	resources := map[string]*Resource{
		"instance server": mk("instance server"),
		"k8s cluster":     mk("k8s cluster"),
		"lb pool":         mk("lb pool"),
	}
	a := &App{
		loc:       NewLocality("fr-par-1"),
		resources: resources,
		stack:     []Screen{NewBrowser(resources["instance server"], nil)},
	}
	a.width, a.height = 100, 30
	a.sizeCurrent()
	typeKey := func(r rune) { a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}) }

	// Open the goto menu.
	typeKey('m')
	if _, ok := a.current().(*Menu); !ok {
		t.Fatalf("current screen = %T, want *Menu", a.current())
	}
	// Type "k8s cl" — every char feeds the fuzzy query.
	for _, r := range "k8s cl" {
		typeKey(r)
	}
	menu := a.current().(*Menu)
	if menu.query != "k8s cl" {
		t.Fatalf("menu query = %q, want 'k8s cl'", menu.query)
	}
	// Enter navigates: the menu yields a cmd producing gotoResourceMsg,
	// which the app turns into a new browser (bubbletea runs the cmd).
	_, cmd := a.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		a.Update(cmd())
	}
	b, ok := a.current().(*Browser)
	if !ok {
		t.Fatalf("current screen = %T, want *Browser", a.current())
	}
	if b.res.Name != "k8s cluster" {
		t.Fatalf("current resource = %s, want k8s cluster", b.res.Name)
	}
	if len(a.stack) != 1 {
		t.Fatalf("stack len = %d, want 1 (goto replaces the stack)", len(a.stack))
	}
}

// TestAppMenuCapturesShortcutKeys verifies that while the goto menu is
// open the global shortcuts (q, r, m, g, :, /) are typed into the fuzzy
// query instead of triggering their action, and that esc clears the
// query before going back.
func TestAppMenuCapturesShortcutKeys(t *testing.T) {
	mk := func(name string) *Resource {
		return &Resource{
			Name:  name,
			Title: name,
			Fetch: func(ctx context.Context) (listResult, error) {
				return listResult{columns: []string{"ID"}, rows: []Row{{ID: name}}}, nil
			},
		}
	}
	resources := map[string]*Resource{
		"instance server": mk("instance server"),
		"k8s cluster":     mk("k8s cluster"),
	}
	a := &App{
		loc:       NewLocality("fr-par-1"),
		resources: resources,
		stack:     []Screen{NewBrowser(resources["instance server"], nil)},
	}
	a.width, a.height = 100, 30
	a.sizeCurrent()
	typeKey := func(r rune) { a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}) }

	typeKey('m')
	menu := a.current().(*Menu)
	// Every "shortcut" char is typed, not interpreted.
	for _, r := range "rqmg/:" {
		typeKey(r)
	}
	if menu.query != "rqmg/:" {
		t.Fatalf("menu query = %q, want 'rqmg/:'", menu.query)
	}
	b := a.stack[0].(*Browser)
	if b.loading {
		t.Fatal("'r' must not refresh the browser while the menu is open")
	}
	if len(a.stack) != 2 {
		t.Fatalf("stack len = %d, want 2 ('m' must not push another menu)", len(a.stack))
	}

	// esc clears the query first (menu still open)...
	a.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if menu.query != "" {
		t.Fatalf("esc should clear the query, got %q", menu.query)
	}
	if _, ok := a.current().(*Menu); !ok {
		t.Fatalf("current screen = %T, want *Menu (esc clears before popping)", a.current())
	}
	// ...a second esc goes back to the browser.
	_, cmd := a.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if cmd != nil {
		a.Update(cmd())
	}
	if _, ok := a.current().(*Browser); !ok {
		t.Fatalf("current screen = %T, want *Browser", a.current())
	}
	if len(a.stack) != 1 {
		t.Fatalf("stack len = %d, want 1", len(a.stack))
	}
}

// TestBrowserViewTitle verifies the browser leads with the browsed
// resource name (title line) and sizes the table below it.
func TestBrowserViewTitle(t *testing.T) {
	res := &Resource{Name: "instance server", Title: "instance server"}
	b := NewBrowser(res, nil)
	b.Size(80, 10)
	if b.table.Height != 8 { // 10 - title line - status line
		t.Fatalf("table height = %d, want 8", b.table.Height)
	}
	plain := ansiRe.ReplaceAllString(b.View(), "")
	if !strings.HasPrefix(plain, "instance server") {
		t.Fatalf("view should lead with the resource title, got:\n%s", plain)
	}
}

func cutTab(s string) (string, string, bool) {
	for i := range len(s) {
		if s[i] == '\t' {
			return s[:i], s[i+1:], true
		}
	}

	return s, "", false
}

// TestFetchMsgRoutesToBrowser verifies the root App forwards a fetchMsg
// to the current screen (the k9s "listener -> screen" path). Without
// this the fetch result is dropped and the table stays "loading…".
func TestFetchMsgRoutesToBrowser(t *testing.T) {
	fake := &Resource{
		Name:  "server",
		Title: "instance/server",
		Fetch: func(ctx context.Context) (listResult, error) { return listResult{}, nil },
	}
	a := &App{
		loc:       NewLocality("fr-par-1"),
		project:   "proj",
		resources: map[string]*Resource{"server": fake},
		stack:     []Screen{NewBrowser(fake, nil)},
	}
	a.width, a.height = 80, 24
	a.sizeCurrent()

	// A fetch result arrives at the root and must reach the browser.
	a.Update(fetchMsg{
		columns: []string{"ID", "NAME"},
		rows:    []Row{{ID: "x", Cells: []string{"x", "foo"}}},
	})

	b := a.current().(*Browser)
	if got := len(b.Table().Rows); got != 1 {
		t.Fatalf("table rows = %d, want 1 (fetchMsg not routed)", got)
	}
	if got := len(b.Table().Columns); got != 2 {
		t.Fatalf("table columns = %d, want 2", got)
	}
	if b.loading {
		t.Fatal("browser still loading after fetchMsg")
	}
	if b.err != "" {
		t.Fatalf("unexpected error: %s", b.err)
	}
}

// TestPromptSwitchResource verifies a `:` command replaces the stack
// with the requested resource (k9s gotoResource), by exact name or
// short alias.
func TestPromptSwitchResource(t *testing.T) {
	mk := func(name string) *Resource {
		return &Resource{
			Name:  name,
			Title: "x/" + name,
			Fetch: func(ctx context.Context) (listResult, error) { return listResult{}, nil },
		}
	}
	server, ip := mk("instance server"), mk("instance ip")
	a := &App{
		loc: NewLocality("fr-par-1"),
		resources: map[string]*Resource{
			"instance server": server,
			"instance ip":     ip,
		},
		stack: []Screen{NewBrowser(server, nil)},
	}
	a.width, a.height = 80, 24
	a.sizeCurrent()

	// Short alias.
	a.prompt = newCmdPrompt()
	a.prompt.Buffer = "i"
	a.execPrompt()
	if got := a.current().(*Browser).res.Name; got != "instance ip" {
		t.Fatalf("current = %s, want instance ip", got)
	}

	// Exact name.
	a.prompt = newCmdPrompt()
	a.prompt.Buffer = "instance server"
	a.execPrompt()
	if got := a.current().(*Browser).res.Name; got != "instance server" {
		t.Fatalf("current = %s, want instance server", got)
	}

	// Unknown command flashes, keeps the current browser.
	a.prompt = newCmdPrompt()
	a.prompt.Buffer = "nope"
	a.execPrompt()
	if a.flash == "" {
		t.Fatal("expected a flash for an unknown command")
	}
	if got := len(a.stack); got != 1 {
		t.Fatalf("stack len = %d, want 1", got)
	}
}
