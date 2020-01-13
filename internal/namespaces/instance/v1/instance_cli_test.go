package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ListServer(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server list -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

}

func Test_GetServer(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server get -h",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: "scw instance server get server-id={{ .Server.id }}",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))
}
