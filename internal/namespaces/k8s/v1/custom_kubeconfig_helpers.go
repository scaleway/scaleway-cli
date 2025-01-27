package k8s

import (
	"context"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
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
	file, err := os.ReadFile(kubeconfigPath)
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

	return os.WriteFile(kubeconfigPath, newKubeconfig, 0o600)
}
