package k8s

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/ghodss/yaml"
	api "github.com/kubernetes-client/go-base/config/api"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_GetKubeconfig(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createClusterAndWaitAndKubeconfig("get-kubeconfig", "Cluster", "Kubeconfig", kapsuleVersion),
		Cmd:        "scw k8s kubeconfig get {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				config, err := yaml.Marshal(ctx.Meta["Kubeconfig"].(api.Config))
				assert.Equal(t, err, nil)
				assert.Equal(t, ctx.Result.(string), string(config))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
	}))
}
