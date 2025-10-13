package k8s

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
)

type k8sKubeconfigUninstallRequest struct {
	ClusterID string
}

func k8sKubeconfigUninstallCommand() *core.Command {
	return &core.Command{
		Short: `Uninstall a kubeconfig`,
		Long: `Remove specified cluster from kubeconfig file specified by the KUBECONFIG env, if empty it will default to $HOME/.kube/config.
If the current context points to this cluster, it will be set to an empty context.`,
		Namespace: "k8s",
		Verb:      "uninstall",
		Resource:  "kubeconfig",
		ArgsType:  reflect.TypeOf(k8sKubeconfigUninstallRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      "Cluster ID from which to uninstall the kubeconfig",
				Required:   true,
				Positional: true,
			},
		},
		Run: k8sKubeconfigUninstallRun,
		Examples: []*core.Example{
			{
				Short:    "Uninstall the kubeconfig for a given cluster",
				ArgsJSON: `{"cluster_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw k8s kubeconfig install",
				Short:   "Install a kubeconfig",
			},
		},

		// avoid calling checkAPIKey (Check if API Key is about to expire)
		DisableAfterChecks: true,
	}
}

// k8sKubeconfigUninstallRun use the specified cluster ID to remove it from the wanted kubeconfig file
// it removes all the users, contexts and clusters that contains this ID from the file
func k8sKubeconfigUninstallRun(ctx context.Context, argsI any) (i any, e error) {
	request := argsI.(*k8sKubeconfigUninstallRequest)
	kubeconfigPath := getKubeconfigPath(ctx)

	// if the file does not exist, the cluster is not there
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return fmt.Sprintf("File %s does not exist.", kubeconfigPath), nil
	}

	existingKubeconfig, err := openAndUnmarshalKubeconfig(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	newClusters := []api.NamedCluster{}
	for _, cluster := range existingKubeconfig.Clusters {
		if !strings.HasSuffix(cluster.Name, request.ClusterID) {
			newClusters = append(newClusters, cluster)
		}
	}

	newContexts := []api.NamedContext{}
	for _, kubeconfigContext := range existingKubeconfig.Contexts {
		if !strings.HasSuffix(kubeconfigContext.Name, request.ClusterID) {
			newContexts = append(newContexts, kubeconfigContext)
		}
	}

	newUsers := []api.NamedAuthInfo{}
	for _, user := range existingKubeconfig.AuthInfos {
		if !strings.HasSuffix(user.Name, request.ClusterID) {
			newUsers = append(newUsers, user)
		}
	}

	// reset the current context
	if strings.HasSuffix(existingKubeconfig.CurrentContext, request.ClusterID) {
		existingKubeconfig.CurrentContext = ""
	}

	// write the modification
	existingKubeconfig.Clusters = newClusters
	existingKubeconfig.Contexts = newContexts
	existingKubeconfig.AuthInfos = newUsers

	err = marshalAndWriteKubeconfig(existingKubeconfig, kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf(
		"Cluster %s successfully deleted from %s",
		request.ClusterID,
		kubeconfigPath,
	), nil
}
