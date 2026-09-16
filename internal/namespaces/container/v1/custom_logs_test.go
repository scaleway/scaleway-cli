package container_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1"
	"github.com/stretchr/testify/assert"
)

func Test_ContainerLogs(t *testing.T) {
	image := "hello-world:latest"

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createNamespace("Namespace"),
			core.ExecStoreBeforeCmd("Container", fmt.Sprintf(
				"scw container container create namespace-id={{ .Namespace.ID }} name=%s image=%s -w",
				core.GetRandomName("test-logs"),
				image,
			)),
		),
		Cmd: "scw container container logs {{ .Container.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNamespace("Namespace"),
		),
	}))

	t.Run("MaxEntries", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createNamespace("Namespace"),
			core.ExecStoreBeforeCmd("Container", fmt.Sprintf(
				"scw container container create namespace-id={{ .Namespace.ID }} name=%s image=%s -w",
				core.GetRandomName("test-logs"),
				image,
			)),
		),
		Cmd: "scw container container logs {{ .Container.ID }} entry-count=5",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				c := ctx.Result.([]container.LogEntry)
				assert.Len(t, c, 5)
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNamespace("Namespace"),
		),
	}))
}
