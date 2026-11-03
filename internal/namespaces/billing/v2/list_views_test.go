package billing_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2"
)

func Test_ListBudget(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: billing.GetCommands(),
		Cmd:      "scw billing budget list",
		Check:    core.TestCheckGolden(),
	}))
}
