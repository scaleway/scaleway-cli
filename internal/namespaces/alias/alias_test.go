package alias_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/stretchr/testify/assert"
)

func Test_Alias(t *testing.T) {
	t.Run("Multiple words", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create i command=instance --yes"),
			core.ExecBeforeCmdArgs([]string{"scw", "alias", "create", "sl", "command=server list", "--yes"}),
		),
		Commands:      commands.GetCommands(),
		Cmd:           "scw i sl -h",
		EnableAliases: true,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, string(ctx.Stderr), "instance server list")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("one word aliases", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create i command=instance --yes"),
			core.ExecBeforeCmd("scw alias create s command=server --yes"),
			core.ExecBeforeCmd("scw alias create l command=list --yes"),
		),
		Commands:      commands.GetCommands(),
		Cmd:           "scw i s l -h",
		EnableAliases: true,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, string(ctx.Stderr), "instance server list")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("list aliases", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create myalias command=iam --yes"),
		),
		Commands:      commands.GetCommands(),
		Cmd:           "scw alias list",
		EnableAliases: true,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, string(ctx.Stdout), "myalias")
				assert.Contains(t, string(ctx.Stdout), "iam")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("delete alias", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecBeforeCmd("scw alias create i command=instance --yes"),
		),
		Commands:      commands.GetCommands(),
		Cmd:           "scw alias delete i --yes",
		EnableAliases: true,
		AfterFunc: core.AfterFuncCombine(
			func(ctx *core.AfterFuncCtx) error {
				res := ctx.ExecuteCmd([]string{"scw", "alias", "list"})
				resString, err := human.Marshal(res, nil)
				if err != nil {
					return err
				}

				if strings.Contains(resString, "instance") {
					return errors.New("alias list should not contain instance")
				}

				return nil
			},
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, string(ctx.Stdout), "Deleted")
			},
		),
		TmpHomeDir: true,
	}))
}
