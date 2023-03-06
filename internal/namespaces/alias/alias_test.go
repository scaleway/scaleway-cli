package alias_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
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

	t.Run("list aliases", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create myalias command=iam"),
		),
		Commands:      namespaces.GetCommands(),
		Cmd:           "scw alias list",
		EnableAliases: true,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Contains(t, string(ctx.Stdout), "myalias")
				assert.Contains(t, string(ctx.Stdout), "iam")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("delete alias", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create i command=instance"),
		),
		Commands:      namespaces.GetCommands(),
		Cmd:           "scw alias delete i",
		EnableAliases: true,
		AfterFunc: core.AfterFuncCombine(
			func(ctx *core.AfterFuncCtx) error {
				res := ctx.ExecuteCmd([]string{"scw", "alias", "list"})
				resString, err := human.Marshal(res, nil)
				if err != nil {
					return err
				}

				if strings.Contains(resString, "instance") {
					return fmt.Errorf("alias list should not contain instance")
				}
				return nil
			},
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Contains(t, string(ctx.Stdout), "Deleted")
			},
		),
		TmpHomeDir: true,
	}))
}
