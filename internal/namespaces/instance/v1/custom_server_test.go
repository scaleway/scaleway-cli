package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	blockCli "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func Test_ServerVolumeUpdate(t *testing.T) {
	t.Run("Attach", func(t *testing.T) {
		t.Run("simple block volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createServerBionic("Server"),
				createVolume("Volume", 10, instance.VolumeVolumeTypeBSSD),
			),
			Cmd: "scw instance server attach-volume server-id={{ .Server.ID }} volume-id={{ .Volume.ID }}",
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 20*scw.GB, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["0"].Size)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["1"].Size)
				assert.Equal(t, instance.VolumeServerVolumeTypeBSSD, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["1"].VolumeType)
			},
			AfterFunc:       deleteServer("Server"),
			DisableParallel: true,
		}))

		t.Run("simple local volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createServerBionic("Server"),
				createVolume("Volume", 10, instance.VolumeVolumeTypeLSSD),
			),
			Cmd: "scw instance server attach-volume server-id={{ .Server.ID }} volume-id={{ .Volume.ID }}",
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 20*scw.GB, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["0"].Size)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["1"].Size)
				assert.Equal(t, instance.VolumeServerVolumeTypeLSSD, ctx.Result.(*instance.AttachVolumeResponse).Server.Volumes["1"].VolumeType)
			},
			AfterFunc:       deleteServer("Server"),
			DisableParallel: true,
		}))

		t.Run("invalid volume UUID", core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: createServerBionic("Server"),
			Cmd:        "scw instance server attach-volume server-id={{ .Server.ID }} volume-id=11111111-1111-1111-1111-111111111111",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(1),
			),
			AfterFunc:       deleteServer("Server"),
			DisableParallel: true,
		}))
	})
	t.Run("Detach", func(t *testing.T) {
		t.Run("simple block volume", core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic additional-volumes.0=block:10G"),
			Cmd:        `scw instance server detach-volume volume-id={{ (index .Server.Volumes "1").ID }}`,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.NotZero(t, ctx.Result.(*instance.DetachVolumeResponse).Server.Volumes["0"])
				assert.Nil(t, ctx.Result.(*instance.DetachVolumeResponse).Server.Volumes["1"])
				assert.Equal(t, 1, len(ctx.Result.(*instance.DetachVolumeResponse).Server.Volumes))
			},
			AfterFunc: core.AfterFuncCombine(
				core.ExecAfterCmd(`scw instance volume delete {{ (index .Server.Volumes "1").ID }}`),
				deleteServer("Server"),
			),
			DisableParallel: true,
		}))

		t.Run("invalid volume UUID", core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: createServerBionic("Server"),
			Cmd:        "scw instance server detach-volume volume-id=11111111-1111-1111-1111-111111111111",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(1),
			),
			AfterFunc:       deleteServer("Server"),
			DisableParallel: true,
		}))
	})
}

