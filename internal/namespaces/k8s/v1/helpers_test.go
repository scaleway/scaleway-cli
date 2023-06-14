package k8s

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	go_api "github.com/kubernetes-client/go-base/config/api"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

const (
	kapsuleVersion = "1.27.1"
)

//
// Clusters
//

// createCluster creates a basic cluster with "poolSize" dev1-m as nodes, the given version and
// register it in the context Meta at metaKey.
func createCluster(clusterNameSuffix string, metaKey string, version string, poolSize int, nodeType string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		metaKey,
		fmt.Sprintf("scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=%s pools.0.size=%d pools.0.name=default", clusterNameSuffix, version, nodeType, poolSize))
}

// createClusterAndWaitAndKubeconfig creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey.
func createClusterAndWaitAndKubeconfig(clusterNameSuffix string, metaKey string, kubeconfigMetaKey string, version string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default --wait", clusterNameSuffix, version)
		res := ctx.ExecuteCmd(strings.Split(cmd, " "))
		cluster := res.(*k8s.Cluster)
		ctx.Meta[metaKey] = cluster
		api := k8s.NewAPI(ctx.Client)
		apiKubeconfig, err := api.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    cluster.Region,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return err
		}

		var kubeconfig go_api.Config

		err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &kubeconfig)
		if err != nil {
			return err
		}

		ctx.Meta[kubeconfigMetaKey] = kubeconfig
		return nil
	}
}

// createClusterAndWaitAndInstallKubeconfig creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey. And install the kubeconfig
func createClusterAndWaitAndInstallKubeconfig(clusterNameSuffix string, metaKey string, kubeconfigMetaKey string, version string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default --wait", clusterNameSuffix, version)
		res := ctx.ExecuteCmd(strings.Split(cmd, " "))
		cluster := res.(*k8s.Cluster)
		ctx.Meta[metaKey] = cluster
		api := k8s.NewAPI(ctx.Client)
		apiKubeconfig, err := api.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    cluster.Region,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return err
		}

		var kubeconfig go_api.Config

		err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &kubeconfig)
		if err != nil {
			return err
		}

		ctx.Meta[kubeconfigMetaKey] = kubeconfig
		cmd = fmt.Sprintf("scw k8s kubeconfig install %s", cluster.ID)
		_ = ctx.ExecuteCmd(strings.Split(cmd, " "))
		return nil
	}
}

// createClusterAndWaitAndKubeconfigAndPopulateFile creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey. It also populates the given file with the given content
func createClusterAndWaitAndKubeconfigAndPopulateFile(clusterNameSuffix string, metaKey string, kubeconfigMetaKey string, version string, file string, content []byte) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default --wait", clusterNameSuffix, version)
		res := ctx.ExecuteCmd(strings.Split(cmd, " "))
		cluster := res.(*k8s.Cluster)
		ctx.Meta[metaKey] = cluster
		api := k8s.NewAPI(ctx.Client)
		apiKubeconfig, err := api.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    cluster.Region,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return err
		}

		var kubeconfig go_api.Config

		err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &kubeconfig)
		if err != nil {
			return err
		}

		ctx.Meta[kubeconfigMetaKey] = kubeconfig
		err = os.WriteFile(file, content, 0644)
		return err
	}
}

// createClusterAndWaitAndKubeconfigAndPopulateFileAndInstall creates a basic cluster with 1 dev1-m as node, the given version and
// register it in the context Meta at metaKey. It also populates the given file with the given content and install the new kubeconfig
func createClusterAndWaitAndKubeconfigAndPopulateFileAndInstall(clusterNameSuffix string, metaKey string, kubeconfigMetaKey string, version string, file string, content []byte) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default --wait", clusterNameSuffix, version)
		res := ctx.ExecuteCmd(strings.Split(cmd, " "))
		cluster := res.(*k8s.Cluster)
		ctx.Meta[metaKey] = cluster
		api := k8s.NewAPI(ctx.Client)
		apiKubeconfig, err := api.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    cluster.Region,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return err
		}

		var kubeconfig go_api.Config

		err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &kubeconfig)
		if err != nil {
			return err
		}

		ctx.Meta[kubeconfigMetaKey] = kubeconfig
		err = os.WriteFile(file, content, 0644)
		if err != nil {
			return err
		}
		cmd = fmt.Sprintf("scw k8s kubeconfig install %s", cluster.ID)
		_ = ctx.ExecuteCmd(strings.Split(cmd, " "))

		return nil
	}
}

// deleteCluster deletes a cluster previously registered in the context Meta at metaKey.
func deleteCluster(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw k8s cluster delete {{ ." + metaKey + ".ID }} --wait")
}
