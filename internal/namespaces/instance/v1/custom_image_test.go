package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_ImageCreate(t *testing.T) {
	t.Run("Create simple image", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			core.ExecStoreBeforeCmd("Snapshot", `scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`),
		),
		Commands: GetCommands(),
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
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", "scw instance server create image=ubuntu_focal root-volume=local:10GB additional-volumes.0=local:10GB -w"),
			core.ExecStoreBeforeCmd("SnapshotA", `scw instance snapshot create -w name=cli-test-image-create-snapshotA volume-id={{ (index .Server.Volumes "0").ID }}`),
			core.ExecStoreBeforeCmd("SnapshotB", `scw instance snapshot create -w name=cli-test-image-create-snapshotB volume-id={{ (index .Server.Volumes "1").ID }}`),
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
	t.Run("simple", core.Test(&core.TestConfig{
		BeforeFunc: createImage("Image"),
		Commands:   GetCommands(),
		Cmd:        "scw instance image delete {{ .Image.Image.ID }} with-snapshots=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				// Assert snapshot are deleted with the image
				api := instance.NewAPI(ctx.Client)
				_, err := api.GetSnapshot(&instance.GetSnapshotRequest{
					SnapshotID: ctx.Meta["Snapshot"].(*instance.CreateSnapshotResponse).Snapshot.ID,
				})
				assert.IsType(t, &scw.ResourceNotFoundError{}, err)
			},
		),
		AfterFunc: deleteServer("Server"),
	}))
}

func createImage(metaKey string) core.BeforeFunc {
	return core.BeforeFuncCombine(
		createServer("Server"),
		core.ExecStoreBeforeCmd("Snapshot", `scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`),
		core.ExecStoreBeforeCmd(metaKey, `scw instance image create snapshot-id={{ .Snapshot.Snapshot.ID }} arch=x86_64`),
	)
}

func deleteImage(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd(`scw instance image delete {{ .` + metaKey + `.Image.ID }} with-snapshots=true`)
}

func Test_ImageList(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		BeforeFunc: createImage("Image"),
		Commands:   GetCommands(),
		Cmd:        "scw instance image list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteImage("Image"),
	}))
}
