package k8s_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
)

func Test_PoolSetTaint(t *testing.T) {
	t.Run("set-empty", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				clusterMetaKey,
				fmt.Sprintf(
					"scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default --wait=true",
					t.Name(),
					kapsuleVersion,
				),
			),
			fetchPoolMetadata(clusterMetaKey, poolMetaKey, "default"),
		),
		Cmd: "scw k8s pool set-taint {{ ." + poolMetaKey + ".ID }} key=foo value=bar effect=NoSchedule",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
}

func Test_PoolRemoveTaint(t *testing.T) {
	t.Run("remove-existing", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				clusterMetaKey,
				fmt.Sprintf(
					"scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default pools.0.taints.0.key=foo pools.0.taints.0.value=bar pools.0.taints.0.effect=NoSchedule --wait=true",
					t.Name(),
					kapsuleVersion,
				),
			),
			fetchPoolMetadata(clusterMetaKey, poolMetaKey, "default"),
		),
		Cmd: "scw k8s pool remove-taint {{ ." + poolMetaKey + ".ID }} key=foo",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
}

func Test_PoolSetStartupTaint(t *testing.T) {
	t.Run("set-empty", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				clusterMetaKey,
				fmt.Sprintf(
					"scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default --wait=true",
					t.Name(),
					kapsuleVersion,
				),
			),
			fetchPoolMetadata(clusterMetaKey, poolMetaKey, "default"),
		),
		Cmd: "scw k8s pool set-startup-taint {{ ." + poolMetaKey + ".ID }} key=foo value=bar effect=NoSchedule",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
}

func Test_PoolRemoveStartupTaint(t *testing.T) {
	t.Run("remove-existing", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				clusterMetaKey,
				fmt.Sprintf(
					"scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default pools.0.taints.0.key=foo pools.0.taints.0.value=bar pools.0.taints.0.effect=NoSchedule --wait=true",
					t.Name(),
					kapsuleVersion,
				),
			),
			fetchPoolMetadata(clusterMetaKey, poolMetaKey, "default"),
		),
		Cmd: "scw k8s pool remove-startup-taint {{ ." + poolMetaKey + ".ID }} key=foo",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
}
