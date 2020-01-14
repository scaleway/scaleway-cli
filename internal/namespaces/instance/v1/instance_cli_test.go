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

	// TODO: Run test when `placement-group=none` is supported
	t.Skip("No initial placement group & placement-group-id=none", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.id }} placement-group=none",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }}")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

	t.Run("No initial placement group & placement-group-id=none", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.id }} placement-group=none",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }}")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.placement_group.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))
}
