package vpc

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
)

func Test_GetPrivateNetwork(t *testing.T) {
	cmds := GetCommands()
	cmds.Merge(instance.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstance(),
			createNIC(),
		),
		Cmd:   "scw vpc private-network get {{ .PN.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteInstance(),
			deletePN(),
		),
	}))
}
