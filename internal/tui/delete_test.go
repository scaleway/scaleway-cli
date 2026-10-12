package tui

import (
	"context"
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type fakeDeleteRequest struct {
	Zone     scw.Zone `json:"-"`
	ServerID string   `json:"-"`
}

func fakeDeleteCommand(called **fakeDeleteRequest) *core.Command {
	return &core.Command{
		Namespace: "instance", Resource: "server", Verb: "delete",
		ArgsType: reflect.TypeFor[fakeDeleteRequest](),
		ArgSpecs: core.ArgSpecs{
			{Name: "server-id", Short: "Server ID", Required: true, Positional: true},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			*called = argsI.(*fakeDeleteRequest)

			return nil, nil
		},
	}
}

func newDeleteTestApp(t *testing.T, res *Resource, confirmDelete bool) *App {
	t.Helper()
	client, err := scw.NewClient(
		scw.WithAuth("SCWAAAAAAAAAAAAAAAAA", "11111111-aaaa-bbbb-cccc-111111111111"),
		scw.WithDefaultOrganizationID("11111111-aaaa-bbbb-cccc-111111111111"),
	)
	if err != nil {
		t.Fatal(err)
	}
	a := NewApp(
		context.Background(), client, nil,
		map[string]*Resource{res.Name: res}, nil, nil, NewLocality("fr-par-1"),
	)
	a.confirmDelete = confirmDelete
	a.width, a.height = 100, 30
	a.sizeCurrent()

	return a
}

func TestResourcesFromCommandsWiresDelete(t *testing.T) {
	mk := func(ns, res, verb string) *core.Command {
		return &core.Command{
			Namespace: ns, Resource: res, Verb: verb,
			Run: func(ctx context.Context, argsI any) (any, error) { return nil, nil },
		}
	}
	list := mk("instance", "server", "list")
	del := fakeDeleteCommand(nil)
	listOnly := mk("lb", "pool", "list")

	res := ResourcesFromCommands([]*core.Command{list, del, listOnly}, nil, NewLocality("fr-par-1"))
	s := res["instance server"]
	if s == nil {
		t.Fatalf("missing 'instance server' in %v", res)
	}
	if s.DeleteCmd != del || !s.Deleteable() {
		t.Fatal("delete command not wired / not deleteable")
	}
	if p := res["lb pool"]; p == nil || p.Deleteable() {
		t.Fatal("'lb pool' (no delete) must not be deleteable")
	}
}

// TestAppDeleteConfirm walks the 'd' key flow with the confirmation
// modal on: request -> modal, n cancels, y runs the delete (with the
// id and session zone on the request) and flashes the outcome.
func TestAppDeleteConfirm(t *testing.T) {
	var got *fakeDeleteRequest
	del := fakeDeleteCommand(&got)
	res := &Resource{
		Name: "instance server", Title: "instance server", DeleteCmd: del,
	}
	a := newDeleteTestApp(t, res, true)

	// d on a selected row asks for confirmation (no call yet).
	_, cmd := a.Update(deleteRequestMsg{res: res, id: "abc"})
	if cmd != nil || a.confirm == nil {
		t.Fatalf("cmd = %v, confirm = %v, want the modal open and no call", cmd, a.confirm)
	}
	// The modal line is in the view.
	if view := a.View(); !strings.Contains(view, "delete instance server abc?") {
		t.Fatalf("missing the confirm line:\n%s", view)
	}
	// n cancels.
	a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if a.confirm != nil || got != nil {
		t.Fatalf("cancel: confirm = %v, got = %v, want closed and no call", a.confirm, got)
	}

	// Ask again, confirm with y: the delete runs.
	a.Update(deleteRequestMsg{res: res, id: "abc"})
	_, cmd = a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if cmd == nil || a.confirm != nil {
		t.Fatalf("cmd = %v, confirm = %v, want the delete running", cmd, a.confirm)
	}
	a.Update(cmd())
	if got == nil || got.ServerID != "abc" || got.Zone != "fr-par-1" {
		t.Fatalf("delete ran with %v, want id abc and zone fr-par-1", got)
	}
	if a.flash != "deleted instance server abc" {
		t.Fatalf("flash = %q", a.flash)
	}

	// A resource without a delete command flashes, does not call.
	plain := &Resource{Name: "lb pool", Title: "lb pool"}
	b := newDeleteTestApp(t, plain, true)
	b.Update(deleteRequestMsg{res: plain, id: "xyz"})
	if b.confirm != nil || b.flash != "delete not available for lb pool" {
		t.Fatalf("flash = %q, confirm = %v", b.flash, b.confirm)
	}
}

// TestAppDeleteNoConfirm verifies tui.delete_confirm: false deletes
// straight away, no modal.
func TestAppDeleteNoConfirm(t *testing.T) {
	var got *fakeDeleteRequest
	del := fakeDeleteCommand(&got)
	res := &Resource{
		Name: "instance server", Title: "instance server", DeleteCmd: del,
	}
	a := newDeleteTestApp(t, res, false)

	_, cmd := a.Update(deleteRequestMsg{res: res, id: "abc"})
	if cmd == nil || a.confirm != nil {
		t.Fatalf("cmd = %v, confirm = %v, want the delete running without a modal", cmd, a.confirm)
	}
	a.Update(cmd())
	if got == nil || got.ServerID != "abc" {
		t.Fatalf("delete ran with %v, want id abc", got)
	}
	if a.flash != "deleted instance server abc" {
		t.Fatalf("flash = %q", a.flash)
	}
}
