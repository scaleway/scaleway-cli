package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func Test_ImageCreate(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Server",
				testServerCommand("stopped=true image=ubuntu-jammy root-volume=l:20G"),
			),
			core.ExecStoreBeforeCmd(
				"Snapshot",
				`scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`,
			),
		),
		Commands: instance.GetCommands(),
		Cmd:      "scw instance image create snapshot-id={{ .Snapshot.Snapshot.ID }} arch=x86_64",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			core.ExecAfterCmd("scw instance image delete {{ .CmdResult.Image.ID }}"),
			deleteSnapshot("Snapshot"),
		),
	}))

	t.Run("Use additional snapshots", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Server",
				"scw instance server create type=DEV1-S ip=none image=ubuntu_focal root-volume=local:10GB additional-volumes.0=local:10GB -w",
			),
			core.ExecStoreBeforeCmd(
				"SnapshotA",
				`scw instance snapshot create -w name=cli-test-image-create-snapshotA volume-id={{ (index .Server.Volumes "0").ID }}`,
			),
			core.ExecStoreBeforeCmd(
				"SnapshotB",
				`scw instance snapshot create -w name=cli-test-image-create-snapshotB volume-id={{ (index .Server.Volumes "1").ID }}`,
			),
		),
		Cmd: "scw instance image create snapshot-id={{ .SnapshotA.ID }} extra-volumes.0.id={{ .SnapshotB.ID }} arch=x86_64",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			core.ExecAfterCmd("scw instance image delete {{ .CmdResult.Image.ID }}"),
			core.ExecAfterCmd("scw instance snapshot delete {{ .SnapshotA.ID }}"),
			core.ExecAfterCmd("scw instance snapshot delete {{ .SnapshotB.ID }}"),
		),
	}))
}

func Test_ImageDelete(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createImage("Image"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance image delete {{ .Image.Image.ID }} with-snapshots=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				// Assert snapshot are deleted with the image
				api := instanceSDK.NewAPI(ctx.Client)
				snapshot := testhelpers.MapValue[*instanceSDK.CreateSnapshotResponse](
					t,
					ctx.Meta,
					"Snapshot",
				)

				_, err := api.GetSnapshot(&instanceSDK.GetSnapshotRequest{
					SnapshotID: snapshot.Snapshot.ID,
				})
				assert.ErrorAs(t, err, new(*scw.ResourceNotFoundError))
			},
		),
		AfterFunc: deleteServer("Server"),
	}))
}

func createImage(metaKey string) core.BeforeFunc {
	return core.BeforeFuncCombine(
		core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true image=ubuntu-jammy root-volume=l:20G"),
		),
		core.ExecStoreBeforeCmd(
			"Snapshot",
			`scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`,
		),
		core.ExecStoreBeforeCmd(
			metaKey,
			`scw instance image create snapshot-id={{ .Snapshot.Snapshot.ID }} arch=x86_64`,
		),
	)
}

func createImageWithSBSRootVolume(metaKey string) core.BeforeFunc {
	return core.BeforeFuncCombine(
		core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true image=ubuntu-jammy root-volume=sbs:20G:5000 --wait"),
		),
		core.ExecStoreBeforeCmd(
			"Snapshot",
			`scw block snapshot create volume-id={{ (index .Server.Volumes "0").ID }} --wait`,
		),
		core.ExecStoreBeforeCmd(
			metaKey,
			`scw instance image create snapshot-id={{ .Snapshot.ID }} arch=x86_64`,
		),
	)
}

func createImageWithSBSAdditionalVolumes(metaKey string) core.BeforeFunc {
	return core.BeforeFuncCombine(
		core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand(
				"stopped=true image=ubuntu-jammy root-volume=local:20G additional-volumes.0=sbs:10GB:5000 additional-volumes.1=sbs:15GB:15000 --wait",
			),
		),
		core.ExecStoreBeforeCmd(
			"SnapshotRoot",
			`scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }} --wait`,
		),
		core.ExecStoreBeforeCmd(
			"SnapshotAdditional1",
			`scw block snapshot create volume-id={{ (index .Server.Volumes "1").ID }} --wait`,
		),
		core.ExecStoreBeforeCmd(
			"SnapshotAdditional2",
			`scw block snapshot create volume-id={{ (index .Server.Volumes "2").ID }} --wait`,
		),
		core.ExecStoreBeforeCmd(
			metaKey,
			`scw instance image create snapshot-id={{ .SnapshotRoot.ID }} additional-volumes.0.id={{ .SnapshotAdditional1.ID }} additional-volumes.1.id={{ .SnapshotAdditional2.ID }} arch=x86_64`,
		),
	)
}

func deleteImage(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd(
		`scw instance image delete {{ .` + metaKey + `.Image.ID }} with-snapshots=true`,
	)
}

func Test_ImageList(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		BeforeFunc: createImage("Image"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance image list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteImage("Image"),
			deleteServer("Server"),
		),
	}))

	t.Run("With SBS root volume", core.Test(&core.TestConfig{
		BeforeFunc: createImageWithSBSRootVolume("ImageSBSRoot"),
		Commands: core.NewCommandsMerge(
			instance.GetCommands(),
			block.GetCommands(),
		),
		Cmd: "scw instance image list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteImage("ImageSBSRoot"),
			deleteServer("Server"),
		),
	}))

	t.Run("With SBS additional volumes", core.Test(&core.TestConfig{
		BeforeFunc: createImageWithSBSAdditionalVolumes("ImageSBSAdditional"),
		Commands: core.NewCommandsMerge(
			instance.GetCommands(),
			block.GetCommands(),
		),
		Cmd: "scw instance image list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteImage("ImageSBSAdditional"),
			deleteServer("Server"),
		),
	}))
}

func Test_ImageUpdate(t *testing.T) {
	t.Run("Change name", core.Test(&core.TestConfig{
		BeforeFunc: createImage("ImageName"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance image update {{ .ImageName.Image.ID }} name=foo",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.NotNil(t, ctx.Result)
				assert.Equal(t, "foo", ctx.Result.(*instanceSDK.UpdateImageResponse).Image.Name)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteImage("ImageName"),
		),
	}))

	t.Run("Change public from default false to true", core.Test(&core.TestConfig{
		BeforeFunc: createImage("ImagePub"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance image update {{ .ImagePub.Image.ID }} public=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.NotNil(t, ctx.Result)
				assert.True(t, ctx.Result.(*instanceSDK.UpdateImageResponse).Image.Public)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteImage("ImagePub"),
		),
	}))

	t.Run("Add extra volume", func(t *testing.T) {
		t.Skip("Skipping 'Add extra volume' test temporarily")
		core.Test(&core.TestConfig{
			BeforeFunc: core.BeforeFuncCombine(
				createNonEmptyLocalVolume("Volume", 10),
				core.ExecStoreBeforeCmd(
					"SnapshotVol",
					`scw instance snapshot create -w name=snapVol volume-id={{ .Volume.ID }}`,
				),
				createImage("ImageExtraVol"),
			),
			Commands: instance.GetCommands(),
			Cmd:      "scw instance image update {{ .ImageExtraVol.Image.ID }} extra-volumes.1.id={{ .SnapshotVol.ID }}",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					assert.Equal(
						t,
						"snapVol",
						ctx.Result.(*instanceSDK.UpdateImageResponse).Image.ExtraVolumes["1"].Name,
					)
				},
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServer("Server"),
				deleteImage("ImageExtraVol"),
				deleteVolume("Volume"),
			),
		})
	})
}
