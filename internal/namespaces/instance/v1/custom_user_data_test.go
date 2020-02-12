package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_UserDataGet(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance user-data get -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Get an existing key", core.Test(&core.TestConfig{
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu_bionic stopped")
			ctx.ExecuteCmd("scw instance user-data set server-id={{.Server.ID}} key=happy content=true")
			return nil
		},
		Commands:  GetCommands(),
		Cmd:       "scw instance user-data get server-id={{.Server.ID}} key=happy",
		AfterFunc: deleteVanillaServer,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Get an nonexistent key", core.Test(&core.TestConfig{
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu_bionic stopped")
			return nil
		},
		Commands:  GetCommands(),
		Cmd:       "scw instance user-data get server-id={{.Server.ID}} key=happy",
		AfterFunc: deleteVanillaServer,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_UserDataSet(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance user-data set -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
