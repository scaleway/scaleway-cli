package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

// These tests needs to be run in sequence
// since they are using the interactive print
func Test_ServerTerminate(t *testing.T) {
	interactive.IsInteractive = true

	t.Run("without IP", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create image=ubuntu-bionic -w"),
		Cmd:        `scw instance server terminate {{ .Server.ID }}`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				api := instance.NewAPI(ctx.Client)
				server := ctx.Meta["Server"].(*instance.Server)
				_, err := api.GetIP(&instance.GetIPRequest{
					IP: server.PublicIP.ID,
				})
				assert.NoError(t, err)
			},
		),
		AfterFunc:       core.ExecAfterCmd(`scw instance ip delete {{ index .Server.PublicIP.ID }}`),
		DisableParallel: true,
	}))

	t.Run("with IP", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create image=ubuntu-bionic -w"),
		Cmd:        `scw instance server terminate {{ .Server.ID }} with-ip=true`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				api := instance.NewAPI(ctx.Client)
				server := ctx.Meta["Server"].(*instance.Server)
				_, err := api.GetIP(&instance.GetIPRequest{
					IP: server.PublicIP.ID,
				})
				require.IsType(t, &scw.ResponseError{}, err)
				assert.Equal(t, 403, err.(*scw.ResponseError).StatusCode)
			},
		),
		DisableParallel: true,
	}))

	t.Run("without block", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create image=ubuntu-bionic additional-volumes.0=block:10G -w"),
		Cmd:        `scw instance server terminate {{ .Server.ID }} with-ip=true with-block=false`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc:       core.ExecAfterCmd(`scw instance volume delete {{ (index .Server.Volumes "1").ID }}`),
		DisableParallel: true,
	}))

	t.Run("with block", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create image=ubuntu-bionic additional-volumes.0=block:10G -w"),
		Cmd:        `scw instance server terminate {{ .Server.ID }} with-ip=true with-block=true -w`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				api := instance.NewAPI(ctx.Client)
				server := ctx.Meta["Server"].(*instance.Server)
				_, err := api.GetVolume(&instance.GetVolumeRequest{
					VolumeID: server.Volumes["0"].ID,
					Zone:     server.Zone,
				})
				require.IsType(t, &scw.ResourceNotFoundError{}, err)
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
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic"),
		Cmd:        `scw instance server backup {{ .Server.ID }} name=backup`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw instance image delete {{ .CmdResult.Image.ID }} with-snapshots=true"),
			core.ExecAfterCmd("scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=local"),
		),
	}))
}

func Test_ServerAction(t *testing.T) {
	t.Run("manual poweron", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu_jammy"),
		Cmd:        `scw instance server action {{ .Server.ID }} action=poweron --wait`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				storedServer := ctx.Meta["Server"].(*instance.Server)
				api := instance.NewAPI(ctx.Client)
				resp, err := api.GetServer(&instance.GetServerRequest{
					Zone:     storedServer.Zone,
					ServerID: storedServer.ID,
				})
				assert.Nil(t, err)
				assert.Equal(t, instance.ServerStateRunning, resp.Server.State)
			},
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=local force-shutdown=true"),
		),
	}))
}

func Test_ServerEnableRoutedIP(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create zone=fr-par-3 type=PRO2-XXS image=ubuntu_jammy ip=new --wait"),
		Cmd:        `scw instance server enable-routed-ip zone=fr-par-3 {{ .Server.ID }} --wait`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				storedServer := ctx.Meta["Server"].(*instance.Server)
				api := instance.NewAPI(ctx.Client)
				server, err := api.GetServer(&instance.GetServerRequest{
					Zone:     storedServer.Zone,
					ServerID: storedServer.ID,
				})
				assert.Nil(t, err)
				assert.Equal(t, true, server.Server.RoutedIPEnabled)
				ip, err := api.GetIP(&instance.GetIPRequest{
					Zone: storedServer.Zone,
					IP:   storedServer.PublicIP.ID,
				})
				assert.Nil(t, err)
				assert.Equal(t, instance.IPTypeRoutedIPv4, ip.IP.Type)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw instance server delete zone=fr-par-3 {{ .Server.ID }} force-shutdown=true with-ip=true with-volumes=local"),
		),
	}))
}
