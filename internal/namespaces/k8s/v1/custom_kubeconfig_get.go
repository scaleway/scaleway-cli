package k8s

import (
	"context"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type k8sKubeconfigGetRequest struct {
	ClusterID  string
	Region     scw.Region
	AuthMethod authMethods
}

func k8sKubeconfigGetCommand() *core.Command {
	return &core.Command{
		Short:     `Retrieve a kubeconfig`,
		Long:      `Retrieve the kubeconfig for a specified cluster.`,
		Namespace: "k8s",
		Verb:      "get",
		Resource:  "kubeconfig",
		ArgsType:  reflect.TypeOf(k8sKubeconfigGetRequest{}),
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
			core.RegionArgSpec(),
		},
		Run: k8sKubeconfigGetRun,
		Examples: []*core.Example{
			{
				Short:    "Get the kubeconfig for a given cluster",
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
				Command: "scw k8s kubeconfig install",
				Short:   "Install a kubeconfig",
			},
		},
	}
}

func k8sKubeconfigGetRun(ctx context.Context, argsI any) (i any, e error) {
	request := argsI.(*k8sKubeconfigGetRequest)

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

	var kubeconfig api.Config
	if err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &kubeconfig); err != nil {
		return nil, err
	}

	if request.AuthMethod != authMethodLegacy {
		namedAuthInfo, err := generateNamedAuthInfo(ctx, request.AuthMethod)
		if err != nil {
			return nil, err
		}
		kubeconfig.AuthInfos[0] = *namedAuthInfo
		kubeconfig.Contexts[0].Context.AuthInfo = namedAuthInfo.Name
	}

	config, err := yaml.Marshal(kubeconfig)
	if err != nil {
		return nil, err
	}

	return string(config), nil
}
