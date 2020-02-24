package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	if !core.UpdateCassettes {
		instance.RetryInterval = 0
	}
}

//
// Server
//
func Test_ListServer(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server list -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_ListServerTypes(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server-type list -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:     GetCommands(),
		Cmd:          "scw instance server-type list",
		UseE2EClient: true,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_GetServer(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server get -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        "scw instance server get server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteServer("Server"),
	}))
}

//
// Volume
//
func Test_CreateVolume(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume create name=test size=20G",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "test", ctx.Result.(*instance.CreateVolumeResponse).Volume.Name)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance volume delete volume-id=" + ctx.CmdResult.(*instance.CreateVolumeResponse).Volume.ID)
			return nil
		},
	}))

	t.Run("Bad size unit", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume create name=test size=20",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_ServerUpdate(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server update -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        "scw instance server update server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run("No initial placement group & placement-group-id=none", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        "scw instance server update server-id={{ .Server.ID }} placement-group=none",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Nil(t, ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run(`No initial placement group & placement-group-id=`, core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        `scw instance server update server-id={{ .Server.ID }} placement-group=`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Nil(t, ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run(`No initial placement group & placement-group-id=<existing pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createPlacementGroup("PlacementGroup"),
			createServer("Server"),
		),
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id={{ .PlacementGroup.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t,
					ctx.Meta["PlacementGroup"].(*instance.PlacementGroup).ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup.ID,
				)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deletePlacementGroup("PlacementGroup"),
		),
	}))

	t.Run(`No initial placement group & placement-group-id=<valid, but non existing pg id>`, core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        `scw instance server update server-id={{ .Server.ID }} placement-group-id=11111111-1111-1111-1111-111111111111`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run(`No initial placement group & placement-group-id=<invalid pg id>`, core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        `scw instance server update server-id={{ .Server.ID }} placement-group-id=1111111`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run(`Initial placement group & placement-group-id=none`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createPlacementGroup("PlacementGroup"),
			core.ExecStoreBeforeCmd("Server", "scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.ID }} stopped"),
		),
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id=none`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Nil(t, ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deletePlacementGroup("PlacementGroup"),
		),
	}))

	t.Run(`Initial placement group & placement-group-id=<current pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createPlacementGroup("PlacementGroup"),
			core.ExecStoreBeforeCmd("Server", "scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.ID }} stopped"),
		),
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id={{ .PlacementGroup.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t,
					ctx.Meta["PlacementGroup"].(*instance.PlacementGroup).ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup.ID,
				)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deletePlacementGroup("PlacementGroup"),
		),
	}))

	t.Run(`Initial placement group & placement-group-id=<new pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createPlacementGroup("PlacementGroup1"),
			createPlacementGroup("PlacementGroup2"),
			core.ExecStoreBeforeCmd("Server", "scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup1.ID }} stopped"),
		),
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id={{ .PlacementGroup2.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.NoError(t, ctx.Err)
				assert.Equal(t,
					ctx.Meta["PlacementGroup2"].(*instance.PlacementGroup).ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup.ID,
				)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deletePlacementGroup("PlacementGroup1"),
			deletePlacementGroup("PlacementGroup2"),
		),
	}))
}

//
// Snapshot
//
func Test_SnapshotCreate(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        `scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			func(ctx *core.AfterFuncCtx) error {
				ctx.ExecuteCmd("scw instance snapshot delete snapshot-id=" + ctx.CmdResult.(*instance.CreateSnapshotResponse).Snapshot.ID)
				return nil
			},
			deleteServer("Server"),
		),
	}))
}
