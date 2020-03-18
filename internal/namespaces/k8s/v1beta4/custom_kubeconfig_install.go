package k8s

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"gopkg.in/yaml.v2"
)

const (
	kubeLocationDir = ".kube"
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

	// get the path to write the wanted kubeconfig on disk
	// either the file pointed by the KUBECONFIG env variable (first one in case of a list)
	// or the $HOME/.kube/config
	var kubeconfigPath string
	kubeconfigEnv := core.ExtractEnv(ctx, "KUBECONFIG")
	if kubeconfigEnv != "" {
		if runtime.GOOS == "windows" {
			kubeconfigPath = strings.Split(kubeconfigEnv, ";")[0] // list is separated by ; on windows
		} else {
			kubeconfigPath = strings.Split(kubeconfigEnv, ":")[0] // list is separated by : on linux/macos
		}
	} else {
		homeDir, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		kubeconfigPath = path.Join(homeDir, kubeLocationDir, "config")
	}

	// create the kubeconfig file if it does not exist
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		f, err := os.OpenFile(kubeconfigPath, os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		f.Close()
	}

	// reading the file
	file, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// merging the wanted kubeconfig into the opened file

	var existingKubeconfig k8s.Kubeconfig

	err = yaml.Unmarshal(file, &existingKubeconfig)
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
		existingKubeconfig.APIVersion = "v1"
	}
	if existingKubeconfig.Kind == "" {
		existingKubeconfig.Kind = "Config"
	}

	// marshal and write the file
	newKubeconfig, err := yaml.Marshal(existingKubeconfig)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(kubeconfigPath, newKubeconfig, 0644)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Kubeconfig for cluster %s successfully written at %s", request.ClusterID, kubeconfigPath), nil
}
