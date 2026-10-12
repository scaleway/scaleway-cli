package tui

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// TestAppConnectFlow verifies the connect wiring: c on the browser sends a
// connect request, the app resolves the registered connector, builds the
// command in a cmd, and hands it to tea.ExecProcess; the hint is shown.
func TestAppConnectFlow(t *testing.T) {
	res := &Resource{
		Name:  "fake box",
		Title: "fake box",
		Fetch: func(ctx context.Context) (listResult, error) {
			return listResult{columns: []string{"ID"}, rows: []Row{{ID: "abc"}}}, nil
		},
	}
	a := &App{
		loc:       NewLocality("fr-par-1"),
		resources: map[string]*Resource{"fake box": res},
		stack:     []Screen{NewBrowser(res, nil)},
	}
	a.width, a.height = 80, 24
	a.sizeCurrent()
	b := a.current().(*Browser)
	b.table.Columns = []string{"ID"}
	b.table.SetRows([]Row{{ID: "abc"}})

	gotRowID := ""
	RegisterConnector("fake box", Connector{
		Hint: "ssh",
		Connect: func(ctx context.Context, client *scw.Client, loc *Locality, row Row) (*exec.Cmd, func(), error) {
			gotRowID = row.ID

			return exec.Command("true"), nil, nil
		},
	})

	// The header hint shows the connector.
	var hintFound bool
	for _, h := range b.Hints() {
		if h == "[ c ] ssh" {
			hintFound = true
		}
	}
	if !hintFound {
		t.Fatalf("hints = %v, want a '[ c ] ssh' hint", b.Hints())
	}

	// c sends the connect request.
	_, cmd := a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if cmd == nil {
		t.Fatal("c must return the connect request cmd")
	}
	// The app resolves the connector and builds the command (async).
	req := cmd().(connectRequestMsg)
	_, cmd = a.Update(req)
	if cmd == nil {
		t.Fatal("connect request must return the build cmd")
	}
	// The build cmd calls the connector with the selected row.
	ready := cmd().(connectReadyMsg)
	if ready.err != nil {
		t.Fatalf("unexpected build error: %s", ready.err)
	}
	if gotRowID != "abc" {
		t.Fatalf("connector row id = %q, want abc", gotRowID)
	}
	// The ready msg is handed to tea.ExecProcess (not run here: it needs
	// a live program and terminal).
	_, cmd = a.Update(ready)
	if cmd == nil {
		t.Fatal("connect ready must return the ExecProcess cmd")
	}
}

// TestAppConnectErrors verifies the error paths: a failing connector and a
// resource without one both flash instead of opening anything.
func TestAppConnectErrors(t *testing.T) {
	res := &Resource{
		Name:  "fake db",
		Title: "fake db",
		Fetch: func(ctx context.Context) (listResult, error) {
			return listResult{columns: []string{"ID"}, rows: []Row{{ID: "def"}}}, nil
		},
	}
	a := &App{
		loc:       NewLocality("fr-par-1"),
		resources: map[string]*Resource{"fake db": res},
		stack:     []Screen{NewBrowser(res, nil)},
	}
	a.width, a.height = 80, 24
	a.sizeCurrent()
	b := a.current().(*Browser)
	b.table.Columns = []string{"ID"}
	b.table.SetRows([]Row{{ID: "def"}})

	// No connector registered for this resource: flash, no cmd.
	_, cmd := a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	req := cmd().(connectRequestMsg)
	a.Update(req)
	if a.flash != "no connector for fake db" {
		t.Fatalf("flash = %q, want 'no connector for fake db'", a.flash)
	}

	// A failing connector flashes its error.
	RegisterConnector("fake db", Connector{
		Hint: "psql",
		Connect: func(ctx context.Context, client *scw.Client, loc *Locality, row Row) (*exec.Cmd, func(), error) {
			return nil, nil, errors.New("boom")
		},
	})
	a.flash = ""
	_, cmd = a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	_, buildCmd := a.Update(cmd().(connectRequestMsg))
	if buildCmd == nil {
		t.Fatal("failing connector must still return the build cmd")
	}
	a.Update(buildCmd().(connectReadyMsg))
	if !strings.Contains(a.flash, "boom") {
		t.Fatalf("flash = %q, want the connector error", a.flash)
	}
}
