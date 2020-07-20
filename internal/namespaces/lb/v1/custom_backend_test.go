package lb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
)

func Test_ImportInstanceBackend(t *testing.T) {
	cmds := GetCommands()
	cmds.Merge(instance.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstance(),
		),
		Cmd:   "scw lb backend import-instance instance-id={{ .Instance.ID }} lb-id={{ .LB.ID }} protocol=tcp port=443",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteInstance(),
		),
	}))
}
