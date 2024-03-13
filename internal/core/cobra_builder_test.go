package core_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_UnknownCommand(t *testing.T) {
	dummyNamespaceCommand := &core.Command{
		Namespace: "instance",
	}

	dummyResourceCommand := &core.Command{
		Namespace: "instance",
		Resource:  "server",
	}

	cmds := core.NewCommands(
		dummyNamespaceCommand,
		dummyResourceCommand,
	)

	t.Run("UnknownResource", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw instance foobar",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("UnknownVerb", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw instance server foobar",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}
