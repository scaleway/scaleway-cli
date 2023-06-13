package container

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

func createNamespace(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw container namespace create -w")
}

func deleteNamespace(metaKey string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		return core.ExecAfterCmd("scw container namespace delete {{ ." + metaKey + ".ID }}")(ctx)
	}
}

func Test_Create(t *testing.T) {
	commands := GetCommands()

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			createNamespace("Namespace"),
		),
		Cmd: fmt.Sprintf("scw container container create namespace-id={{ .Namespace.ID }} name=%s deploy=true", core.GetRandomName("test")),
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				c := ctx.Result.(*container.Container)
				assert.Equal(t, container.ContainerStatusPending, c.Status)
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNamespace("Namespace"),
		),
		DisableParallel: true,
	}))
}
