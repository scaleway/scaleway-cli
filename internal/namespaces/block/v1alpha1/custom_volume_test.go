package block_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	blockSDK "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/stretchr/testify/require"
)

func Test_VolumeWait(t *testing.T) {
	t.Run("Wait command", core.Test(&core.TestConfig{
		Commands: block.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Volume",
				"scw block volume create perf-iops=5000 from-empty.size=20GB",
			),
		),
		Cmd: "scw block volume wait {{ .Volume.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteVolume("Volume"),
	}))

	t.Run("Wait flag", core.Test(&core.TestConfig{
		Commands: block.GetCommands(),
		Cmd:      "scw block volume create perf-iops=5000 from-empty.size=20GB -w",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				vol := testhelpers.Value[*blockSDK.Volume](t, ctx.Result)
				require.Equal(t, blockSDK.VolumeStatusAvailable, vol.Status)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteVolume("CmdResult"),
	}))
}
