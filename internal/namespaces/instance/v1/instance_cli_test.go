package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

//
// Server
//

func Test_ListServer(t *testing.T) {

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

func Test_ListServerTypes(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server-type list -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:     GetCommands(),
		Cmd:          "scw instance server-type list",
		UseE2EClient: true,
		Check:        core.TestCheckGolden(),
	}))

}

func Test_GetServer(t *testing.T) {

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
			ctx.Meta["server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: "scw instance server get server-id={{ .server.id }}",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .server.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

}

//
// Volume
//

func Test_CreateVolume(t *testing.T) {

	deleteVolumeAfterFunc := func(ctx *core.AfterFuncCtx) error {
		// Get ID of the created volume.
		volumeID, err := ctx.ExtractResourceID()
		if err != nil {
			return err
		}

		// Delete the test volume.
		ctx.ExecuteCmd("scw instance volume delete volume-id=" + volumeID)
		return nil
	}

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:  GetCommands(),
		Cmd:       "scw instance volume create name=test size=20G",
		AfterFunc: deleteVolumeAfterFunc,
		Check:     core.TestCheckGolden(),
	}))

	t.Run("Bad size unit", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume create name=test size=20",
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
		Cmd: "scw -o json instance server update server-id={{ .Server.id }} placement-group=none",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual("{{.Result.server.placement_group}}", "<no value>"),
		),
	}))

	t.Run(`No initial placement group & placement-group-id=`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: `scw -o json instance server update server-id={{ .Server.id }} placement-group=`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }} delete-ip=true delete-volumes=true")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual("{{.Result.server.placement_group}}", "<no value>"),
		),
	}))

	t.Run(`No initial placement group & placement-group-id=<existing pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: `scw -o json instance server update server-id={{ .Server.id }} placement-group={{ .PlacementGroup.placement_group.id }}`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual("{{.Result.server.placement_group.id}}", "{{ .PlacementGroup.placement_group.id }}"),
		),
	}))

	t.Run(`No initial placement group & placement-group-id=<valid, but non existing pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.id }} placement-group=11111111-1111-1111-1111-111111111111`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }} delete-ip=true delete-volumes=true")
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
		Cmd: `scw instance server update server-id={{ .Server.id }} placement-group=1111111`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }} delete-ip=true delete-volumes=true")
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
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Cmd: `scw -o json instance server update server-id={{ .Server.id }} placement-group=none`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual("{{.Result.server.placement_group}}", "<no value>"),
		),
	}))

	t.Run(`Initial placement group & placement-group-id=<current pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Cmd: `scw -o json instance server update server-id={{ .Server.id }} placement-group={{ .PlacementGroup.placement_group.id }}`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual("{{.Result.server.placement_group.id}}", "{{.PlacementGroup.placement_group.id}}"),
		),
	}))

	t.Run(`Initial placement group & placement-group-id=<new pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["PlacementGroup2"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Cmd: `scw -o json instance server update server-id={{ .Server.id }} placement-group={{ .PlacementGroup2.placement_group.id }}`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckEqual("{{.Result.server.placement_group.id}}", "{{ .PlacementGroup2.placement_group.id }}"),
		),
	}))
}
