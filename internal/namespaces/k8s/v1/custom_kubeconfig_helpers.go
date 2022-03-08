package k8s

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"
	api "github.com/kubernetes-client/go-base/config/api"
	"github.com/scaleway/scaleway-cli/internal/core"
)

// get the path to the wanted kubeconfig on disk
// either the file pointed by the KUBECONFIG env variable (first one in case of a list)
// or the $HOME/.kube/config
func getKubeconfigPath(ctx context.Context) (string, error) {
	var kubeconfigPath string
	kubeconfigEnv := core.ExtractEnv(ctx, "KUBECONFIG")
	if kubeconfigEnv != "" {
		if runtime.GOOS == "windows" {
			kubeconfigPath = strings.Split(kubeconfigEnv, ";")[0] // list is separated by ; on windows
		} else {
			kubeconfigPath = strings.Split(kubeconfigEnv, ":")[0] // list is separated by : on linux/macos
		}
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		kubeconfigPath = path.Join(homeDir, kubeLocationDir, "config")
	}

	return kubeconfigPath, nil
}

func openAndUnmarshalKubeconfig(kubeconfigPath string) (*api.Config, error) {
	// getting the existing file
	file, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	var kubeconfig api.Config

	err = yaml.Unmarshal(file, &kubeconfig)
	if err != nil {
		return nil, err
	}

	return &kubeconfig, nil
}

func marshalAndWriteKubeconfig(kubeconfig *api.Config, kubeconfigPath string) error {
	newKubeconfig, err := yaml.Marshal(*kubeconfig)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(kubeconfigPath, newKubeconfig, 0600)
}
