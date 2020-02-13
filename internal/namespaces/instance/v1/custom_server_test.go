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

func Test_ServerVolumeUpdate(t *testing.T) {
	t.Run("Attach", func(t *testing.T) {
		t.Run("help", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server attach-volume -h",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
			),
		}))

		t.Run("simple block volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["Response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=10G volume-type=b_ssd")
				return createVanillaServer(ctx)
			},
			Cmd: "scw instance server attach-volume server-id={{ .Server.ID }} volume-id={{ .Response.Volume.ID }}",
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 20*scw.GB, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["0"].Size)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["1"].Size)
				assert.Equal(t, instance.VolumeTypeBSSD, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["1"].VolumeType)
			},
			AfterFunc: deleteVanillaServer,
		}))

		t.Run("simple local volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["Response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=10G volume-type=l_ssd")
				return createVanillaServer(ctx)
			},
			Cmd: "scw instance server attach-volume server-id={{ .Server.ID }} volume-id={{ .Response.Volume.ID }}",
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 20*scw.GB, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["0"].Size)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["1"].Size)
				assert.Equal(t, instance.VolumeTypeLSSD, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["1"].VolumeType)
			},
			AfterFunc: deleteVanillaServer,
		}))

		t.Run("invalid volume UUID", core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: createVanillaServer,
			Cmd:        "scw instance server attach-volume server-id={{ .Server.ID }} volume-id=11111111-1111-1111-1111-111111111111",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(1),
			),
			AfterFunc: deleteVanillaServer,
		}))
	})
	t.Run("Detach", func(t *testing.T) {
		t.Run("help", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server detach-volume -h",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
			),
		}))

		t.Run("simple block volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create stopped=true image=ubuntu-bionic additional-volumes.0=block:10G")
				return nil
			},
			Cmd: `scw instance server detach-volume volume-id={{ (index .Server.Volumes "1").ID }}`,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.NotZero(t, ctx.Result.(*instance.DetachVolumeResponse).Server.Volumes["0"])
				assert.Nil(t, ctx.Result.(*instance.DetachVolumeResponse).Server.Volumes["1"])
				assert.Equal(t, 1, len(ctx.Result.(*instance.DetachVolumeResponse).Server.Volumes))
			},
			AfterFunc: func(ctx *core.AfterFuncCtx) error {
				ctx.ExecuteCmd(`scw instance volume delete volume-id={{ (index .Server.Volumes "1").ID }}`)
				return deleteVanillaServer(ctx)
			},
		}))

		t.Run("invalid volume UUID", core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: createVanillaServer,
			Cmd:        "scw instance server detach-volume volume-id=11111111-1111-1111-1111-111111111111",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(1),
			),
			AfterFunc: deleteVanillaServer,
		}))
	})
}

func Test_ServerUpdateCustom(t *testing.T) {

	// IP cases.
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
			core.TestCheckExitCode(0),
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
			core.TestCheckExitCode(0),
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
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance ip delete ip={{ .CreateIPResponse.IP.Address }}")
			return nil
		},
	}))

	// Placement group cases.
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

	// Security group cases.
	t.Run("Update server security-group-id from server with security-group-id", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["SecurityGroupResponse"] = ctx.ExecuteCmd("scw instance security-group create")
			ctx.Meta["SecurityGroupResponse2"] = ctx.ExecuteCmd("scw instance security-group create")
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create stopped=true image=ubuntu-bionic security-group-id={{ .SecurityGroupResponse.SecurityGroup.ID }}")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} security-group-id={{ .SecurityGroupResponse2.SecurityGroup.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t,
					ctx.Meta["SecurityGroupResponse2"].(*instance.CreateSecurityGroupResponse).SecurityGroup.ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.SecurityGroup.ID)
			},
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			ctx.ExecuteCmd("scw instance security-group delete security-group-id={{ .SecurityGroupResponse.SecurityGroup.ID }}")
			ctx.ExecuteCmd("scw instance security-group delete security-group-id={{ .SecurityGroupResponse2.SecurityGroup.ID }}")
			return nil
		},
	}))

	// Volumes cases.
	t.Run("Volumes", func(t *testing.T) {
		t.Run("valid simple block volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["Response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=10G volume-type=b_ssd")
				return createVanillaServer(ctx)
			},
			Cmd: `scw instance server update server-id={{ .Server.ID }} volume-ids.0={{ (index .Server.Volumes "0").ID }} volume-ids.1={{ .Response.Volume.ID }}`,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 20*scw.GB, ctx.Result.(*instance.UpdateServerResponse).Server.Volumes["0"].Size)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.UpdateServerResponse).Server.Volumes["1"].Size)
			},
			AfterFunc: deleteVanillaServer,
		}))
	})
}
