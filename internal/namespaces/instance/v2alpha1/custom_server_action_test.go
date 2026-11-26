package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ServerAction(t *testing.T) {
	cmds := instanceV1AndV2Commands()

	t.Run("Start", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true image=ubuntu_jammy"),
		),
		Cmd: `scw instance server start {{ .Server.ID }} --wait`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				server := ctx.Result.(*instanceSDK.Server)

				assert.Equal(t, instanceSDK.ServerStatusStarted, server.Status)
			},
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all force-shutdown=true",
			),
		),
	}))

	t.Run("Stop and delete", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("image=ubuntu_jammy --wait"),
		),
		Cmd: `scw instance server stop-and-delete {{ .Server.ID }}`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				server := ctx.Result.(*instanceSDK.Server)
				assert.Equal(t, instanceSDK.ServerStatusStopping, server.Status)

				_, err := instanceSDK.NewAPI(ctx.Client).
					WaitForServer(&instanceSDK.WaitForServerRequest{
						Zone:          server.Zone,
						ServerID:      server.ID,
						RetryInterval: core.DefaultRetryInterval,
					}, scw.WithContext(t.Context()))

				require.Error(t, err)
				assert.ErrorAs(t, err, new(*scw.ResourceNotFoundError))
			},
		),
	}))
}
