package marketplace

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_marketplaceImageList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw marketplace image list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
