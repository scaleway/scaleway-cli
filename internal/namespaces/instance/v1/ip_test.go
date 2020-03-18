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
		Check:    core.TestCheckGolden(),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance ip delete ip=" + ctx.CmdResult.(*instance.CreateIPResponse).IP.ID)
			return nil
		},
	}))
}

func Test_IpDelete(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   GetCommands(),
		Cmd:        "scw instance ip delete ip={{ .Ip.ID }}",
		Check:      core.TestCheckGolden(),
		//AfterFunc:  deleteIP("Ip"),
	}))
}

func Test_IpGet(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   GetCommands(),
		Cmd:        "scw instance ip get ip={{ .Ip.ID }}",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteIP("Ip"),
	}))
}

func Test_IpList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   GetCommands(),
		Cmd:        "scw instance ip list",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteIP("Ip"),
	}))
}

func Test_IpUpdate(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   GetCommands(),
		Cmd:        "scw instance ip update ip={{ .Ip.ID }}",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteIP("Ip"),
	}))
}
