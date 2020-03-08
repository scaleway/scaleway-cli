package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func Test_SnapshotList(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance snapshot list -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance snapshot list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_SnapshotGet(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance snapshot get -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance snapshot get",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_SnapshotDelete(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance snapshot delete -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance snapshot delete",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_SnapshotCreate(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance snapshot create -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance snapshot create",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        `scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			func(ctx *core.AfterFuncCtx) error {
				ctx.ExecuteCmd("scw instance snapshot delete snapshot-id=" + ctx.CmdResult.(*instance.CreateSnapshotResponse).Snapshot.ID)
				return nil
			},
			deleteServer("Server"),
		),
	}))
}
