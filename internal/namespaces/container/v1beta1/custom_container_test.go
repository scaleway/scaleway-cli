package container_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	container "github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1beta1"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/stretchr/testify/assert"
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
	commands := container.GetCommands()

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			createNamespace("Namespace"),
		),
		Cmd: fmt.Sprintf(
			"scw container container create namespace-id={{ .Namespace.ID }} name=%s deploy=true",
			core.GetRandomName("test"),
		),
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				c := ctx.Result.(*containerSDK.Container)
				assert.Equal(t, containerSDK.ContainerStatusPending, c.Status)
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
