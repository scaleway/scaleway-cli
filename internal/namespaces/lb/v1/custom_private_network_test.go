package lb

import (
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
)

func Test_ListLBPrivateNetwork(t *testing.T) {
	cmds := GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createPN(),
			attachPN(),
		),
		Cmd:   "scw lb private-network list {{ .LB.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			detachPN(),
			deleteLB(),
			core.AfterFuncWhenUpdatingCassette(
				func(_ *core.AfterFuncCtx) error {
					time.Sleep(10 * time.Second)
					return nil
				},
			),
			deletePN(),
		),
	}))
}
