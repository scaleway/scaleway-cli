package lb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_GetFrontend(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
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
