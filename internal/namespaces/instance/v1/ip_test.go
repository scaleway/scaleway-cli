package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func Test_IpCreate(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance ip create",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance ip delete " + ctx.CmdResult.(*instance.CreateIPResponse).IP.ID)
			return nil
		},
	}))
}

func Test_IpDelete(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   GetCommands(),
		Cmd:        "scw instance ip delete {{ .Ip.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_IpGet(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   GetCommands(),
		Cmd:        "scw instance ip get {{ .Ip.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteIP("Ip"),
	}))
}

func Test_IpList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   GetCommands(),
		Cmd:        "scw instance ip list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteIP("Ip"),
	}))
}

func Test_IpUpdate(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   GetCommands(),
		Cmd:        "scw instance ip update {{ .Ip.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteIP("Ip"),
	}))
}
