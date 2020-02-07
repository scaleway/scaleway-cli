package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func createVanillaServer(ctx *core.BeforeFuncCtx) error {
	ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create stopped=true image=ubuntu-bionic")
	return nil
}

func deleteVanillaServer(ctx *core.AfterFuncCtx) error {
	ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
	return nil
}

func Test_ServerUpdateCustom(t *testing.T) {

	////
	// IP use cases
	////
	t.Run("Try to remove ip from server without ip", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createVanillaServer,
		Cmd:        "scw instance server update server-id={{ .Server.ID }} ip=none",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, (*instance.ServerIP)(nil), ctx.Result.(*instance.UpdateServerResponse).Server.PublicIP)
			},
		),
		AfterFunc: deleteVanillaServer,
	}))

	t.Run("Update server ip from server without ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["CreateIPResponse"] = ctx.ExecuteCmd("scw instance ip create")
			return createVanillaServer(ctx)
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} ip={{ .CreateIPResponse.IP.Address }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, ctx.Meta["CreateIPResponse"].(*instance.CreateIPResponse).IP.Address, ctx.Result.(*instance.UpdateServerResponse).Server.PublicIP.Address)
			},
		),
		AfterFunc: deleteVanillaServer,
	}))

	t.Run("Update server ip from server with ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			createVanillaServer(ctx)
			ctx.Meta["CreateIPResponse"] = ctx.ExecuteCmd("scw instance ip create")
			ctx.Meta["ServerUpdated"] = ctx.ExecuteCmd("scw instance server update server-id={{ .Server.ID }} ip={{ .CreateIPResponse.IP.Address }}")
			ctx.Meta["CreateIPResponse2"] = ctx.ExecuteCmd("scw instance ip create")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} ip={{ .CreateIPResponse2.IP.Address }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t,
					ctx.Meta["CreateIPResponse"].(*instance.CreateIPResponse).IP.Address,
					ctx.Meta["ServerUpdated"].(*instance.UpdateServerResponse).Server.PublicIP.Address)
				assert.Equal(t,
					ctx.Meta["CreateIPResponse2"].(*instance.CreateIPResponse).IP.Address,
					ctx.Result.(*instance.UpdateServerResponse).Server.PublicIP.Address)
			},
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			deleteVanillaServer(ctx)
			ctx.ExecuteCmd("scw instance ip delete ip={{ .CreateIPResponse.IP.Address }}")
			return nil
		},
	}))

	////
	// Placement group use cases
	////
	t.Run("Update server placement-group-id from server with placement-group-id", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["PlacementGroupResponse"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["PlacementGroupResponse2"] = ctx.ExecuteCmd("scw instance placement-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create stopped=true image=ubuntu-bionic placement-group-id={{ .PlacementGroupResponse.PlacementGroup.ID }}")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} placement-group-id={{ .PlacementGroupResponse2.PlacementGroup.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t,
					ctx.Meta["PlacementGroupResponse2"].(*instance.CreatePlacementGroupResponse).PlacementGroup.ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup.ID)
			},
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroupResponse.PlacementGroup.ID }}")
			ctx.ExecuteCmd("scw instance placement-group delete placement-group-id={{ .PlacementGroupResponse2.PlacementGroup.ID }}")
			return nil
		},
	}))

	////
	// Volume use cases
	////
	t.Run("Volumes", func(t *testing.T) {
		t.Run("valid simple block volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["Response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=10G volume-type=b_ssd")
				ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create stopped=true image=ubuntu-bionic")
				return nil
			},
			Cmd: "scw instance server update server-id={{ .Server.ID }} additional-volume-ids.0={{ .Response.Volume.ID }}",
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 20*scw.GB, ctx.Result.(*instance.UpdateServerResponse).Server.Volumes["0"].Size)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.UpdateServerResponse).Server.Volumes["1"].Size)
			},
			AfterFunc: deleteVanillaServer,
		}))

		t.Run("valid double local volumes", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["Response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=10G volume-type=l_ssd")
				ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create stopped=true image=ubuntu-bionic root-volume=local:10GB additional-volumes.0=l:10G")
				return nil
			},
			Cmd: "scw instance server update server-id={{ .Server.ID }} additional-volume-ids.0={{ .Response.Volume.ID }}",
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.UpdateServerResponse).Server.Volumes["0"].Size)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.UpdateServerResponse).Server.Volumes["1"].Size)
			},
			AfterFunc: deleteVanillaServer,
		}))
	})
}
