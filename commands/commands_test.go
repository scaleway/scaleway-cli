package commands_test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
)

// listViewExceptions documents the list commands that must not have an
// explicit view, with the reason why. A view is only meaningful when the
// command's result is a slice of structs rendered by the default reflective
// human marshaler: views are ignored by custom slice marshalers and do not
// apply to non-slice results.
var listViewExceptions = map[string]string{
	"annotations key-value list":          "result is a struct (key-value map), not a slice",
	"edge-services route-rules list":      "result is a struct, not a slice",
	"environmental-footprint report list": "result is a struct, not a slice",
	"function runtime list":               "result is a struct, not a slice",
	"iam saml-certificates list":          "result is a struct, not a slice",
	"instance image list":                 "custom slice marshaler, views are ignored",
	"instance server list":                "custom slice marshaler, views are ignored",
	"jobs secret list":                    "result is a struct, not a slice",
	"lb private-network list":             "custom slice marshaler, views are ignored",
	"object bucket list":                  "custom slice marshaler (S3 types), views are ignored",
	"rdb log list":                        "result is a struct, not a slice",
	"registry docker-helper list":         "result is a slice of scalars, not structs",
	"tem offers list":                     "result is a struct, not a slice",
	"webhosting backup-item list":         "result is a struct, not a slice",
	"webhosting freedomain list":          "result is a slice of scalars, not structs",
	"webhosting progress list":            "result is a struct, not a slice",
}

// Test_ListCommandsHaveView ensures every list command in the merged command
// tree has an explicit view, so table columns stay stable when new fields are
// added to the API models (e.g. the SRN field). When this test fails, either
// add a view with core.ListViewBuilder(...) in the owning namespace package,
// or document the command in listViewExceptions.
func Test_ListCommandsHaveView(t *testing.T) {
	ctx := core.InjectMeta(context.Background(), &core.Meta{BetaMode: true})
	cmds := commands.GetCommands(ctx)

	var missing []string
	for _, cmd := range cmds.GetAll() {
		if cmd.Verb != "list" {
			continue
		}
		name := fmt.Sprintf("%s %s list", cmd.Namespace, cmd.Resource)
		if _, ok := listViewExceptions[name]; ok {
			continue
		}
		if cmd.View == nil {
			missing = append(missing, name)
		}
	}

	sort.Strings(missing)
	require.Empty(t, missing, "list commands without an explicit view; add a view or a documented exception: %v", missing)
}
