package k8s

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_WaitCluster(t *testing.T) {
	t.Run("wait for pools", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createCluster("Cluster", kapsuleVersion, 3),
		Cmd:        "scw k8s cluster wait {{ .Cluster.ID }} wait-for-pools=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
	}))
}
