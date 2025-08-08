package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateSnapshot(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		t.Run("Change tags", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createNonEmptyLocalVolume("Volume", 10),
				core.ExecStoreBeforeCmd(
					"CreateSnapshot",
					"scw instance snapshot create volume-id={{ .Volume.ID }} name=cli-test-snapshot-update-tags tags.0=foo tags.1=bar",
				),
			),
			Cmd: "scw instance snapshot update -w {{ .CreateSnapshot.Snapshot.ID }} tags.0=bar tags.1=foo",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					snapshot := ctx.Result.(*instanceSDK.Snapshot)
					assert.Equal(t, "cli-test-snapshot-update-tags", snapshot.Name)
					assert.Len(t, snapshot.Tags, 2)
					assert.Equal(t, "bar", snapshot.Tags[0])
					assert.Equal(t, "foo", snapshot.Tags[1])
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteSnapshot("CreateSnapshot"),
				deleteVolume("Volume"),
			),
		}))
		t.Run("Change name", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createNonEmptyLocalVolume("Volume", 10),
				core.ExecStoreBeforeCmd(
					"CreateSnapshot",
					"scw instance snapshot create volume-id={{ .Volume.ID }} name=cli-test-snapshot-update-name tags.0=foo tags.1=bar",
				),
			),
			Cmd: "scw instance snapshot update -w {{ .CreateSnapshot.Snapshot.ID }} name=cli-test-snapshot-update-name-updated",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					snapshot := ctx.Result.(*instanceSDK.Snapshot)
					assert.Equal(t, "cli-test-snapshot-update-name-updated", snapshot.Name)
					assert.Len(t, snapshot.Tags, 2)
					assert.Equal(t, "foo", snapshot.Tags[0])
					assert.Equal(t, "bar", snapshot.Tags[1])
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteSnapshot("CreateSnapshot"),
				deleteVolume("Volume"),
			),
		}))
	})
}
