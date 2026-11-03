package account_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
)

func Test_ListProject(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: account.GetCommands(),
		Cmd:      "scw account project list",
		Check:    core.TestCheckGolden(),
	}))
}
