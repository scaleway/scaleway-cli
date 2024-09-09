package marketplace_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/marketplace/v2"
)

func Test_marketplaceImageList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: marketplace.GetCommands(),
		Cmd:      "scw marketplace image list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
