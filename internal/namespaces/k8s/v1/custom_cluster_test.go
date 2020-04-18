package k8s

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

// deleteClusterAfterFunc deletes the created cluster.
func deleteClusterAfterFunc() core.AfterFunc {
	return core.ExecAfterCmd("scw k8s cluster delete {{ .CmdResult.ID }} --wait")
}

func Test_ClusterCreate(t *testing.T) {
	t.Run("cluster with pool size zero", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw k8s cluster create name=test version=1.18.2 pools.0.size=1 pools.0.node-type=DEV1-M pools.0.name=default cni=cilium pools.1.size=0 pools.1.name=empty pools.1.node-type=DEV1-M",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				poolResMapSize := map[string]uint32{
					"default": 1,
					"empty":   0,
				}
				api := k8s.NewAPI(ctx.Client)
				pools, err := api.ListPools(&k8s.ListPoolsRequest{
					ClusterID: ctx.Result.(*k8s.Cluster).ID,
				})
				assert.Equal(t, nil, err)
				for _, pool := range pools.Pools {
					expectedSize, ok := poolResMapSize[pool.Name]
					assert.Equal(t, true, ok)
					assert.Equal(t, expectedSize, pool.Size)
				}
			},
		),
		AfterFunc: deleteClusterAfterFunc(),
	}))
}
