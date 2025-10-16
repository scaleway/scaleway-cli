package k8s_test

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	go_api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
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
	wait bool,
) core.BeforeFunc {
	format := "scw k8s cluster create name=cli-test-%s version=%s cni=cilium pools.0.node-type=DEV1-M pools.0.size=1 pools.0.name=default"
	if wait {
		format += " --wait"
	}

	return core.ExecStoreBeforeCmd(
		clusterMetaKey,
		fmt.Sprintf(
			format,
			clusterNameSuffix,
			kapsuleVersion,
		),
	)
}

// fetchClusterKubeconfigMetadata fetch kubeconfig of previously created cluster.
func fetchClusterKubeconfigMetadata(
	redacted bool,
) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cluster := ctx.Meta[clusterMetaKey].(*k8s.Cluster)

		apiKubeconfig, err := k8s.NewAPI(ctx.Client).
			GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
				Region:    cluster.Region,
				ClusterID: cluster.ID,
				Redacted:  scw.BoolPtr(redacted),
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

func writeKubeconfigFile(kubeconfigRaw []byte) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		// Populate $HOME/.kube/config with existing data
		kubeconfigPath := path.Join(ctx.Meta["HOME"].(string), ".kube", "config")
		if err := os.MkdirAll(path.Dir(kubeconfigPath), 0o755); err != nil {
			return err
		}

		return os.WriteFile(kubeconfigPath, kubeconfigRaw, 0o600)
	}
}

func cliInstallKubeconfig() core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cluster := ctx.Meta[clusterMetaKey].(*k8s.Cluster)
		cmd := "scw k8s kubeconfig install " + cluster.ID
		installOut := ctx.ExecuteCmd(strings.Split(cmd, " "))
		if !strings.Contains(installOut.(string), "successfully written") {
			return errors.New("kubeconfig install failed")
		}

		return nil
	}
}

// deleteCluster deletes a cluster previously registered in the context Meta at metaKey.
func deleteCluster() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw k8s cluster delete {{ ." + clusterMetaKey + ".ID }} with-additional-resources=true --wait",
	)
}
