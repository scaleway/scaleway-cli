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
