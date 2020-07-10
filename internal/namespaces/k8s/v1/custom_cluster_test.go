package k8s

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_GetCluster(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createCluster("Cluster", kapsuleVersion, 1),
		Cmd:        "scw k8s cluster get {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
	}))
}

func Test_WaitCluster(t *testing.T) {
	t.Run("wait for pools", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Cluster",
			"scw k8s cluster create name=cli-test version=1.18.5 cni=cilium pools.0.node-type=DEV1-M pools.0.size=3 pools.0.name=default",
		),
		Cmd: "scw k8s cluster wait {{ .Cluster.ID }} wait-for-pools=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
	}))
}
