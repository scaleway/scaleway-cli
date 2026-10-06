package main

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
)

func Test_MainUsage(t *testing.T) {
	t.Run("usage", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(context.Background()),
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

	for _, cmd := range commands.GetCommands(context.Background()).GetAll() {
		commandLine := cmd.GetCommandLine("scw")
		commandLine = strings.TrimPrefix(commandLine, "scw ")
		if _, exists := excludedCommands[commandLine]; exists || cmd.Hidden {
			continue
		}

		t.Run(commandLine+" usage", core.Test(&core.TestConfig{
			Commands: commands.GetCommands(context.Background()),
			Cmd:      "scw " + commandLine + " -h",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGoldenAndReplacePatterns(core.GoldenReplacement{
					Pattern:       regexp.MustCompile("[0-9a-f-]{36}"),
					Replacement:   "11111111-1111-1111-1111-111111111111",
					OptionalMatch: true,
				}),
			),
		}))
	}
}
