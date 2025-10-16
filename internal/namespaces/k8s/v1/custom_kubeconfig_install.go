package k8s

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	kubeLocationDir = ".kube"
)

type k8sKubeconfigInstallRequest struct {
	ClusterID          string
	Region             scw.Region
	AuthMethod         authMethods
	KeepCurrentContext bool
}

func k8sKubeconfigInstallCommand() *core.Command {
	return &core.Command{
		Short: `Install a kubeconfig`,
		Long: `Retrieve the kubeconfig for a specified cluster and write it on disk. 
It will merge the new kubeconfig in the file pointed by the KUBECONFIG variable. If empty it will default to $HOME/.kube/config.`,
		Namespace: "k8s",
		Verb:      "install",
		Resource:  "kubeconfig",
		ArgsType:  reflect.TypeOf(k8sKubeconfigInstallRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      "Cluster ID from which to retrieve the kubeconfig",
				Required:   true,
				Positional: true,
			},
			{
				Name:       "auth-method",
				Short:      `Which method to use to authenticate using kubelet`,
				Default:    core.DefaultValueSetter(defaultAuthMethod),
				EnumValues: enumAuthMethods,
			},
			{
				Name:  "keep-current-context",
				Short: "Whether or not to keep the current kubeconfig context unmodified",
			},
			core.RegionArgSpec(),
		},
		Run: k8sKubeconfigInstallRun,
		Examples: []*core.Example{
			{
				Short:    "Install the kubeconfig for a given cluster and using the new context",
				ArgsJSON: `{"cluster_id": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Get the kubeconfig for a given cluster by copying current secret_key to it",
				ArgsJSON: `{"cluster_id": "11111111-1111-1111-1111-111111111111", "auth_method": "copy-cli-token"}`,
			},

			{
				Short:    "Get the kubeconfig for a given cluster and use legacy authentication",
				ArgsJSON: `{"cluster_id": "11111111-1111-1111-1111-111111111111", "auth_method": "legacy"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw k8s kubeconfig uninstall",
				Short:   "Uninstall a kubeconfig",
			},
		},
	}
}

func k8sKubeconfigInstallRun(ctx context.Context, argsI any) (i any, e error) {
	request := argsI.(*k8sKubeconfigInstallRequest)
	kubeconfigPath := getKubeconfigPath(ctx)

	// get cluster kubeconfig
	apiKubeconfig, err := k8s.NewAPI(core.ExtractClient(ctx)).
		GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
			Region:    request.Region,
			ClusterID: request.ClusterID,
			Redacted: scw.BoolPtr(
				request.AuthMethod != authMethodLegacy,
			), // put true after legacy deprecation
		})
	if err != nil {
		return nil, err
	}

	var clusterKubeconfig api.Config
	if err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &clusterKubeconfig); err != nil {
		return nil, err
	}

	namedAuthInfo := clusterKubeconfig.AuthInfos[0]
	namedAuthInfo.Name = fmt.Sprintf("%s-%s", clusterKubeconfig.Clusters[0].Name, request.ClusterID)
	if request.AuthMethod != authMethodLegacy {
		namedAuthInfoPtr, err := generateNamedAuthInfo(ctx, request.AuthMethod)
		if err != nil {
			return nil, err
		}
		namedAuthInfo = *namedAuthInfoPtr
	}

	kubeconfigManager := NewKubeMapConfig()
	if _, err := os.Stat(kubeconfigPath); err == nil {
		kubeconfigManager, err = LoadKubeMapConfig(ctx, kubeconfigPath)
		if err != nil {
			return nil, err
		}
	}

	err = kubeconfigManager.SetUser(namedAuthInfo.Name, namedAuthInfo.AuthInfo, true)
	if err != nil {
		return nil, err
	}

	clusterNameWithID := fmt.Sprintf("%s-%s", clusterKubeconfig.Clusters[0].Name, request.ClusterID)
	err = kubeconfigManager.SetCluster(
		clusterNameWithID,
		clusterKubeconfig.Clusters[0].Cluster,
		true,
	)
	if err != nil {
		return nil, err
	}

	err = kubeconfigManager.SetContext(
		clusterNameWithID,
		api.Context{Cluster: clusterNameWithID, AuthInfo: namedAuthInfo.Name},
		true,
	)
	if err != nil {
		return nil, err
	}

	if !request.KeepCurrentContext {
		kubeconfigManager.CurrentContext = clusterNameWithID
	}

	err = kubeconfigManager.Save(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf(
		"Kubeconfig for cluster %s successfully written at %s",
		request.ClusterID,
		kubeconfigPath,
	), nil
}
