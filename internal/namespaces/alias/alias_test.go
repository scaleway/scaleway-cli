package alias_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces"
)

func Test_Alias(t *testing.T) {
	t.Run("Multiple words", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create i command=instance"),
			core.ExecBeforeCmdArgs([]string{"scw", "alias", "create", "sl", "command=server list"}),
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

	t.Run("one word aliases", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create i command=instance"),
			core.ExecBeforeCmd("scw alias create s command=server"),
			core.ExecBeforeCmd("scw alias create l command=list"),
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
