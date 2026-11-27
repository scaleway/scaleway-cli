package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/stretchr/testify/assert"
)

func Test_Server(t *testing.T) {
	t.Run("Start Spot", core.Test(&core.TestConfig{
		Commands: instanceV1AndV2Commands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true type=POP2-2C-8G zone=pl-waw-3"),
		),
		Cmd: "scw instance server start-spot {{ .Server.ID }} zone=pl-waw-3 -w",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				server := ctx.Result.(*instance.Server)
				assert.NotNil(t, server.RuntimeInfo)
				assert.True(t, server.RuntimeInfo.Spot)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw instance server stop {{ .Server.ID }} zone=pl-waw-3 -w"),
			deleteServerV1("Server"),
		),
	}))
}
