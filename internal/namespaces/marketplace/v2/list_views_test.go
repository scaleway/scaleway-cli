package marketplace_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/marketplace/v2"
)

func Test_ListCategory(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: marketplace.GetCommands(),
		Cmd:      "scw marketplace category list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: marketplace.GetCommands(),
		Cmd:      "scw marketplace version list",
		Check:    core.TestCheckGolden(),
	}))
}
