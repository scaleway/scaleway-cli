package core_test

import (
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func TestCobraCompletion(t *testing.T) {
	run := func(words []string, expected ...string) func(*testing.T) {
		return core.Test(&core.TestConfig{
			Commands: testAutocompleteGetCommands(),
			Args:     append([]string{"scw", "__complete"}, words...),
			Check: core.TestCheckCombine(
				core.TestCheckStdout(strings.Join(expected, "\n")+"\n"),
				core.TestCheckExitCode(0),
			),
		})
	}

	t.Run("subcommands", run(
		[]string{"test", "flower", ""},
		"create", "delete", ":4",
	))
	t.Run("arg names", run(
		[]string{"test", "flower", "create", "s"},
		"size=", "species=", ":6",
	))
	t.Run("completed args are not suggested again", run(
		[]string{"test", "flower", "create", "species=rose", "s"},
		"size=", ":6",
	))
	t.Run("enum values", run(
		[]string{"test", "flower", "create", "species=v"},
		"species=violet", "species=virginia bluebell", ":4",
	))
	t.Run("positional and arg names", run(
		[]string{"test", "flower", "delete", ""},
		"anemone", "hibiscus", "with-leaves=", ":6",
	))
	t.Run("flags are skipped", run(
		[]string{"test", "flower", "create", "-p", "default", "species=v"},
		"species=violet", "species=virginia bluebell", ":4",
	))
	t.Run("output flag value", run(
		[]string{"test", "flower", "create", "-o", ""},
		"human", "json", "yaml", "template", ":4",
	))
}
