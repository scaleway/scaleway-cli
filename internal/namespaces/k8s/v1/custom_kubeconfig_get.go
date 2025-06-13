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
	ClusterID string
	Region    scw.Region
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
			core.RegionArgSpec(),
		},
		Run: k8sKubeconfigGetRun,
		Examples: []*core.Example{
			{
				Short:    "Get the kubeconfig for a given cluster",
				ArgsJSON: `{"cluster_id": "11111111-1111-1111-1111-111111111111"}`,
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

	kubeconfigRequest := &k8s.GetClusterKubeConfigRequest{
		Region:    request.Region,
		ClusterID: request.ClusterID,
	}

	client := core.ExtractClient(ctx)
	apiK8s := k8s.NewAPI(client)

	apiKubeconfig, err := apiK8s.GetClusterKubeConfig(kubeconfigRequest)
	if err != nil {
		return nil, err
	}

	var kubeconfig api.Config

	err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &kubeconfig)
	if err != nil {
		return nil, err
	}

	config, err := yaml.Marshal(kubeconfig)
	if err != nil {
		return nil, err
	}

	return string(config), nil
}
