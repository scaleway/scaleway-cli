package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
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

	t.Run("List images", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			core.ExecStoreBeforeCmd("Snapshot", `scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`),
			core.ExecStoreBeforeCmd("DeletedServer", `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all`),
		),
		Commands: GetCommands(),
		Cmd:      "scw instance image create snapshot-id={{ .Snapshot.Snapshot.ID }} arch=x86_64",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			func(ctx *core.AfterFuncCtx) error {
				ctx.ExecuteCmd("scw instance image delete " + ctx.CmdResult.(*instance.CreateImageResponse).Image.ID)
				return nil
			},
			deleteSnapshot("Snapshot"),
		),
	}))
}
