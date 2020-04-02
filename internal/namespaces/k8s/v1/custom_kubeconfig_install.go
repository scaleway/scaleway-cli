package k8s

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	kubeLocationDir      = ".kube"
	kubeconfigAPIVersion = "v1"
	kubeconfigKind       = "Config"
)

type k8sKubeconfigInstallRequest struct {
	ClusterID string
	Region    scw.Region
}

func k8sKubeconfigInstallCommand() *core.Command {
	return &core.Command{
		Short:     `Install a kubeconfig`,
		Long:      `Retrieve the kubeconfig for a specified cluster and write it on disk. It will merge the new kubeconfig in the file pointed by the KUBECONFIG variable. If empty it will default to $HOME/.kube/config.`,
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
			core.RegionArgSpec(),
		},
		Run: k8sKubeconfigInstallRun,
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
	kubeconfig, err := apiK8s.GetClusterKubeConfig(kubeconfigRequest)
	if err != nil {
		return nil, err
	}

	kubeconfigPath, err := getKubeconfigPath(ctx)
	if err != nil {
		return nil, err
	}

	// create the kubeconfig file if it does not exist
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		f, err := os.OpenFile(kubeconfigPath, os.O_CREATE, 0644)
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
		existingKubeconfig.Clusters = append(existingKubeconfig.Clusters, &k8s.KubeconfigClusterWithName{
			Name:    kubeconfig.Clusters[0].Name + "-" + request.ClusterID,
			Cluster: kubeconfig.Clusters[0].Cluster,
		})
	}

	// loop through all contexts and insert the wanted one if it does not exist
	contextFoundInExistingKubeconfig := false
	for _, kubeconfigContext := range existingKubeconfig.Contexts {
		if kubeconfigContext.Name == kubeconfig.Contexts[0].Name+"-"+request.ClusterID {
			contextFoundInExistingKubeconfig = true
			kubeconfigContext.Context = k8s.KubeconfigContext{
				Cluster: kubeconfig.Clusters[0].Name + "-" + request.ClusterID,
				User:    kubeconfig.Users[0].Name + "-" + request.ClusterID,
			}
			break
		}
	}
	if !contextFoundInExistingKubeconfig {
		existingKubeconfig.Contexts = append(existingKubeconfig.Contexts, &k8s.KubeconfigContextWithName{
			Name: kubeconfig.Contexts[0].Name + "-" + request.ClusterID,
			Context: k8s.KubeconfigContext{
				Cluster: kubeconfig.Clusters[0].Name + "-" + request.ClusterID,
				User:    kubeconfig.Users[0].Name + "-" + request.ClusterID,
			},
		})
	}

	// loop through all users and insert the wanted one if it does not exist
	userFoundInExistingKubeconfig := false
	for _, user := range existingKubeconfig.Users {
		if user.Name == kubeconfig.Users[0].Name+"-"+request.ClusterID {
			userFoundInExistingKubeconfig = true
			user.User = kubeconfig.Users[0].User
			break
		}
	}
	if !userFoundInExistingKubeconfig {
		existingKubeconfig.Users = append(existingKubeconfig.Users, &k8s.KubeconfigUserWithName{
			Name: kubeconfig.Users[0].Name + "-" + request.ClusterID,
			User: kubeconfig.Users[0].User,
		})
	}

	// set the current context to the new one
	existingKubeconfig.CurrentContext = kubeconfig.Contexts[0].Name + "-" + request.ClusterID

	// if it's a new file, set the correct config in the file
	if existingKubeconfig.APIVersion == "" {
		existingKubeconfig.APIVersion = kubeconfigAPIVersion
	}
	if existingKubeconfig.Kind == "" {
		existingKubeconfig.Kind = kubeconfigKind
	}

	err = marshalAndWriteKubeconfig(existingKubeconfig, kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Kubeconfig for cluster %s successfully written at %s", request.ClusterID, kubeconfigPath), nil
}
