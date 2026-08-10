package k8s_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
)

func Test_PoolSetLabel(t *testing.T) {
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
		Cmd: "scw k8s pool set-label {{ ." + poolMetaKey + ".ID }} key=foo value=bar",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
}

func Test_PoolRemoveLabel(t *testing.T) {
	t.Run("remove-existing", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				clusterMetaKey,
				fmt.Sprintf(
					"scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default pools.0.labels.foo=bar --wait=true",
					t.Name(),
					kapsuleVersion,
				),
			),
			fetchPoolMetadata(clusterMetaKey, poolMetaKey, "default"),
		),
		Cmd: "scw k8s pool remove-label {{ ." + poolMetaKey + ".ID }} key=foo",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
}
