package block_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	blockSDK "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/stretchr/testify/require"
)

func Test_SnapshotWait(t *testing.T) {
	t.Run("Wait command", core.Test(&core.TestConfig{
		Commands: block.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Volume",
				"scw block volume create perf-iops=5000 from-empty.size=20GB -w",
			),
			core.ExecStoreBeforeCmd(
				"Snapshot",
				"scw block snapshot create volume-id={{ .Volume.ID }}",
			),
		),
		Cmd: "scw block snapshot wait {{ .Snapshot.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteSnapshot("Snapshot"),
			deleteVolume("Volume"),
		),
	}))

	t.Run("Wait flag", core.Test(&core.TestConfig{
		Commands: block.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Volume",
			"scw block volume create perf-iops=5000 from-empty.size=20GB -w",
		),
		Cmd: "scw block snapshot create volume-id={{ .Volume.ID }} -w",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				snap := testhelpers.Value[*blockSDK.Snapshot](t, ctx.Result)
				require.Equal(t, blockSDK.SnapshotStatusAvailable, snap.Status)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteSnapshot("CmdResult"),
			deleteVolume("Volume"),
		),
	}))
}
