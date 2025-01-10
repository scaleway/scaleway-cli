package lb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
)

func Test_GetFrontend(t *testing.T) {
	cmds := lb.GetCommands()
	cmds.Merge(instance.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createInstance(),
			createBackend(80),
			createFrontend(8888),
		),
		Cmd:   "scw lb frontend get {{ .Frontend.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteInstance(),
			deleteLBFlexibleIP(),
		),
	}))
}
