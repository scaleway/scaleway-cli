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
	Region    scw.Region
	ClusterID string
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
			core.RegionArgSpec(),
			{
				Name:       "cluster-id",
				Short:      "Cluster ID from which to retrieve the kubeconfig",
				Required:   true,
				Positional: true,
			},
		},
		Run: k8sKubeconfigInstallRun,
	}
}

func k8sKubeconfigInstallRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*k8sKubeconfigInstallRequest)

	kubeconfigRequest := &k8s.GetClusterKubeConfigRequest{
		Region:    args.Region,
		ClusterID: args.ClusterID,
	}

	client := core.ExtractClient(ctx)
	apiK8s := k8s.NewAPI(client)

	kubeconfig, err := apiK8s.GetClusterKubeConfig(kubeconfigRequest)
	if err != nil {
		return nil, err
	}

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
		f, err := os.OpenFile(kubeconfigPath, os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		f.Close()
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

	found := false
	for _, cluster := range existingKubeconfig.Clusters {
		if cluster.Name == kubeconfig.Clusters[0].Name+"-"+args.ClusterID {
			found = true
			cluster.Cluster = kubeconfig.Clusters[0].Cluster
			break
		}
	}
	if !found {
		existingKubeconfig.Clusters = append(existingKubeconfig.Clusters, &k8s.KubeconfigClusterWithName{
			Name:    kubeconfig.Clusters[0].Name + "-" + args.ClusterID,
			Cluster: kubeconfig.Clusters[0].Cluster,
		})
	}

	found = false
	for _, context := range existingKubeconfig.Contexts {
		if context.Name == kubeconfig.Contexts[0].Name+"-"+args.ClusterID {
			found = true
			context.Context = k8s.KubeconfigContext{
				Cluster: kubeconfig.Clusters[0].Name + "-" + args.ClusterID,
				User:    kubeconfig.Users[0].Name + "-" + args.ClusterID,
			}
			break
		}
	}

	if !found {
		existingKubeconfig.Contexts = append(existingKubeconfig.Contexts, &k8s.KubeconfigContextWithName{
			Name: kubeconfig.Contexts[0].Name + "-" + args.ClusterID,
			Context: k8s.KubeconfigContext{
				Cluster: kubeconfig.Clusters[0].Name + "-" + args.ClusterID,
				User:    kubeconfig.Users[0].Name + "-" + args.ClusterID,
			},
		})
	}

	found = false
	for _, user := range existingKubeconfig.Users {
		if user.Name == kubeconfig.Users[0].Name+"-"+args.ClusterID {
			found = true
			user.User = kubeconfig.Users[0].User
			break
		}
	}

	if !found {
		existingKubeconfig.Users = append(existingKubeconfig.Users, &k8s.KubeconfigUserWithName{
			Name: kubeconfig.Users[0].Name + "-" + args.ClusterID,
			User: kubeconfig.Users[0].User,
		})
	}

	existingKubeconfig.CurrentContext = kubeconfig.Contexts[0].Name + "-" + args.ClusterID
	if existingKubeconfig.APIVersion == "" {
		existingKubeconfig.APIVersion = "v1"
	}
	if existingKubeconfig.Kind == "" {
		existingKubeconfig.Kind = "Config"
	}

	newKubeconfig, err := yaml.Marshal(existingKubeconfig)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(kubeconfigPath, newKubeconfig, 0644)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Kubeconfig for cluster %s successfully written at %s", args.ClusterID, kubeconfigPath), nil
}
