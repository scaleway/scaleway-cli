package k8s

import (
	"fmt"
	"io/ioutil"

	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
)

const (
	kapsuleVersion = "1.17.3"
)

//
// Clusters
//

// createClusterAndWaitAndKubeconfig creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey.
func createClusterAndWaitAndKubeconfig(metaKey string, kubeconfigMetaKey string, version string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw k8s cluster create name=cli-test version=%s cni=cilium default-pool-config.node-type=DEV1-M default-pool-config.size=1 --wait", version)
		res := ctx.ExecuteCmd(cmd)
		cluster := res.(*k8s.Cluster)
		ctx.Meta[metaKey] = cluster
		api := k8s.NewAPI(ctx.Client)
		kubeconfig, err := api.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    cluster.Region,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return err
		}
		ctx.Meta[kubeconfigMetaKey] = kubeconfig
		return nil
	}
}

// createClusterAndWaitAndInstallKubeconfig creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey. And install the kubeconfig
func createClusterAndWaitAndInstallKubeconfig(metaKey string, kubeconfigMetaKey string, version string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw k8s cluster create name=cli-test version=%s cni=cilium default-pool-config.node-type=DEV1-M default-pool-config.size=1 --wait", version)
		res := ctx.ExecuteCmd(cmd)
		cluster := res.(*k8s.Cluster)
		ctx.Meta[metaKey] = cluster
		api := k8s.NewAPI(ctx.Client)
		kubeconfig, err := api.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    cluster.Region,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return err
		}
		ctx.Meta[kubeconfigMetaKey] = kubeconfig
		cmd = fmt.Sprintf("scw k8s kubeconfig install %s", cluster.ID)
		_ = ctx.ExecuteCmd(cmd)
		return nil
	}
}

// createClusterAndWaitAndKubeconfigAndPopulateFile creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey. It also populates the given file with the given content
func createClusterAndWaitAndKubeconfigAndPopulateFile(metaKey string, kubeconfigMetaKey string, version string, file string, content []byte) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw k8s cluster create name=cli-test version=%s cni=cilium default-pool-config.node-type=DEV1-M default-pool-config.size=1 --wait", version)
		res := ctx.ExecuteCmd(cmd)
		cluster := res.(*k8s.Cluster)
		ctx.Meta[metaKey] = cluster
		api := k8s.NewAPI(ctx.Client)
		kubeconfig, err := api.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    cluster.Region,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return err
		}
		ctx.Meta[kubeconfigMetaKey] = kubeconfig
		err = ioutil.WriteFile(file, content, 0644)
		return err
	}
}

// createClusterAndWaitAndKubeconfigAndPopulateFileAndInstall creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey. It also populates the given file with the given content and install the new kubeconfig
func createClusterAndWaitAndKubeconfigAndPopulateFileAndInstall(metaKey string, kubeconfigMetaKey string, version string, file string, content []byte) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw k8s cluster create name=cli-test version=%s cni=cilium default-pool-config.node-type=DEV1-M default-pool-config.size=1 --wait", version)
		res := ctx.ExecuteCmd(cmd)
		cluster := res.(*k8s.Cluster)
		ctx.Meta[metaKey] = cluster
		api := k8s.NewAPI(ctx.Client)
		kubeconfig, err := api.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    cluster.Region,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return err
		}
		ctx.Meta[kubeconfigMetaKey] = kubeconfig
		err = ioutil.WriteFile(file, content, 0644)
		if err != nil {
			return err
		}
		cmd = fmt.Sprintf("scw k8s kubeconfig install %s", cluster.ID)
		_ = ctx.ExecuteCmd(cmd)

		return nil
	}
}

// deleteCluster deletes a cluster previously registered in the context Meta at metaKey.
func deleteCluster(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw k8s cluster delete cluster-id={{ ." + metaKey + ".ID }} --wait")
}
