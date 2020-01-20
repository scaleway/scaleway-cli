package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ServerList(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server list -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server list",
		Check:    core.TestCheckGolden(),
	}))

}

func Test_ServerReboot(t *testing.T) {

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server reboot",
		Check:    core.TestCheckGolden(),
	}))

}

func Test_ServerStandby(t *testing.T) {

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server standby",
		Check:    core.TestCheckGolden(),
	}))

}

func Test_ServerGet(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server get -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: "scw instance server get server-id={{ .Server.ID }}",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

}

func Test_ServerDelete(t *testing.T) {

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server delete",
		Check:    core.TestCheckGolden(),
	}))

}

func Test_ServerStart(t *testing.T) {

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server start",
		Check:    core.TestCheckGolden(),
	}))

}

func Test_ServerStop(t *testing.T) {

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server stop",
		Check:    core.TestCheckGolden(),
	}))

}

func Test_ServerUpdate(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server update -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

	t.Run("No initial placement group & placement-group-id=none", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} placement-group=none",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckNil(".Server.PlacementGroup"),
		),
	}))

	t.Run(`No initial placement group & placement-group-id=`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group=`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckNil(".Server.PlacementGroup"),
		),
	}))

	t.Run(`No initial placement group & placement-group-id=<existing pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group={{ .PlacementGroup.PlacementGroup.ID }}`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual(".PlacementGroup.PlacementGroup.ID", ".Server.PlacementGroup.ID"),
		),
	}))

	t.Run(`No initial placement group & placement-group-id=<valid, but non existing pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group=11111111-1111-1111-1111-111111111111`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

	t.Run(`No initial placement group & placement-group-id=<invalid pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group=1111111`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

	t.Run(`Initial placement group & placement-group-id=none`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group=none`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckNil(".Server.PlacementGroup"),
		),
	}))

	t.Run(`Initial placement group & placement-group-id=<current pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group={{ .PlacementGroup.PlacementGroup.ID }}`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual(".PlacementGroup.PlacementGroup.ID", ".Server.PlacementGroup.ID"),
		),
	}))

	t.Run(`Initial placement group & placement-group-id=<new pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["PlacementGroup2"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group={{ .PlacementGroup2.PlacementGroup.ID }}`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual(".PlacementGroup2.PlacementGroup.ID", ".Server.PlacementGroup.ID"),
		),
	}))
}
