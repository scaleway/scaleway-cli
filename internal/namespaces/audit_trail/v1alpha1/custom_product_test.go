package audit_trail_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	audit_trail "github.com/scaleway/scaleway-cli/v2/internal/namespaces/audit_trail/v1alpha1"
)

func Test_ProductList(t *testing.T) {
	t.Run("base", core.Test(&core.TestConfig{
		Commands: audit_trail.GetCommands(),
		Cmd:      "scw audit-trail product list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
