package rdb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_ListEngineSettings(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw rdb engine settings name=MySQL version=8",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
