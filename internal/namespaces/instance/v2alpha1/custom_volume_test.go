package instance_test

import (
	"fmt"
	"strings"
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

	t.Run("Create bad size unit", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance volume create name=test volume-type=l_ssd size=20",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Update", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Volume",
				"scw instance volume create volume-type=l_ssd size=20GB tags.0=test",
			),
		),
		Cmd: "scw instance volume update {{ .Volume.ID }} size=30GB name=cli-test-instance-volume-update " +
			"tags.0=cli-test tags.1=instance-volume tags.2=update",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				volume := ctx.Result.(*instanceSDK.Volume)

				assert.Equal(t, instanceSDK.VolumeVolumeTypeLSSD, volume.VolumeType)
				assert.Equal(t, 30*scw.GB, volume.Size)
				assert.Equal(t, "cli-test-instance-volume-update", volume.Name)
				assert.Len(t, volume.Tags, 3)
				assert.Equal(t, "cli-test", volume.Tags[0])
				assert.Equal(t, "instance-volume", volume.Tags[1])
				assert.Equal(t, "update", volume.Tags[2])
			},
		),
		AfterFunc: core.ExecAfterCmd("scw instance volume delete {{ .Volume.ID }}"),
	}))

	t.Run("List", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Vol1",
				"scw instance volume create volume-type=l_ssd size=25GB name=list-test-vol1",
			),
			core.ExecStoreBeforeCmd(
				"Vol2",
				"scw instance volume create volume-type=l_ssd size=20GB tags.0=cli-test name=list-test-vol2",
			),
			core.ExecStoreBeforeCmd(
				"Vol3",
				"scw instance volume create volume-type=scratch size=200GB name=list-test-vol3",
			),
		),
		Commands: instance.GetCommands(),
		Cmd:      "scw instance volume list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				volumes := ctx.Result.([]*instanceSDK.Volume)
				volumesFound := map[int]bool{1: false, 2: false, 3: false}

				for _, volume := range volumes {
					switch volume.Name {
					case "list-test-vol1":
						assert.Equal(t, instanceSDK.VolumeVolumeTypeLSSD, volume.VolumeType)
						assert.Equal(t, 25*scw.GB, volume.Size)
						volumesFound[1] = true

					case "list-test-vol2":
						assert.Equal(t, instanceSDK.VolumeVolumeTypeLSSD, volume.VolumeType)
						assert.Equal(t, 20*scw.GB, volume.Size)
						assert.Len(t, volume.Tags, 1)
						assert.Equal(t, "cli-test", volume.Tags[0])
						volumesFound[2] = true

					case "list-test-vol3":
						assert.Equal(t, instanceSDK.VolumeVolumeTypeScratch, volume.VolumeType)
						assert.Equal(t, 200*scw.GB, volume.Size)
						volumesFound[3] = true
					}
				}

				errorMsg := strings.Builder{}
				printMsg := false

				errorMsg.WriteString("Error: could not find some of the expected volumes:\n")

				for index, found := range volumesFound {
					if !found {
						errorMsg.WriteString(
							fmt.Sprintf("- \"list-test-vol%d\" was missing\n", index),
						)
						printMsg = true
					}
				}

				if printMsg {
					t.Error(errorMsg.String())
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteVolume("Vol1"),
			deleteVolume("Vol2"),
			deleteVolume("Vol3"),
		),
	}))
}
