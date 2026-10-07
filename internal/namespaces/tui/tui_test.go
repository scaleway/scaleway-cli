package tui_test

import (
	"context"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/internal/tui"
)

// TestAllListCommandsSupported guards that every `scw <ns> <res> list`
// command in the CLI registry is exposed as a browsable TUI resource,
// with a unique name and a working fetch.
func TestAllListCommandsSupported(t *testing.T) {
	cmds := commands.GetCommands(context.Background()).GetAll()

	// Every list command maps to exactly one resource name.
	want := map[string]bool{}
	for _, c := range cmds {
		if !c.IsList() || c.Run == nil {
			continue
		}
		name := c.Namespace
		if c.Resource != "" {
			name += " " + c.Resource
		}
		want[name] = true
	}

	resources := tui.ResourcesFromCommands(cmds, nil, tui.NewLocality(""))
	if len(resources) != len(want) {
		t.Fatalf("resources = %d, want %d (one per list command)", len(resources), len(want))
	}
	for name := range want {
		res, ok := resources[name]
		if !ok {
			t.Errorf("list command %q has no TUI resource", name)

			continue
		}
		if res.Name != name || res.Title == "" {
			t.Errorf("resource %q: bad name/title (%q/%q)", name, res.Name, res.Title)
		}
		if res.Fetch == nil {
			t.Errorf("resource %q: nil Fetch", name)
		}
	}
	if len(resources) < 100 {
		t.Errorf("only %d resources, expected full CLI coverage", len(resources))
	}
	// The resources of the initial view must exist.
	for _, name := range []string{"instance server", "instance volume", "instance ip"} {
		if _, ok := resources[name]; !ok {
			t.Errorf("missing expected resource %q", name)
		}
	}
}