func Test_ServerUpdateCustom(t *testing.T) {
	// IP cases.
	t.Run("Try to remove ip from server without ip", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServerBionic("Server"),
		Cmd:        "scw instance server update {{ .Server.ID }} ip=none",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, (*instance.ServerIP)(nil), ctx.Result.(*instance.UpdateServerResponse).Server.PublicIP)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run("Update server ip from server without ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServerBionic("Server"),
			createIP("IP"),
		),
		Cmd: "scw instance server update {{ .Server.ID }} ip={{ .IP.Address }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, ctx.Meta["IP"].(*instance.IP).Address, ctx.Result.(*instance.UpdateServerResponse).Server.PublicIP.Address)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run("Update server ip from server with ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServerBionic("Server"),
			createIP("IP1"),
			createIP("IP2"),

			// Attach IP1 to Server.
			core.ExecStoreBeforeCmd("UpdatedServer", "scw instance server update {{ .Server.ID }} ip={{ .IP1.Address }}"),
		),
		Cmd: "scw instance server update {{ .Server.ID }} ip={{ .IP2.Address }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				// Test that the Server WAS attached to IP1.
				assert.Equal(t,
					ctx.Meta["IP1"].(*instance.IP).Address,
					ctx.Meta["UpdatedServer"].(*instance.UpdateServerResponse).Server.PublicIP.Address)
				// Test that the Server IS attached to IP2.
				assert.Equal(t,
					ctx.Meta["IP2"].(*instance.IP).Address,
					ctx.Result.(*instance.UpdateServerResponse).Server.PublicIP.Address)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteIP("IP1"),
		),
	}))

	// Placement group cases.
	t.Run("Update server placement-group-id from server with placement-group-id", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createPlacementGroup("PlacementGroup1"),
			createPlacementGroup("PlacementGroup2"),
			core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic placement-group-id={{ .PlacementGroup1.ID }}"),
		),
		Cmd: "scw instance server update {{ .Server.ID }} placement-group-id={{ .PlacementGroup2.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t,
					ctx.Meta["PlacementGroup2"].(*instance.PlacementGroup).ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.PlacementGroup.ID)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deletePlacementGroup("PlacementGroup1"),
			deletePlacementGroup("PlacementGroup2"),
		),
	}))

	// Security group cases.
	t.Run("Update server security-group-id from server with security-group-id", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createSecurityGroup("SecurityGroup1"),
			createSecurityGroup("SecurityGroup2"),
			core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic security-group-id={{ .SecurityGroup1.ID }}"),
		),
		Cmd: "scw instance server update {{ .Server.ID }} security-group-id={{ .SecurityGroup2.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t,
					ctx.Meta["SecurityGroup2"].(*instance.SecurityGroup).ID,
					ctx.Result.(*instance.UpdateServerResponse).Server.SecurityGroup.ID)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteSecurityGroup("SecurityGroup1"),
			deleteSecurityGroup("SecurityGroup2"),
		),
	}))

	// Volumes cases.
	t.Run("Volumes", func(t *testing.T) {
		t.Run("valid simple block volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createServerBionic("Server"),
				createVolume("Volume", 10, instance.VolumeVolumeTypeBSSD),
			),
			Cmd: `scw instance server update {{ .Server.ID }} volume-ids.0={{ (index .Server.Volumes "0").ID }} volume-ids.1={{ .Volume.ID }}`,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 20*scw.GB, ctx.Result.(*instance.UpdateServerResponse).Server.Volumes["0"].Size)
				assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.UpdateServerResponse).Server.Volumes["1"].Size)
			},
			AfterFunc: deleteServer("Server"),
		}))

		t.Run("detach all volumes", core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic additional-volumes.0=block:10G"),
			Cmd:        `scw instance server update {{ .Server.ID }} volume-ids=none`,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				require.NoError(t, ctx.Err)
				assert.Equal(t, 0, len(ctx.Result.(*instance.UpdateServerResponse).Server.Volumes))
			},
			AfterFunc: core.AfterFuncCombine(
				core.ExecAfterCmd(`scw instance volume delete {{ (index .Server.Volumes "0").ID }}`),
				core.ExecAfterCmd(`scw instance volume delete {{ (index .Server.Volumes "1").ID }}`),
				deleteServer("Server"),
			),
		}))
	})
}

// These tests needs to be run in sequence
// since they are using the interactive print
func Test_ServerDelete(t *testing.T) {
	interactive.IsInteractive = true

	t.Run("with all volumes", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic additional-volumes.0=block:10G"),
		Cmd:        `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		DisableParallel: true,
	}))

	t.Run("only block volumes", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic additional-volumes.0=block:10G"),
		Cmd:        `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=block`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc:       core.ExecAfterCmd(`scw instance volume delete {{ (index .Server.Volumes "0").ID }}`),
		DisableParallel: true,
	}))

	t.Run("only local volumes", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic additional-volumes.0=block:10G"),
		Cmd:        `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=local`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc:       core.ExecAfterCmd(`scw instance volume delete {{ (index .Server.Volumes "1").ID }}`),
		DisableParallel: true,
	}))

	t.Run("with none volumes", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic additional-volumes.0=block:10G"),
		Cmd:        `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=none`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				api := instance.NewAPI(ctx.Client)
				server := ctx.Meta["Server"].(*instance.Server)
				_, err := api.GetVolume(&instance.GetVolumeRequest{
					VolumeID: server.Volumes["0"].ID,
				})
				assert.NoError(t, err)
			},
		),
		AfterFunc:       core.ExecAfterCmd(`scw instance volume delete {{ (index .Server.Volumes "0").ID }}`),
		DisableParallel: true,
	}))

	t.Run("with sbs volumes", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			GetCommands(),
			blockCli.GetCommands(),
		),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("BlockVolume", "scw block volume create perf-iops=5000 from-empty.size=10G name=cli-test-server-delete-with-sbs-volumes"),
			core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-jammy"),
			core.ExecBeforeCmd("scw instance server attach-volume server-id={{ .Server.ID }} volume-id={{ .BlockVolume.ID }}"),
		),
		Cmd: `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				api := block.NewAPI(ctx.Client)
				blockVolume := ctx.Meta["BlockVolume"].(*block.Volume)
				resp, err := api.GetVolume(&block.GetVolumeRequest{
					Zone:     blockVolume.Zone,
					VolumeID: blockVolume.ID,
				})
				assert.Error(t, err, "%v", resp)
			},
		),
		DisableParallel: true,
	}))

	interactive.IsInteractive = false
}
