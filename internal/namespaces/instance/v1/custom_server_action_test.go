package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	blockSDK "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests needs to be run in sequence
// since they are using the interactive print
func Test_ServerTerminate(t *testing.T) {
	interactive.IsInteractive = true

	t.Run("without IP", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("image=ubuntu-jammy ip=new -w"),
		),
		Cmd: `scw instance server terminate {{ .Server.ID }} with-block=true`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				api := instanceSDK.NewAPI(ctx.Client)
				server := testhelpers.MapValue[*instance.ServerWithWarningsResponse](
					t,
					ctx.Meta,
					"Server",
				).Server
				assert.NotNil(t, server.PublicIP)
				_, err := api.GetIP(&instanceSDK.GetIPRequest{
					IP: server.PublicIP.ID,
				})
				assert.NoError(t, err)
			},
		),
		AfterFunc: core.ExecAfterCmd(
			`scw instance ip delete {{ index .Server.PublicIP.ID }}`,
		),
		DisableParallel: true,
	}))

	t.Run("with IP", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("image=ubuntu-jammy ip=new -w"),
		),
		Cmd: `scw instance server terminate {{ .Server.ID }} with-ip=true with-block=true`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				api := instanceSDK.NewAPI(ctx.Client)
				server := testhelpers.MapValue[*instance.ServerWithWarningsResponse](
					t,
					ctx.Meta,
					"Server",
				).Server
				assert.NotNil(t, server.PublicIP)

				_, err := api.GetIP(&instanceSDK.GetIPRequest{
					IP: server.PublicIP.ID,
				})
				require.ErrorAs(t, err, new(*scw.ResourceNotFoundError))
			},
		),
		DisableParallel: true,
	}))

	t.Run("without block", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			instance.GetCommands(),
			block.GetCommands(),
		),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("image=ubuntu-jammy additional-volumes.0=block:10G -w"),
		),
		Cmd: `scw instance server terminate {{ .Server.ID }} with-ip=true with-block=false`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				`scw block volume wait terminal-status=available {{ (index .Server.Volumes "1").ID }}`,
			),
			core.ExecAfterCmd(`scw block volume delete {{ (index .Server.Volumes "1").ID }}`),
			core.ExecAfterCmd(
				`scw block volume wait terminal-status=available {{ (index .Server.Volumes "0").ID }}`,
			),
			core.ExecAfterCmd(`scw block volume delete {{ (index .Server.Volumes "0").ID }}`),
		),
		DisableParallel: true,
	}))

	t.Run("with block", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("image=ubuntu-jammy additional-volumes.0=block:10G -w"),
		),
		Cmd: `scw instance server terminate {{ .Server.ID }} with-ip=true with-block=true -w`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				api := blockSDK.NewAPI(ctx.Client)
				server := testhelpers.MapValue[*instance.ServerWithWarningsResponse](
					t,
					ctx.Meta,
					"Server",
				).Server
				rootVolume := testhelpers.MapTValue(t, server.Volumes, "0")

				_, err := api.GetVolume(&blockSDK.GetVolumeRequest{
					VolumeID: rootVolume.ID,
					Zone:     server.Zone,
				})
				require.ErrorAs(t, err, new(*scw.ResourceNotFoundError))

				additionalVolume := testhelpers.MapTValue(t, server.Volumes, "1")
				_, err = api.GetVolume(&blockSDK.GetVolumeRequest{
					VolumeID: additionalVolume.ID,
					Zone:     server.Zone,
				})
				require.ErrorAs(t, err, new(*scw.ResourceNotFoundError))
			},
		),
		DisableParallel: true,
	}))

	interactive.IsInteractive = false
}

// These tests needs to be run in sequence
// since they are using the interactive print
func Test_ServerBackup(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true image=ubuntu-jammy"),
		),
		Cmd: `scw instance server backup {{ .Server.ID }} name=backup`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance image delete {{ .CmdResult.Image.ID }} with-snapshots=true",
			),
			core.ExecAfterCmd(
				"scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all",
			),
		),
	}))

	t.Run("With SBS volumes", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand(
				"root-volume=sbs:20G additional-volumes.0=sbs:10G additional-volumes.1=sbs:15G stopped=true image=ubuntu-jammy",
			),
		),
		Cmd: `scw instance server backup {{ .Server.ID }} name=backup`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance image delete {{ .CmdResult.Image.ID }} with-snapshots=true",
			),
			core.ExecAfterCmd(
				"scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all",
			),
		),
	}))
}

func Test_ServerAction(t *testing.T) {
	t.Run("manual poweron", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true image=ubuntu_jammy"),
		),
		Cmd: `scw instance server action {{ .Server.ID }} action=poweron --wait`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				storedServer := testhelpers.MapValue[*instance.ServerWithWarningsResponse](
					t,
					ctx.Meta,
					"Server",
				).Server
				api := instanceSDK.NewAPI(ctx.Client)
				resp, err := api.GetServer(&instanceSDK.GetServerRequest{
					Zone:     storedServer.Zone,
					ServerID: storedServer.ID,
				})
				require.NoError(t, err)
				assert.Equal(t, instanceSDK.ServerStateRunning, resp.Server.State)
			},
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all force-shutdown=true",
			),
		),
	}))
}
