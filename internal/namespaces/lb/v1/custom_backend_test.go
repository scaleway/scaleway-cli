package lb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
)

func Test_GetBackend(t *testing.T) {
	cmds := GetCommands()
	cmds.Merge(instance.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createInstance(),
			createBackend(80),
			addIP2Backend("{{ .Instance.PublicIP.Address }}"),
		),
		Cmd:   "scw lb backend get {{ .Backend.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteInstance(),
		),
	}))
}
