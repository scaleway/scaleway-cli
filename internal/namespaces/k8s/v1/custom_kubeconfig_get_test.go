package k8s_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetKubeconfig(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: createClusterAndWaitAndKubeconfig(
			"get-kubeconfig",
			"Cluster",
			"Kubeconfig",
			kapsuleVersion,
		),
		Cmd: "scw k8s kubeconfig get {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				config, err := yaml.Marshal(ctx.Meta["Kubeconfig"].(api.Config))
				require.NoError(t, err)
				assert.Equal(t, ctx.Result.(string), string(config))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
	}))
}
