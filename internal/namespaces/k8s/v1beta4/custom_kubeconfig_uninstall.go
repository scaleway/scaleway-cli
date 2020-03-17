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
	"gopkg.in/yaml.v2"
)

type k8sKubeconfigUninstallRequest struct {
	ClusterID string
}

func k8sKubeconfigUninstallCommand() *core.Command {
	return &core.Command{
		Short:     `Uninstall a kubeconfig`,
		Long:      `Remove specified cluster from kubeconfig file specified by the KUBECONFIG env, if empty it will default to $HOME/.kube/config.`,
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
	}
}

func k8sKubeconfigUninstallRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	request := argsI.(*k8sKubeconfigUninstallRequest)

	// get the path to write the wanted kubeconfig on disk
	// either the file pointed by the KUBECONFIG env variable (first one in case of a list)
	// or the $HOME/.kube/config
	var kubeconfigPath string
	kubeconfigEnv := core.ExtractEnv(ctx, "KUBECONFIG")
	if kubeconfigEnv != "" {
		if runtime.GOOS == "windows" {
			kubeconfigPath = strings.Split(kubeconfigEnv, ";")[0]
		} else {
			kubeconfigPath = strings.Split(kubeconfigEnv, ":")[0]
		}
	} else {
		homeDir, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		kubeconfigPath = path.Join(homeDir, kubeLocationDir, "config")
	}

	// if the file does not exist, the cluster is not there
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return fmt.Sprintf("File %s does not exists.", kubeconfigPath), nil
	}

	// getting the existing file
	file, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	var existingKubeconfig k8s.Kubeconfig

	err = yaml.Unmarshal(file, &existingKubeconfig)
	if err != nil {
		return nil, err
	}

	// delete the wanted cluster from the file
	newClusters := []*k8s.KubeconfigClusterWithName{}
	for _, cluster := range existingKubeconfig.Clusters {
		if !strings.HasSuffix(cluster.Name, request.ClusterID) {
			newClusters = append(newClusters, cluster)
		}
	}

	// delete the wanted context from the file
	newContexts := []*k8s.KubeconfigContextWithName{}
	for _, kubeconfigContext := range existingKubeconfig.Contexts {
		if !strings.HasSuffix(kubeconfigContext.Name, request.ClusterID) {
			newContexts = append(newContexts, kubeconfigContext)
		}
	}

	// delete the wanted user from the file
	newUsers := []*k8s.KubeconfigUserWithName{}
	for _, user := range existingKubeconfig.Users {
		if !strings.HasSuffix(user.Name, request.ClusterID) {
			newUsers = append(newUsers, user)
		}
	}

	// reset the current context
	existingKubeconfig.CurrentContext = ""

	// write the modification
	existingKubeconfig.Clusters = newClusters
	existingKubeconfig.Contexts = newContexts
	existingKubeconfig.Users = newUsers

	newKubeconfig, err := yaml.Marshal(existingKubeconfig)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(kubeconfigPath, newKubeconfig, 0644)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Cluster %s successfully deleted from %s", request.ClusterID, kubeconfigPath), nil
}
