package tui

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
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

	// GetCmd, UpdateCmd and DeleteCmd (may be nil) are the `get`,
	// `update` and `delete` commands of the resource: they make it
	// editable with the 'e' key (see edit.go) and deletable with the
	// 'd' key (see delete.go).
	GetCmd    *core.Command
	UpdateCmd *core.Command
	DeleteCmd *core.Command

	// Sub, if non-nil, is the resource's subresource (k9s
	// drill-down): the browser shows a drill key that lists Sub.Res
	// for the selected row, when it is non-empty.
	Sub *SubResource
	// Parent is the name of the resource this one belongs to ("" if
	// top-level): subresources are hidden from the goto menu and the
	// `:` prompt — they only make sense for one parent row.
	Parent string
}

// SubResource is a child list of a resource, scoped to one row of the
// parent (k9s drill-down).
type SubResource struct {
	// Res is the child resource (e.g. "k8s pool").
	Res *Resource
	// Scoped returns the child list fetch pinned to one parent row.
	Scoped func(parentID string) func(context.Context) (listResult, error)
}

// Editable reports whether the resource can be edited with the 'e' key.
func (r *Resource) Editable() bool {
	return r.GetCmd != nil && r.UpdateCmd != nil && r.UpdateCmd.ArgsType != nil
}

// Deleteable reports whether the resource can be deleted with the 'd'
// key.
func (r *Resource) Deleteable() bool {
	return r.DeleteCmd != nil && r.DeleteCmd.ArgsType != nil
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
