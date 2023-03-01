package alias_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces"
)

func Test_Alias(t *testing.T) {
	t.Run("raw alias", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create i command.0=instance"),
			core.ExecBeforeCmd("scw alias create sl command.0=server command.1=list"),
		),
		Commands:      namespaces.GetCommands(),
		Cmd:           "scw i sl -h",
		EnableAliases: true,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Contains(t, string(ctx.Stderr), "instance server list")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("resource alias", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create i resource=instance"),
			core.ExecBeforeCmd("scw alias create s resource=instance.server"),
			core.ExecBeforeCmd("scw alias create l resource=instance.server.list"),
		),
		Commands:      namespaces.GetCommands(),
		Cmd:           "scw i s l -h",
		EnableAliases: true,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Contains(t, string(ctx.Stderr), "instance server list")
			},
		),
		TmpHomeDir: true,
	}))
}
