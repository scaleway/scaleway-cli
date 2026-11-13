package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func Test_Volume(t *testing.T) {
	t.Run("Create", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance volume create volume-type=l_ssd size=20GB tags.0=cli-test",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				volume := ctx.Result.(*instanceSDK.Volume)

				assert.Equal(t, instanceSDK.VolumeVolumeTypeLSSD, volume.VolumeType)
				assert.Equal(t, 20*scw.GB, volume.Size)
				assert.Len(t, volume.Tags, 1)
				assert.Equal(t, "cli-test", volume.Tags[0])
			},
		),
		AfterFunc: core.ExecAfterCmd("scw instance volume delete {{ .CmdResult.ID }}"),
	}))

	cmds := instanceV1AndV2Commands()

	t.Run("Create from snapshot", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createNonEmptyLocalVolume("Volume", 10),
			createSnapshot("Snap", "Volume"),
		),
		Cmd: "scw instance volume create base-snapshot={{ .Snap.ID }} volume-type=l_ssd size=25GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				volume := ctx.Result.(*instanceSDK.Volume)
				snapID := ctx.Meta.Render("{{.Snap.ID}}")

				assert.Equal(t, instanceSDK.VolumeVolumeTypeLSSD, volume.VolumeType)
				assert.Equal(t, 25*scw.GB, volume.Size)
				assert.NotNil(t, volume.BaseSnapshotID)
				assert.Equal(t, snapID, *volume.BaseSnapshotID)
			},
		),
		AfterFunc: core.ExecAfterCmd("scw instance volume delete {{ .CmdResult.ID }}"),
	}))

	t.Run("List", core.Test(&core.TestConfig{
		// TODO: display ok but test more stuff, create volumes to be listed, etc.
		Commands: instance.GetCommands(),
		Cmd:      "scw instance volume list",
		Check:    core.TestCheckGolden(),
	}))
}
