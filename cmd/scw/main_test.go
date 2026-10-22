package main

import (
	"context"
	"os"
	"path/filepath"
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

func Test_DispatchPlugin(t *testing.T) {
	dir := t.TempDir()
	argsOut := filepath.Join(dir, "args")
	pluginPath := filepath.Join(dir, "scw-hello")
	pluginContent := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$SCW_TEST_OUT\"\nexit 3\n"
	if err := os.WriteFile(pluginPath, []byte(pluginContent), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("SCW_TEST_OUT", argsOut)

	exitCode, handled := dispatchPlugin([]string{"scw", "hello", "arg1", "--flag", "arg2"})
	if !handled {
		t.Fatal("expected scw-hello to be handled")
	}
	if exitCode != 3 {
		t.Fatalf("exit code = %d, want 3", exitCode)
	}
	data, err := os.ReadFile(argsOut)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(data)); got != "arg1\n--flag\narg2" {
		t.Fatalf("args passed to plugin = %q, want %q", got, "arg1\n--flag\narg2")
	}

	if _, handled := dispatchPlugin([]string{"scw", "scw-not-a-plugin-xyz"}); handled {
		t.Fatal("expected unknown scw-* command to not be handled")
	}
	if _, handled := dispatchPlugin([]string{"scw", "instance"}); handled {
		t.Fatal("expected non-plugin command to not be handled")
	}
	if _, handled := dispatchPlugin([]string{"scw"}); handled {
		t.Fatal("expected bare scw to not be handled")
	}
}
