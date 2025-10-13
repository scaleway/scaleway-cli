package k8s_test

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	go_api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

const (
	kapsuleVersion    = "1.32.3"
	clusterMetaKey    = "Cluster"
	kubeconfigMetaKey = "Kubeconfig"
)

//
// Clusters
//

// createCluster creates a basic cluster with "poolSize" dev1-m as nodes, the given version and
// register it in the context Meta at metaKey.
func createCluster(
	clusterNameSuffix string,
	poolSize int,
	nodeType string,
) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		clusterMetaKey,
		fmt.Sprintf(
			"scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=%s pools.0.size=%d pools.0.name=default",
			clusterNameSuffix,
			kapsuleVersion,
			nodeType,
			poolSize,
		),
	)
}

// createClusterAndWaitAndKubeconfig creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey.
func createClusterAndWaitAndKubeconfig(
	clusterNameSuffix string,
) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf(
			"scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default --wait",
			clusterNameSuffix,
			kapsuleVersion,
		)
		res := ctx.ExecuteCmd(strings.Split(cmd, " "))
		cluster := res.(*k8s.Cluster)
		ctx.Meta[clusterMetaKey] = cluster

		apiKubeconfig, err := k8s.NewAPI(ctx.Client).
			GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
				Region:    cluster.Region,
				ClusterID: cluster.ID,
			})
		if err != nil {
			return err
		}

		var kubeconfig go_api.Config
		if err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &kubeconfig); err != nil {
			return err
		}
		ctx.Meta[kubeconfigMetaKey] = kubeconfig

		return nil
	}
}

func populateKubeconfigAndCreateCluster(
	kubeconfigRaw []byte,
	clusterNameSuffix string,
) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		// Populate $HOME/.kube/config with existing data
		kubeconfigPath := path.Join(ctx.Meta["HOME"].(string), ".kube", "config")
		if err := os.MkdirAll(path.Dir(kubeconfigPath), 0o755); err != nil {
			return err
		}

		if err := os.WriteFile(kubeconfigPath, kubeconfigRaw, 0o600); err != nil {
			return err
		}

		// Then create a cluster
		return createClusterAndWaitAndKubeconfig(
			clusterNameSuffix,
		)(ctx)
	}
}

// deleteCluster deletes a cluster previously registered in the context Meta at metaKey.
func deleteCluster() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw k8s cluster delete {{ ." + clusterMetaKey + ".ID }} with-additional-resources=true --wait",
	)
}
