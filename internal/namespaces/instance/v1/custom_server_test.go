package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func Test_ServerUpdateCustom(t *testing.T) {
	t.Run("Try to remove ip from server without ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} ip=none",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, (*instance.ServerIP)(nil), ctx.Result.(*instance.UpdateServerResponse).Server.PublicIP)
			},
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
	}))

	t.Run("Update server ip from server without ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
			ctx.Meta["CreateIPResponse"] = ctx.ExecuteCmd("scw instance ip create")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} ip={{ .CreateIPResponse.IP.Address }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, ctx.Meta["CreateIPResponse"].(*instance.CreateIPResponse).IP.Address, ctx.Result.(*instance.UpdateServerResponse).Server.PublicIP.Address)
			},
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
	}))

	t.Run("Update server ip from server with ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic stopped")
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
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance ip delete ip={{ .CreateIPResponse.IP.Address }}")
			return nil
		},
	}))

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
			core.TestCheckExitCode(0),
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
}
