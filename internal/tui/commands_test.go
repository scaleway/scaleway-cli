package tui

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type fakeServer struct {
	ID           string
	Name         string
	State        string
	Zone         scw.Zone
	CreationDate time.Time
	Tags         []string
}

func TestToRowsSlice(t *testing.T) {
	resp := []*fakeServer{{
		ID:           "12345678-aaaa-bbbb-cccc-000000000000",
		Name:         "web-1",
		State:        "running",
		Zone:         "fr-par-1",
		CreationDate: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		Tags:         []string{"a", "b"},
	}}
	lr, err := toRows(resp)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"ID", "NAME", "STATE", "ZONE", "CREATED"}; !slices.Equal(lr.columns, want) {
		t.Fatalf("columns = %v, want %v", lr.columns, want)
	}
	row := lr.rows[0]
	if row.ID != "12345678-aaaa-bbbb-cccc-000000000000" {
		t.Fatalf("row.ID = %q (want the full id, for copy)", row.ID)
	}
	wantCells := []string{"12345678", "web-1", "running", "fr-par-1", "2026-01-02 03:04"}
	if !slices.Equal(row.Cells, wantCells) {
		t.Fatalf("cells = %v, want %v", row.Cells, wantCells)
	}
	if len(row.CellStyles) != 1 {
		t.Fatalf("expected a STATE cell style, got %v", row.CellStyles)
	}
}

type fakeCluster struct {
	ID           string
	Name         string
	State        string
	Region       scw.Region
	CreationDate time.Time
}

func TestToRowsWrappedStruct(t *testing.T) {
	resp := &struct {
		Total    int
		Clusters []*fakeCluster
	}{Clusters: []*fakeCluster{{
		ID: "c1", Name: "k8s", State: "ready", Region: "fr-par",
		CreationDate: time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC),
	}}}
	lr, err := toRows(resp)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{
		"ID",
		"NAME",
		"STATE",
		"REGION",
		"CREATED",
	}; !slices.Equal(
		lr.columns,
		want,
	) {
		t.Fatalf("columns = %v, want %v", lr.columns, want)
	}
	if lr.rows[0].Cells[3] != "fr-par" {
		t.Fatalf("region cell = %q", lr.rows[0].Cells[3])
	}
}

func TestToRowsPlainSlice(t *testing.T) {
	lr, err := toRows([]string{"fr-par-1", "nl-ams-1"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"VALUE"}; !slices.Equal(lr.columns, want) {
		t.Fatalf("columns = %v, want %v", lr.columns, want)
	}
	if len(lr.rows) != 2 || lr.rows[0].Cells[0] != "fr-par-1" {
		t.Fatalf("rows = %v", lr.rows)
	}
}

func TestToRowsNoList(t *testing.T) {
	if _, err := toRows(42); err == nil {
		t.Fatal("expected an error for a non-list response")
	}
	if _, err := toRows("hello"); err == nil {
		t.Fatal("expected an error for a string response")
	}
}

type zoneArgs struct {
	Zone   scw.Zone
	Region scw.Region
}

func TestFillLocalities(t *testing.T) {
	args := &zoneArgs{}
	fillLocalities(args, "nl-ams-1", "")
	if args.Zone != "nl-ams-1" {
		t.Fatalf("Zone = %q", args.Zone)
	}
	if args.Region != "nl-ams" {
		t.Fatalf("Region = %q, want nl-ams", args.Region)
	}
	// An explicit region takes precedence over the zone-derived one.
	args3 := &zoneArgs{}
	fillLocalities(args3, "nl-ams-1", "fr-par")
	if args3.Region != "fr-par" {
		t.Fatalf("Region = %q, want the explicit fr-par", args3.Region)
	}
	// Explicit values are kept.
	args2 := &zoneArgs{Zone: "fr-par-2"}
	fillLocalities(args2, "nl-ams-1", "")
	if args2.Zone != "fr-par-2" {
		t.Fatalf("Zone = %q, want the explicit fr-par-2", args2.Zone)
	}
}

func TestResourcesFromCommands(t *testing.T) {
	var calls int
	first := &core.Command{
		Namespace: "instance", Resource: "server", Verb: "list",
		ArgsType: reflect.TypeFor[zoneArgs](),
		Run: func(ctx context.Context, argsI any) (any, error) {
			calls++

			return []*fakeServer{{ID: "abc", Name: "n"}}, nil
		},
	}
	// Same path: must be ignored (cobra builder: first command wins).
	second := &core.Command{
		Namespace: "instance", Resource: "server", Verb: "list",
		Run: func(ctx context.Context, argsI any) (any, error) {
			t.Fatal("duplicate command must not be registered")

			return nil, nil
		},
	}
	// Not a list command: skipped.
	other := &core.Command{
		Namespace: "tui", Resource: "server", Verb: "get",
		Run: func(ctx context.Context, argsI any) (any, error) { return nil, nil },
	}
	res := ResourcesFromCommands(
		[]*core.Command{first, second, other},
		nil,
		NewLocality("fr-par-1"),
	)
	if len(res) != 1 {
		t.Fatalf("resources = %d, want 1", len(res))
	}
	r, ok := res["instance server"]
	if !ok {
		t.Fatalf("missing 'instance server' in %v", res)
	}
	lr, err := r.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("fetch calls = %d, want 1", calls)
	}
	if lr.rows[0].ID != "abc" {
		t.Fatalf("rows = %v", lr.rows)
	}
}

