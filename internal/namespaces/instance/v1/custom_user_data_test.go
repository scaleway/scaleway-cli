package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_UserDataGet(t *testing.T) {
	t.Run("Get an existing key", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			core.ExecBeforeCmd("scw instance user-data set server-id={{.Server.ID}} key=happy content=true"),
		),
		Commands:  GetCommands(),
		Cmd:       "scw instance user-data get server-id={{.Server.ID}} key=happy",
		AfterFunc: deleteServer("Server"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Get an nonexistent key", core.Test(&core.TestConfig{
		BeforeFunc: createServer("Server"),
		Commands:   GetCommands(),
		Cmd:        "scw instance user-data get server-id={{.Server.ID}} key=happy",
		AfterFunc:  deleteServer("Server"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}
