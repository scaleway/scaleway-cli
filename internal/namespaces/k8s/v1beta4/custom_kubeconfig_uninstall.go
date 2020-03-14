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
			core.RegionArgSpec(),
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
	args := argsI.(*k8sKubeconfigUninstallRequest)

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

	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return fmt.Sprintf("File %s does not exists.", kubeconfigPath), nil
	}

	file, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	var existingKubeconfig k8s.Kubeconfig

	err = yaml.Unmarshal(file, &existingKubeconfig)
	if err != nil {
		return nil, err
	}

	newClusters := []*k8s.KubeconfigClusterWithName{}
	for _, cluster := range existingKubeconfig.Clusters {
		if !strings.HasSuffix(cluster.Name, args.ClusterID) {
			newClusters = append(newClusters, cluster)
		}
	}

	newContexts := []*k8s.KubeconfigContextWithName{}
	for _, context := range existingKubeconfig.Contexts {
		if !strings.HasSuffix(context.Name, args.ClusterID) {
			newContexts = append(newContexts, context)
		}
	}

	newUsers := []*k8s.KubeconfigUserWithName{}
	for _, user := range existingKubeconfig.Users {
		if !strings.HasSuffix(user.Name, args.ClusterID) {
			newUsers = append(newUsers, user)
		}
	}

	existingKubeconfig.CurrentContext = ""
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

	return fmt.Sprintf("Cluster %s successfully deleted from %s", args.ClusterID, kubeconfigPath), nil
}
