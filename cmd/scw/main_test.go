package main

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_MainUsage(t *testing.T) {
	t.Run("usage", core.Test(&core.TestConfig{
		Commands: getCommands(),
		Cmd:      "scw -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}
