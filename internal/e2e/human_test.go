package e2e_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/test/v1"
)

func TestTestCommand(t *testing.T) {
	t.Skip("Test API not available")

	t.Run("usage", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		Cmd:             "scw test -h",
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanCreate(t *testing.T) {
	t.Skip("Test API not available")

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human create -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test human create",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("with args", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test human create height=170.5 shoe-size=35.1 altitude-in-meter=-12 altitude-in-millimeter=-12050 fingers-count=21 hair-count=9223372036854775808 is-happy=true eyes-color=amber organization-id=b3ba839a-dcf2-4b0a-ac81-fc32370052a0",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("invalid boolean", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test human create is-happy=so-so",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanList(t *testing.T) {
	t.Skip("Test API not available")

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human list -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw test human create"),
			core.ExecBeforeCmd("scw test human create"),
			core.ExecBeforeCmd("scw test human create"),
		),
		Cmd: "scw test human list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanUpdate(t *testing.T) {
	t.Skip("Test API not available")

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human update -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("single arg", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		BeforeFunc: core.ExecBeforeCmd(
			"scw test human create height=170.5 shoe-size=35.1 altitude-in-meter=-12 altitude-in-millimeter=-12050 fingers-count=21 hair-count=9223372036854775808 is-happy=true eyes-color=amber organization-id=b3ba839a-dcf2-4b0a-ac81-fc32370052a0",
		),
		Cmd: "scw test human update 0194fdc2-fa2f-fcc0-41d3-ff12045b73c8 is-happy=false",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("multiple args", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		BeforeFunc:      core.ExecBeforeCmd("scw test human create"),
		Cmd:             "scw test human update 0194fdc2-fa2f-fcc0-41d3-ff12045b73c8 height=155.666 shoe-size=36.0 altitude-in-meter=2147483647 altitude-in-millimeter=2147483647285 fingers-count=20 hair-count=9223372036854775809 is-happy=true eyes-color=blue",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanGet(t *testing.T) {
	t.Skip("Test API not available")

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human get -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		BeforeFunc: core.ExecBeforeCmd(
			"scw test human create height=155.666 shoe-size=36.0 altitude-in-meter=2147483647 altitude-in-millimeter=2147483647285 fingers-count=20 hair-count=9223372036854775809 is-happy=true eyes-color=blue organization-id=b3ba839a-dcf2-4b0a-ac81-fc32370052a0",
		),
		Cmd: "scw test human get 0194fdc2-fa2f-fcc0-41d3-ff12045b73c8",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("pass human UUID without arg", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test human get 11111111-1111-1111-1111-111111111111",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))

	t.Run("invalid arg name", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test human get 11111111-1111-1111-1111-111111111111 invalid=true",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanDelete(t *testing.T) {
	t.Skip("Test API not available")

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human delete -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		BeforeFunc: core.ExecBeforeCmd(
			"scw test human create height=155.666 shoe-size=36.0 altitude-in-meter=2147483647 altitude-in-millimeter=2147483647285 fingers-count=20 hair-count=9223372036854775809 is-happy=true eyes-color=blue organization-id=b3ba839a-dcf2-4b0a-ac81-fc32370052a0",
		),
		Cmd: "scw test human delete 0194fdc2-fa2f-fcc0-41d3-ff12045b73c8",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanRun(t *testing.T) {
	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human run -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}
