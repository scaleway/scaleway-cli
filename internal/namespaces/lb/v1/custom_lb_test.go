package lb

import (
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
)

func Test_ListLB(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createLB(),
		Cmd:        "scw lb lb list",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteLB(),
	}))
}

func Test_CreateLB(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:  GetCommands(),
		Cmd:       "scw lb lb create name=foobar description=foobar --wait",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw lb lb delete {{ .CmdResult.ID }}"),
	}))
}

func Test_GetLB(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createLB(),
		Cmd:        "scw lb lb get {{ .LB.ID }}",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteLB(),
	}))
}

func Test_WaitLB(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"LB",
			"scw lb lb create name=cli-test description=cli-test",
		),
		Cmd:       "scw lb lb wait {{ .LB.ID }}",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteLB(),
	}))
}

func Test_GetStats(t *testing.T) {
	commands := GetCommands()
	commands.Merge(instance.GetCommands())
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createInstance(),
			createBackend(80),
			createBackend(81),
			addIP2Backend("{{ .Instance.PublicIP.Address }}"),
			createFrontend(8888),
			// We let enough time for the health checks to come through
			core.BeforeFuncWhenUpdatingCassette(
				func(ctx *core.BeforeFuncCtx) error {
					time.Sleep(10 * time.Second)
					return nil
				},
			),
		),
		Cmd:   "scw lb lb get-stats {{ .LB.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteInstance(),
		),
	}))
}

func Test_GetK8sTaggedLB(t *testing.T) {
	cmds := GetCommands()
	cmds.Merge(k8s.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:        cmds,
		DisableParallel: true,
		BeforeFunc: core.BeforeFuncCombine(
			createClusterAndWaitAndInstallKubeconfig("Cluster", "Kubeconfig", "1.24.7"),
			retrieveLBID("LBID"),
		),
		Cmd:   "scw lb lb get {{ .LBID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteCluster("Cluster"),
			core.ExecAfterCmd("scw lb lb delete {{ .LBID }}"),
		),
	}))
}
