package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/stretchr/testify/assert"
)

func Test_Snapshot(t *testing.T) {
	cmds := instanceV1AndV2Commands()

	t.Run("Create", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: createNonEmptyLocalVolumeV1("Volume", 10),
		Cmd:        "scw instance snapshot create base-volume-id={{ .Volume.ID }} name=cli-test-snapshot-create tags.0=foo tags.1=bar",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.NotNil(t, ctx.Result)

				snapshot := ctx.Result.(*instanceSDK.Snapshot)
				assert.Equal(t, "cli-test-snapshot-create", snapshot.Name)
				assert.Len(t, snapshot.Tags, 2)
				assert.Equal(t, "foo", snapshot.Tags[0])
				assert.Equal(t, "bar", snapshot.Tags[1])
				assert.False(t, snapshot.Public)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteSnapshot("CmdResult"),
			deleteVolumeV1("Volume"),
		),
	}))

	t.Run("List", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createNonEmptyLocalVolumeV1("Vol1", 10),
			createNonEmptyLocalVolumeV1("Vol2", 15),
			createSnapshot("Snap1", "Vol1"),
			createSnapshot("Snap2", "Vol2"),
			createSnapshot("Snap3", "Vol2"),
		),
		Cmd: "scw instance snapshot list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.NotNil(t, ctx.Result)

				snapshots := ctx.Result.([]*instanceSDK.Snapshot)
				assert.Len(t, snapshots, 3)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteSnapshot("Snap1"),
			deleteSnapshot("Snap2"),
			deleteSnapshot("Snap3"),
			deleteVolumeV1("Vol1"),
			deleteVolumeV1("Vol2"),
		),
	}))

	t.Run("Update", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createNonEmptyLocalVolumeV1("Volume", 10),
			createSnapshot("Snapshot", "Volume"),
		),
		Cmd: "scw instance snapshot update {{ .Snapshot.ID }} name=cli-test-snapshot-update tags.0=foo tags.1=bar public=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.NotNil(t, ctx.Result)

				snapshot := ctx.Result.(*instanceSDK.Snapshot)
				assert.Equal(t, "cli-test-snapshot-update", snapshot.Name)
				assert.Len(t, snapshot.Tags, 2)
				assert.Equal(t, "foo", snapshot.Tags[0])
				assert.Equal(t, "bar", snapshot.Tags[1])
				assert.True(t, snapshot.Public)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteSnapshot("Snapshot"),
			deleteVolumeV1("Volume"),
		),
	}))
}
