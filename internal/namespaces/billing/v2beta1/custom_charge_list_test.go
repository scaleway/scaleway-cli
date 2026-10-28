package billing_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
)

func Test_ChargeList(t *testing.T) {
	t.Run("list all pages", core.Test(&core.TestConfig{
		Commands: billing.GetCommands(),
		Cmd:      "scw billing charge list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}
