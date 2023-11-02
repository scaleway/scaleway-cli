package applesilicon

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_ServerTypeList(t *testing.T) {
	t.Run("base", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw apple-silicon server-type list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
