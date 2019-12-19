package e2e

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces/test/v1"
)

func TestTestCommand(t *testing.T) {

	t.Run("usage", core.Test(&core.TestConfig{
		Commands:     test.GetCommands(),
		Cmd:          "scw test -h",
		UseE2EClient: true,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanCreate(t *testing.T) {

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human create -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:     test.GetCommands(),
		Cmd:          "scw test human create",
		UseE2EClient: true,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("with args", core.Test(&core.TestConfig{
		Commands:     test.GetCommands(),
		Cmd:          "scw test human create height=170.5 shoe-size=35.1 altitude-in-meter=-12 altitude-in-millimeter=-12050 fingers-count=21 hair-count=9223372036854775808 is-happy=true eyes-color=amber organization-id=b3ba839a-dcf2-4b0a-ac81-fc32370052a0",
		UseE2EClient: true,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanList(t *testing.T) {

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human list -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.ExecuteCmd("scw test human create")
			ctx.ExecuteCmd("scw test human create")
			ctx.ExecuteCmd("scw test human create")
			return nil
		},
		Cmd:          "scw test human list",
		UseE2EClient: true,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

}

func TestHumanUpdate(t *testing.T) {

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human update -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("single arg", core.Test(&core.TestConfig{
		Commands:     test.GetCommands(),
		UseE2EClient: true,
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.ExecuteCmd("scw test human create height=170.5 shoe-size=35.1 altitude-in-meter=-12 altitude-in-millimeter=-12050 fingers-count=21 hair-count=9223372036854775808 is-happy=true eyes-color=amber organization-id=b3ba839a-dcf2-4b0a-ac81-fc32370052a0")
			return nil
		},
		Cmd: "scw test human update human-id=0194fdc2-fa2f-fcc0-41d3-ff12045b73c8 is-happy=false",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("multiple args", core.Test(&core.TestConfig{
		Commands:     test.GetCommands(),
		UseE2EClient: true,
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.ExecuteCmd("scw test human create")
			return nil
		},
		Cmd: "scw test human update human-id=0194fdc2-fa2f-fcc0-41d3-ff12045b73c8 height=155.666 shoe-size=36.0 altitude-in-meter=2147483647 altitude-in-millimeter=2147483647285 fingers-count=20 hair-count=9223372036854775809 is-happy=true eyes-color=blue",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

}

func TestHumanGet(t *testing.T) {

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human get -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:     test.GetCommands(),
		UseE2EClient: true,
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.ExecuteCmd("scw test human create height=155.666 shoe-size=36.0 altitude-in-meter=2147483647 altitude-in-millimeter=2147483647285 fingers-count=20 hair-count=9223372036854775809 is-happy=true eyes-color=blue organization-id=b3ba839a-dcf2-4b0a-ac81-fc32370052a0")
			return nil
		},
		Cmd: "scw test human get human-id=0194fdc2-fa2f-fcc0-41d3-ff12045b73c8",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
}

func TestHumanDelete(t *testing.T) {

	t.Run("usage", core.Test(&core.TestConfig{
		Commands: test.GetCommands(),
		Cmd:      "scw test human delete -h",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:     test.GetCommands(),
		UseE2EClient: true,
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.ExecuteCmd("scw test human create height=155.666 shoe-size=36.0 altitude-in-meter=2147483647 altitude-in-millimeter=2147483647285 fingers-count=20 hair-count=9223372036854775809 is-happy=true eyes-color=blue organization-id=b3ba839a-dcf2-4b0a-ac81-fc32370052a0")
			return nil
		},
		Cmd: "scw test human delete human-id=0194fdc2-fa2f-fcc0-41d3-ff12045b73c8",
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
