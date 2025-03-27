package k8s

import (
	"context"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	kubeLocationDir      = ".kube"
	KubeconfigAPIVersion = "v1"
	KubeconfigKind       = "Config"
)

type k8sKubeconfigInstallRequest struct {
	ClusterID          string
	Region             scw.Region
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
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw k8s kubeconfig uninstall",
				Short:   "Uninstall a kubeconfig",
			},
		},
	}
}

func k8sKubeconfigInstallRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	request := argsI.(*k8sKubeconfigInstallRequest)

	kubeconfigRequest := &k8s.GetClusterKubeConfigRequest{
		Region:    request.Region,
		ClusterID: request.ClusterID,
	}

	client := core.ExtractClient(ctx)
	apiK8s := k8s.NewAPI(client)

	// get the wanted kubeconfig
	apiKubeconfig, err := apiK8s.GetClusterKubeConfig(kubeconfigRequest)
	if err != nil {
		return nil, err
	}
	var kubeconfig api.Config

	err = yaml.Unmarshal(apiKubeconfig.GetRaw(), &kubeconfig)
	if err != nil {
		return nil, err
	}

	kubeconfigPath, err := getKubeconfigPath(ctx)
	if err != nil {
		return nil, err
	}

	// create the kubeconfig file if it does not exist
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		// make sure the directory exists
		err = os.MkdirAll(path.Dir(kubeconfigPath), 0o755)
		if err != nil {
			return nil, err
		}

		// create the file
		f, err := os.OpenFile(kubeconfigPath, os.O_CREATE, 0o600)
		if err != nil {
			return nil, err
		}
		f.Close()
	}

	existingKubeconfig, err := openAndUnmarshalKubeconfig(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// loop through all clusters and insert the wanted one if it does not exist
	clusterFoundInExistingKubeconfig := false
	for _, cluster := range existingKubeconfig.Clusters {
		if cluster.Name == kubeconfig.Clusters[0].Name+"-"+request.ClusterID {
			clusterFoundInExistingKubeconfig = true
			cluster.Cluster = kubeconfig.Clusters[0].Cluster

			break
		}
	}
	if !clusterFoundInExistingKubeconfig {
		existingKubeconfig.Clusters = append(existingKubeconfig.Clusters, api.NamedCluster{
			Name:    kubeconfig.Clusters[0].Name + "-" + request.ClusterID,
			Cluster: kubeconfig.Clusters[0].Cluster,
		})
	}

	// loop through all contexts and insert the wanted one if it does not exist
	contextFoundInExistingKubeconfig := false
	for _, kubeconfigContext := range existingKubeconfig.Contexts {
		if kubeconfigContext.Name == kubeconfig.Contexts[0].Name+"-"+request.ClusterID {
			contextFoundInExistingKubeconfig = true
			kubeconfigContext.Context = api.Context{
				Cluster:  kubeconfig.Clusters[0].Name + "-" + request.ClusterID,
				AuthInfo: kubeconfig.AuthInfos[0].Name,
			}

			break
		}
	}
	if !contextFoundInExistingKubeconfig {
		existingKubeconfig.Contexts = append(existingKubeconfig.Contexts, api.NamedContext{
			Name: kubeconfig.Contexts[0].Name + "-" + request.ClusterID,
			Context: api.Context{
				Cluster:  kubeconfig.Clusters[0].Name + "-" + request.ClusterID,
				AuthInfo: kubeconfig.AuthInfos[0].Name + "-" + request.ClusterID,
			},
		})
	}

	// loop through all users and insert the wanted one if it does not exist
	userFoundInExistingKubeconfig := false
	for _, user := range existingKubeconfig.AuthInfos {
		if user.Name == kubeconfig.AuthInfos[0].Name+"-"+request.ClusterID {
			userFoundInExistingKubeconfig = true
			user.AuthInfo = kubeconfig.AuthInfos[0].AuthInfo

			break
		}
	}
	if !userFoundInExistingKubeconfig {
		existingKubeconfig.AuthInfos = append(existingKubeconfig.AuthInfos, api.NamedAuthInfo{
			Name:     kubeconfig.AuthInfos[0].Name + "-" + request.ClusterID,
			AuthInfo: kubeconfig.AuthInfos[0].AuthInfo,
		})
	}

	// set the current context to the new one
	if !request.KeepCurrentContext {
		existingKubeconfig.CurrentContext = kubeconfig.Contexts[0].Name + "-" + request.ClusterID
	}

	// if it's a new file, set the correct config in the file
	if existingKubeconfig.APIVersion == "" {
		existingKubeconfig.APIVersion = KubeconfigAPIVersion
	}
	if existingKubeconfig.Kind == "" {
		existingKubeconfig.Kind = KubeconfigKind
	}

	err = marshalAndWriteKubeconfig(existingKubeconfig, kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf(
		"Kubeconfig for cluster %s successfully written at %s",
		request.ClusterID,
		kubeconfigPath,
	), nil
}
