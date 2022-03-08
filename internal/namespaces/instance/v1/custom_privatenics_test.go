package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ListNICs(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		// Temporary in waiting for the private network support in the CLI
		Cmd: "scw instance private-nic list server-id=4fe24c2a-3c65-4530-b274-574b22ba3d14",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))
}
