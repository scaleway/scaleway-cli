package tui

import (
	"context"
)

// listResult is the outcome of a Resource fetch: the table columns
// (derived from the response, so every resource gets sensible columns)
// and the rows.
type listResult struct {
	columns []string
	rows    []Row
}

// Resource describes a listable Scaleway resource: its name (the
// `namespace resource` command path), its header title and how to fetch
// its rows. It is the k9s model.Registry entry (GVR -> DAO + Renderer).
//
// Resources are generated from the CLI command registry (see
// ResourcesFromCommands), so every `scw <ns> <res> list` command is
// browsable in the TUI.
type Resource struct {
	Name  string
	Title string
	Fetch func(ctx context.Context) (listResult, error)
}

// resourceAliases maps short `:` prompt words to resource names.
var resourceAliases = map[string]string{
	"s": "instance server",
	"v": "instance volume",
	"i": "instance ip",
}

func shortID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}

	return id
}
