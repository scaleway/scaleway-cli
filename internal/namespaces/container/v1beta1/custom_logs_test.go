package container_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	container "github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1beta1"
)

func Test_ContainerLogs(t *testing.T) {
	image := "hello-world:latest"

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createNamespace("Namespace"),
			core.ExecStoreBeforeCmd("Container", fmt.Sprintf(
				"scw container container create namespace-id={{ .Namespace.ID }} name=%s registry-image=%s -w",
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
}
