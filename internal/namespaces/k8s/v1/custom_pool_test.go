package k8s

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

func Test_ClusterPoolCreate(t *testing.T) {
	t.Run("pool size zero", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createClusterAndWait("Cluster", "1.18.2"),
		Cmd:        "scw k8s pool create {{ .Cluster.ID }} name=empty size=0 node-type=DEV1-M min-size=0",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "empty", ctx.Result.(*k8s.Pool).Name)
				assert.Equal(t, uint32(0), ctx.Result.(*k8s.Pool).Size)
			},
		),
		AfterFunc: deleteCluster("Cluster"),
	}))
}
