package k8s

import (
	"context"
	"io/ioutil"
	"path"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"gopkg.in/yaml.v2"
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
		homeDir, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		kubeconfigPath = path.Join(homeDir, kubeLocationDir, "config")
	}

	return kubeconfigPath, nil
}

func openAndUnmarshalKubeconfig(kubeconfigPath string) (*k8s.Kubeconfig, error) {
	// getting the existing file
	file, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	var kubeconfig k8s.Kubeconfig

	err = yaml.Unmarshal(file, &kubeconfig)
	if err != nil {
		return nil, err
	}

	return &kubeconfig, nil
}

func marshalAndWriteKubeconfig(kubeconfig *k8s.Kubeconfig, kubeconfigPath string) error {
	newKubeconfig, err := yaml.Marshal(*kubeconfig)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(kubeconfigPath, newKubeconfig, 0644)
}
