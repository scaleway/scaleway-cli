package k8s_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
)

func Test_GetCluster(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		BeforeFunc: createCluster("get-cluster", false),
		Cmd:        "scw k8s cluster get {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
}

func Test_WaitCluster(t *testing.T) {
	t.Run("wait for pools", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		BeforeFunc: createCluster("wait-cluster", false),
		Cmd:        "scw k8s cluster wait {{ .Cluster.ID }} wait-for-pools=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
}
