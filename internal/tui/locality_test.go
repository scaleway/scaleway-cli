package tui

import (
	"context"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// TestLocality verifies the session locality holder: an explicit region
// takes precedence, and changing the zone clears the explicit region.
func TestLocality(t *testing.T) {
	loc := NewLocality("fr-par-1")
	if z, r := loc.Get(); z != "fr-par-1" || r != "" {
		t.Fatalf("initial = (%q,%q), want (fr-par-1,\"\")", z, r)
	}
	loc.SetRegion("nl-ams")
	if z, r := loc.Get(); z != "fr-par-1" || r != "nl-ams" {
		t.Fatalf("after SetRegion = (%q,%q), want (fr-par-1,nl-ams)", z, r)
	}
	loc.SetZone("nl-ams-1")
	if z, r := loc.Get(); z != "nl-ams-1" || r != "" {
		t.Fatalf("after SetZone = (%q,%q), want (nl-ams-1,\"\")", z, r)
	}
	// A nil locality yields empty values (defensive).
	var nilLoc *Locality
	if z, r := nilLoc.Get(); z != "" || r != "" {
		t.Fatalf("nil locality = (%q,%q), want empty", z, r)
	}
}

// TestResourceFetchUsesLiveLocality is the core of the feature: a
// resource fetch must read the session locality at fetch time, so
// changing it mid-session changes what the next fetch targets.
func TestResourceFetchUsesLiveLocality(t *testing.T) {
	loc := NewLocality("fr-par-1")
	var (
		gotZone   scw.Zone
		gotRegion scw.Region
	)
	cmd := &core.Command{
		Namespace: "instance", Resource: "server", Verb: "list",
		ArgsType: reflect.TypeFor[zoneArgs](),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*zoneArgs)
			gotZone, gotRegion = args.Zone, args.Region

			return []string{"ok"}, nil
		},
	}
	res := ResourcesFromCommands([]*core.Command{cmd}, nil, loc)
	r := res["instance server"]

	fetch := func() {
		if _, err := r.Fetch(context.Background()); err != nil {
			t.Fatal(err)
		}
	}

	fetch()
	if gotZone != "fr-par-1" || gotRegion != "fr-par" {
		t.Fatalf("fetch1 zone=%q region=%q, want fr-par-1/fr-par", gotZone, gotRegion)
	}
	// Change the zone: the next fetch must target it (region re-derived).
	loc.SetZone("nl-ams-1")
	fetch()
	if gotZone != "nl-ams-1" || gotRegion != "nl-ams" {
		t.Fatalf("fetch2 zone=%q region=%q, want nl-ams-1/nl-ams", gotZone, gotRegion)
	}
	// An explicit region takes precedence over the zone-derived one.
	loc.SetRegion("pl-waw")
	fetch()
	if gotRegion != "pl-waw" {
		t.Fatalf("fetch3 region=%q, want pl-waw", gotRegion)
	}
}

// TestAppSetLocality verifies the app applies a selection: it updates
// the session locality, closes the picker, and re-fetches the screen
// underneath.
func TestAppSetLocality(t *testing.T) {
	res := &Resource{Name: "instance server", Title: "instance server"}
	loc := NewLocality("fr-par-1")
	a := &App{
		loc:       loc,
		resources: map[string]*Resource{"instance server": res},
		stack: []Screen{
			NewBrowser(res, nil),
			NewLocalityMenu(true, loc, nil),
		},
	}
	a.width, a.height = 100, 30
	a.sizeCurrent()

	a.Update(setLocalityMsg{isZone: true, zone: "nl-ams-2"})

	if loc.GetZone() != "nl-ams-2" {
		t.Fatalf("zone = %q, want nl-ams-2", loc.GetZone())
	}
	if len(a.stack) != 1 {
		t.Fatalf("stack len = %d, want 1 (picker closed)", len(a.stack))
	}
	if _, ok := a.current().(*Browser); !ok {
		t.Fatalf("current = %T, want Browser", a.current())
	}
	if a.flash != "zone set to nl-ams-2" {
		t.Fatalf("flash = %q", a.flash)
	}
}

// TestAppLocalityKeysOpenPickers verifies the default ctrl+z / ctrl+r
// bindings open the zone / region pickers.
func TestAppLocalityKeysOpenPickers(t *testing.T) {
	res := &Resource{Name: "instance server", Title: "instance server"}
	loc := NewLocality("fr-par-1")
	newApp := func() *App {
		a := &App{
			loc:       loc,
			resources: map[string]*Resource{"instance server": res},
			stack:     []Screen{NewBrowser(res, nil)},
		}
		a.width, a.height = 100, 30
		a.sizeCurrent()

		return a
	}

	a := newApp()
	a.Update(tea.KeyMsg{Type: tea.KeyCtrlZ})
	if m, ok := a.current().(*LocalityMenu); !ok || !m.isZone {
		t.Fatalf("ctrl+z: expected the zone picker, got %T", a.current())
	}

	a = newApp()
	a.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	if m, ok := a.current().(*LocalityMenu); !ok || m.isZone {
		t.Fatalf("ctrl+r: expected the region picker, got %T", a.current())
	}
}

// TestLocalityMenuSelect verifies the picker lists every zone/region,
// starts on the current one, and emits a setLocalityMsg for it.
func TestLocalityMenuSelect(t *testing.T) {
	loc := NewLocality("fr-par-1")

	m := NewLocalityMenu(true, loc, nil)
	m.Size(100, 20)
	if len(m.items) != len(scw.AllZones) {
		t.Fatalf("zone items = %d, want %d", len(m.items), len(scw.AllZones))
	}
	msg, ok := m.selectVisible()().(setLocalityMsg)
	if !ok || !msg.isZone || msg.zone != "fr-par-1" {
		t.Fatalf("zone select = %+v, want fr-par-1", msg)
	}

	mr := NewLocalityMenu(false, loc, nil)
	mr.Size(100, 20)
	if len(mr.items) != len(scw.AllRegions) {
		t.Fatalf("region items = %d, want %d", len(mr.items), len(scw.AllRegions))
	}
	rmsg, ok := mr.selectVisible()().(setLocalityMsg)
	if !ok || rmsg.isZone || rmsg.region == "" {
		t.Fatalf("region select = %+v, want a non-empty region", rmsg)
	}
}
