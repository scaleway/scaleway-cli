package core

import (
	"testing"
)

func Test_UnknownCommand(t *testing.T) {
	dummyNamespaceCommand := &Command{
		Namespace: "instance",
	}

	dummyResourceCommand := &Command{
		Namespace: "instance",
		Resource:  "server",
	}

	cmds := NewCommands(
		dummyNamespaceCommand,
		dummyResourceCommand,
	)

	t.Run("UnknownResource", Test(&TestConfig{
		Commands: cmds,
		Cmd:      "scw instance foobar",
		Check: TestCheckCombine(
			TestCheckGolden(),
			TestCheckExitCode(1),
		),
	}))

	t.Run("UnknownVerb", Test(&TestConfig{
		Commands: cmds,
		Cmd:      "scw instance server foobar",
		Check: TestCheckCombine(
			TestCheckGolden(),
			TestCheckExitCode(1),
		),
	}))
}
