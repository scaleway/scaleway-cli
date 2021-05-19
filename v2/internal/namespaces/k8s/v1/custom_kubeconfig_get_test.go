package k8s

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

func Test_GetKubeconfig(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createClusterAndWaitAndKubeconfig("Cluster", "Kubeconfig", kapsuleVersion),
		Cmd:        "scw k8s kubeconfig get {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, ctx.Result.(string), string(ctx.Meta["Kubeconfig"].(*k8s.Kubeconfig).GetRaw()))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
	}))
}
