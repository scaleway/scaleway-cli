package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ServerUpdateCustom(t *testing.T) {
	t.Run("Remove ip from server", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} ip=none",
		Check: core.TestCheckCombine(
			core.TestCheckNil(".Server.IP"),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
	}))

	t.Run("Update server ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create image=ubuntu-bionic")
			ctx.Meta["IP"] = ctx.ExecuteCmd("scw instance ip create")
			return nil
		},
		Cmd: "scw instance server update server-id={{ .Server.ID }} ip={{ .IP.IP }}",
		Check: core.TestCheckCombine(
			core.TestCheckEqual("IP.ip", ".Server.IP"),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .Server.ID }} delete-ip=true delete-volumes=true")
			return nil
		},
	}))
}
