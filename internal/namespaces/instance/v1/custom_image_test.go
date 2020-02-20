package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ImageCreate(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu_bionic stopped")
			//ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }}")
			ctx.Meta["SnapshotResponse"] = ctx.ExecuteCmd(`scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`)
			return nil
		},
		Cmd: "scw instance image create snapshot-id={{ .SnapshotResponse.Snapshot.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {

				//				println("ee")
				//				assert.Equal(t, ctx.Meta["Server"].(*instance.Server).Image.ID, ctx.Result.(*instance.CreateImageResponse).Image.ID)
			},
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance snapshot delete snapshot-id={{ .SnapshotResponse.Snapshot.ID }}")
			//			ctx.ExecuteCmd("scw instance image delete image-id={{ .Result.Image.ID }}")
			return nil
		},
	}))
}
