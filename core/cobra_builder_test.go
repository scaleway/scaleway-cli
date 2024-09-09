package core_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
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

func Test_DeprecatedCommand(t *testing.T) {
	dummyNamespaceCommand := &core.Command{
		Namespace: "instance",
	}

	dummyResourceCommand1 := &core.Command{
		Namespace:  "instance",
		Resource:   "server",
		Short:      "short server",
		Deprecated: false,
	}

	dummyResourceCommand2 := &core.Command{
		Namespace:  "instance",
		Resource:   "engine",
		Short:      "short engine",
		Deprecated: true,
	}

	dummyResourceCommand3 := &core.Command{
		Namespace:  "instance",
		Resource:   "a",
		Short:      "short server",
		Deprecated: false,
	}

	dummyResourceCommand4 := &core.Command{
		Namespace:  "instance",
		Resource:   "b",
		Short:      "short server",
		Deprecated: true,
	}

	cmds := core.NewCommands(
		dummyNamespaceCommand,
		dummyResourceCommand1,
		dummyResourceCommand2,
		dummyResourceCommand3,
		dummyResourceCommand4,
	)
	t.Run("", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw instance -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
