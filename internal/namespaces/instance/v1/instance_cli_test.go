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
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			return nil
		},
		Cmd: "scw instance server get server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }}")
			return nil
		},
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
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }}")
			return nil
		},
	}))

	t.Run("No initial placement group & placement-group-id=none", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} placement-group=none",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Nil(t, ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }}")
			return nil
		},
	}))

	t.Run(`No initial placement group & placement-group-id=`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group=`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Nil(t, ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
	}))

	t.Run(`No initial placement group & placement-group-id=<existing pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t,
					ctx.Meta["PlacementGroup"].(*instance.CreatePlacementGroupResponse).PlacementGroup.ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup.ID,
				)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
	}))

	t.Run(`No initial placement group & placement-group-id=<valid, but non existing pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id=11111111-1111-1111-1111-111111111111`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
	}))

	t.Run(`No initial placement group & placement-group-id=<invalid pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id=1111111`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
	}))

	t.Run(`Initial placement group & placement-group-id=none`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.PlacementGroup.ID }} stopped")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id=none`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Nil(t, ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
	}))

	t.Run(`Initial placement group & placement-group-id=<current pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.PlacementGroup.ID }} stopped")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t,
					ctx.Meta["PlacementGroup"].(*instance.CreatePlacementGroupResponse).PlacementGroup.ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup.ID,
				)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			return nil
		},
	}))

	t.Run(`Initial placement group & placement-group-id=<new pg id>`, core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroup"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["PlacementGroup2"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic placement-group-id={{ .PlacementGroup.PlacementGroup.ID }} stopped")
			return nil
		},
		Cmd: `scw instance server update server-id={{ .Server.ID }} placement-group-id={{ .PlacementGroup2.PlacementGroup.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.NoError(t, ctx.Err)
				assert.Equal(t,
					ctx.Meta["PlacementGroup2"].(*instance.CreatePlacementGroupResponse).PlacementGroup.ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup.ID,
				)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup.PlacementGroup.ID }}")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroup2.PlacementGroup.ID }}")
			return nil
		},
	}))
}

func Test_ImageCreate(t *testing.T) {
	t.Run("test", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance image create -D",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

//scw instance image create -D
