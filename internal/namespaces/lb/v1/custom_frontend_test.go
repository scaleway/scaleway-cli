package lb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
)

func Test_GetFrontend(t *testing.T) {
	cmds := GetCommands()
	cmds.Merge(instance.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createInstance(),
			createBackend(),
			createFrontend(),
		),
		Cmd:   "scw lb frontend get {{ .LB.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteInstance(),
		),
	}))
}