func TestFuzzyMatch(t *testing.T) {
	if pos := fuzzyMatch("instance server", "srv"); len(pos) != 3 {
		t.Fatalf("fuzzyMatch = %v, want 3 positions", pos)
	}
	if fuzzyMatch("k8s cluster", "xyz") != nil {
		t.Fatal("fuzzyMatch should not match")
	}
	if fuzzyMatch("K8S Cluster", "k8s") == nil {
		t.Fatal("fuzzyMatch must be case-insensitive")
	}
}

func TestMenuFuzzyFilter(t *testing.T) {
	mk := func(name string) *Resource {
		return &Resource{
			Name: name, Title: name,
			Fetch: func(ctx context.Context) (listResult, error) { return listResult{}, nil },
		}
	}
	m := NewMenu(map[string]*Resource{
		"instance server": mk("instance server"),
		"k8s cluster":     mk("k8s cluster"),
		"lb pool":         mk("lb pool"),
	}, nil)
	m.Size(100, 10)
	for _, ch := range []string{"k", "8", "s", " ", "c", "l"} {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(ch)})
	}
	if m.query != "k8s cl" {
		t.Fatalf("query = %q, want 'k8s cl'", m.query)
	}
	vis := m.visible()
	if len(vis) != 1 || m.items[vis[0]].label != "k8s cluster" {
		t.Fatalf("visible = %v, want only 'k8s cluster'", vis)
	}
	// Enter opens the highlighted resource.
	cmd := m.selectVisible()
	gm, ok := cmd().(gotoResourceMsg)
	if !ok || gm.res.Name != "k8s cluster" {
		t.Fatalf("selectVisible = %v", cmd())
	}
	// Backspace + esc trim the query.
	m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if m.query != "k8s c" {
		t.Fatalf("query = %q, want 'k8s c'", m.query)
	}
	// esc clears the query; a second esc pops the menu.
	m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if m.query != "" {
		t.Fatalf("esc should clear the query, got %q", m.query)
	}
}

func TestMenuScroll(t *testing.T) {
	items := map[string]*Resource{}
	for i := range 240 {
		name := fmt.Sprintf("ns resource-%03d", i)
		items[name] = &Resource{Name: name, Title: name}
	}
	m := NewMenu(items, nil)
	m.Size(100, 10) // 9 visible rows
	m.move(20)
	if m.cursor != 20 || m.off != 12 {
		t.Fatalf("cursor = %d off = %d, want 20/12", m.cursor, m.off)
	}
	m.move(-100)
	if m.cursor != 0 || m.off != 0 {
		t.Fatalf("cursor = %d off = %d, want 0/0", m.cursor, m.off)
	}
}
