package audit_trail_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	audit_trail "github.com/scaleway/scaleway-cli/v2/internal/namespaces/audit_trail/v1alpha1"
)

func Test_ListProducts(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: audit_trail.GetCommands(),
		Cmd:      "scw audit-trail product list",
		Check:    core.TestCheckGolden(),
	}))
}
