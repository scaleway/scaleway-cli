package billing_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
)

func Test_ListConsumption(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: billing.GetCommands(),
		Cmd:      "scw billing consumption list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDiscount(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: billing.GetCommands(),
		Cmd:      "scw billing discount list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListInvoice(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: billing.GetCommands(),
		Cmd:      "scw billing invoice list",
		Check:    core.TestCheckGolden(),
	}))
}
