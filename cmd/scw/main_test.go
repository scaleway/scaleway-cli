package main

import (
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces"
)

func Test_MainUsage(t *testing.T) {
	t.Run("usage", core.Test(&core.TestConfig{
		Commands: namespaces.GetCommands(),
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

	for _, cmd := range namespaces.GetCommands().GetAll() {
		commandLine := cmd.GetCommandLine("scw")
		commandLine = strings.TrimPrefix(commandLine, "scw ")
		if _, exists := excludedCommands[commandLine]; exists || cmd.Hidden {
			continue
		}

		t.Run(commandLine+" usage", core.Test(&core.TestConfig{
			Commands: namespaces.GetCommands(),
			Cmd:      "scw " + commandLine + " -h",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
			),
		}))
	}
}
