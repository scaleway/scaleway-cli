package lb_test

import (
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/ipam/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
)

func Test_ListLBPrivateNetwork(t *testing.T) {
	cmds := lb.GetCommands()
	cmds.Merge(vpc.GetCommands())
	cmds.Merge(ipam.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createPN(),
			createIPAMIP(),
			attachPN(),
		),
		Cmd:   "scw lb private-network list {{ .LB.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			detachPN(),
			deleteLB(),
			core.AfterFuncWhenUpdatingCassette(
				func(_ *core.AfterFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
			deleteIPAMIP(),
			deletePN(),
			deleteLBFlexibleIP(),
		),
	}))
}
