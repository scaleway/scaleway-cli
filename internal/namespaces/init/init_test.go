package init

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_Init(t *testing.T) {
	t.Run("Helper", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw init --help",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
