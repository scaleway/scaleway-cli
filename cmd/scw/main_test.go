package main

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/command"
	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_MainUsage(t *testing.T) {
	t.Run("usage", core.Test(&core.TestConfig{
		Commands: command.GetCommands(),
		Cmd:      "scw -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}

func Test_AllUsage(t *testing.T) {
	// The help for these commands can not be tested because it depends on the environment
	excludedCommands := map[string]bool{
		"init":                true,
		"config":              true,
		"autocomplete script": true,
	}

	for _, cmd := range command.GetCommands().GetAll() {
		commandLine := cmd.GetCommandLine()
		if _, exists := excludedCommands[commandLine]; exists {
			continue
		}

		t.Run(commandLine+" usage", core.Test(&core.TestConfig{
			Commands: command.GetCommands(),
			Cmd:      "scw " + commandLine + " -h",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
			),
		}))
	}
}
