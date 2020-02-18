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

func Test_AllUsage(t *testing.T) {
	for _, command := range getCommands().GetAll() {
		commandLine := command.GetCommandLine()
		t.Run(commandLine+" usage", core.Test(&core.TestConfig{
			Commands: getCommands(),
			Cmd:      "scw " + commandLine + " -h",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
			),
		}))
	}
}
